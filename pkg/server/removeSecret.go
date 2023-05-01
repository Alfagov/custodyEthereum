package server

import (
	"custodyEthereum/internal/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func (s *Server) RemoveSecret() gin.HandlerFunc {
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

		action := models.StoreAction{
			Action:     "remove",
			SecretPath: path,
			Role:       userRole,
		}

		store.ActionChannel <- &action

		s.Mux.Lock()
		for _, r := range s.Roles {
			delete(s.RoleSecrets[r], path)
		}
		s.Mux.Unlock()
	}
}
