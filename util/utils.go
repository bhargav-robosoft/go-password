package util

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

func GenerateOtp(numberOfDigits int) (otp int) {
	max := math.Pow(10, float64(numberOfDigits)-1) * 9
	return rand.Intn(int(max)) + int(math.Pow(10, float64(numberOfDigits)-1))
}

func TimestampToUnix(timestampMilliseconds int64) (unixTime time.Time) {
	timestampSeconds := timestampMilliseconds / 1000
	timestampNanoseconds := (timestampMilliseconds % 1000) * 1e6
	return time.Unix(timestampSeconds, timestampNanoseconds)
}

func GenerateCookie(cookieName string, cookieValue string, age int) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   age,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
		HttpOnly: false,
	}
}

func GetImage(domain string) (imageUrl string) {
	imageUrl = fmt.Sprintf("https://logo.clearbit.com/%v", domain)
	response, err := http.Get(imageUrl)
	if err != nil {
		return ""
	}
	if response.Status != "200 OK" {
		return ""
	}
	return imageUrl
}
