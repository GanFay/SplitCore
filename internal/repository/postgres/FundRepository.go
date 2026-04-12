package postgres

import (
	"SplitCore/internal/domain"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FundRepository struct {
	DB *pgxpool.Pool
}

func NewFundRepository(pool *pgxpool.Pool) *FundRepository {
	slog.Info("init fundRepository")

	return &FundRepository{DB: pool}
}

func (r *FundRepository) Create(ctx context.Context, fund *domain.Fund) (*domain.Fund, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			slog.Debug("rollback", "err", err)
			return
		}
	}(tx, ctx)

	err = r.DB.QueryRow(ctx, `INSERT INTO funds
    (name, author_id, invite_code) 
	VALUES ($1, $2, $3) 
	ON CONFLICT DO NOTHING
	RETURNING id, created_at`, fund.Name, fund.AuthorID, fund.InviteCode).Scan(&fund.ID, &fund.CreatedAt)
	if err != nil {
		return nil, err
	}
	queryMember := `INSERT INTO fund_members (fund_id, user_id) VALUES ($1, $2)`
	_, err = tx.Exec(ctx, queryMember, fund.ID, fund.AuthorID)
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}
	return fund, err
}

func (r *FundRepository) GetByInviteCode(ctx context.Context, code string) (*domain.Fund, error) {
	var fund domain.Fund
	err := r.DB.QueryRow(ctx, `SELECT id, name, author_id, invite_code, created_at FROM funds WHERE invite_code = $1`, code).Scan(
		&fund.ID, &fund.Name, &fund.AuthorID, &fund.InviteCode, &fund.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &fund, nil
}

func (r *FundRepository) GetByUserID(ctx context.Context, userID int64, limit string, offset string) ([]domain.Fund, error) {
	var funds []domain.Fund
	query, err := r.DB.Query(ctx, `SELECT id, name, author_id, invite_code, created_at FROM funds WHERE author_id = $1 LIMIT $2 OFFSET $3`, userID, limit, offset)
	if err != nil {
		slog.Debug(err.Error())
		return nil, err
	}
	for query.Next() {
		var fund domain.Fund
		err = query.Scan(&fund.ID, &fund.Name, &fund.AuthorID, &fund.InviteCode, &fund.CreatedAt)
		if err != nil {
			return nil, err
		}
		funds = append(funds, fund)
	}
	defer query.Close()
	return funds, nil
}

func (r *FundRepository) AddMember(ctx context.Context, fund *domain.Fund, userID int64) error {
	queryMember := `INSERT INTO fund_members (fund_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.DB.Exec(ctx, queryMember, fund.ID, userID)
	return err
}
