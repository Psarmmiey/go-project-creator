package templates

import (
	"fmt"
	"os"
)

var envFile = `
STELLAR_DB_CONNECTION_STRING=postgres://stellar:n0p2d3kUAcFwhi8e@localhost:5431/core
STELLAR_DB_TYPE=postgres

HORIZON_DB_CONNECTION_STRING=postgres://stellar:n0p2d3kUAcFwhi8e@localhost:5431/horizon
HORIZON_DB_TYPE=postgres
`

func CreateEnvFile() {
	// Create the file
	file, err := os.Create(".env")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Write the file
	_, err = file.WriteString(envFile)
	if err != nil {
		fmt.Println(err)
	}

	// Create env.sample file
	file, err = os.Create(".env.sample")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Write the file
	_, err = file.WriteString(envFile)
	if err != nil {
		fmt.Println(err)
	}

}
