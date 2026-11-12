package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/Baroneee/EBVN_Assignment_1/internal/service"
	"github.com/gin-gonic/gin"
)

// ShortenLink handles HTTP requests for creating shortened links.
type ShortenLink interface {
	CreateShorternLink(c *gin.Context)
	Redirect(c *gin.Context)
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
	URL        string `json:"url" binding:"required,url"`
	ExpireTime int    `json:"expireTime" binding:"gte=0,lte=604800"` 
}

// @Summary Generate shortern link
// @Description Generate shortern link
// @Tags link
// @Accept application/json
// @Produce application/jsons
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

	code, err := h.svc.CreateShorternLink(c, input.URL, time.Duration(input.ExpireTime)*time.Second)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": code})
}


// Redirect retrieves the original URL for a given code and redirects the client to that URL.
// @Tags link
// @Accept application/json
// @Produce application/json
// @Param code path string true "Shortened code"
// @Success 302
// @Router /v1/links/redirect/{code} [get]
func (h *shortenURL) Redirect(c *gin.Context) {
	code := c.Param("code")

	if code == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}

	url, err := h.svc.GetUrlFromCode(c, code)
	if err != nil {
		if errors.Is(err, service.ErrUrlDoesNotExist) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Code Not Exist"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}