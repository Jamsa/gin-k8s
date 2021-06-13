package user_service

import (
	//"encoding/json"

	"github.com/jamsa/gin-k8s/models"
	//"github.com/jamsa/gin-k8s/pkg/gredis"
	//"github.com/jamsa/gin-k8s/pkg/logging"
	//"github.com/jamsa/gin-k8s/service/cache_service"
)

type User struct {
	ID      int
	Username string
	Fullname    string
	Desc string
	State   int

	PageNum  int
	PageSize int
}

func (a *User) Add() error {
	user := map[string]interface{}{

		"fullname":    a.Fullname,
		"state": a.State,
	}

	if err := models.AddUser(user); err != nil {
		return err
	}

	return nil
}

func (a *User) Edit() error {
	return models.EditUser(a.ID, map[string]interface{}{
		"fullname":    a.Fullname,
		"username":    a.Username,
		"state":   a.State,
		//"modified_by":     a.ModifiedBy,
	})
}

func (a *User) Get() (*models.User, error) {
	// var cacheUser *models.User

	// cache := cache_service.User{ID: a.ID}
	// key := cache.GetUserKey()
	// if gredis.Exists(key) {
	// 	data, err := gredis.Get(key)
	// 	if err != nil {
	// 		logging.Info(err)
	// 	} else {
	// 		json.Unmarshal(data, &cacheUser)
	// 		return cacheUser, nil
	// 	}
	// }

	user, err := models.GetUser(a.ID)
	if err != nil {
		return nil, err
	}

	// gredis.Set(key, user, 3600) 
	return user, nil
}

func (a *User) GetAll() ([]*models.User, error) {
	var (
		users /*, cacheUsers*/ []*models.User
	)

	// cache := cache_service.User{
	// 	State: a.State,

	// 	PageNum:  a.PageNum,
	// 	PageSize: a.PageSize,
	// }
	// key := cache.GetUsersKey()
	// if gredis.Exists(key) {
	// 	data, err := gredis.Get(key)
	// 	if err != nil {
	// 		logging.Info(err)
	// 	} else {
	// 		json.Unmarshal(data, &cacheUsers)
	// 		return cacheUsers, nil
	// 	}
	// }

	users, err := models.GetUsers(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	// gredis.Set(key, users, 3600)
	return users, nil
}

func (a *User) Delete() error {
	return models.DeleteUser(a.ID)
}

func (a *User) ExistByID() (bool, error) {
	return models.ExistUserByID(a.ID)
}

func (a *User) Count() (int, error) {
	return models.GetUserTotal(a.getMaps())
}

func (a *User) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	if a.State != -1 {
		maps["state"] = a.State
	}
	if len(a.Fullname)>0 {
		maps["fullname"] = a.Fullname
	}
	if len(a.Username)>0 {
		maps["Username"] = a.Username
	}
	/*if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}*/

	return maps
}
