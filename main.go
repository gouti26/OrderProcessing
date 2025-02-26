package main

import (
	"OrderProcessing/handler"
	"OrderProcessing/migration"
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("orders.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	handlerService := handler.NewHandler(db)
	server := migration.NewServer(db, handlerService)
	server.InitializeServer()

	router := server.SetupRouter()

	fmt.Println("--------------")
	fmt.Println("Server is running on http://localhost:8080")
	fmt.Println("--------------")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
