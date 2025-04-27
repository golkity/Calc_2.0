package handler

import (
	"auth_service/internal/models"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"auth_service/internal/db"
	"auth_service/internal/kafka"
	"auth_service/internal/tokens"
)

type Handler struct {
	store     *db.Storage
	tm        *tokens.Manager
	refreshDB *tokens.Store
	kafka     *kafka.Producer
}

func New(store *db.Storage, tm *tokens.Manager, rdb *tokens.Store, k *kafka.Producer) *Handler {
	return &Handler{store, tm, rdb, k}
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	u := models.User{Email: c.Email, Password: string(hash)}
	if err := h.store.CreateUser(r.Context(), &u); err != nil {
		http.Error(w, err.Error(), 409)
		return
	}

	h.kafka.PublishUserCreated(r.Context(), u.ID, u.Email)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var c credentials
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	u, err := h.store.GetUserByEmail(r.Context(), c.Email)
	if err != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Password)) != nil {
		http.Error(w, "invalid credentials", 401)
		return
	}

	tks, err := h.tm.Generate(u.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_ = h.refreshDB.Save(r.Context(), tks.RefreshToken, time.Duration(tks.RefreshExpiresIn)*time.Second)
	_ = h.kafka.PublishLogin(r.Context(), u.ID)

	json.NewEncoder(w).Encode(tks)
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Refresh string `json:"refresh_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	cl, err := h.tm.Parse(body.Refresh)
	if err != nil {
		http.Error(w, "invalid refresh", 401)
		return
	}

	ok, _ := h.refreshDB.Exists(r.Context(), body.Refresh)
	if !ok {
		http.Error(w, "token revoked", 401)
		return
	}

	_ = h.refreshDB.Delete(r.Context(), body.Refresh)

	tks, err := h.tm.Generate(cl.UserID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_ = h.refreshDB.Save(r.Context(), tks.RefreshToken, time.Duration(tks.RefreshExpiresIn)*time.Second)
	json.NewEncoder(w).Encode(tks)
}
