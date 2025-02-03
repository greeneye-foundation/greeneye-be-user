// internal/services/token_service.go
package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/greeneye-foundation/greeneye-be-user/internal/models"
	"github.com/greeneye-foundation/greeneye-be-user/internal/pkg/errors"
)

type TokenService struct {
	jwtSecret []byte
}

func (s *TokenService) GenerateAuthToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":       user.ID,
		"mobile_number": user.MobileNumber,
		"country_code":  user.CountryCode,
		"exp":           time.Now().Add(time.Hour * 24).Unix(),
		"issued_at":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) GeneratePasswordResetToken(mobileNumber string, countryCode string) (string, error) {
	claims := jwt.MapClaims{
		"mobile_number": mobileNumber,
		"country_code":  countryCode,
		"exp":           time.Now().Add(time.Hour).Unix(), // 1 hour validity
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *TokenService) ValidatePasswordResetToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New(errors.InvalidToken, "Invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return "", errors.New(errors.InvalidToken, "Invalid mobile number in token")
	}

	return email, nil
}
