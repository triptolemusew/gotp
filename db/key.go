package db

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Key struct {
	ID        uint `gorm:"primaryKey"`
	Type      string
	Issuer    string
	Account   string
	Secret    string
	Algorithm string
	Counter   uint64
	Digits    int
	Period    int
	CreatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
}

func GetAll(db *gorm.DB) ([]Key, error) {
	var keys []Key

	result := db.Find(&keys)
	if result.Error != nil {
		return keys, result.Error
	}

	return keys, nil
}

func Get(db *gorm.DB, id uint) (Key, error) {
	var key Key

	result := db.First(&key, id)
	if result.Error != nil {
		return key, result.Error
	}

	return key, nil
}

func Create(db *gorm.DB, key *Key) error {
	if result := db.Create(&key); result.Error != nil {
		return result.Error
	}
	return nil
}

func FilterByIssuerAndAccount(keys []Key, term string) []Key {
	var out []Key
	for _, each := range keys {
		if strings.Contains(each.Issuer, term) || strings.Contains(each.Account, term) {
			out = append(out, each)
		}
	}

	return out
}
