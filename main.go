/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/sirupsen/logrus"
	"webFinal/config"
	"webFinal/database"
	"webFinal/routers"

	_ "webFinal/docs"
)

// @title webFinal API
// @version 1.0
// @license.name GPL 3.0
// @host localhost:8080
// @BasePath /
func main() {
	// read config
	if err := config.ReadConfig("./config.yaml"); err != nil {
		logrus.Error("Failed to read config file")
		logrus.Fatal(err)
	}

	if err := database.Init(); err != nil {
		logrus.Error("Failed to init database")
		logrus.Fatal(err)
	}

	app := fiber.New(fiber.Config{})

	app.Use(cors.New(cors.Config{
		AllowHeaders:  "Authorization,Origin,Content-Type,Accept",
		ExposeHeaders: "Content-Disposition",
	}))

	routers.InitRouter(app)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Static("/img", "./static")

	if err := app.Listen(config.CONFIG.BindAddr); err != nil {
		logrus.Error("Failed to listen")
		logrus.Panic(err)
	}
}
