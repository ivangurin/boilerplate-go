package jwt

import (
	"boilerplate/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	KeyUserID   = "user_id"
	KeyUserName = "user_name"
	KeyExp      = "exp"
)

func ParseToken(token string, config *model.ConfigAPI) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(config.AccessPrivateKey), nil
	})
	if err != nil {
		return nil, nil, err
	}

	return parsedToken, claims, nil
}

func ValidateToken(token string, config *model.ConfigAPI) (jwt.MapClaims, error) {
	parsedToken, claims, err := ParseToken(token, config)
	if err != nil {
		return nil, errors.New("недействительный токен")
	}

	if !parsedToken.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}

func GenerateAccessToken(userID int, userName string, config *model.ConfigAPI) (string, error) {
	claims := jwt.MapClaims{
		KeyUserID:   userID,
		KeyUserName: userName,
		KeyExp:      time.Now().UTC().Add(time.Second * time.Duration(config.AccessTokenTTL)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AccessPrivateKey))
}

func GenerateRefreshToken(userID int, userName string, config *model.ConfigAPI) (string, error) {
	claims := jwt.MapClaims{
		KeyUserID:   userID,
		KeyUserName: userName,
		KeyExp:      time.Now().UTC().Add(time.Second * time.Duration(config.RefreshTokenTTL)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AccessPrivateKey))
}

func GetUserID(claims jwt.MapClaims) (int, bool) {
	userID, ok := claims[KeyUserID].(float64)
	if !ok {
		return 0, false
	}
	return int(userID), true
}

func GetUserName(claims jwt.MapClaims) (string, bool) {
	userName, ok := claims[KeyUserName].(string)
	if !ok {
		return "", false
	}
	return userName, true
}
