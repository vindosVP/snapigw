package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"

	"github.com/vindosVP/snapigw/internal/utils/response"
)

type Claims struct {
	jwt.RegisteredClaims
	Email   string `json:"email,omitempty"`
	Id      int    `json:"id"`
	IsAdmin *bool  `json:"isAdmin,omitempty"`
}

func Authorize(secret string, requireAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			response.AbortErr(c, http.StatusUnauthorized, err.Error())
			return
		}
		token, err := parseToken(jwtToken, secret)
		if err != nil {
			response.AbortErr(c, http.StatusUnauthorized, err.Error())
			return
		}
		claims, OK := token.Claims.(*Claims)
		if !OK {
			response.AbortErr(c, http.StatusInternalServerError, err.Error())
			return
		}
		if requireAdmin && !*claims.IsAdmin {
			response.AbortErr(c, http.StatusUnauthorized, "You are not authorized for this operation")
		}
		c.Set("userId", claims.Id)
		c.Set("isAdmin", claims.IsAdmin)
		c.Next()
	}
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func parseToken(jwtToken string, secret string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, OK := token.Method.(*jwt.SigningMethodHMAC); !OK {
			return nil, errors.New("bad signed method received")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New("bad jwt token")
	}

	return token, nil
}
