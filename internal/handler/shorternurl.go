package handler

import (
	"net/http"

	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	"github.com/gin-gonic/gin"
)

// ShortenLink handles HTTP requests for creating shortened links.
type ShortenLink interface {
	CreateShorternLink(c *gin.Context)
}

type shortenURL struct {
	svc service.ShortenURL
}

// NewShortenURL creates a shortened-link handler backed by svc.
func NewShortenURL(svc service.ShortenURL) ShortenLink {
	return &shortenURL{svc: svc}
}

// shortenLinkInput contains the URL submitted for shortening.
type shortenLinkInput struct {
	URL string `json:"url"`
}

// @Summary Generate shortern link
// @Description Generate shortern link
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param input body shortenLinkInput true "Input required"
// @Success 200 {object} map[string]string
// @Router /v1/links/shortern [post]
// CreateShorternLink generates a shortened link and writes the result as JSON.
func (h *shortenURL) CreateShorternLink(c *gin.Context) {
	input := &shortenLinkInput{}

	if err := c.ShouldBindBodyWithJSON(input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	code, err := h.svc.CreateShorternLink(c, input.URL)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}
