/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package services

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"webFinal/database"
	"webFinal/models"
	"webFinal/models/web"
)

func DeletePicViaUid(uid uint) error {
	var p []models.Pic
	//delete pic from local
	if err := database.DB.Where("user_id = ?", uid).Find(&p).Error; err != nil {
		for _, v := range p {
			err := os.Remove("./static/" + v.UUID)
			if err != nil {
				return err
			}
		}
	}
	return database.DB.Where("user_id = ?", uid).Delete(&p).Error
}

//func FindUserPic(uid uint) ([]models.Pic, error) {
//	var p []models.Pic
//	if err := database.DB.Where("user_id = ?", uid).Find(&p).Error; err != nil {
//		return nil, err
//	}
//	return p, nil
//}

// @Summary get user info
// @Tags user
// @Produce json
// @Param id path string true "user id"
// @router /user/info/{id} [get]
// @Success 200 {object} web.User{data=web.UserInfo}
func GetUserInfo(c *fiber.Ctx) error {
	targetId := c.Params("id")

	var user models.User

	var result web.UserInfo

	if err := database.DB.Where("id = ?", targetId).Or("name = ?", targetId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.User{
			Status: fiber.StatusNotFound,
			Errors: "User not found",
		})
	}

	result.Role = (user.Role)
	result.Username = user.Name
	result.ID = user.ID

	return c.Status(fiber.StatusOK).JSON(web.User{
		Status: fiber.StatusOK,
		Data:   result,
	})
}

// @Summary Update user info
// @Tags user
// @Produce json
// @Param id path string true "user id"
// @router /user/info/{id} [post]
// @Success 200 {object} web.User{data=web.UserInfo}
func UpdateUserInfo(c *fiber.Ctx) error {
	targetId := c.Params("id")

	uId := c.Locals("sub").(uint)

	if targetId != strconv.FormatUint(uint64(uId), 10) && c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusUnauthorized).JSON(web.User{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	var user models.User

	var result web.UserInfo

	if err := database.DB.Where("id = ?", targetId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.User{
			Status: fiber.StatusNotFound,
			Errors: "User not found",
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.User{
			Status: fiber.StatusBadRequest,
			Errors: "Bad request",
		})
	}
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.User{
			Status: fiber.StatusInternalServerError,
			Errors: "Database error",
		})
	}

	result.Role = (user.Role)
	result.Username = user.Name
	result.ID = user.ID

	return c.Status(fiber.StatusOK).JSON(web.User{
		Status: fiber.StatusOK,
		Data:   result,
	})
}

// @Summary list all user info (admin only)
// @Tags user
// @Produce json
// @router /user/list [get]
// @Success 200 {object} []web.User{data=web.UserInfo}
func GetUserList(c *fiber.Ctx) error {
	var users []models.User
	//var result []web.UserInfo

	if c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusUnauthorized).JSON(web.User{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.User{
			Status: fiber.StatusInternalServerError,
			Errors: "Database error",
		})
	}
	//fmt.Printf("", users)

	return c.Status(fiber.StatusOK).JSON(web.User{
		Status: fiber.StatusOK,
		Data:   users,
	})
}

// @Summary delete user (admin only)
// @Tags user
// @Produce json
// @Param id path string true "user id"
// @router /user/delete/{id} [get]
// @Success 200 {string} string "success"
func DeleteUser(c *fiber.Ctx) error {
	targetId := c.Params("id")

	if c.Locals("role").(models.UserRole) != models.Admin {
		return c.Status(fiber.StatusUnauthorized).JSON(web.User{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	var user models.User

	if err := database.DB.Where("id = ?", targetId).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.User{
			Status: fiber.StatusNotFound,
			Errors: "User not found",
		})
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.User{
			Status: fiber.StatusInternalServerError,
			Errors: "Database error",
		})
	}

	if err := DeletePicViaUid(user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.User{
			Status: fiber.StatusInternalServerError,
			Errors: "Delete Pic error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(web.User{
		Status: fiber.StatusOK,
		Data:   "Success",
	})
}
