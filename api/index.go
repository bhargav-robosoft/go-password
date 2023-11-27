package api

import (
	"net/http"
	"password-manager/controller"
	"password-manager/logger"
	"password-manager/middleware"
	"password-manager/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	logger.Init()

	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	config := cors.DefaultConfig()
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowOrigins = []string{"https://react-password-manager.vercel.app", "http://localhost:3000"}
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowCredentials = true
	config.ExposeHeaders = []string{"Authorization", "Set-Cookie"}
	server.Use(cors.New(config))

	authController := controller.NewAuthController(service.NewAuthService())
	siteController := controller.NewSiteController(service.NewSiteService())

	server.POST("/generate-otp", authController.GenerateOtp)
	server.POST("/verify-otp", authController.VerifyOtp)
	server.POST("/sign-up", authController.SignUp)
	server.POST("/sign-in", authController.SignIn)
	server.PUT("/forgot-password", authController.ForgotPassword)
	server.PUT("/reset-password", middleware.TokenAuthMiddleware(), authController.ResetPassword)
	server.GET("/sign-out", middleware.TokenAuthMiddleware(), authController.SignOut)
	server.GET("/check-token", middleware.TokenAuthMiddleware(), authController.CheckToken)

	server.POST("/save-site", middleware.TokenAuthMiddleware(), siteController.SaveSite)
	server.GET("/get-sites", middleware.TokenAuthMiddleware(), siteController.GetSites)
	server.PATCH("/edit-site", middleware.TokenAuthMiddleware(), siteController.EditSite)
	server.DELETE("/delete-site", middleware.TokenAuthMiddleware(), siteController.DeleteSite)

	server.ServeHTTP(w, r)
}
