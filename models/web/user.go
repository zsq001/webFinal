/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package web

import "webFinal/models"

type User struct {
	Status int         `json:"status,omitempty"`
	Errors string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type UserInfo struct {
	ID       uint            `json:"id"`
	Username string          `json:"name"`
	Role     models.UserRole `json:"role"`
}
