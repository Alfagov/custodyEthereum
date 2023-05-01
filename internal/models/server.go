package models

import (
	"github.com/gin-gonic/gin"
)

type ServerInterface interface {
	NewStore() gin.HandlerFunc
	Unlock(bool) gin.HandlerFunc
	AddSecret() gin.HandlerFunc
	RemoveSecret() gin.HandlerFunc
	UpdateSecret() gin.HandlerFunc
	UpdateStore(storeName string, override bool)
}
