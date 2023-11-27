package service

import (
	"net/http"
	"password-manager/db"
	"password-manager/logger"
	"password-manager/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	GenerateOtp(email string, otpType string) (id string, expiresAt string, err error)
	VerifyOtp(dbId string, email string, otp string) (expiresAt string, err error)
	SignUp(dbId string, email string, password string) (err error)
	SignIn(email string, password string) (token string, err error)
	ForgotPassword(dbId string, email string, password string) (err error)
	ResetPassword(userId string, oldPassword string, newPassword string) (token string, err error)
	SignOut(token string, expirationTime time.Time) (err error)
}

type authService struct{}

func NewAuthService() AuthService {
	return &authService{}
}

func (service *authService) GenerateOtp(email string, otpType string) (id string, expiresAt string, err error) {
	registerationStatus, err := db.CheckUserRegistered(email)
	var purpose string
	if otpType == "reset" {
		purpose = "forgot password"
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", err
		}
		if !registerationStatus {
			message := "Email is not registered"
			logger.ErrorLogger.Println(message)
			return "", "", &util.CustomError{Message: message, Status: http.StatusConflict}
		}
	} else {
		purpose = "registration"
		if err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", err
		}
		if registerationStatus {
			message := "Email is already registered"
			logger.ErrorLogger.Println(message)
			return "", "", &util.CustomError{Message: message, Status: http.StatusConflict}
		}
	}

	id, err = db.CheckOtpGenerated(email)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", "", &util.CustomError{Message: "OTP generation failed", Status: http.StatusInternalServerError}
	}

	otp := util.GenerateOtp(6)
	if id == "" {
		if id, expiresAt, err = db.GenerateOtp(email, strconv.Itoa(otp)); err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", &util.CustomError{Message: "OTP generation failed", Status: http.StatusInternalServerError}
		}

		if err = util.SendEmailOtp(email, otp, purpose); err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", &util.CustomError{Message: "OTP generation failed", Status: http.StatusInternalServerError}
		}
	} else {
		if expiresAt, err = db.ReGenerateOtp(email, strconv.Itoa(otp)); err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", &util.CustomError{Message: "OTP generation failed", Status: http.StatusInternalServerError}
		}

		if err = util.SendEmailOtp(email, otp, purpose); err != nil {
			logger.ErrorLogger.Println(err.Error())
			return "", "", &util.CustomError{Message: "OTP generation failed", Status: http.StatusInternalServerError}
		}
	}

	return id, expiresAt, nil
}

func (service *authService) VerifyOtp(dbId string, email string, otp string) (expiresAt string, err error) {
	if err = db.VerifyOtp(dbId, email, otp); err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	if expiresAt, err = db.OtpVerified(email); err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: err.Error(), Status: http.StatusInternalServerError}
	}

	return expiresAt, nil
}

func (service *authService) SignUp(dbId string, email string, password string) error {
	verificationStatus, err := db.RemoveVerifiedUser(dbId, email)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return err
	}
	if !verificationStatus {
		message := "Email is not verified"
		logger.ErrorLogger.Println(message)
		return &util.CustomError{Message: message, Status: http.StatusBadRequest}
	}

	_, err = db.RegisterUser(email, password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	return err
}

func (service *authService) SignIn(email string, password string) (string, error) {
	registerationStatus, err := db.CheckUserRegistered(email)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", err
	}
	if !registerationStatus {
		message := "Email is not registered"
		logger.ErrorLogger.Println(message)
		return "", &util.CustomError{Message: message, Status: http.StatusNotFound}
	}

	validCredentials, userId, passwordSetAt, err := db.CheckUserCredentials(email, password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", err
	}
	if !validCredentials {
		message := "Email or password is wrong"
		logger.ErrorLogger.Println(message)
		return "", &util.CustomError{Message: message, Status: http.StatusNotFound}
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["passwordSetAt"] = passwordSetAt
	claims["id"] = userId

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return t, nil
}

func (service *authService) ForgotPassword(dbId string, email string, password string) error {
	verificationStatus, err := db.CheckUserVerified(dbId, email)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return err
	}
	if !verificationStatus {
		message := "Email is not verified"
		logger.ErrorLogger.Println(message)
		return &util.CustomError{Message: message, Status: http.StatusBadRequest}
	}

	_, err = db.ResetPassword(email, password)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
	}

	return err
}

func (service *authService) ResetPassword(userId string, oldPassword string, newPassword string) (string, error) {
	userEmail, err := db.CheckUserCredentialsWithId(userId, oldPassword)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", err
	}
	if userEmail == "" {
		message := "Password is wrong"
		logger.ErrorLogger.Println(message)
		return "", &util.CustomError{Message: message, Status: http.StatusNotFound}
	}

	passwordSetAt, err := db.ResetPassword(userEmail, newPassword)
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", err
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 60).Unix()
	claims["passwordSetAt"] = passwordSetAt
	claims["id"] = userId

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.ErrorLogger.Println(err.Error())
		return "", &util.CustomError{Message: "Internal Server Error", Status: http.StatusInternalServerError}
	}

	return t, err
}

func (service *authService) SignOut(token string, expirationTime time.Time) error {
	db.BlacklistToken(token, expirationTime)
	return nil
}
