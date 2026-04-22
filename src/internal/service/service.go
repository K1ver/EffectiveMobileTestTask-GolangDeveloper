package service

import (
	"context"

	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/model"
	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/repository"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type SubscriptionService interface {
	Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	GetAll(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTotalCost(ctx context.Context, q model.CostQuery) (*model.CostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	sub := &model.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	log.WithField("subscription_id", sub.ID).Info("Creating subscription")
	if err := s.repo.Create(ctx, sub); err != nil {
		log.WithError(err).Error("Failed to create subscription")
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	log.WithField("id", id).Info("Getting subscription by ID")
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Failed to get subscription")
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) GetAll(ctx context.Context) ([]model.Subscription, error) {
	log.Info("Getting all subscriptions")
	return s.repo.GetAll(ctx)
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	log.WithField("id", id).Info("Updating subscription")
	sub, err := s.repo.Update(ctx, id, req)
	if err != nil {
		log.WithError(err).WithField("id", id).Error("Failed to update subscription")
		return nil, err
	}
	return sub, nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	log.WithField("id", id).Info("Deleting subscription")
	if err := s.repo.Delete(ctx, id); err != nil {
		log.WithError(err).WithField("id", id).Error("Failed to delete subscription")
		return err
	}
	return nil
}

func (s *subscriptionService) GetTotalCost(ctx context.Context, q model.CostQuery) (*model.CostResponse, error) {
	log.Info("Calculating total cost")
	total, err := s.repo.GetTotalCost(ctx, q)
	if err != nil {
		log.WithError(err).Error("Failed to calculate total cost")
		return nil, err
	}
	return &model.CostResponse{TotalCost: total}, nil
}
