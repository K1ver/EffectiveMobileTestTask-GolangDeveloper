package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ServiceName string    `json:"service_name" db:"service_name" binding:"required"`
	Price       int       `json:"price" db:"price" binding:"required,gte=0"`
	UserID      uuid.UUID `json:"user_id" db:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" db:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty" db:"end_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSubscriptionRequest struct {
	ServiceName string    `json:"service_name" binding:"required"`
	Price       int       `json:"price" binding:"required,gte=0"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	StartDate   string    `json:"start_date" binding:"required"`
	EndDate     *string   `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int    `json:"price,omitempty" binding:"omitempty,gte=0"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type CostQuery struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	From        string  `form:"from" binding:"required"`
	To          string  `form:"to" binding:"required"`
}

func (q *CostQuery) ParseUserID() (*uuid.UUID, error) {
	if q.UserID == nil || *q.UserID == "" {
		return nil, nil
	}
	uid, err := uuid.Parse(*q.UserID)
	if err != nil {
		return nil, err
	}
	return &uid, nil
}

type CostResponse struct {
	TotalCost int `json:"total_cost"`
}
