package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TaskRepo struct{ db *pgx.Conn }

func NewTaskRepo(db *pgx.Conn) *TaskRepo { return &TaskRepo{db} }

func (r *TaskRepo) Create(ctx context.Context, exprID int64, arg1, arg2 float64, op string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		insert into tasks (expression_id,arg1,arg2,op,status) 
		values ($1,$2,$3,$4,'pending') returning id`,
		exprID, arg1, arg2, op).Scan(&id)
	return id, err
}
