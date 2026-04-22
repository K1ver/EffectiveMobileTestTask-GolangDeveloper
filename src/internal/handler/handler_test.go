package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service
type MockService struct {
	mock.Mock
}

func (m *MockService) Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockService) GetAll(ctx context.Context) ([]model.Subscription, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func (m *MockService) Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) (*model.Subscription, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(*model.Subscription), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) GetTotalCost(ctx context.Context, q model.CostQuery) (*model.CostResponse, error) {
	args := m.Called(ctx, q)
	return args.Get(0).(*model.CostResponse), args.Error(1)
}

func setupRouter(h *Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h.RegisterRoutes(r)
	return r
}

func TestCreate(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	userID := uuid.New()
	subID := uuid.New()

	req := model.CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	expectedSub := &model.Subscription{
		ID:          subID,
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      userID,
		StartDate:   "07-2025",
	}

	mockSvc.On("Create", mock.Anything, req).Return(expectedSub, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.Subscription
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Yandex Plus", resp.ServiceName)
	assert.Equal(t, 400, resp.Price)
	mockSvc.AssertExpectations(t)
}

func TestGetByID(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	id := uuid.New()
	expectedSub := &model.Subscription{
		ID:          id,
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   "01-2025",
	}

	mockSvc.On("GetByID", mock.Anything, id).Return(expectedSub, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/subscriptions/"+id.String(), nil)
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetAll(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	subs := []model.Subscription{
		{ID: uuid.New(), ServiceName: "Spotify", Price: 199},
		{ID: uuid.New(), ServiceName: "Netflix", Price: 999},
	}

	mockSvc.On("GetAll", mock.Anything).Return(subs, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/subscriptions", nil)
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []model.Subscription
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	mockSvc.AssertExpectations(t)
}

func TestDelete(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	id := uuid.New()
	mockSvc.On("Delete", mock.Anything, id).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/v1/subscriptions/"+id.String(), nil)
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetTotalCost(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	expected := &model.CostResponse{TotalCost: 1598}

	mockSvc.On("GetTotalCost", mock.Anything, mock.AnythingOfType("model.CostQuery")).Return(expected, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/subscriptions/cost?from=01-2025&to=12-2025", nil)
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp model.CostResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 1598, resp.TotalCost)
	mockSvc.AssertExpectations(t)
}

func TestCreateInvalidBody(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/v1/subscriptions", bytes.NewBufferString(`{}`))
	httpReq.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByIDInvalidUUID(t *testing.T) {
	mockSvc := new(MockService)
	h := NewHandler(mockSvc)
	r := setupRouter(h)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/v1/subscriptions/not-a-uuid", nil)
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
