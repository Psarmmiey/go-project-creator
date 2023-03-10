package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Psarmmiey/go-project-creator/templates"
	"github.com/Psarmmiey/go-project-creator/utils"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

type tomlConfig struct {
	Project struct {
		Name        string
		Module      string
		Description string
	}
	Folders struct {
		Internal []string
	}
	Components map[string][]string
	Models     map[string][]model
	Options    struct {
		Db              bool
		Errors          bool
		Middleware      bool
		GithubWorkflows bool
		Docker          bool
		MakeFile        bool
		Env             bool
	}
}

type model struct {
	Name   string
	Fields []string
	CRUD   bool
}

func main() {
	app := cli.NewApp()
	app.Name = "go-project-create"
	app.Usage = "Create a new go project with Interstellar format"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "config, c",
			Usage:    "Path to load the config file",
			Required: false,
		},
		cli.StringFlag{
			Name:  "component, comp",
			Usage: "Create a new component",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("Creating project structure...")
		configFile := c.String("config")
		if configFile != "" {

			var config tomlConfig
			if _, err := toml.DecodeFile(configFile, &config); err != nil {
				fmt.Println(err)
				return err
			}
			createEntireStructure(config)
			runCommands(config)
		}

		if c.String("component") != "" {
			// create the component
			createComponent(c.String("component"))
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func createEntireStructure(config tomlConfig) {
	// create the internal folder
	err := os.MkdirAll("internal", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, folder := range config.Folders.Internal {
		err := os.MkdirAll(filepath.Join("internal", folder), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		if folder == "db" {
			if config.Options.Db {
				// create the db folder
				err := os.MkdirAll(filepath.Join("internal", folder), os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
				// create the db code
				templates.CreateDB()
			}
		}

		if folder == "errors" {
			if config.Options.Errors {
				// create the responses folder
				err := os.MkdirAll(filepath.Join("internal", folder), os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
				// create the responses code
				templates.CreateResponses("internal/errors")
			}
		}

		if folder == "components" {
			for k, v := range config.Components {
				err := os.MkdirAll(filepath.Join("internal", folder, k), os.ModePerm)
				if err != nil {
					fmt.Println(err)
					return
				}
				for _, subfolder := range v {
					err := os.MkdirAll(filepath.Join("internal", folder, k, subfolder), os.ModePerm)
					if err != nil {
						fmt.Println(err)

						return
					}
					if subfolder == "models" {
						file, err := os.OpenFile(filepath.Join("internal", folder, k, subfolder, "main.go"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							fmt.Println(err)
							return
						}
						defer file.Close()

						// check if package name is already written
						f, _ := os.Open(filepath.Join("internal", folder, k, subfolder, "main.go"))
						defer f.Close()
						packageCheck, _ := io.ReadAll(f)
						if !strings.Contains(string(packageCheck), "package "+subfolder) {
							_, err = file.WriteString(fmt.Sprintf("package %s\n\n", subfolder))
							if err != nil {
								fmt.Println(err)
								return
							}
						}
						for _, model := range config.Models[k] {

							_, err = file.WriteString(fmt.Sprintf("type %s struct{\n", model.Name))
							if err != nil {
								fmt.Println(err)
								return
							}
							for _, field := range model.Fields {
								_, err = file.WriteString(fmt.Sprintf("\t%s string `json:\"%s\" gorm:\"column:%s\"`\n", utils.ToTitle(field), utils.ToLowerCamel(field), utils.ToUnderscore(field)))
								if err != nil {
									fmt.Println(err)
									return
								}
							}
							_, err = file.WriteString("}\n")
							if err != nil {
								fmt.Println(err)
								return
							}
						}
					} else if subfolder == "services" {
						file, err := os.Create(filepath.Join("internal", folder, k, subfolder, "main.go"))
						if err != nil {
							fmt.Println(err)
							return
						}
						defer file.Close()
						_, err = file.WriteString(fmt.Sprintf("package %s\n\n", subfolder))
						if err != nil {
							fmt.Println(err)
							return
						}
						for _, model := range config.Models[k] {
							if model.CRUD {
								_, err = file.WriteString(fmt.Sprintf("func Create%s(){\n\t//TODO: Add Create Logic\n}\n", model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Read%s(){\n\t//TODO: Add Read Logic\n}\n", model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Update%s(){\n\t//TODO: Add Update Logic\n}\n", model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Delete%s(){\n\t//TODO: Add Delete Logic\n}\n", model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}
					} else if subfolder == "controllers" {
						file, err := os.Create(filepath.Join("internal", folder, k, subfolder, "main.go"))
						if err != nil {
							fmt.Println(err)
							return
						}
						defer file.Close()
						var copyright = `/*
Copyright Interstellar, Inc - All Rights Reserved.
Unauthorized copying of this file, via any medium is strictly prohibited.
Proprietary and confidential.
*/`

						var imports = `
import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"{{.Project.Module}}/internal/components/products/services"

)`
						_, err = file.WriteString(fmt.Sprintf("%s\n", copyright))
						if err != nil {
							fmt.Println(err)
							return
						}
						_, err = file.WriteString(fmt.Sprintf("package %s\n\n", subfolder))
						if err != nil {
							fmt.Println(err)
							return
						}

						t, err := template.New("imports").Parse(imports)
						if err != nil {
							fmt.Println(err)
							return
						}
						err = t.Execute(file, config)
						if err != nil {
							fmt.Println(err)
							return
						}

						// Create Init Function with router and db as parameters and add routes
						_, err = file.Write([]byte("\n\nfunc Init(router *gin.Engine, db *gorm.DB){\n"))
						if err != nil {
							fmt.Println(err)
							return
						}
						_, err = file.Write([]byte("\t//TODO: Add Routes\n}\n"))
						if err != nil {
							fmt.Println(err)
							return
						}
						for _, model := range config.Models[k] {

							if model.CRUD {
								_, err = file.WriteString(fmt.Sprintf("func Create%s(){\n\t//TODO: Add Create Logic\n\tservices.Create%s()\n}\n", model.Name, model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Read%s(){\n\t//TODO: Add Read Logic\n\tservices.Read%s()\n}\n", model.Name, model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Update%s(){\n\t//TODO: Add Update Logic\n\tservices.Update%s()\n}\n", model.Name, model.Name))
								if err != nil {
									fmt.Println(err)

									return
								}
								_, err = file.WriteString(fmt.Sprintf("func Delete%s(){\n\t//TODO: Add Delete Logic\n\tservices.Delete%s()\n}\n", model.Name, model.Name))
								if err != nil {
									fmt.Println(err)
									return
								}
							}
						}
					}
				}
			}
		}

		// Check if middleware is enabled
		if config.Options.Middleware {
			// check if components folder exists
			if _, err := os.Stat("internal/components"); os.IsNotExist(err) {
				// create components folder
				err = os.Mkdir("internal/components", os.ModePerm)
				if err != nil {
					fmt.Println(err)
				}
			}
			// create middleware folder
			err = os.Mkdir("internal/components/middleware", os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
			// create middleware file
			templates.CreateMiddleware("internal/components/middleware")
		}
	}

	// Check if the "db" option is enabled
	if config.Options.Db {
		// copy the db code
		templates.CreateDB()
	}

	// Check if GithubWorkflow is enabled
	if config.Options.GithubWorkflows {
		// Create the .github/workflows/main.yml file
		templates.CreateGithubWorkflow()
	}

	// check if makefile is enabled
	if config.Options.MakeFile {
		// create makefile
		templates.CreateMakeFile()
	}

	// check if dockerfile is enabled
	if config.Options.Docker {
		// create dockerfile
		templates.CreateDockerFile(config.Project.Name)
	}

	// check if env file is enabled
	if config.Options.Env {
		// create env file
		templates.CreateEnvFile()
	}

	templates.CreateMainGoFile(config.Project.Name, config.Project.Description, config.Project.Module)
}

func runCommands(config tomlConfig) {

	// go mod init
	fmt.Println("Initializing go module...")
	cmd := exec.Command("go", "mod", "init", config.Project.Module)
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// swagger init
	fmt.Println("Initializing swagger...")
	cmd = exec.Command("go", "run", "github.com/swaggo/swag/cmd/swag@latest", "init")
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// git init
	fmt.Println("Initializing git...")
	cmd = exec.Command("git", "init")
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// go get
	fmt.Println("Getting dependencies...")
	cmd = exec.Command("go", "get")
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// go mod tidy
	fmt.Println("Tidying dependencies...")
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// go fmt
	fmt.Println("Formatting code...")
	cmd = exec.Command("go", "fmt")
	cmd.Dir = "./"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	// output
	fmt.Println("Done!")

}

func createComponent(componentName string) {
	// check if internal folder exists
	if _, err := os.Stat("internal"); os.IsNotExist(err) {
		// create internal folder
		err = os.Mkdir("internal", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	// check if components folder exists
	if _, err := os.Stat("internal/components"); os.IsNotExist(err) {
		// create components folder
		err = os.Mkdir("internal/components", os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	// create component directory
	err := os.Mkdir("internal/components/"+componentName, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	// create model directory
	err = os.Mkdir("internal/components/"+componentName+"/model", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	// create service directory
	err = os.Mkdir("internal/components/"+componentName+"/service", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	// create controller directory
	err = os.Mkdir("internal/components/"+componentName+"/controller", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	// create model file
	file, err := os.Create("internal/components/" + componentName + "/model/" + "main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	var modelFile = `
	package model

	import "time"

	type ` + componentName + ` struct {
		` + utils.ToTitle(componentName) + `ID string ` + "`json:\"" + componentName + "ID\"`" + `
		CreatedAt time.Time ` + "`json:\"createdAt\"`" + `
		UpdatedAt time.Time ` + "`json:\"updatedAt\"`" + `
	}
	`

	// write to model file
	_, err = file.WriteString(modelFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// create service file
	file, err = os.Create("internal/components/" + componentName + "/service/" + "main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// write to service file
	_, err = file.WriteString("package service \r \r")
	if err != nil {
		fmt.Println(err)
		return
	}

	// create controller file
	file, err = os.Create("internal/components/" + componentName + "/controller/" + "main.go")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// create response
	var routeSample = `
	package controller

	import (
		"net/http"

		"github.com/gin-gonic/gin"
		"github.com/jinzhu/gorm"
	)


	func Init(router *gin.Engine, db *gorm.DB) {
		group := router.Group("/{{.componentName}}") 
		group.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Hello World!"})
		})
	}
	`
	tmpl, err := template.New("main").Parse(routeSample)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = tmpl.Execute(file, map[string]string{"componentName": componentName})
	if err != nil {
		fmt.Println(err)
		return
	}

	// write to controller file

}
