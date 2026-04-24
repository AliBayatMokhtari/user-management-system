package main

import (
	"log"
	"net/http"
	"ums/handler"
	"ums/model"
	"ums/repository"
	"ums/router"
	"ums/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=secret dbname=mydb port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	db.AutoMigrate(&model.User{})

	repo := repository.NewGormUserRepository(db)
	userSvc := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userSvc)

	r := router.NewRouter(userHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
