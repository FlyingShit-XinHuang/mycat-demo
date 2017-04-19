package space

import (
	"time"
)

/*
CREATE TABLE `whispir_spaces` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

type Space struct {
	Id        uint   `gorm:"primary_key;AUTO_INCREMENT;not null"`
	Name      string `gorm:"type:varchar(100);not null"`
	// This field could be automatically initialized by gorm when inserting a record
	CreatedAt time.Time
	// This field could be automatically set by gorm when soft deleteing a record
	DeletedAt *time.Time
}

// Specify table name with this method
func (Space) TableName() string {
	return "whispir_spaces"
}

func NewSpace(name string) *Space {
	return &Space{
		Name: name,
	}
}
