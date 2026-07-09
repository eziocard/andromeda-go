package main

import (
	"github.com/eziocard/andromeda-go/initializers"
	"github.com/eziocard/andromeda-go/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(
		&models.Product{},
		&models.Sale{},
		&models.SaleDetail{},
		&models.SalePayment{},
		&models.Business{})
}
