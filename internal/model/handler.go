package model

import "github.com/gin-gonic/gin"

type Handler interface {
	Mount(router *gin.RouterGroup)
}

type HandlerError struct {
	Error string `json:"error"`
}

type HandlerMessage struct {
	Message string `json:"message"`
}
