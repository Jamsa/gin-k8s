package v1

import (
	"context"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/jamsa/gin-k8s/pkg/app"
	"github.com/jamsa/gin-k8s/pkg/e"
	"github.com/jamsa/gin-k8s/pkg/logging"
	"github.com/jamsa/gin-k8s/service/k8s"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// @ Param page body int false "Page"
// @Summary Get multiple deployments
// @Produce  json
// @Param cluster path string true "ClusterID"
// @Param namespace query string true "Namespace"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/k8s/{cluster}/deployments [get]
func GetDeployments(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	cluster := c.Param("cluster")
	namespace := c.Query("namespace")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	clientset, err := k8s.GetClient(cluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NO_K8S_CLIENT, nil)
		return
	}

	deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NO_K8S_RESOURCE, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, deployments)
}

// @Summary Get single deployment
// @Produce json
// @Param cluster path string true "ClusterID"
// @Param namespace path string true "Namespace"
// @Param deploymentName path string true "Deployment Name"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/k8s/{cluster}/deployments/{namespace}/{deploymentName} [get]
func GetDeployment(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}
	cluster := c.Param("cluster")
	namespace := c.Param("namespace")
	deploymentName := c.Param("deploymentName")

	if valid.HasErrors() {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	clientset, err := k8s.GetClient(cluster)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_NO_K8S_CLIENT, nil)
		return
	}

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logging.Error(err)
		appG.Response(http.StatusInternalServerError, e.ERROR_NO_K8S_RESOURCE, nil)
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, deployment)
}
