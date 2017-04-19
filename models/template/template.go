package template

import (
	"time"
	"whispir/mycat-demo/models"
	"whispir/mycat-demo/models/space"
)

/*
CREATE TABLE `whispir_templates` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `name` varchar(100) NOT NULL,
  `content` blob NOT NULL,
  `space_id` int(10) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `space_id` (`space_id`),
  CONSTRAINT `whispir_templates_ibfk_1` FOREIGN KEY (`space_id`) REFERENCES `whispir_spaces` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
*/

type Template struct {
	Id        uint `gorm:"primary_key;AUTO_INCREMENT;not null"`
	CreatedAt time.Time
	// This field could be automatically set by gorm when updating a record
	UpdatedAt time.Time
	DeletedAt *time.Time

	Name    string `gorm:"type:varchar(100);not null"`
	Content []byte `gorm:"type:blob;not null"`

	SpaceId uint `gorm:"not null"`
}

// Specify table name with this method
func (Template) TableName() string {
	return "whispir_templates"
}

func NewTemplate(name, content string, space *space.Space) *Template {
	return &Template{
		Name:    name,
		Content: []byte(content),
		SpaceId: space.Id,
	}
}

func ListInSpace(space *space.Space) (tmpls []Template, err error) {
	err = models.GetDB().Find(&tmpls, Template{SpaceId: space.Id}).Error
	return
}
