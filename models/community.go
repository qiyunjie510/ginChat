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

func JoinCommunity(ownerId uint, dstId uint) error {
	// 是否存在该群

	community := Community{}
	community.ID = dstId
	if err := utils.DB.Where("id = ? ", dstId).Find(&community).Error; err != nil || community.Name == "" {
		return errors.New("未找到该群")
	}
	user := UserBasic{}
	user.ID = ownerId
	if err := utils.DB.Where("id = ? ", ownerId).Find(&user).Error; err != nil || user.Name == "" {
		return errors.New("未找到当前登录人信息")
	}
	contact := Contact{}
	contact.OwnerId = ownerId
	contact.TargetId = dstId
	contact.Type = 2
	if err := utils.DB.Save(&contact).Error; err != nil {
		return err
	}
	return nil
}

func LoadCommunity(contact Contact) ([]Community, error) {
	var contactList []Contact
	if err := utils.DB.Where("owner_id = ? and type = 2", contact.OwnerId).Find(&contactList).Error; err != nil {
		return nil, errors.New("查询好友关系出错")
	}
	var groupIDs []uint
	for _, v := range contactList {
		groupIDs = append(groupIDs, v.TargetId)
	}
	var community []Community
	if err := utils.DB.Where("id in ?", groupIDs).Find(&community).Error; err != nil {
		return nil, errors.New("查询群信息出错")
	}
	return community, nil
}

func getGroupNumberList(groupId uint) ([]Contact, error) {
	var contactList []Contact
	if err := utils.DB.Where("target_id = ? and type = 2", groupId).Find(&contactList).Error; err != nil {
		return nil, errors.New("查询出错")
	}
	return contactList, nil
}
