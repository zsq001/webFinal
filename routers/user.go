/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package routers

import (
	"github.com/gofiber/fiber/v2"
	"webFinal/services"
)

func setUserRouterPub(r fiber.Router) {
	r.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	user := r.Group("/user")
	user.Post("/login", services.Login)
	user.Post("/register", services.Register)
}

func setupUserRouter(r fiber.Router) {
	user := r.Group("/user")
	user.Get("/info/:id", services.GetUserInfo)
	user.Get("/list", services.GetUserList)
	user.Post("/update/:id", services.UpdateUserInfo)
	user.Get("/delete/:id", services.DeleteUser)
	user.Get("/whoami", services.WhoAmI)
}
