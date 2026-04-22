package handler

import (
	"net/http"

	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/model"
	"github.com/K1ver/EffectiveMobileTestTask-GolangDeveloper/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	svc service.SubscriptionService
}

func NewHandler(svc service.SubscriptionService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.POST("/subscriptions", h.Create)
		api.GET("/subscriptions", h.GetAll)
		api.GET("/subscriptions/:id", h.GetByID)
		api.PUT("/subscriptions/:id", h.Update)
		api.DELETE("/subscriptions/:id", h.Delete)
		api.GET("/subscriptions/cost", h.GetTotalCost)
	}
}

// Create godoc
// @Summary Create subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param input body model.CreateSubscriptionRequest true "Subscription"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /api/v1/subscriptions [post]
func (h *Handler) Create(c *gin.Context) {
	var req model.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Warn("Invalid create request")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

// GetAll godoc
// @Summary List all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Router /api/v1/subscriptions [get]
func (h *Handler) GetAll(c *gin.Context) {
	subs, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// GetByID godoc
// @Summary Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Failure 404 {object} map[string]string
// @Router /api/v1/subscriptions/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	sub, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Update godoc
// @Summary Update subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param input body model.UpdateSubscriptionRequest true "Fields to update"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Router /api/v1/subscriptions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// Delete godoc
// @Summary Delete subscription
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/subscriptions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetTotalCost godoc
// @Summary Get total cost of subscriptions
// @Tags subscriptions
// @Produce json
// @Param from query string true "Start period (MM-YYYY)"
// @Param to query string true "End period (MM-YYYY)"
// @Param user_id query string false "User ID (UUID)"
// @Param service_name query string false "Service name"
// @Success 200 {object} model.CostResponse
// @Router /api/v1/subscriptions/cost [get]
func (h *Handler) GetTotalCost(c *gin.Context) {
	var q model.CostQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.svc.GetTotalCost(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
