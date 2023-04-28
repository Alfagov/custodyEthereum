package jwtHelper

import (
	"custodyEthereum/configs"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type AccessTokenCustomClaims struct {
	Id      string `json:"id"`
	User    bool   `json:"user"`
	keyType string
	Role    string
	jwt.StandardClaims
}

var (
	rootRole     = "root"
	adminRole    = "admin"
	allowedRoles = []string{rootRole, adminRole}
)

func GenerateAccessToken(isUser bool, role string) (string, error) {

	tokenType := "access"
	cfg := configs.GlobalViper

	claims := &AccessTokenCustomClaims{ // PayLoads
		"",
		isUser,
		tokenType,
		role,
		jwt.StandardClaims{
			// 3 month
			ExpiresAt: time.Now().Add(time.Hour * 256).Unix(),
			//ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			//ExpiresAt: time.Now().Add(time.Second * 20).Unix(),
			Issuer:   "custodyEthereum",
			IssuedAt: time.Now().Unix(), // ! Can use to reject in case of certain incident or attack
		},
	}
	//fmt.printf("Reading File\n")
	signBytes, err := os.ReadFile(cfg.GetString("jwt.private_key"))
	if err != nil {
		return "", err
	}
	//fmt.printf("File Read\n")
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	return token.SignedString(signKey)
}
