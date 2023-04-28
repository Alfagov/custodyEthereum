package server

import (
	logger "custodyEthereum/internal/logger"
	"custodyEthereum/internal/models"
	"custodyEthereum/pkg/encryptedStore"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

type Server struct {
	AvailableStores map[string]bool `json:"availableStores"`
	Stores          map[string]*models.ServerStore
	Users           []models.User                       `json:"users"`
	Roles           []string                            `json:"roles"`
	Secrets         []models.Secret                     `json:"secrets"`
	RoleSecrets     map[string]map[string]models.Secret `json:"roleSecrets"`
	ACL             map[string]map[string]models.Secret
	Logger          *zap.Logger
}

func NewServer() *Server {
	lo, _ := logger.NewLogger("server_logger.log")

	return &Server{
		AvailableStores: make(map[string]bool),
		Stores:          make(map[string]*models.ServerStore),
		Users:           []models.User{},
		Roles:           []string{"root", "admin", "user"},
		RoleSecrets:     make(map[string]map[string]models.Secret),
		ACL:             make(map[string]map[string]models.Secret),
		Logger:          lo,
	}
}

func (s *Server) NewStore() gin.HandlerFunc {
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

		shares := encryptedStore.CreateNewBoxStore(req.StoreName, req.Threshold, req.Total, "")

		log.Println("Shares: ", shares)
		c.JSON(200, gin.H{"message": "Store created successfully", "shares": shares})
		return
	}
}
