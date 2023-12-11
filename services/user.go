/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"webFinal/config"
	"webFinal/database"
	"webFinal/models"
	wb "webFinal/models/web"
	"webFinal/utils"
)

// jwt token process

func jwtSign(user models.User) (string, error) {
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"name": user.Name,
		"role": user.Role,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(time.Hour * 72).Unix(),
	})
	return rawToken.SignedString([]byte(config.CONFIG.JwtKey))
}

func SaveJWTtoLocal(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	if user.Valid {
		claims := user.Claims.(jwt.MapClaims)
		c.Locals("sub", uint(claims["sub"].(float64)))
		c.Locals("role", models.UserRole(claims["role"].(float64)))
	}
	return c.Next()
}

// @Summary Login
// @description Login api
// @description user can login with name
// @Param user formData string true "user name"
// @Param pass formData string true "user password"
// @Produce json
// @Success 200 {object} wb.User{data=string}
// @Failure 400 {object} wb.User{data=int}
// @Failure 401 {object} wb.User{data=int}
// @Failure 500 {object} wb.User{data=int}
// @Router /api/v1/user/login [post]
// @tags user
func Login(c *fiber.Ctx) error {
	user := c.FormValue("user")
	pass := c.FormValue("pass")
	if user == "" || pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user or pass is empty",
		})
	}
	var UserObject models.User

	// check if user exist
	if err := database.DB.Where("name = ?", user).First(&UserObject).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(wb.User{
			Status: fiber.StatusUnauthorized,
			Errors: "user not found or password error",
			Data:   nil,
		})
	}

	// check if password is correct
	if err := utils.CompareUserAndPass(UserObject.Password, pass); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(wb.User{
			Status: fiber.StatusUnauthorized,
			Errors: "user not found or password error",
			Data:   nil,
		})
	}

	token, err := jwtSign(UserObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wb.User{
			Status: fiber.StatusInternalServerError,
			Errors: err.Error(),
			Data:   nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(wb.User{
		Status: fiber.StatusOK,
		Data:   token,
	})
}

// @Summary Register
// @description Register api
// @description user can register with name
// @Param user formData string true "user name"
// @Param pass formData string true "user password"
// @Produce json
// @Success 200 {object} wb.User{data=string}
// @Failure 400 {object} wb.User{data=int}
// @Failure 401 {object} wb.User{data=int}
// @Failure 500 {object} wb.User{data=int}
// @Router /api/v1/user/register [post]
// @tags user
func Register(c *fiber.Ctx) error {
	name := c.FormValue("user")
	pass := c.FormValue("pass")
	if name == "" || pass == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user or pass is empty",
		})
	}
	var UserObject models.User
	UserObject.Name = name
	UserObject.Password = pass // For homework only, DO NOT USE IN PRODUCTION
	UserObject.Role = models.Normal
	if err := database.DB.Create(&UserObject).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wb.User{
			Status: fiber.StatusInternalServerError,
			Errors: "Create User Internal Error",
			Data:   nil,
		})
	}

	token, err := jwtSign(UserObject)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(wb.User{
			Status: fiber.StatusInternalServerError,
			Errors: err.Error(),
			Data:   nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(wb.User{
		Status: fiber.StatusOK,
		Data:   token,
	})
}
