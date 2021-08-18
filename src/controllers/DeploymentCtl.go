package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/client-go/kubernetes"
	"k8sapi/src/services"
)

type DeploymentCtl struct {
	K8sClient *kubernetes.Clientset `inject:"-"`
	DepService *services.DeploymentService  `inject:"-"`  //首字母一定要大写
}

func NewDeploymentCtl() *DeploymentCtl {
	return &DeploymentCtl{}
}
func(this *DeploymentCtl) GetList(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default") // GET /deployments?ns=xxx
	return gin.H{
		"code":20000,
		"data":this.DepService.ListAll(ns),
	}

}
func(this *DeploymentCtl)  Build(goft *goft.Goft){
	//路由
	goft.Handle("GET","/deployments",this.GetList)
}
func(*DeploymentCtl) Name() string{
	return "DeploymentCtl"
}