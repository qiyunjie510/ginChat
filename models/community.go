package models

import (
	"errors"
	"ginchat/utils"
	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string
	OwnerId uint
	Img     string
	Desc    string
}

func (table *Community) TableName() string {
	return "community"
}

func CreateCommunity(community Community) error {
	if len(community.Name) == 0 {
		return errors.New("群名不能为空")
	}
	if community.OwnerId == 0 {
		return errors.New("所有者不能为空")
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		return err
	}
	return nil
}
