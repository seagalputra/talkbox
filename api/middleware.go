package api

import (
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if hasPrefix := strings.Contains(authHeader, "Bearer"); !hasPrefix {
			ctx.JSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
			ctx.Abort()
		}

		auth := strings.Split(authHeader, " ")
		token, err := jwt.Parse(auth[1], func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("failed get signing method")
			} else if method != JwtSigningMethod {
				return nil, fmt.Errorf("signing method not match")
			}

			return JwtSecretKey, nil
		})

		if err != nil {
			log.Printf("[AuthenticateUser] %v", err)
			ctx.JSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
			ctx.Abort()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			ctx.JSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
			ctx.Abort()
		}

		userID := claims["id"].(string)
		user, err := FindUserByID(userID)
		if err != nil {
			log.Printf("[AuthenticateUser] %v", err)
			ctx.JSON(401, gin.H{
				"status":  "error",
				"message": "Unauthorized",
			})
			ctx.Abort()
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}
