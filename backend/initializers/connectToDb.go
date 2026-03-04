package initializers

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB, uygulama genelinde kullanılan paylaşımlı GORM veritabanı bağlantı nesnesidir.
// Tüm controller ve middleware'ler bu değişken üzerinden veritabanına erişir.
var DB *gorm.DB

// ConnectToDb, .env dosyasındaki DB değişkenini kullanarak PostgreSQL veritabanına bağlanır.
// Bağlantı başarısız olursa uygulama panic ile durdurulur.
func ConnectToDb() {
	var err error

	// DB ortam değişkeninden PostgreSQL bağlantı dizgisini (DSN) al
	// Örnek: "host=localhost user=postgres password=1234 dbname=staj_db port=5432 sslmode=disable"
	dsn := os.Getenv("DB")

	// GORM aracılığıyla PostgreSQL bağlantısını kur ve global DB değişkenine ata
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		// Veritabanına bağlanılamıyorsa uygulama çalışamaz; panic ile durdur
		panic("Failed to connect to DB")
	} else {
		fmt.Println("Veritabanına bağlandı!")
	}
}
