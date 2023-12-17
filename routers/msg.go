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

func SetupMsgRouterv1(r fiber.Router) {
	base := r.Group("/msg")
	msg := base.Group("/v1")
	msg.Get("/get/:tid", services.GetMsgs)
	msg.Post("/send/:tid", services.SendMsg)
	msg.Get("/list", services.ListMsgs)
	msg.Get("/delete/:mid", services.DeleteMsg)
	msg.Get("/recall/:mid", services.RecallMsg)
}
