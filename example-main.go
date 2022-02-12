package main

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
   	"net/http"
	"api/user"
)

func main()  {
	// connect to mysql
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection to database is good")

	get data from database
	var users []user.User
	length := len(users) // masih 0 data
	fmt.Println(length)
	
	// secara otomatis var users dapat terkoneksi dengan db melalui file user entity.go
	db.Find(&users)

	length = len(users) // sudah berisi data dari database
	fmt.Println(length)

	for _, user := range users {
		fmt.Println(user.Name)
		fmt.Println(user.Email)
		fmt.Println("==============")
	}

	// membuat routing
	// router := gin.Default() // membuat router di gin
	// router.GET("/handler", handler) // object router
	// router.Run() // menjalankan router
}

// menampilkan data dalam bentuk json pada web
function untuk dijadikan handler / controller dalam laravel
func handler(c *gin.Context)  {
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
  	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	var users []user.User
	db.Find(&users)

	c.JSON(http.StatusOK, users)
}

/*
Step API:
-> input
-> handler - mapping input ke struct
-> service - mapping ke struct User
-> repository - menyimpan struct User ke db
-> db
*/