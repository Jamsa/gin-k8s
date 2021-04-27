package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jamsa/gin-k8s/pkg/app"
	"github.com/jamsa/gin-k8s/pkg/e"
	"github.com/jamsa/gin-k8s/pkg/setting"
	"github.com/jamsa/gin-k8s/pkg/util"
	"github.com/jamsa/gin-k8s/service/cluster_service"
	"github.com/unknwon/com"
)

// @Summary Get a single cluster
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/clusters/{id} [get]
func GetCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	clusterService := cluster_service.Cluster{ID: id}
	exists, err := clusterService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	cluster, err := clusterService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, cluster)
}

// @Summary Get multiple clusters
// @Produce  json
// 不要的Param tag_id body int false "TagID"
// @Param state body int false "State"
// @Param created_by body int false "CreatedBy"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/clusters [get]
func GetClusters(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state")
	}

	// tagId := -1
	// if arg := c.PostForm("tag_id"); arg != "" {
	// 	tagId = com.StrTo(arg).MustInt()
	// 	valid.Min(tagId, 1, "tag_id")
	// }

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	clusterService := cluster_service.Cluster{
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := clusterService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	clusters, err := clusterService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = clusters
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddClusterForm struct {
	Name    string `form:"name" valid:"Required;MaxSize(100)"`
	Desc    string `form:"desc" valid:"Required;MaxSize(255)"`
	Content string `form:"content" valid:"Required;MaxSize(65535)"`
	State   int    `form:"state" valid:"Range(0,1)"`
}

// @Summary Add cluster
// @Produce  json
// Param tag_id body int true "TagID"
// @Param title body string true "Title"
// @Param desc body string true "Desc"
// @Param content body string true "Content"
// @Param created_by body string true "CreatedBy"
// @Param state body int true "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/clusters [post]
func AddCluster(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddClusterForm
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	/*tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}*/

	clusterService := cluster_service.Cluster{
		Name:    form.Name,
		Desc:    form.Desc,
		Content: form.Content,
		State:   form.State,
	}
	if err := clusterService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditClusterForm struct {
	ID      int    `form:"id" valid:"Required;Min(1)"`
	Name    string `form:"name" valid:"Required;MaxSize(100)"`
	Desc    string `form:"desc" valid:"Required;MaxSize(255)"`
	Content string `form:"content" valid:"Required;MaxSize(65535)"`
	//ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State int `form:"state" valid:"Range(0,1)"`
}

// @Summary Update cluster
// @Produce  json
// @Param id path int true "ID"
// @Param tag_id body string false "TagID"
// @Param title body string false "Title"
// @Param desc body string false "Desc"
// @Param content body string false "Content"
// @Param modified_by body string true "ModifiedBy"
// @Param state body int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/clusters/{id} [put]
func EditCluster(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditClusterForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	clusterService := cluster_service.Cluster{
		ID:      form.ID,
		Name:    form.Name,
		Desc:    form.Desc,
		Content: form.Content,
		//ModifiedBy:    form.ModifiedBy,
		State: form.State,
	}
	exists, err := clusterService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	// tagService := tag_service.Tag{ID: form.TagID}
	// exists, err = tagService.ExistByID()
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
	// 	return
	// }

	// if !exists {
	// 	appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
	// 	return
	// }

	err = clusterService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete cluster
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/clusters/{id} [delete]
func DeleteCluster(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	clusterService := cluster_service.Cluster{ID: id}
	exists, err := clusterService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = clusterService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

/*
const (
	QRCODE_URL = "https://github.com/EDDYCJY/blog#gin%E7%B3%BB%E5%88%97%E7%9B%AE%E5%BD%95"
)

func GenerateClusterPoster(c *gin.Context) {
	appG := app.Gin{C: c}
	cluster := &cluster_service.Cluster{}
	qr := qrcode.NewQrCode(QRCODE_URL, 300, 300, qr.M, qr.Auto)
	posterName := cluster_service.GetPosterFlag() + "-" + qrcode.GetQrCodeFileName(qr.URL) + qr.GetQrCodeExt()
	clusterPoster := cluster_service.NewClusterPoster(posterName, cluster, qr)
	clusterPosterBgService := cluster_service.NewClusterPosterBg(
		"bg.jpg",
		clusterPoster,
		&cluster_service.Rect{
			X0: 0,
			Y0: 0,
			X1: 550,
			Y1: 700,
		},
		&cluster_service.Pt{
			X: 125,
			Y: 298,
		},
	)

	_, filePath, err := clusterPosterBgService.Generate()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GEN_ARTICLE_POSTER_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"poster_url":      qrcode.GetQrCodeFullUrl(posterName),
		"poster_save_url": filePath + posterName,
	})
}
*/
