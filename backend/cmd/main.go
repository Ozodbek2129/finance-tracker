package main

import (
	"log"

	connectiondb "finance-tracker/internal/connectiondb"
	handler "finance-tracker/handler"
	postgrescrud "finance-tracker/internal/connectiondb/postgres_crud"
)

func main() {
	db, err := connectiondb.ConnectionDb()
	if err != nil {
		log.Fatal("DB ga ulanishda xato:", err)
	}
	defer db.Close()

	ft := postgrescrud.NewFinanceTracker(db)

	h := handler.NewHandler(ft)

	r := h.SetupRouter()

	// Server ishga tushirish
	log.Println("Server :8080 portida ishlamoqda...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server ishga tushmadi:", err)
	}
}