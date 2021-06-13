package cluster_service

import (
	"encoding/json"

	"github.com/jamsa/gin-k8s/models"
	"github.com/jamsa/gin-k8s/pkg/gredis"
	"github.com/jamsa/gin-k8s/pkg/logging"
	"github.com/jamsa/gin-k8s/service/cache_service"
)

type Cluster struct {
	ID      int
	Name    string
	Desc    string
	Content string
	State   int

	PageNum  int
	PageSize int
}

func (a *Cluster) Add() error {
	cluster := map[string]interface{}{

		"name":    a.Name,
		"desc":    a.Desc,
		"content": a.Content,
		//"created_by":      a.CreatedBy,
		"state": a.State,
	}

	if err := models.AddCluster(cluster); err != nil {
		return err
	}

	return nil
}

func (a *Cluster) Edit() error {
	return models.EditCluster(a.ID, map[string]interface{}{
		"name":    a.Name,
		"desc":    a.Desc,
		"content": a.Content,
		"state":   a.State,
		//"modified_by":     a.ModifiedBy,
	})
}

func (a *Cluster) Get() (*models.Cluster, error) {
	var cacheCluster *models.Cluster

	cache := cache_service.Cluster{ID: a.ID}
	key := cache.GetClusterKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheCluster)
			return cacheCluster, nil
		}
	}

	cluster, err := models.GetCluster(a.ID)
	if err != nil {
		return nil, err
	}

	gredis.Set(key, cluster, 3600)
	return cluster, nil
}

func (a *Cluster) GetAll() ([]*models.Cluster, error) {
	var (
		clusters/*, cacheClusters*/ []*models.Cluster
	)

	/*
	cache := cache_service.Cluster{
		State: a.State,

		PageNum:  a.PageNum,
		PageSize: a.PageSize,
	}
	key := cache.GetClustersKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Info(err)
		} else {
			json.Unmarshal(data, &cacheClusters)
			return cacheClusters, nil
		}
	}
	*/

	clusters, err := models.GetClusters(a.PageNum, a.PageSize, a.getMaps())
	if err != nil {
		return nil, err
	}

	//gredis.Set(key, clusters, 3600)
	return clusters, nil
}

func (a *Cluster) Delete() error {
	return models.DeleteCluster(a.ID)
}

func (a *Cluster) ExistByID() (bool, error) {
	return models.ExistClusterByID(a.ID)
}

func (a *Cluster) Count() (int, error) {
	return models.GetClusterTotal(a.getMaps())
}

func (a *Cluster) getMaps() map[string]interface{} {
	maps := make(map[string]interface{})
	maps["deleted_on"] = 0
	/*if a.State != -1 {
		maps["state"] = a.State
	}*/
	if len(a.Name)>0 {
		maps["name"] = a.Name
	}
	if len(a.Desc)>0 {
		maps["desc"] = a.Desc
	}
	/*if a.TagID != -1 {
		maps["tag_id"] = a.TagID
	}*/

	return maps
}
