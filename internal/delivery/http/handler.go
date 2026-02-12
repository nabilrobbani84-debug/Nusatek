package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"nusatek-backend/internal/domain"
)

type PropertyHandler struct {
	PropertyUsecase domain.PropertyUsecase
}

func NewPropertyHandler(r *gin.Engine, us domain.PropertyUsecase) {
	handler := &PropertyHandler{
		PropertyUsecase: us,
	}

	api := r.Group("/api/v1")
	{
		api.GET("/properties", handler.Fetch)
		api.GET("/properties/:id", handler.GetByID)
		api.POST("/properties", handler.Store)
		api.PUT("/properties/:id", handler.Update)
		api.DELETE("/properties/:id", handler.Delete)
	}
}

func (h *PropertyHandler) Fetch(c *gin.Context) {
	limit := 10
	offset := 0

	// Basic query param parsing
	if l, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = l
	}
	if o, err := strconv.Atoi(c.Query("offset")); err == nil {
		offset = o
	}

	properties, err := h.PropertyUsecase.Fetch(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, properties)
}

func (h *PropertyHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	property, err := h.PropertyUsecase.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Property not found"})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) Store(c *gin.Context) {
	var property domain.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.PropertyUsecase.Store(c.Request.Context(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, property)
}

func (h *PropertyHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var property domain.Property
	if err := c.ShouldBindJSON(&property); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	property.ID = int64(id)

	if err := h.PropertyUsecase.Update(c.Request.Context(), &property); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, property)
}

func (h *PropertyHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.PropertyUsecase.Delete(c.Request.Context(), int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
