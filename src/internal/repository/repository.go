package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, s *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, q model.CostQuery) (int, error)
}

type subscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, s *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES (:id, :service_name, :price, :user_id, :start_date, :end_date)
		RETURNING created_at, updated_at
	`

	startDate, err := normalizeDate(s.StartDate)
	if err != nil {
		return err
	}
	endDate, err := normalizeDate(*s.EndDate)
	if err != nil {
		return err
	}

	s.StartDate = startDate
	s.EndDate = &endDate

	rows, err := r.db.NamedQueryContext(ctx, query, s)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.Scan(&s.CreatedAt, &s.UpdatedAt)
	}
	return nil
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	var s model.Subscription
	err := r.db.GetContext(ctx, &s, "SELECT * FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepo) GetAll(ctx context.Context) ([]model.Subscription, error) {
	var subs []model.Subscription
	err := r.db.SelectContext(ctx, &subs, "SELECT * FROM subscriptions ORDER BY created_at DESC")
	return subs, err
}

func (r *subscriptionRepo) Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	setClauses := []string{}
	args := map[string]interface{}{"id": id}

	if req.ServiceName != nil {
		setClauses = append(setClauses, "service_name = :service_name")
		args["service_name"] = *req.ServiceName
	}
	if req.Price != nil {
		setClauses = append(setClauses, "price = :price")
		args["price"] = *req.Price
	}
	if req.StartDate != nil {
		setClauses = append(setClauses, "start_date = :start_date")
		startDate, err := normalizeDate(*req.StartDate)
		if err != nil {
			return nil, err
		}
		args["start_date"] = startDate
	}
	if req.EndDate != nil {
		endDate, err := normalizeDate(*req.EndDate)
		if err != nil {
			return nil, err
		}
		setClauses = append(setClauses, "end_date = :end_date")
		args["end_date"] = endDate
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	query := "UPDATE subscriptions SET "
	for i, c := range setClauses {
		if i > 0 {
			query += ", "
		}
		query += c
	}
	query += ", updated_at = NOW() WHERE id = :id"

	_, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subscription not found")
	}
	return nil
}

func (r *subscriptionRepo) GetTotalCost(ctx context.Context, q model.CostQuery) (int, error) {
	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date >= $1::date
		  AND (end_date IS NULL OR end_date <= $2::date)`

	from, err := normalizeDate(q.From)
	if err != nil {
		return 0, err
	}
	to, err := normalizeDate(q.To)
	if err != nil {
		return 0, err
	}
	args := []interface{}{from, to}
	argIdx := 3

	if q.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, uuid.MustParse(*q.UserID))
		argIdx++
	}
	if q.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, *q.ServiceName)
	}

	var total int
	err = r.db.GetContext(ctx, &total, query, args...)
	return total, err
}

func normalizeDate(input string) (string, error) {
	t, err := time.Parse("01-2006", input)
	if err != nil {
		return "", err
	}
	return t.Format("2006-01-02"), nil
}
