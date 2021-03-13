package handler

import (
	"github.com/gin-gonic/gin"
	"log"
)

type Error struct {
	Msg string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string){
	log.Print(message)
	c.AbortWithStatusJSON(statusCode, Error{message})
}
