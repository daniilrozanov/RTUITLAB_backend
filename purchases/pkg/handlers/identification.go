package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	userIdCtx = "userId"
)

func (h *Handler) identifyUser (c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerPaths := strings.Split(header, " ")
	if len(headerPaths) != 2 || headerPaths[0] != "Bearer" {
		newErrorResponse(c, http.StatusUnauthorized, "bad auth header")
	}
	if len(headerPaths[1]) == 0 {
		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}
	userId, err := h.buis.ParseToken(headerPaths[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userIdCtx, userId)
}

func (h *Handler) getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userIdCtx)
	if !ok {
		return 0, errors.New("user not found")
	}
	intId , ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}
	return intId, nil
}