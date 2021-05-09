package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jamsa/gin-k8s/pkg/app"
	"github.com/jamsa/gin-k8s/pkg/e"
	"github.com/jamsa/gin-k8s/pkg/setting"
	"github.com/jamsa/gin-k8s/pkg/util"
	"github.com/jamsa/gin-k8s/service/user_service"
	"github.com/unknwon/com"
)

// @Summary Get a single user
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/users/{id} [get]
func GetUser(c *gin.Context) {
	appG := app.Gin{C: c}
	id := com.StrTo(c.Param("id")).MustInt()
	valid := validation.Validation{}
	valid.Min(id, 1, "id")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	userService := user_service.User{ID: id}
	exists, err := userService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	user, err := userService.Get()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, user)
}

// @ Param tag_id body int false "TagID"
// @ Param created_by body int false "CreatedBy"

// @Summary Get multiple users
// @Produce  json
// @Param state query int false "State"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/users [get]
func GetUsers(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	state := -1
	//if arg := c.PostForm("state"); arg != "" {
	if arg := c.Query("state"); arg != "" {
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

	userService := user_service.User{
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := userService.Count()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	users, err := userService.GetAll()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = users
	data["total"] = total

	appG.Response(http.StatusOK, e.SUCCESS, data)
}

type AddUserForm struct {
	Username    string `form:"name" valid:"Required;MaxSize(100)"`
	Fullname    string `form:"name" valid:"Required;MaxSize(100)"`
	Desc    string `form:"desc" valid:"Required;MaxSize(255)"`
	State   int    `form:"state" valid:"Range(0,1)"`
}

// @ Param tag_id body int true "TagID"
// @ Param title body string true "Title"
// @ Param desc body string true "Desc"
// @ Param content body string true "Content"
// @ Param created_by body string true "CreatedBy"
// @ Param state body int true "State"

// @Summary Add user
// @Produce  json
// @Param user body v1.AddUserForm true "user集群"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/users [post]
func AddUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form AddUserForm
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

	userService := user_service.User{
		ID:       0,
		Username: form.Username,
		Fullname: form.Fullname,
		Desc: form.Desc,
		State:    form.State,
		PageNum:  0,
		PageSize: 0,
	}
	if err := userService.Add(); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditUserForm struct {
	ID      int    `form:"id" valid:"Required;Min(1)"`
	Fullname    string `form:"fullname" valid:"Required;MaxSize(100)"`
	Username string `form:"username" valid:"Required;MaxSize(100)"`
	Desc    string `form:"desc" valid:"Required;MaxSize(255)"`
	//ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State int `form:"state" valid:"Range(0,1)"`
}

// @ Param id path int true "ID"
// @ Param tag_id body string false "TagID"
// @ Param title body string false "Title"
// @ Param desc body string false "Desc"
// @ Param content body string false "Content"
// @ Param modified_by body string true "ModifiedBy"
// @ Param state body int false "State"


// @Summary Update user
// @Produce  json
// @Param user body v1.EditUserForm true "user集群"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/users/{id} [put]
func EditUser(c *gin.Context) {
	var (
		appG = app.Gin{C: c}
		form = EditUserForm{ID: com.StrTo(c.Param("id")).MustInt()}
	)

	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	userService := user_service.User{
		ID:      form.ID,
		Username: form.Username,
		Fullname:    form.Fullname,
		Desc:    form.Desc,
		//ModifiedBy:    form.ModifiedBy,
		State: form.State,
	}
	exists, err := userService.ExistByID()
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

	err = userService.Edit()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

// @Summary Delete user
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	id := com.StrTo(c.Param("id")).MustInt()
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	userService := user_service.User{ID: id}
	exists, err := userService.ExistByID()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		appG.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = userService.Delete()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, nil)
}

