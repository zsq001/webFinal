/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package web

type Pic struct {
	Status int         `json:"status"`
	Errors string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type ListPic struct {
	UserId uint   `json:"user_id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
}
