package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	// "net/http"
	// "strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(h.logger, c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(h.logger, c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	h.logger.Info("-------------------------------------------", headerParts)

	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(h.logger, c, http.StatusUnauthorized, err.Error())
	}
	h.logger.Info("-------------------------------------------", userId)
	c.Set(userCtx, userId)
}
