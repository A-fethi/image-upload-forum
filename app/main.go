package main

import (
	"database/sql"
	"forum/app/config"
	database "forum/app/db"
	"forum/app/handlers"
	"forum/app/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	config.InitLogger()
	config.InitTemplates("./templates/*.html")
}

func main() {

	db, err := sql.Open("sqlite3", "./app/db/forum.db")
	if err != nil {
		config.Logger.Println("error opening database: ", err)
		config.Logger.Println("Shutting down...")
		return
	}
	database.CreateTables(db)
	handlers.RegisterRoutes(db)

	go func() {
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			config.Logger.Println("Error starting server: ", err)
			config.Logger.Println("Shutting down...")
			config.CloseLogger()
			db.Close()
			os.Exit(1)
		}
	}()
	utils.Print()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	sig := <-sigChan
	config.Logger.Println("\nReceived termination signal:", sig)
	config.Logger.Println("Shutting down...")
	db.Close()
	config.Logger.Println("Database connection closed.")
	config.CloseLogger()

}
