package templates

import (
	"log"
	"os"
)

var responses = `package errors

import "net/http"

type JsonResult struct {
	Status     string      ` + "`json:\"status\"`" + `
	Message    string      ` + "`json:\"message\"`" + `
	Data       interface{} ` + "`json:\"data\"`" + `
	HttpStatus int         ` + "`json:\"-\"`" + `
}

func NewSuccessResult(label string, data interface{}) *JsonResult {
	var jsonResult JsonResult
	jsonResult.Data = map[string]interface{}{label: data}
	jsonResult.Status = "success"
	jsonResult.HttpStatus = http.StatusOK
	return &jsonResult
}

func NewFailResult(err error, data interface{}, httpStatusCode int) *JsonResult {
	var jsonResult JsonResult
	jsonResult.Data = data
	jsonResult.Status = "fail"
	jsonResult.Message = err.Error()
	jsonResult.HttpStatus = httpStatusCode
	return &jsonResult
}

func NewFailBadRequest(err error, key string, value string) *JsonResult {
	var jsonResult JsonResult
	jsonResult.Data = map[string]string{key: value}
	jsonResult.Status = "fail"
	jsonResult.Message = err.Error()
	jsonResult.HttpStatus = http.StatusBadRequest
	return &jsonResult
}

func NewErrorResult(err error, data interface{}) *JsonResult {
	var jsonResult JsonResult
	jsonResult.Data = data
	jsonResult.Status = "error"
	jsonResult.Message = err.Error()
	jsonResult.HttpStatus = http.StatusInternalServerError
	return &jsonResult
}
`

func CreateResponses(filepath string) {

	// Create the responses/main.go file
	mainFilePath := filepath + "/main.go"
	_, err := os.Create(mainFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Open the responses/main.go file for writing
	f, err := os.OpenFile(mainFilePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}

	// Write the code to the file
	_, err = f.WriteString(responses)
	if err != nil {
		log.Fatal(err)
	}

	// Close the file
	f.Close()
}
