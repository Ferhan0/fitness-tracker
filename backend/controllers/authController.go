package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Ferhan0/fitness-tracker/initializers"
	"github.com/Ferhan0/fitness-tracker/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// SignUp, POST /signup rotasının handler fonksiyonudur.
// İstek gövdesinden e-posta ve şifreyi alır, şifreyi hashler ve yeni kullanıcı oluşturur.
func SignUp(c *gin.Context) {
	// İstek gövdesinden e-posta ve şifre alanlarını almak için anonim struct tanımla
	var body struct {
		Email    string
		Password string
	}

	// JSON gövdesini struct'a bağla; hata varsa 400 döndür
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	// Şifreyi bcrypt ile hashle (DefaultCost = 10 round); ham şifre asla saklanmaz
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Hashlenmiş şifreyle yeni User kaydı oluştur ve veritabanına kaydet
	user := models.User{Email: body.Email, Password: string(hashedPassword)}
	result := initializers.DB.Create(&user)

	// Veritabanı hatası varsa (örn. e-posta zaten kayıtlıysa) 400 döndür
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	// Kullanıcı başarıyla oluşturuldu; 201 Created yanıtı döndür
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

// Login, POST /login rotasının handler fonksiyonudur.
// Kullanıcı kimlik bilgilerini doğrular, JWT token üretir ve HttpOnly cookie olarak set eder.
func Login(c *gin.Context) {
	// İstek gövdesinden e-posta ve şifre alanlarını almak için anonim struct tanımla
	var body struct {
		Email    string
		Password string
	}

	// JSON gövdesini struct'a bağla; hata varsa 400 döndür
	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	// Veritabanında e-posta adresine göre kullanıcıyı ara
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)

	// Kullanıcı bulunamazsa (ID == 0) genel hata mesajı döndür (email enumeration'ı önler)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Gönderilen şifreyi veritabanındaki bcrypt hash'i ile karşılaştır
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		// Şifre eşleşmedi; güvenlik için e-posta ile aynı hata mesajı kullan
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Access token oluştur: sub = kullanıcı ID, exp = 24 saat, type = "access"
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
		"type": "access",
	})

	// Refresh token oluştur: 7 günlük ömür, type = "refresh"
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
		"type": "refresh",
	})

	// Access token'ı imzala
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create access token"})
		return
	}

	// Refresh token'ı imzala
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create refresh token"})
		return
	}

	// Cookie güvenlik ayarı: SameSite=Lax (CSRF koruması için)
	c.SetSameSite(http.SameSiteLaxMode)

	// Access token: 24 saatlik HttpOnly cookie
	c.SetCookie("Authorization", accessTokenString, 3600*24, "", "", false, true)

	// Refresh token: 7 günlük HttpOnly cookie
	c.SetCookie("RefreshToken", refreshTokenString, 3600*24*7, "", "", false, true)

	// Başarılı giriş
	c.JSON(http.StatusOK, gin.H{})
}

// RefreshToken, POST /refresh rotasının handler fonksiyonudur.
// "RefreshToken" cookie'sindeki token'ı doğrular ve yeni bir access token üretir.
func RefreshToken(c *gin.Context) {
	// Refresh token cookie'sini al
	refreshTokenString, err := c.Cookie("RefreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	// Refresh token'ı parse et ve doğrula
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	// Token tipinin "refresh" olduğunu doğrula; access token ile bu endpoint kullanılamaz
	if claims["type"] != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
		return
	}

	// Süre kontrolü
	if float64(time.Now().Unix()) > claims["exp"].(float64) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
		return
	}

	// Kullanıcıyı veritabanından çek
	var user models.User
	initializers.DB.First(&user, claims["sub"])
	if user.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Yeni access token oluştur
	newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"exp":  time.Now().Add(time.Minute * 1).Unix(),
		"type": "access",
	})

	newAccessTokenString, err := newAccessToken.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", newAccessTokenString, 3600*24, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Token refreshed"})
}

// Validate, GET /validate rotasının handler fonksiyonudur.
// RequireAuth middleware'i tarafından doğrulandıktan sonra çalışır;
// context'e eklenen kullanıcı bilgisini döndürür.
func Validate(c *gin.Context) {
	// RequireAuth middleware'inin context'e eklediği kullanıcı nesnesini al
	user, _ := c.Get("user")

	// Oturum açmış kullanıcı bilgisini JSON olarak döndür
	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
