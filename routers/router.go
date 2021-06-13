package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jamsa/gin-k8s/api"
	admv1 "github.com/jamsa/gin-k8s/api/v1/admin"
	k8sv1 "github.com/jamsa/gin-k8s/api/v1/k8s"
	_ "github.com/jamsa/gin-k8s/docs"
	"github.com/jamsa/gin-k8s/pkg/logging"
	"github.com/jamsa/gin-k8s/pkg/setting"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	//r.Use(gin.Logger())
	if setting.AppSetting.LogGin {
		r.Use(logging.LogToLogrus())
	}

	r.Use(gin.Recovery())

	/*
	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))
	r.POST("/upload", api.UploadImage)
	*/
	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/k8s/:cluster/pods", k8sv1.GetPods)
		apiv1.GET("/k8s/:cluster/pods/:namespace/:podname", k8sv1.GetPod)

		apiv1.GET("/k8s/:cluster/deployments", k8sv1.GetDeployments)
		apiv1.GET("/k8s/:cluster/deployments/:namespace/:deploymentName", k8sv1.GetDeployment)

		apiv1.GET("/k8s/:cluster/services", k8sv1.GetServices)
		apiv1.GET("/k8s/:cluster/services/:namespace/:serviceName", k8sv1.GetService)

		apiv1.GET("/k8s/:cluster/statefulsets", k8sv1.GetStatefulSets)
		apiv1.GET("/k8s/:cluster/statefulsets/:namespace/:statefulsetName", k8sv1.GetStatefulSet)

		apiv1.GET("/k8s/:cluster/ingresses", k8sv1.GetIngresses)
		apiv1.GET("/k8s/:cluster/ingresses/:namespace/:ingressName", k8sv1.GetIngress)

		apiv1.GET("/k8s/:cluster/configmaps", k8sv1.GetConfigMaps)
		apiv1.GET("/k8s/:cluster/configmaps/:namespace/:configmapName", k8sv1.GetConfigMap)

		apiv1.GET("/k8s/:cluster/persistentvolumeclaims", k8sv1.GetPersistentVolumeClaims)
		apiv1.GET("/k8s/:cluster/persistentvolumeclaims/:namespace/:persistentvolumeclaimName", k8sv1.GetPersistentVolumeClaim)

		//不区分namespace的
		apiv1.GET("/k8s/:cluster/persistentvolumes", k8sv1.GetPersistentVolumes)
		apiv1.GET("/k8s/:cluster/persistentvolumes/:persistentvolumeName", k8sv1.GetPersistentVolume)

		apiv1.GET("/k8s/:cluster/nodes", k8sv1.GetNodes)
		apiv1.GET("/k8s/:cluster/nodes/:nodeName", k8sv1.GetNode)

		apiv1.GET("/k8s/:cluster/namespaces", k8sv1.GetNamespaces)
		apiv1.GET("/k8s/:cluster/namespaces/:namespaceName", k8sv1.GetNamespace)


		//集群
		apiv1.GET("/admin/clusters", admv1.GetClusters)
		apiv1.GET("/admin/clusters/:id", admv1.GetCluster)
		apiv1.POST("/admin/clusters", admv1.AddCluster)
		apiv1.PUT("/admin/clusters/:id", admv1.EditCluster)
		apiv1.DELETE("/admin/clusters/:id", admv1.DeleteCluster)

		//用户
		apiv1.GET("/admin/users", admv1.GetUsers)
		apiv1.GET("/admin/users/:id", admv1.GetUser)
		apiv1.POST("/admin/users", admv1.AddUser)
		apiv1.PUT("/admin/users/:id", admv1.EditUser)
		apiv1.DELETE("/admin/users/:id", admv1.DeleteUser)
	}
	/*{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
		//生成文章海报
		apiv1.POST("/articles/poster/generate", v1.GenerateArticlePoster)
	}*/

	return r
}
