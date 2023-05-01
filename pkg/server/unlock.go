package server

import (
	"custodyEthereum/internal/models"
	"custodyEthereum/pkg/encryptedStore"
	"encoding/base64"
	"encoding/json"
	"github.com/awnumar/memguard"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/vault/shamir"
	"go.uber.org/zap"
	"log"
)

func (s *Server) Unlock(reload bool) gin.HandlerFunc {
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

		if s.AvailableStores[req.StoreName] && !reload {
			c.AbortWithStatusJSON(400, gin.H{"message": "Store already opened"})
			return
		}

		keyShares := make([][]byte, len(req.Shares))
		for _, share := range req.Shares {
			if share == "" {
				c.AbortWithStatusJSON(400, gin.H{"message": "Share cannot be empty"})
				return
			}
			k, err := base64.StdEncoding.DecodeString(share)
			if err != nil {
				s.Logger.Error("Error in decoding share", zap.Error(err))
				c.AbortWithStatusJSON(400, gin.H{"message": "Error in decoding share"})
				return
			}
			keyShares = append(keyShares, k)
		}

		recoveredKey, err := shamir.Combine(keyShares)
		if err != nil {
			s.Logger.Error("Error in combining shares", zap.Error(err))
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in combining shares"})
			return
		}

		var usableKey [models.KeySize]byte
		copy(usableKey[:], recoveredKey)

		memguard.WipeBytes(recoveredKey)

		store, err := encryptedStore.OpenBoxStore(req.StoreName, usableKey)

		log.Println("Decrypted Store: ", store.ID)

		if reload {
			s.removeStoreFromServer(s.Stores[req.StoreName])
		}

		_ = s.importStoreWithACL(store, usableKey)

		memguard.WipeBytes(usableKey[:])

		c.JSON(200, gin.H{"message": "Store unlocked successfully", "store": store.ID})
		return
	}
}

func (s *Server) UpdateStore(storeName string, override bool) {
	s.UpdateStoreChan <- &models.UpdateStoreAction{
		StoreName: storeName,
		Override:  override,
	}
}
