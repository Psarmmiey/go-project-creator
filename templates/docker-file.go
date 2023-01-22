package templates

import (
	"fmt"
	"html/template"
	"os"
)

var docker = `
# syntax=docker/dockerfile:1
## Build
FROM golang:1.18-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o /{{.ProjectName}}

## Deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /{{.ProjectName}} /{{.ProjectName}}

EXPOSE 8080

USER nonroot:nonroot

ENV STELLAR_DB_CONNECTION_STRING=postgresql://horizon_plus:horizon_plus@localhost:1111/stellar
ENV STELLAR_DB_TYPE=postgres

ENTRYPOINT ["{{.ProjectName}}"]
`

func CreateDockerFile(projectName string) {
	// Create the file
	file, err := os.Create("Dockerfile")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// Write the file
	t := template.Must(template.New("docker").Parse(docker))
	err = t.Execute(file, struct{ ProjectName string }{projectName})
	if err != nil {
		fmt.Println(err)
	}
}
