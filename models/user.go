/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package models

import "gorm.io/gorm"

type UserRole int64

type User struct {
	gorm.Model
	Name     string
	Role     UserRole
	Pics     []Pic  `gorm:"foreignKey:UserID"`
	Password string `gorm:"not null"`
}

const (
	Normal UserRole = iota
	Admin
)

func (u *User) IsAdmin() bool {
	return u.Role == Admin
}
