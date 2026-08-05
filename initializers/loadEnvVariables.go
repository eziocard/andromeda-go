package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró .env, se usará el entorno del sistema (esperado en Docker/producción)")
	}
}
