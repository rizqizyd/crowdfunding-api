package main

import (
	"api/handler"
	"api/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// connect to mysql
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// passing db ke NewRepository pada file repository
	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	// membuat router
	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("/api/v1")

	// register handler untuk dapat diakses pada api "/users"
	api.POST("/users", userHandler.RegisterUser)
	router.Run()

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Tes simpan dari service"
	// userInput.Email = "cth@gmail.com"
	// userInput.Occupation = "programmer"
	// userInput.Password = "pass"

	// userService.RegisterUser(userInput)

	// user := user.User{
	// 	Name: "Test simpan",
	// }

	// userRepository.Save(user)
}

/*
Step API:
-> input
-> handler - mapping input dari user ke struct input
-> service - mapping dari struct input ke struct User
-> repository - menyimpan struct User ke db
-> db
*/
