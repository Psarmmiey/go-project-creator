package templates

import (
	"log"
	"os"
)

var middleware = `/*
Copyright Interstellar, Inc - All Rights Reserved.
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
Written by Fritz Ekwoge (fritz@interstellar.cm), March 2021.
*/
package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-BANTUPAY-PUBLIC-KEY, X-BANTUPAY-SIGNATURE")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
`

func CreateMiddleware(filepath string) {

	// Create the middleware/main.go file
	mainFilePath := filepath + "/cors.go"
	_, err := os.Create(mainFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Open the middleware/main.go file for writing
	f, err := os.OpenFile(mainFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	// Write the code to the file
	_, err = f.WriteString(middleware)
	if err != nil {
		log.Fatal(err)
	}

	// Close the file
	f.Close()
}
