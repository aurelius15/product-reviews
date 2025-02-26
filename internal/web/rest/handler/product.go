package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aurelius15/product-reviews/internal/service"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type ProductHandler struct {
	ser *service.APIService
}

func NewProductHandler(ser *service.APIService) *ProductHandler {
	return &ProductHandler{
		ser: ser,
	}
}

func (h *ProductHandler) Retrieve(c *gin.Context) {
	id, err := h.parseProductID(c)
	if err != nil {
		return
	}

	p, err := h.ser.RetrieveProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product apimodel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	p, err := h.ser.SaveProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id, err := h.parseProductID(c)
	if err != nil {
		return
	}

	var product apimodel.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	product.ID = id

	p, err := h.ser.SaveProduct(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := h.parseProductID(c)
	if err != nil {
		return
	}

	if serErr := h.ser.DeleteProduct(id); serErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": serErr.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ProductHandler) parseProductID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return 0, err
	}

	return id, nil
}
