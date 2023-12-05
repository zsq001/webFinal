/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package routers

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"webFinal/config"
	"webFinal/services"
)

func InitRouter(fiber *fiber.App) {
	r := fiber.Group("/api/v1")

	setUserRouterPub(r)

	SetUpJwtTokenMiddleware(r)

	SetupPicRouter(r)
	setupUserRouter(r)
	SetupMsgRouterv1(r)
}

func SetUpJwtTokenMiddleware(r fiber.Router) {
	r.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(config.CONFIG.JwtKey),
	}))
	r.Use(services.SaveJWTtoLocal)
}
