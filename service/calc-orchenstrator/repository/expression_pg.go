package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Expression struct {
	ID     int64    `json:"id"`
	UserID int64    `json:"user_id"`
	Raw    string   `json:"raw"`
	Result *float64 `json:"result,omitempty"`
	Status string   `json:"status"`
}

type ExpressionRepo struct{ db *pgx.Conn }

func NewExpr(db *pgx.Conn) *ExpressionRepo { return &ExpressionRepo{db} }

func (r *ExpressionRepo) Create(ctx context.Context, userID int64, raw string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		`insert into expressions(user_id,raw,status) values ($1,$2,'pending') returning id`,
		userID, raw).Scan(&id)
	return id, err
}

func (r *ExpressionRepo) Done(ctx context.Context, id int64, res float64) error {
	_, err := r.db.Exec(ctx,
		`update expressions set result=$1, status='done' where id=$2`, res, id)
	return err
}

func (r *ExpressionRepo) Get(ctx context.Context, id int64) (Expression, error) {
	var e Expression
	err := r.db.QueryRow(ctx,
		`select id, user_id, raw, result, status from expressions where id=$1`, id).
		Scan(&e.ID, &e.UserID, &e.Raw, &e.Result, &e.Status)
	return e, err
}
