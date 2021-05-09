package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	Model
	//Auth

	Username string `json:"username"`
	Password string `json:"password"`
	Fullname    string `json:"fullname"`
	Desc    string `json:"desc"`
	State   int    `json:"state"`
}

// ExistUserByID checks if an User exists based on ID
func ExistUserByID(id int) (bool, error) {
	var user User
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if user.ID > 0 {
		return true, nil
	}

	return false, nil
}

// GetUserTotal gets the total number of users based on the constraints
func GetUserTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&User{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetUsers gets a list of users based on paging constraints
func GetUsers(pageNum int, pageSize int, maps interface{}) ([]*User, error) {
	var users []*User
	//err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&users).Error
	err := db.Model(&User{}).Where(maps).Offset(pageNum).Limit(pageSize).Find(&users).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return users, nil
}

// GetUser Get a single user based on ID
func GetUser(id int) (*User, error) {
	var user User
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// err = db.Model(&user).Related(&user.Tag).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return nil, err
	// }

	return &user, nil
}

// EditUser modify a single user
func EditUser(id int, data interface{}) error {
	if err := db.Model(&User{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// AddUser add a single user
func AddUser(data map[string]interface{}) error {
	user := User{
		Username:    data["username"].(string),
		Desc:    data["desc"].(string),
		Password: data["password"].(string),
		Fullname: data["fullname"].(string),
		State:   data["state"].(int),
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// DeleteUser delete a single user
func DeleteUser(id int) error {
	if err := db.Where("id = ?", id).Delete(User{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllUser clear all user
func CleanAllUser() error {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}
