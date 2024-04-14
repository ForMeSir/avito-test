package handler

import (
	"github.com/gin-gonic/gin"
)

type SecondError struct {
	Description string `json:"description"`
}

type Error struct{
	Description string `json:"description"`
	Error       string `json:"error"`
}

func ErrorResponse(c *gin.Context, statusCode int, description string, err string){
  c.AbortWithStatusJSON(statusCode, Error{Description: description, Error: err})
}

func newErrorResponse(c *gin.Context, statusCode int, description string){
  c.AbortWithStatusJSON(statusCode, SecondError{Description: description})
}