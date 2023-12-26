/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
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

func isNumeric(str string) bool {
	_, errInt := strconv.Atoi(str)             // 尝试将字符串转为整数
	_, errFloat := strconv.ParseFloat(str, 64) // 尝试将字符串转为浮点数

	return errInt == nil || errFloat == nil
}

// @Summary Login
// @description Login api
// @description user can login with name
// @Param user formData string true "user name"
// @Param pass formData string true "user password"
// @Produce json
// @Success 200 {object} wb.User{data=wb.User}
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
// @Success 200 {object} wb.User{data=wb.User}
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
	if isNumeric(name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user name can't be number",
		})
	}

	// check if user exist
	var UserObject models.User

	database.DB.Where("name = ?", name).First(&UserObject)
	if UserObject.Name != "" {
		return c.Status(fiber.StatusBadRequest).JSON(wb.User{
			Status: fiber.StatusBadRequest,
			Errors: "user already exist",
			Data:   nil,
		})
	}

	// create user
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

// @Summary check id
// @description check id api
// @description user can check id of themselves
// @Produce json
// @Success 200 {object} wb.User{data=wb.User}
// @Router /api/v1/user/whoami [get]
// @tags user
func WhoAmI(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)
	return c.Status(fiber.StatusOK).JSON(wb.User{
		Status: fiber.StatusOK,
		Data:   uid,
	})
}
