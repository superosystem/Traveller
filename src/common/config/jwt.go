package config

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/superosystem/trainingsystem-backend/src/common/helper"
	"strings"
	"time"
)

type JWTConfig struct {
	// JWT signature secret
	JWTSecret string
}

func NewJWTConfig(secret string) *JWTConfig {
	return &JWTConfig{
		JWTSecret: secret,
	}
}

type JWTCustomClaims struct {
	UserId   string `json:"user_id"`
	MenteeId string `json:"mentee_id,omitempty"`
	MentorId string `json:"mentor_id,omitempty"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

func (config *JWTConfig) GenerateToken(userId string, actorId string, role string, exp time.Time) (string, error) {
	claims := JWTCustomClaims{
		UserId: userId,
		Role:   role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	switch role {
	case "mentee":
		claims.MenteeId = actorId
	case "mentor":
		claims.MentorId = actorId
	default:
		return "", helper.ErrInvalidJWTPayload
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := config.JWTSecret

	return token.SignedString([]byte(jwtSecret))
}

func (config *JWTConfig) ExtractToken(c echo.Context) (*JWTCustomClaims, error) {
	tokenString := c.Request().Header.Get("Authorization")

	sanitizedTokenBearer := strings.Replace(tokenString, "Bearer ", "", 1)

	token, err := jwt.Parse(sanitizedTokenBearer, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, helper.ErrInvalidTokenHeader
		}

		jwtSecret := config.JWTSecret

		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if token.Valid {
		customClaim := JWTCustomClaims{}

		claims := token.Claims.(jwt.MapClaims)
		customClaim.UserId = claims["user_id"].(string)
		customClaim.Role = claims["role"].(string)

		switch customClaim.Role {
		case "mentee":
			customClaim.MenteeId = claims["mentee_id"].(string)
		case "mentor":
			customClaim.MentorId = claims["mentor_id"].(string)
		}

		return &customClaim, nil
	}

	return nil, err
}
