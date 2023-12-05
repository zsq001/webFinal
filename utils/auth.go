/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package utils

import "errors"

func CompareUserAndPass(Password string, userPassword string) error {
	if Password != userPassword {
		return errors.New("password is not correct")
	}
	return nil
}
