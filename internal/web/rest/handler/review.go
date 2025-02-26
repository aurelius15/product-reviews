package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/aurelius15/product-reviews/internal/service"
	"github.com/aurelius15/product-reviews/internal/web/rest/apimodel"
)

type ReviewHandler struct {
	ser *service.APIService
}

func NewReviewHandler(ser *service.APIService) *ReviewHandler {
	return &ReviewHandler{
		ser: ser,
	}
}

func (h *ReviewHandler) Retrieve(c *gin.Context) {
	id, err := h.parseReviewID(c)
	if err != nil {
		return
	}

	r, err := h.ser.RetrieveReview(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, r)
}

func (h *ReviewHandler) Create(c *gin.Context) {
	var review apimodel.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	r, err := h.ser.SaveReview(&review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, r)
}

func (h *ReviewHandler) Update(c *gin.Context) {
	id, err := h.parseReviewID(c)
	if err != nil {
		return
	}

	var review apimodel.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review.ID = id

	r, err := h.ser.SaveReview(&review)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, r)
}

func (h *ReviewHandler) Delete(c *gin.Context) {
	id, err := h.parseReviewID(c)
	if err != nil {
		return
	}

	if serErr := h.ser.DeleteReview(id); serErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": serErr.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (_ *ReviewHandler) parseReviewID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return 0, err
	}

	return id, nil
}
