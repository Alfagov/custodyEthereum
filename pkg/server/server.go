package server

import (
	"custodyEthereum/internal/models"
	"custodyEthereum/pkg/encryptedStore"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/vault/shamir"
	"go.uber.org/zap"
	"log"
)

type Server struct {
	AvailableStores map[string]bool `json:"availableStores"`
	Logger          *zap.Logger
}

func (s *Server) Initialize() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.InitializeRequest

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

		if s.AvailableStores[req.StoreName] {
			c.AbortWithStatusJSON(400, gin.H{"message": "Store already exists"})
			return
		}

		shares := encryptedStore.CreateNewStore(req.StoreName, req.Threshold, req.Total)

		log.Println("Shares: ", shares)
		c.JSON(200, gin.H{"message": "Store created successfully", "shares": shares})
		return
	}
}

func (s *Server) Unlock() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.UnlockRequest

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

		keyShares := make([][]byte, len(req.Shares))
		for _, share := range req.Shares {
			if share == "" {
				c.AbortWithStatusJSON(400, gin.H{"message": "Share cannot be empty"})
				return
			}
			keyShares = append(keyShares, []byte(share))
		}

		recoveredKey, err := shamir.Combine(keyShares)
		if err != nil {
			s.Logger.Error("Error in combining shares", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in combining shares"})
			return
		}

		store := encryptedStore.LoadEncryptedKVStore(req.StoreName)

		decryptedStore := encryptedStore.DecryptKVStore(store, recoveredKey)

		log.Println("Decrypted Store: ", decryptedStore)

		c.JSON(200, gin.H{"message": "Store unlocked successfully", "store": decryptedStore})
		return
	}
}
