/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package models

import "gorm.io/gorm"

type Msg struct {
	gorm.Model
	Content        string
	SendUserID     uint
	ReceiveUserID  uint
	SendTime       int64 `gorm:"autoCreateTime"`
	SendVisible    bool  `gorm:"default:true"`
	ReceiveVisible bool  `gorm:"default:true"`
	IsRecall       bool  `gorm:"default:false"`
}
