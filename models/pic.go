/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package models

import (
	"gorm.io/gorm"
)

type Pic struct {
	gorm.Model
	Name       string
	UUID       string
	UserID     uint
	UploadTime int64 `gorm:"autoCreateTime"`
}
