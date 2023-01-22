package templates

import (
	"fmt"
	"html/template"
	"os"
)

var main = `
package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"{{.ProjectModule}}/internal/db"
	"{{.ProjectModule}}/internal/components/middleware"
	_ "{{.ProjectModule}}/docs"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)




// @title {{.ProjectName}} API Reference
// @description {{.ProjectDesc}}.
// @version 1.0.0
// @BasePath /
// @schemes http https
func main() {
	rand.Seed(time.Now().UnixNano())

	log.Info().Msg("hello world")

	errEnv := godotenv.Load()
	if errEnv != nil {
		path, _ := os.Getwd()
		log.Printf("could not find or load any .env file from %v...skipping...\n", path)
	}

	//var stellarDatabase *gorm.DB
	var horizonDatabase *gorm.DB

	{

		{

			exit := false
			requiredEnvironmentVariables := []string{"HORIZON_DB_CONNECTION_STRING", "HORIZON_DB_TYPE"}

			for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
				if len(os.Getenv(requiredEnvironmentVariable)) == 0 {
					log.Printf("Required environment variable is missing %v", requiredEnvironmentVariable)
					exit = true
				}
			}

			if exit {
				panic("missing env variables")
			}

		}

		var err error
		// stellarDatabase, err = db.OpenDb("STELLAR")

		// if err != nil {
		// 	log.Fatal().Err(err).Msg("Error opening Stellar DB")
		// 	return
		// }

		horizonDatabase, err = db.OpenDb("HORIZON")

		if err != nil {
			log.Fatal().Err(err).Msg("Error opening Horizon DB")
			return
		}

	}

	var router *gin.Engine = gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.CORSMiddleware())
	// Todo: add other routes
	router.Run()

}
`

func CreateMainGoFile(projectName, projectDesc, projectModule string) {
	// Create the file
	file, err := os.Create("main.go")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Write the file
	t := template.Must(template.New("main").Parse(main))
	err = t.Execute(file, struct{ ProjectName, ProjectDesc, ProjectModule string }{projectName, projectDesc, projectModule})
	if err != nil {
		fmt.Println(err)
	}

}
