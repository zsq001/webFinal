/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package services

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
	db "webFinal/database"
	"webFinal/models"
	mw "webFinal/models/web"
)

func FindUserNameById(id uint) string {
	var user models.User
	db.DB.First(&user, id)
	return user.Name
}

// @Summary Send
// @description message api
// @description user can send msg with name
// @Produce json
// @Success 200 {object} wb.Msg{data=string}
// @Router /api/v1/msg/v1/send/{tid} [post]
// @tags msg

func SendMsg(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)
	tid, err := strconv.ParseUint(c.Params("tid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mw.Msg{
			Status: fiber.StatusBadRequest,
			Errors: "target Id invalid",
		})
	}

	var message models.Msg

	msg := c.FormValue("Content")

	flag := c.Locals("role").(models.UserRole) >= models.Normal
	if !flag {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	message.SendUserID = uid
	message.ReceiveUserID = uint(tid)
	message.Content = msg

	if err := db.DB.Create(&models.Msg{
		SendUserID:    message.SendUserID,
		ReceiveUserID: message.ReceiveUserID,
		Content:       message.Content,
	}).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mw.Msg{
			Status: fiber.StatusInternalServerError,
			Errors: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mw.Msg{
		Status: fiber.StatusOK,
		Data:   message,
	})
}

// @Summary Get
// @description message api
// @description user can get exist msg with target id
// @Produce json
// @Success 200 {object} wb.Msg{data=string}
// @Router /api/v1/msg/v1/get/{tid} [get]
// @tags msg

func GetMsgs(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)
	tid, err := strconv.ParseUint(c.Params("tid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mw.Msg{
			Status: fiber.StatusBadRequest,
			Errors: "target Id invalid",
		})
	}

	var messages []models.Msg

	flag := c.Locals("role").(models.UserRole) >= models.Normal
	if !flag {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if err := db.DB.Where("send_user_id = ? AND receive_user_id = ?", uid, tid).Or("send_user_id = ? AND receive_user_id = ?", tid, uid).Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mw.Msg{
			Status: fiber.StatusInternalServerError,
			Errors: "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mw.Msg{
		Status: fiber.StatusOK,
		Data:   messages,
	})
}

// @Summary List
// @description message api
// @description user can get exist msg list
// @Produce json
// @Success 200 {object} wb.Msg{data=string}
// @Router /api/v1/msg/v1/list [get]
// @tags msg

func ListMsgs(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)

	var messages []models.Msg

	flag := c.Locals("role").(models.UserRole) >= models.Normal
	if !flag {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if err := db.DB.Where("send_user_id = ? OR receive_user_id = ?", uid, uid).Last(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mw.Msg{
			Status: fiber.StatusInternalServerError,
			Errors: "Internal Server Error",
		})
	}

	var result []mw.MsgList
	for _, v := range messages {
		var tmp mw.MsgList
		if v.IsRecall == true {
			continue
		}
		if v.SendUserID == uid {
			if v.SendVisible == false {
				continue
			}
			tmp.Name = FindUserNameById(v.ReceiveUserID)
		} else {
			if v.ReceiveVisible == false {
				continue
			}
			tmp.Name = FindUserNameById(v.SendUserID)
		}
		tmp.LatestMsg = v.Content
		tmp.LatestMsgTime = strconv.FormatInt(v.SendTime, 10)
		result = append(result, tmp)
	}

	return c.Status(fiber.StatusOK).JSON(mw.Msg{
		Status: fiber.StatusOK,
		Data:   result,
	})
}

// @Summary Delete
// @description message api
// @description user can delete exist msg with message id
// @Produce json
// @Success 200 {object} wb.Msg{data=string}
// @Router /api/v1/msg/v1/delete/:tid [get]
// @tags msg

func DeleteMsg(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)
	mid, err := strconv.ParseUint(c.Params("mid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mw.Msg{
			Status: fiber.StatusBadRequest,
			Errors: "message Id invalid",
		})
	}

	var message models.Msg

	flag := c.Locals("role").(models.UserRole) >= models.Normal
	if !flag {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if err := db.DB.Where("id = ?", mid).First(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mw.Msg{
			Status: fiber.StatusInternalServerError,
			Errors: "Internal Server Error",
		})
	}

	if message.SendUserID == uid {
		message.SendVisible = false
	} else if message.ReceiveUserID == uid {
		message.ReceiveVisible = false
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	return c.Status(fiber.StatusOK).JSON(mw.Msg{
		Status: fiber.StatusOK,
		Data:   message,
	})
}

// @Summary Recall
// @description message api
// @description user can recall msg sent in 10 minutes
// @Produce json
// @Success 200 {object} wb.Msg{data=string}
// @Router /api/v1/msg/v1/recall/:tid [get]
// @tags msg

func RecallMsg(c *fiber.Ctx) error {
	uid := c.Locals("sub").(uint)
	mid, err := strconv.ParseUint(c.Params("mid"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(mw.Msg{
			Status: fiber.StatusBadRequest,
			Errors: "message Id invalid",
		})
	}

	var message models.Msg

	flag := c.Locals("role").(models.UserRole) >= models.Normal
	if !flag {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if err := db.DB.Where("id = ?", mid).First(&message).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(mw.Msg{
			Status: fiber.StatusInternalServerError,
			Errors: "Internal Server Error",
		})
	}

	if message.SendUserID != uid {
		return c.Status(fiber.StatusUnauthorized).JSON(mw.Msg{
			Status: fiber.StatusUnauthorized,
			Errors: "Unauthorized",
		})
	}

	if message.SendTime+600 < time.Now().Unix() {
		return c.Status(fiber.StatusBadRequest).JSON(mw.Msg{
			Status: fiber.StatusBadRequest,
			Errors: "Recall time expired",
		})
	}

	message.IsRecall = true

	return c.Status(fiber.StatusOK).JSON(mw.Msg{
		Status: fiber.StatusOK,
		Data:   message,
	})
}
