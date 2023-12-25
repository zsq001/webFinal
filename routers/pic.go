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

func SetupPicRouter(r fiber.Router) {
	pic := r.Group("/pic")
	pic.Get("/list", services.GetUserPic)
	pic.Post("/upload", services.UploadPic)
	pic.Get("/delete/:uuid", services.DeleteUserPic)
	pic.Get("/download/:uuid", services.DownloadPic)
	pic.Get("/list/all", services.ListAllUserPic)
}
