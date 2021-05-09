package api

import (
	//"fmt"
	"net/http"

	//"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"github.com/jamsa/gin-k8s/pkg/app"
	"github.com/jamsa/gin-k8s/pkg/e"
	"github.com/jamsa/gin-k8s/pkg/util"
	"github.com/jamsa/gin-k8s/service/auth_service"
)

type Auth struct {
	Username string `valid:"Required; MaxSize(50)" json:"username"`
	Password string `valid:"Required; MaxSize(50)" json:"password"`
}

// @Summary Get Auth
// @Produce  json
// @Param auth body api.Auth true "认证信息"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [post]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	auth := Auth{};
	// c.BindJSON(auth);
	// fmt.Println(auth);
	// a := auth;
	// valid := validation.Validation{}
	// ok, _ := valid.Valid(&a)

	// if !ok {
	// 	app.MarkErrors(valid.Errors)
	// 	appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
	// 	return
	// }

	httpCode, errCode := app.BindAndValid(c, &auth)
	if errCode != e.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}

	authService := auth_service.Auth{Username: auth.Username, Password: auth.Password}
	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(auth.Username, auth.Password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
