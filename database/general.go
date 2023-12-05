/*
 * Copyright (c) 2023 zsq001
 * Licensed under the GNU General Public License v3.0
 * See https://www.gnu.org/licenses/gpl-3.0.html for details.
 */

package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"webFinal/config"
	m "webFinal/models"
)

var DB *gorm.DB

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(m.User{})
	err = db.AutoMigrate(m.Pic{})
	err = db.AutoMigrate(m.Msg{})
	if err != nil {
		return err
	}
	return nil
}

func Init() error {
	var err error
	DB, err = sqlInit()
	if err != nil {
		return err
	}
	return nil
}

func sqlInit() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.CONFIG.Database.Path), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return db, nil
}
