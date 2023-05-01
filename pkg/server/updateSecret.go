package server

import (
	"custodyEthereum/internal/models"
	"encoding/base64"
	"encoding/json"
	"github.com/awnumar/memguard"
	"github.com/gin-gonic/gin"
)

func (s *Server) UpdateSecret() gin.HandlerFunc {
	return func(c *gin.Context) {

		var req models.ReqStoreAction

		reqData, err := c.GetRawData()
		if err != nil {
			s.Logger.Error("Error in getting request data")
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in getting request data"})
			return
		}

		err = json.Unmarshal(reqData, &req)
		if err != nil {
			s.Logger.Error("Error in unmarshalling request data")
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in unmarshalling request data"})
			return
		}

		userRole := c.GetString("role")
		if userRole == "" {
			c.AbortWithStatusJSON(400, gin.H{"message": "User role not found"})
			return
		}

		store, err := s.preliminaryRequestChecks(req.StoreName, req.SecretName, userRole, false)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"message": err.Error()})
			return
		}

		path := req.StoreName + "/" + req.SecretName

		encodedData, err := base64.StdEncoding.DecodeString(req.SecretData)
		if err != nil {
			s.Logger.Error("Error in decoding secret data")
			c.AbortWithStatusJSON(400, gin.H{"message": "Error in decoding secret data"})
			return
		}

		action := models.StoreAction{
			Action:      "update",
			SecretPath:  path,
			SecretData:  memguard.NewBufferFromBytes(encodedData),
			Description: req.Description,
			Role:        userRole,
		}

		store.ActionChannel <- &action

	}
}
