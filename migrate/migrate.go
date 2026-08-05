package main

import (
	"log"

	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/internal/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	err := initializers.DB.AutoMigrate(
		&models.Product{},
		&models.Sale{},
		&models.SaleDetail{},
		&models.SalePayment{},
		&models.Business{},
		&models.Role{},
		&models.User{},
		&models.RefreshToken{},
		&models.VoidedSale{},
		&models.VoidedSaleDetail{},
		&models.VoidedSalePayment{},
	)

	if err != nil {
		log.Fatalf("Error al migrar: %v", err)
	}

	log.Println("Migraciones aplicadas correctamente")
}
