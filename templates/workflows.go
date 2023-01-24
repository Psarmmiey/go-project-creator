package templates

import (
	"fmt"
	"os"
)

var dockerImageYml = `# This workflow uses actions that are not certified by GitHub.
# They are provided by a third-party and are governed by
# separate terms of service, privacy policy, and support
# documentation.

# GitHub recommends pinning actions to a commit SHA.
# To get a newer version, you will need to update the SHA.
# You can also reference a tag or branch, but the action may change without warning.

name: Create and publish a Docker image

on:
  push:
    branches: ['release']

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  build-and-push-image:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Set variables
        run: |
          VER=$(cat version.txt)
          echo "VERSION=$VER" >> $GITHUB_ENV

      - name: Log in to the Container registry
        uses: docker/login-action@f054a8b539a109f9f41c372932f1ae047eff08c9
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Docker image
        uses: docker/build-push-action@ad44023a93711e3deb337508980b4b5e9bcdc5dc
        with:
          context: .
          push: true
          tags: |
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ env.VERSION }}
            ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
`

var releasePleaseYml = `
on:
  push:
    branches:
      - main
name: release-please
jobs:
  release-please:
    runs-on: ubuntu-latest
    steps:
      - uses: google-github-actions/release-please-action@v3
        with:
          release-type: simple
          package-name: release-please-action
`

var docValidationYml = `
name: Go doc checker

on:
  push:
    branches:
      - api-reference
      - main
      - dev
  pull_request:
    branches:
      - api-reference
      - main
      - dev

jobs:
  check-comments:
    runs-on: ubuntu-latest

    steps:
    - name: Checkout code
      uses: actions/checkout@v3
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.19.x'
    - name: Check Doc
      run: |
        go run github.com/Psarmmiey/check-comment@latest --path=$GITHUB_WORKSPACE
        if [ $? -eq 0 ]; then
          echo "All functions passed the comment check."
        else
          echo "Some functions failed the comment check."
          exit 1
        fi
`

func CreateDockerImageYml() string {
	return dockerImageYml
}

func CreateReleasePleaseYml() string {
	return releasePleaseYml
}

func CreateGithubWorkflow() {
	// Create the .github/workflows folder
	err := os.Mkdir(".github", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Mkdir(".github/workflows", os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}

	// Create the .github/workflows/release-please.yml file
	_, err = os.Create(".github/workflows/release-please.yml")
	if err != nil {
		fmt.Println(err)
	}

	// Open the .github/workflows/release-please.yml file for writing
	f, err := os.OpenFile(".github/workflows/release-please.yml", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {

		fmt.Println(err)
	}

	// Write the code to the file
	_, err = f.WriteString(releasePleaseYml)
	if err != nil {
		fmt.Println(err)
	}

	// Close the file
	f.Close()

	// Create the .github/workflows/docker-image.yml file
	_, err = os.Create(".github/workflows/docker-image.yml")
	if err != nil {
		fmt.Println(err)
	}

	// Open the .github/workflows/docker-image.yml file for writing
	f, err = os.OpenFile(".github/workflows/docker-image.yml", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}

	// Write the code to the file
	_, err = f.WriteString(dockerImageYml)
	if err != nil {
		fmt.Println(err)
	}

	// Close the file
	f.Close()

	// Create the .github/workflows/doc-validation.yml file
	_, err = os.Create(".github/workflows/doc-validation.yml")
	if err != nil {
		fmt.Println(err)
	}

	// Open the .github/workflows/doc-validation.yml file for writing
	f, err = os.OpenFile(".github/workflows/doc-validation.yml", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}

	// Write the code to the file
	_, err = f.WriteString(docValidationYml)
	if err != nil {
		fmt.Println(err)
	}

	// Close the file
	f.Close()

}
