package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func Connect(user, password, host, database string) (err error) {
	db, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", user, password, host, database))
	return
}

func GetDB() *gorm.DB {
	if nil == db {
		panic(fmt.Errorf("database handle "))
	}
	return db
}

func Create(obj interface{}) (int64, error) {
	d := db.Create(obj)
	return d.RowsAffected, d.Error
}

func FindByPK(obj interface{}, pk interface{}) error {
	return db.First(obj, pk).Error
}

func Update(obj interface{}, query interface{}, args ...interface{}) (int64, error) {
	d := db.Model(obj).Where(query, args...).Updates(obj)
	return d.RowsAffected, d.Error
}

func Delete(objType interface{}, query interface{}, args ...interface{}) (int64, error) {
	d := db.Where(query, args...).Delete(objType)
	return d.RowsAffected, d.Error
}
