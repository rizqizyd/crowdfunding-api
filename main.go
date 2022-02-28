package main

import (
	"api/auth"
	"api/campaign"
	"api/handler"
	"api/helper"
	"api/payment"
	"api/transaction"
	"api/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
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
	// buat instance dari campaign repository
	campaignRepository := campaign.NewRepository(db)
	// instansiasi transaction repository untuk bisa passing db
	transactionRepository := transaction.NewRepository(db)

	// panggil semua data campaign dari database (cek manual)
	// campaigns, err := campaignRepository.FindAll()
	// panggil data campaign by ID
	// campaigns, err := campaignRepository.FindByUserID(1)
	// fmt.Println("debug")
	// fmt.Println(len(campaigns)) // menampilkan jumlah campaign
	// // tampilkan nama setiap campaign
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	// cek campaign memiliki gambar atau tidak
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println("jumlah gambar yg di load:", len(campaign.CampaignImages))
	// 		// akses campaign images
	// 		fmt.Println(campaign.CampaignImages[0].FileName)
	// 	}
	// }

	// akses terhadap user repository
	userService := user.NewService(userRepository)

	// menampilkan data campaign
	campaignService := campaign.NewService(campaignRepository)
	// campaigns, _ := campaignService.GetCampaigns(0)
	// fmt.Println(len(campaigns))

	// memanggil service auth
	authService := auth.NewService()
	// payment service
	paymentService := payment.NewService()
	// memanggil transaction service
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	// handler
	// membuat router. authService yang udah dibuat, kita passing ke dalam userHandler
	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)

	// tes transactionService (manual)
	// user, _ := userService.GetUserByID(1)
	// input := transaction.CreateTransactionInput{
	// 	CampaignID: 6,
	// 	Amount:     3000000,
	// 	User:       user,
	// }
	// transactionService.CreateTransaction(input)

	// panggil service function CreateCampaign untuk tes (manual)
	// input := campaign.CreateCampaignInput{}
	// input.Name = "Penggalangan Dana Startup"
	// input.ShortDescription = "short"
	// input.Description = "long description"
	// input.GoalAmount = 1000000000
	// input.Perks = "hadiah satu, dua, dan tiga"

	// // menggunakan user 1
	// inputUser, _ := userService.GetUserByID(3)
	// input.User = inputUser

	// // panggil CreateCampaign
	// _, err = campaignService.CreateCampaign(input)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// tes validate token (manual)
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.BZcBBLhOhjA9ojwmRNMLx7x0IR83QyTeiH48psbhKLI")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// }

	// if token.Valid {
	// 	fmt.Println("VALID")
	// } else {
	// 	fmt.Println("INVALID")
	// }

	// tes hasil kembalian dari function generate token (manual)
	// fmt.Println(authService.GenerateToken(1001))

	// save avatar (manual)
	// userService.SaveAvatar(1, "images/1-profile.png")

	// login user (tes service manual)
	// input := user.LoginInput{
	// 	Email:    "masonmount@gmail.com",
	// 	Password: "1234a5678",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi Kesalahan")
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(user.Email)
	// fmt.Println(user.Name)

	// menampilkan user by email (manual)
	// userByEmail, err := userRepository.FindByEmail("masonmount@gmail.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if userByEmail.ID == 0 {
	// 	fmt.Println("User Tidak Ditemukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	router := gin.Default()
	// CORS (allow cors)
	router.Use(cors.Default())
	// set router untuk mengambil gambar user melalui folder images
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	// register handler untuk dapat diakses pada api "/users"
	// daftarkan endpoint
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvaliability)
	// jika kita melakukan request ke avatars, kita perlu mengirimkan jwt token sebelum menuju ke userHandler
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	// create new campaign
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaign/:id", campaignHandler.GetCampaign)
	api.GET("/campaign/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransactions)
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTransactions)

	// update campaign
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)

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

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil nilai header Authorization: Bearer tokentokentoken
		authHeader := c.GetHeader("Authorization")
		// apakah di dalam string authHeader terdapat kata Bearer
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized 1", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // hentikan status jika tidak ada bearer token
			return
		}
		// dari header Authorization, kita ambil nilai tokennya saja
		// bearer token
		var tokenString string
		// tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}
		// validasi token
		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized 2", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // hentikan status jika tidak ada bearer token
			return
		}
		// ambil data yang ada di dalam token
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized 3", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // hentikan status jika tidak ada bearer token
			return
		}
		// jika berhasil maka ambil userID
		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		// jika user tidak ditemukan
		if err != nil {
			response := helper.APIResponse("Unauthorized 4", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response) // hentikan status jika tidak ada bearer token
			return
		}
		// jika tidak ada error yang lain kita set context isinya user yang melakukan request
		c.Set("currentUser", user) // context-nya sudah di set, dengan key yang namanya "currentUser"
	}
}

/*
Authentication Middleware Steps:
-> ambil nilai header Authorization: Bearer tokentokentoken
-> dari header Authorization, kita ambil nilai tokennya saja
-> kita validasi token
-> kita ambil user_id
-> ambil user dari db berdasarkan user_id lewat service
-> kita set context isinya user
*/
