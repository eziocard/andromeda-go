package initializers

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("DB_URL no está definida")
	}

	var err error
	maxRetries := 10
	wait := time.Second

	for i := 1; i <= maxRetries; i++ {
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Conectado a la base de datos")
			return
		}

		log.Printf("No se pudo conectar a la DB (intento %d/%d): %v", i, maxRetries, err)

		if i == maxRetries {
			break
		}

		time.Sleep(wait)
		if wait < 30*time.Second {
			wait *= 2 // backoff exponencial, con techo en 30s
		}
	}

	log.Fatalf("No se pudo conectar a la base de datos después de %d intentos: %v", maxRetries, err)
}
