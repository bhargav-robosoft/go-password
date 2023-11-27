package controller

import (
	"net/http"
	"password-manager/entity"
	"password-manager/logger"
	"password-manager/service"
	"password-manager/util"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	GenerateOtp(ctx *gin.Context)
	VerifyOtp(ctx *gin.Context)
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
	ForgotPassword(ctx *gin.Context)
	ResetPassword(ctx *gin.Context)
	SignOut(ctx *gin.Context)
	CheckToken(ctx *gin.Context)
}

type authController struct {
	service service.AuthService
}

func NewAuthController(service service.AuthService) AuthController {
	return &authController{
		service: service,
	}
}

func (controller *authController) GenerateOtp(ctx *gin.Context) {
	var request entity.GenerateOtpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.ErrorLogger.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Email and Type are required and cannot be empty",
		})
		return
	}

	if !(request.Type == "register" || request.Type == "reset") {
		message := "Type can either be 'register' or 'reset'"
		logger.ErrorLogger.Println(message)
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	id, expiresAt, err := controller.service.GenerateOtp(request.Email, request.Type)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		parsedTime, _ := time.Parse(time.RFC3339, expiresAt)
		message := "OTP is generated and sent to mail"
		logger.InfoLogger.Println(message)
		cookie := util.GenerateCookie("id", id, int(time.Until(parsedTime).Seconds()))
		http.SetCookie(ctx.Writer, cookie)
		ctx.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"message":   message,
			"expiresAt": expiresAt,
		})
	}

}

func (controller *authController) VerifyOtp(ctx *gin.Context) {
	var request entity.VerifyOtpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message := "Email and OTP are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	id, err := ctx.Cookie("id")
	if err != nil {
		message := "OTP not generated"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	expiresAt, err := controller.service.VerifyOtp(id, request.Email, request.Otp)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "OTP is verified"
		logger.InfoLogger.Println(message)
		parsedTime, _ := time.Parse(time.RFC3339, expiresAt)
		cookie := util.GenerateCookie("id", id, int(time.Until(parsedTime).Seconds()))
		http.SetCookie(ctx.Writer, cookie)
		ctx.JSON(http.StatusOK, gin.H{
			"status":    http.StatusOK,
			"message":   message,
			"expiresAt": expiresAt,
		})
	}

}

func (controller *authController) SignUp(ctx *gin.Context) {
	var request entity.AuthRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		message := "Email and Password are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	id, err := ctx.Cookie("id")
	if err != nil {
		message := "Email not verified"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	err = controller.service.SignUp(id, request.Email, request.Password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "User successfully registered"
		logger.InfoLogger.Println(message)
		cookie := util.GenerateCookie("id", id, -1)
		http.SetCookie(ctx.Writer, cookie)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}

func (controller *authController) SignIn(ctx *gin.Context) {
	var request entity.AuthRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		message := "Email and Password are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	token, err := controller.service.SignIn(request.Email, request.Password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Sign in successful"
		logger.InfoLogger.Println(message)
		ctx.Header("Authorization", "Bearer "+token)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}

func (controller *authController) ForgotPassword(ctx *gin.Context) {
	var request entity.AuthRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		message := "Email and Password are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	id, err := ctx.Cookie("id")
	if err != nil {
		message := "Email not verified"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	err = controller.service.ForgotPassword(id, request.Email, request.Password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Password reset successful"
		logger.InfoLogger.Println(message)
		cookie := util.GenerateCookie("id", id, -1)
		http.SetCookie(ctx.Writer, cookie)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}

func (controller *authController) ResetPassword(ctx *gin.Context) {
	var request entity.ResetPasswordRequest
	err := ctx.ShouldBindJSON(&request)

	if err != nil {
		message := "Password and New Password are required and cannot be empty"
		logger.ErrorLogger.Println(err.Error())
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": message,
		})
		return
	}

	userId, _ := ctx.Get("userId")
	token, err := controller.service.ResetPassword(userId.(string), request.OldPassword, request.NewPassword)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Password reset successful"
		ctx.Header("Authorization", "Bearer "+token)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}

func (controller *authController) SignOut(ctx *gin.Context) {
	token, _ := ctx.Get("token")
	expirationTime, _ := ctx.Get("expirationTime")

	err := controller.service.SignOut(token.(string), expirationTime.(time.Time))

	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		if customErr, ok := err.(*util.CustomError); ok {
			logger.InfoLogger.Println(customErr.Message)
			ctx.JSON(customErr.Status, gin.H{
				"status":  customErr.Status,
				"message": customErr.Message,
			})
		}
	} else {
		message := "Sign out successful"
		logger.InfoLogger.Println(message)
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": message,
		})
	}
}

func (controller *authController) CheckToken(ctx *gin.Context) {
	message := "Token is valid"
	logger.InfoLogger.Println(message)
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": message,
	})
}
