package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aurelius15/product-reviews/internal/service"
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
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	p, err := h.ser.RetrieveProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (_ *ProductHandler) Create(_ *gin.Context) {

}

func (_ *ProductHandler) Update(_ *gin.Context) {

}

func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	if serErr := h.ser.DeleteProduct(id); serErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": serErr.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
