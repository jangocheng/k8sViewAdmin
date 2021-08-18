package main

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8sapi/pkg/deployoments"
	"k8sapi/pkg/rbac"
	"k8sapi/pkg/resources"
	"k8sapi/src/configs"
	"k8sapi/src/controllers"
	"net/http"
)
func cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

	}
}
func main() {
	 server:=goft.Ignite(cross()).Config(
		configs.NewK8sHandler(),  //1
		configs.NewK8sConfig(), //2
		configs.NewK8sMaps(), //3
		configs.NewServiceConfig(), //4
		).
		Mount("",
			deployoments.NewDeploymentCtlV2(),
			controllers.NewDeploymentCtl(),
			controllers.NewPodCtl(),
			controllers.NewUserCtl(),
			controllers.NewWsCtl(),
			controllers.NewNsCtl(),
			controllers.NewIngressCtl(),
			controllers.NewSvcCtl(),
			controllers.NewSecretCtl(),
			controllers.NewConfigMapCtl(),
			controllers.NewPodLogsCtl(),
			controllers.NewNodeCtl(),
			rbac.NewRBACCtl(),
			resources.NewResourcesCtl(),
		).
		Attach(
			//middlewares.NewCrosMiddleware(),//跨域中间件
			)
   //server.GET("/admin/*filepath", func(c *gin.Context) {
	//  http.FileServer(FS(false)).ServeHTTP(c.Writer,c.Request)
   //})
   server.Launch()






}
