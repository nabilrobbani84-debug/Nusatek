package http

import (
	"net/http"
	"strconv"

	"nusatek-backend/internal/domain"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	CUsecase domain.CustomerUsecase
}

func NewCustomerHandler(r *gin.Engine, us domain.CustomerUsecase) {
	handler := &CustomerHandler{
		CUsecase: us,
	}
	r.GET("/api/v1/customers", handler.Fetch)
	r.POST("/api/v1/customers", handler.Store)
	r.GET("/api/v1/customers/:id", handler.GetByID)
	r.DELETE("/api/v1/customers/:id", handler.Delete)
}

func (h *CustomerHandler) Fetch(c *gin.Context) {
	limit := 10
	offset := 0
	
	list, err := h.CUsecase.Fetch(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
		return
	}

	cust, err := h.CUsecase.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *CustomerHandler) Store(c *gin.Context) {
	var cust domain.Customer
	if err := c.BindJSON(&cust); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.CUsecase.Store(c.Request.Context(), &cust); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cust)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.CUsecase.Delete(c.Request.Context(), int64(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}
