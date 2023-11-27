package middleware

import (
	"errors"
	"net/http"
	"password-manager/db"
	"password-manager/logger"
	"password-manager/util"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := processAuthHeader(c)
		if err != nil {
			return
		}

		token, err := parseToken(c, tokenString)
		if err != nil {
			return
		}

		claims, err := checkClaims(c, token)
		if err != nil {
			return
		}

		var wg sync.WaitGroup

		wg.Add(1)
		var err1 *util.CustomError
		go func() {
			defer wg.Done()
			err1 = checkBlacklisted(c, tokenString)
		}()

		wg.Add(1)
		var err2 *util.CustomError
		go func() {
			defer wg.Done()
			err2 = checkPasswordTimestamp(c, claims["id"].(string), claims["passwordSetAt"].(string))
		}()

		wg.Wait()

		if err1 != nil {
			response := gin.H{
				"status":  err1.Status,
				"message": err1.Message}
			if err1.Message == "Token is blacklisted" {
				response["sessionTimedOut"] = true
			}
			c.JSON(err1.Status, response)
			c.Abort()
			return
		}

		if err2 != nil {
			response := gin.H{
				"status":  err2.Status,
				"message": err2.Message}
			if err2.Message == "Token is blacklisted" {
				response["sessionTimedOut"] = true
			}
			c.JSON(err2.Status, response)
			c.Abort()
			return
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		c.Set("token", tokenString)
		c.Set("userId", claims["id"])
		c.Set("expirationTime", expirationTime)
		c.Next()
	}
}

func processAuthHeader(c *gin.Context) (token string, err error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		message := "Authorization header is missing"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": message,
		})
		c.Abort()
		return "", errors.New(strings.ToLower(message))
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
		message := "Invalid or missing Bearer token"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": message,
		})
		c.Abort()
		return "", errors.New(strings.ToLower(message))
	}

	return authHeaderParts[1], nil
}

func parseToken(c *gin.Context, tokenString string) (token *jwt.Token, err error) {
	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		message := err.Error()
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":          http.StatusUnauthorized,
					"message":         "Token is expired",
					"sessionTimedOut": true,
				})
				c.Abort()
				return nil, err
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": message,
		})
		c.Abort()
		return nil, err
	}
	return token, nil
}

func checkClaims(c *gin.Context, token *jwt.Token) (claims jwt.MapClaims, err error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		if claims["id"] == nil || claims["passwordSetAt"] == nil || claims["exp"] == nil {
			message := "Invalid token. Required claims not found."
			logger.ErrorLogger.Println(message)
			logger.InfoLogger.Println(message)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": message,
			})
			c.Abort()
			return nil, errors.New(strings.ToLower(message))
		}

		if !(reflect.TypeOf(claims["id"]).Kind() == reflect.String) && (reflect.TypeOf(claims["passwordSetAt"]).Kind() == reflect.String) && (reflect.TypeOf(claims["exp"]).Kind() == reflect.Float64) {
			message := "Invalid token. Required claims types invalid."
			logger.ErrorLogger.Println(message)
			logger.InfoLogger.Println(message)
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": message,
			})
			c.Abort()
			return nil, errors.New(strings.ToLower(message))
		}
	} else {
		message := "Token is invalid"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": message,
		})
		c.Abort()
		return nil, errors.New(strings.ToLower(message))
	}
	return claims, nil
}

func checkBlacklisted(c *gin.Context, tokenString string) (err *util.CustomError) {
	blacklisted, er := db.CheckBlacklist(tokenString)

	if er != nil {
		message := "Internal Server Error"
		logger.ErrorLogger.Println(er.Error())
		logger.InfoLogger.Println(message)
		return &util.CustomError{Status: http.StatusInternalServerError, Message: message}
	}

	if blacklisted {
		message := "Token is blacklisted"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		return &util.CustomError{Status: http.StatusUnauthorized, Message: message}
	}
	return nil
}

func checkPasswordTimestamp(c *gin.Context, id string, passwordSetAt string) (err *util.CustomError) {
	passwordTime, er := db.CheckPasswordReset(id)

	if er != nil {
		message := "Internal Server Error"
		logger.ErrorLogger.Println(er.Error())
		logger.InfoLogger.Println(message)
		return &util.CustomError{Status: http.StatusInternalServerError, Message: message}
	}

	tokenPasswordTime := passwordSetAt
	formattedPasswordTime := passwordTime

	if tokenPasswordTime != formattedPasswordTime {
		message := "Token is blacklisted"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		return &util.CustomError{Status: http.StatusUnauthorized, Message: message}
	}

	return nil
}
