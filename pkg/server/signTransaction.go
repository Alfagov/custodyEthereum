package server

import (
	"custodyEthereum/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (s *Server) SignTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.RequestSignTransaction

		reqData, err := c.GetRawData()
		if err != nil {
			s.Logger.Error("Error in getting request data", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in getting request data"})
			return
		}

		err = json.Unmarshal(reqData, &req)
		if err != nil {
			s.Logger.Error("Error in unmarshalling request data", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in unmarshalling request data"})
			return
		}

		if s.AvailableStores[req.TransactionStoreName] {
			c.AbortWithStatusJSON(400, gin.H{"message": "Store already opened"})
			return
		}

		store := s.Stores[req.TransactionStoreName]

		return
	}
}
