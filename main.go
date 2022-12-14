package main

import (
	"fmt"

	"github.com/kamalshkeir/kenv"
	"github.com/kamalshkeir/muzzsol/router"
	"github.com/kamalshkeir/muzzsol/services"
	"github.com/kamalshkeir/muzzsol/settings"

	"github.com/labstack/echo/v4"
)



func main() {
	// load config from .env and fill config strct in settings
	kenv.Load(".env")
	err := kenv.Fill(settings.Config)
	if err != nil {
		fmt.Println(err)
	}
	
	// init echo router
	e := echo.New()

	// init urls
	router.InitUrls(e)

	// init database and migrate
	err = services.NewDatabaseService().Init().Migrate()
	if err != nil {
		fmt.Println("error migrate:",err)
	}
	
	// Run Server
	fmt.Println("running on http://localhost:"+settings.Config.Port)
	e.Logger.Fatal(e.Start(":"+settings.Config.Port))
}