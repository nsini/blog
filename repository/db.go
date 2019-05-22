package repository

import "github.com/jinzhu/gorm"

func NewDb() (*gorm.DB, error) {
	db, err := gorm.Open("mysql", "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	return db, nil
}
