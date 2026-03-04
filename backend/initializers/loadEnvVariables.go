package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

// LoadEnvVariables, proje kök dizinindeki .env dosyasını okuyarak
// PORT, DB ve SECRET_KEY gibi çevre değişkenlerini os.Getenv ile erişilebilir hale getirir.
// .env dosyası bulunamazsa veya okunamazsa uygulama başlatılmaz (log.Fatal).
func LoadEnvVariables() {
	// godotenv.Load, .env dosyasını okur ve değişkenleri process ortamına yükler
	err := godotenv.Load(".env")
	if err != nil {
		// .env dosyası yoksa veya okunamıyorsa uygulamayı durdur
		log.Fatal("error loading .env file")
	}
}
