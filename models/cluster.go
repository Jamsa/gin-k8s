package models

import (
	"github.com/jinzhu/gorm"
)

type Cluster struct {
	Model

	Name    string `json:"name"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
	State   int    `json:"state"`
}

// ExistClusterByID checks if an Cluster exists based on ID
func ExistClusterByID(id int) (bool, error) {
	var cluster Cluster
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&cluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if cluster.ID > 0 {
		return true, nil
	}

	return false, nil
}

// GetClusterTotal gets the total number of clusters based on the constraints
func GetClusterTotal(maps interface{}) (int, error) {
	var count int
	if err := db.Model(&Cluster{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

// GetClusters gets a list of clusters based on paging constraints
func GetClusters(pageNum int, pageSize int, maps interface{}) ([]*Cluster, error) {
	var clusters []*Cluster
	//err := db.Preload("Tag").Where(maps).Offset(pageNum).Limit(pageSize).Find(&clusters).Error
	err := db.Model(&Cluster{}).Where(maps).Offset(pageNum).Limit(pageSize).Find(&clusters).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return clusters, nil
}

// GetCluster Get a single cluster based on ID
func GetCluster(id int) (*Cluster, error) {
	var cluster Cluster
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&cluster).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// err = db.Model(&cluster).Related(&cluster.Tag).Error
	// if err != nil && err != gorm.ErrRecordNotFound {
	// 	return nil, err
	// }

	return &cluster, nil
}

// EditCluster modify a single cluster
func EditCluster(id int, data interface{}) error {
	if err := db.Model(&Cluster{}).Where("id = ? AND deleted_on = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

// AddCluster add a single cluster
func AddCluster(data map[string]interface{}) error {
	cluster := Cluster{
		Name:    data["name"].(string),
		Desc:    data["desc"].(string),
		Content: data["content"].(string),
		State:   data["state"].(int),
	}
	if err := db.Create(&cluster).Error; err != nil {
		return err
	}

	return nil
}

// DeleteCluster delete a single cluster
func DeleteCluster(id int) error {
	if err := db.Where("id = ?", id).Delete(Cluster{}).Error; err != nil {
		return err
	}

	return nil
}

// CleanAllCluster clear all cluster
func CleanAllCluster() error {
	if err := db.Unscoped().Where("deleted_on != ? ", 0).Delete(&Cluster{}).Error; err != nil {
		return err
	}

	return nil
}
