package jwtHelper

import (
	"custodyEthereum/configs"
	"errors"
	"github.com/golang-jwt/jwt"
	"os"
)

func ValidateAccessToken(tokenString string) (string, string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method in auth token")
		}

		//fmt.println("IS OK! ")
		verifyBytes, err := os.ReadFile(configs.GlobalViper.GetString("jwt.public_key"))
		if err != nil {
			return nil, err
		}
		//fmt.println("verifyBytes : ", verifyBytes)

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			return nil, err
		}
		//fmt.println("verifyKey : ", verifyKey)

		return verifyKey, nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*AccessTokenCustomClaims)
	if !ok || !token.Valid || claims.Id == "" {
		//fmt.println("ok : ", ok)
		//fmt.println("token.valid : ", token.Valid)
		//fmt.println("claim.id : ", claims.Id)
		return "", "", errors.New("invalid token: authentication failed")
	}
	// ! CHECK OTHER FILES
	return claims.Id, claims.Role, nil
}
