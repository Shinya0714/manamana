package handler

import (
	"github.com/Shinya0714/manamana/go/app/general"
	"github.com/Shinya0714/manamana/go/app/rooting"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Handler() {

	general.LoadEnv()

	e := echo.New()
	e.Use(middleware.CORS())

	e.GET("/balance", rooting.GetBalance)
	e.GET("/schedule", rooting.GetSchedule)
	e.GET("/sbiBookBuilding/:tickerSymbol", rooting.SbiBookBuilding)
	e.GET("/mizuhoBookBuilding/:tickerSymbol", rooting.MizuhoBookBuilding)

	e.Logger.Fatal(e.Start(":8000"))
}
