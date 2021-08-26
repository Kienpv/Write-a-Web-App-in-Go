package main

import (
	"net/http"
	"github.com/Write-a-Web-App-in-Go/models"
	"github.com/Write-a-Web-App-in-Go/utils"
	"github.com/Write-a-Web-App-in-Go/routes"
)

func main() {
	models.Init()
	utils.LoadTemplates("templates/*.html")
	r := routes.NewRoute() 
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
