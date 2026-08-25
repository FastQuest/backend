package main

import (
	"flashquest/internal/ai"
	"flashquest/internal/platform/database"
	"fmt"
	"log"
)

// @title FastQuest API
// @version 1.0
// @description Coloque sua descrição aqui
// @termsOfService http://swagger.io/terms/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	fmt.Println("Running backend")
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	ai.InitGemini()

	srv := NewServer(db)
	log.Fatal(srv.ListenAndServe())
}
