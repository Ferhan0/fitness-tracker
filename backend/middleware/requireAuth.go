package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Ferhan0/fitness-tracker/initializers"
	"github.com/Ferhan0/fitness-tracker/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequireAuth, korumalı rotalara erişimden önce çalışan JWT doğrulama middleware'idir.
// Cookie'deki token'ı doğrular, kullanıcıyı veritabanından çeker ve context'e ekler.
func RequireAuth(c *gin.Context) {
	fmt.Println("Middleware executed")

	// "Authorization" adlı cookie'den JWT token string'ini al
	// Cookie yoksa 401 Unauthorized döndür ve isteği durdur
	tokenString, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(tokenString)

	// JWT token'ı parse et ve imzasını doğrula
	// WithValidMethods ile yalnızca HS256 algoritması kabul edilir (algorithm confusion saldırılarına karşı)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Token imzasını doğrulamak için SECRET_KEY'i döndür
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// Token claim'lerini MapClaims tipine dönüştür ve işle
	if claims, ok := token.Claims.(jwt.MapClaims); ok {

		// Yalnızca "access" tipindeki token'lara izin ver; refresh token ile bu middleware geçilemez
		if claims["type"] != "access" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Token'ın süresinin dolup dolmadığını kontrol et
		// exp claim'i Unix timestamp olarak saklıdır
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Token'ın "sub" claim'indeki kullanıcı ID'sini kullanarak veritabanından kullanıcıyı çek
		var user models.User
		initializers.DB.First(&user, claims["sub"])

		// Kullanıcı veritabanında bulunamazsa (silinmiş ya da geçersiz ID) 401 döndür
		if user.ID == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// Doğrulanmış kullanıcıyı Gin context'ine ekle; controller'lar c.Get("user") ile erişebilir
		c.Set("user", user)
	} else {
		// Claims dönüşümü başarısızsa hatayı logla
		fmt.Println(err)
	}

	// Bir sonraki handler veya middleware'e geç
	c.Next()
}
