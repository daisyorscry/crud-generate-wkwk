package jwt

import (
	"daisy/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwt(userId, roleId string) (string, error) {
	if userId == "" || roleId == "" {
		return "", fmt.Errorf("userid or roleid is missing")
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"role_id": roleId,
		"exp":     time.Now().Add(time.Minute * 5).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedToken, err := token.SignedString(config.Get().JWT.SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}

func GenerateRefreshJwt(userId, roleId string) (string, error) {
	if userId == "" || roleId == "" {
		return "", fmt.Errorf("userid or roleid is missing")
	}

	claims := jwt.MapClaims{
		"user_id": userId,
		"role_id": roleId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(config.Get().JWT.SecretKeyRefresh)

}

func ValidateJwt(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return config.Get().JWT.SecretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}

func ValidateJwtRefresh(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return config.Get().JWT.SecretKeyRefresh, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}
