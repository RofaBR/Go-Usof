package handler

import (
	"net/http"

	"github.com/RofaBR/Go-Usof/internal/domain"
	"github.com/RofaBR/Go-Usof/internal/dto/request"
	"github.com/RofaBR/Go-Usof/internal/services"
	"github.com/RofaBR/Go-Usof/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
	log             *logger.Logger
}

func NewCategoryHandler(service *services.CategoryService, log *logger.Logger) *CategoryHandler {
	return &CategoryHandler{categoryService: service, log: log}
}

func (h *CategoryHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	h.log.Info("handling category create")

	var req request.CreateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := &domain.Category{
		Title: req.Title,
		Desc:  req.Desc,
		Slug:  slug.Make(req.Title),
	}

	err := h.categoryService.Create(ctx, category)
	if err != nil {
		h.log.Warn("error creating category", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	h.log.Info("category created", "title", req.Title, "desc", req.Desc)
	c.JSON(http.StatusOK, gin.H{
		"message": "category created",
	})
}
