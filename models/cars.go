package models

import "gorm.io/gorm"

type Cars struct {
	ID    uint    `gorm:"primary key;autoIncrement" json:"id"`
	Merk  *string `json:"merk"`
	Tipe  *string `json:"tipe"`
	Warna *string `josn:"warna"`
}

func MigrateCars(db *gorm.DB) error{
	err := db.AutoMigrate (&Cars{})
	return err
}