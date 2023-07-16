package models

import (
	"fmt"
	"ginchat/utils"
	"gorm.io/gorm"
)

// Contact 人员关系
type Contact struct {
	gorm.Model
	OwnerId  uint // 谁的关系
	TargetId uint // 对应的谁
	Type     int  // 关系类型 1 好友，2 群
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId string) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("owner_id = ? and type = 1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, v.TargetId)
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in (?)", objIds).Find(&users)
	return users
}

func AddFriend(userId uint, targetId uint) bool {
	user := UserBasic{}
	if targetId != 0 {
		utils.DB.Where("id = ?", targetId).First(&user)
		if user.ID == 0 {
			contract := Contact{}
			contract.OwnerId = userId
			contract.TargetId = targetId
			contract.Type = 1
			utils.DB.Create(&contract)
			return true
		}
	}
	return false
}
