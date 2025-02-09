package middleware

import (
	"context"
	"hr-app-diploma-backend/pkg/protobuf/hr-app-diploma-backend/auth-service/proto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authClient proto.AuthServiceClient
}

func NewAuthMiddleware(authClient proto.AuthServiceClient) *AuthMiddleware {
	return &AuthMiddleware{authClient: authClient}
}

func (m *AuthMiddleware) VerifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		res, err := m.authClient.VerifyToken(context.Background(), &proto.VerifyTokenRequest{Token: tokenString})
		if err != nil || !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", res.UserId)

		c.Next()
	}
}
