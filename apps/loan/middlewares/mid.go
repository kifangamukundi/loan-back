package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kifangamukundi/gm/loan/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var (
	TokenPrefix     = "Bearer"
	ErrUnauthorized = errors.New("not authorized to access this resource")
	ErrForbidden    = errors.New("you do not have the required permissions to access this resource")
)

func ExtractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, TokenPrefix) {
		return "", ErrUnauthorized
	}
	return strings.TrimPrefix(authHeader, TokenPrefix+" "), nil
}

func VerifyToken(tokenString string) (jwt.MapClaims, error) {
	service := os.Getenv("ISSUE_NAME")
	recipient := os.Getenv("RECIPIENT_NAME")
	accessSecret := os.Getenv("ACCESS_SECRET")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(accessSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, errors.New("token has expired")
		}
	} else {
		return nil, errors.New("missing exp claim")
	}

	if claims["aud"] != recipient || claims["iss"] != service {
		return nil, errors.New("invalid audience or issuer")
	}

	return claims, nil
}

func AdvancedAuth(db *gorm.DB, requiredPermissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ExtractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		claims, err := VerifyToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		var user models.User
		err = db.Preload("Roles.Permissions").First(&user, "id = ?", int(userID)).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
			}
			c.Abort()
			return
		}

		userPermissions := map[string]bool{}
		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				userPermissions[permission.PermissionName] = true
			}
		}

		for _, requiredPermission := range requiredPermissions {
			if !userPermissions[requiredPermission] {
				c.JSON(http.StatusForbidden, gin.H{"error": ErrForbidden.Error()})
				c.Abort()
				return
			}
		}

		c.Set("user", user)

		c.Next()
	}
}
