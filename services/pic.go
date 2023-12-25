/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http"
	"os"
	db "webFinal/database"
	"webFinal/models"
	"webFinal/models/web"
)

//func CheckUserPerm(uid uint, pid uint) bool {
//	var pic models.Pic
//	db.DB.Where("id = ?", pid).First(&pic)
//	return pic.UserID == uid
//}

// @Summary upload pic
// @Tags pic
// @Accept mpfd
// @Produce json
// @router /pic/create [post]
// @Success 200 {object} web.Pic{data=models.Pic}
func UploadPic(c *fiber.Ctx) error {

	uId := c.Locals("sub")

	if uId == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(web.Pic{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	var Pic models.Pic

	file, err := c.FormFile("document")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: err.Error(),
		})
	}
	fileUUID, err := uuid.NewUUID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "UUid Gen error",
		})
	}

	if err := c.SaveFile(file, "./static/"+fileUUID.String()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "Save file error",
		})
	}
	targetFile, err := os.Open("./static/" + fileUUID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "Open file error",
		})
	}
	defer targetFile.Close()

	buffer := make([]byte, 512)
	n, err := targetFile.Read(buffer)
	contentType := http.DetectContentType(buffer[:n])
	if contentType != "image/jpeg" && contentType != "image/png" {
		os.Remove("./static/" + fileUUID.String())
		return c.Status(fiber.StatusForbidden).JSON(web.Pic{
			Status: fiber.StatusForbidden,
			Errors: "File type error",
		})
	}
	Pic.Name = file.Filename
	Pic.UUID = fileUUID.String()
	if err := db.DB.Create(&models.Pic{
		UserID: uId.(uint),
		UUID:   Pic.UUID,
		Name:   Pic.Name,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(web.Pic{
		Status: fiber.StatusOK,
		Data:   Pic,
	})
}

// @Summary delete user pic
// @Tags pic
// @Accept json
// @Produce json
// @Param uuid path string true "picture uuid"
// @router /pic/delete/{uuid} [get]
// @Success 200 {object} web.Pic{data=models.Pic}
func DeleteUserPic(c *fiber.Ctx) error {
	uId := c.Locals("sub").(uint)
	pId := c.Params("uuid")

	var pic models.Pic

	if err := db.DB.Where("uuid = ?", pId).First(&pic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "Image not found",
		})
	}

	if pic.UserID != uId && c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusForbidden).JSON(web.Pic{
			Status: fiber.StatusForbidden,
			Errors: "Permission denied",
		})
	}

	if err := db.DB.Delete(&pic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "Deletion error",
		})
	}

	if err := os.Remove("./static/" + pic.UUID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "File Delete error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(web.Pic{
		Status: fiber.StatusOK,
		Data:   nil,
	})
}

// @Summary Get user pic
// @Tags pic
// @Accept json
// @Produce json
// @router /pic/list/ [get]
// @Success 200 {object} web.Pic{data=[]models.Pic}
func GetUserPic(c *fiber.Ctx) error {
	uId := c.Locals("sub").(uint)

	var pictures []models.Pic

	db.DB.Where("user_id = ?", uId).Find(&pictures)

	return c.Status(fiber.StatusOK).JSON(web.Pic{
		Status: fiber.StatusOK,
		Data:   pictures,
	})
}

// @Summary Download pic
// @Tags pic
// @Accept json
// @router /pic/download/{uuid} [get]
// @Success 200 {formData} file
func DownloadPic(c *fiber.Ctx) error {
	uId := c.Locals("sub").(uint)
	pId := c.Params("uuid")

	var pic models.Pic

	if err := db.DB.Where("uuid = ?", pId).First(&pic).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.Pic{
			Status: fiber.StatusInternalServerError,
			Errors: "Image not found",
		})
	}

	if pic.UserID != uId && c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusForbidden).JSON(web.Pic{
			Status: fiber.StatusForbidden,
			Errors: "Permission denied",
		})
	}

	c.Set("Content-Disposition", "attachment; filename="+pic.Name)
	c.Set("Content-Type", "application/octet-stream")

	return c.SendFile("./static/" + pic.UUID)
}

// @Summary List all user pic (admin only)
// @Tags pic
// @Produce json
// @router /pic/list/all [get]
// @Success 200 {object} web.Pic{data=[]models.Pic}
func ListAllUserPic(c *fiber.Ctx) error {
	//uId := c.Locals("sub").(uint)

	var pictures []models.Pic

	if c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusForbidden).JSON(web.Pic{
			Status: fiber.StatusForbidden,
			Errors: "Permission denied",
		})
	}
	db.DB.Find(&pictures)

	return c.Status(fiber.StatusOK).JSON(web.Pic{
		Status: fiber.StatusOK,
		Data:   pictures,
	})
}
