package database

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	connStr := "root:123456@tcp(10.11.1.181:3306)/teenpatti"
	masterDB := mysql.Open(connStr)
	db, err := gorm.Open(masterDB)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(db)
}
