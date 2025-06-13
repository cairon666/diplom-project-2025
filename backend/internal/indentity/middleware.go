package indentity

import (
	"context"
	"net/http"
	"strings"

	"github.com/cairon666/vkr-backend/pkg/api_key"
	"github.com/gin-gonic/gin"
)

// Контекстный ключ для хранения claims.
type contextKey string

const AuthClaimsContextKey = contextKey("auth_claims")

func GetAuthClaims(ctx context.Context) (*AuthClaims, bool) {
	claims, ok := ctx.Value(AuthClaimsContextKey).(*AuthClaims)

	return claims, ok
}

type IdentityService struct {
	jwtService    *JWTService
	apiKeyService *ApiKeyService
}

func NewIdentityService(jwtService *JWTService, apiKeyService *ApiKeyService) *IdentityService {
	return &IdentityService{
		jwtService:    jwtService,
		apiKeyService: apiKeyService,
	}
}

func (is *IdentityService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Попытка получить JWT из заголовка Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				token := parts[1]
				claims, err := is.jwtService.ParseToken(token)
				if err == nil {
					ctx := context.WithValue(c.Request.Context(), AuthClaimsContextKey, claims)
					c.Request = c.Request.WithContext(ctx)
					c.Next()

					return
				}
				// если ошибка - не прерываем сразу, попробуем api_key
			}
		}

		// 2. Попытка получить API key из заголовка X-Api-Key
		apiKey := c.GetHeader("X-Api-Key")
		if apiKey != "" {
			apiKeyHash := api_key.HashAPIKey(apiKey)

			claims, err := is.apiKeyService.GetAuthClaimsByAPIKeyHash(apiKeyHash)
			if err == nil && claims != nil {
				ctx := context.WithValue(c.Request.Context(), AuthClaimsContextKey, claims)
				c.Request = c.Request.WithContext(ctx)
				c.Next()

				return
			}
		}

		// Если ни JWT, ни API key не сработали
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
