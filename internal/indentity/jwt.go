package indentity

import (
	"fmt"
	"time"

	"github.com/cairon666/vkr-backend/internal/apperrors"
	"github.com/cairon666/vkr-backend/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtClaims struct {
	AuthClaims AuthClaims
	jwt.RegisteredClaims
}

type JWTService struct {
	config *config.Config
}

func (s *JWTService) newClaims(claims AuthClaims, duration time.Duration) *JwtClaims {
	id := uuid.New()

	return &JwtClaims{
		claims,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.JWT.Issuer,
			Subject:   "",
			ID:        id.String(),
			Audience:  []string{""},
		},
	}
}

func NewJWTService(config *config.Config) *JWTService {
	return &JWTService{config: config}
}

func (s *JWTService) getSecret(token *jwt.Token) (interface{}, error) {
	return []byte(s.config.JWT.Secret), nil
}

func (s *JWTService) ParseToken(tokenString string) (*AuthClaims, error) {
	var claims JwtClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, s.getSecret, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	fmt.Println()
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", apperrors.ErrInvalidToken)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", apperrors.ErrInvalidToken)
	}

	if claims, ok := token.Claims.(*JwtClaims); ok {
		return &claims.AuthClaims, nil
	} else {
		return nil, fmt.Errorf("invalid token claims: %w", apperrors.ErrInvalidToken)
	}
}

func (s *JWTService) createToken(claims *JwtClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.config.JWT.Secret))
}

func (s *JWTService) CreateAccessToken(authClaims AuthClaims) (string, error) {
	accessClaims := s.newClaims(authClaims, s.config.JWT.AccessTokenDuration)

	token, err := s.createToken(accessClaims)
	if err != nil {
		return "", fmt.Errorf("failed to create access token: %w", err)
	}

	return token, nil
}

func (s *JWTService) CreateRefreshToken(authClaims AuthClaims) (string, error) {
	refreshClaims := s.newClaims(authClaims, s.config.JWT.RefreshTokenDuration)

	token, err := s.createToken(refreshClaims)
	if err != nil {
		return "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return token, nil
}

func (s *JWTService) SetRefreshCookie(c *gin.Context, value string) {
	c.SetCookie("refresh_token", value, 3600, "/", "*", false, true)
}

func (s *JWTService) GetRefreshCookie(c *gin.Context) (string, error) {
	return c.Cookie("refresh_token")
}
