package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/yaml"

	"k8sapi/src/models"
	"k8sapi/src/services"

)

type IngressCtl struct{
	IngressMap *services.IngressMapStruct `inject:"-"`
	IngressService *services.IngressService `inject:"-"`
	Client *kubernetes.Clientset  `inject:"-"`
}
func NewIngressCtl() *IngressCtl{
  return &IngressCtl{}
}
func(*IngressCtl)  Name() string{
	 return "IngressCtl"
}

func(this *IngressCtl) PostIngress(c *gin.Context) goft.Json{
	postModel:=&models.IngressPost{}
	goft.Error(c.BindJSON(postModel))
	goft.Error(this.IngressService.PostIngress(postModel))
	return gin.H{
		"code":20000,
		"data":postModel,
	}
}
//DELETE /ingress?ns=xx&name=xx
func(this *IngressCtl) RmIngress(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	name:=c.DefaultQuery("name","")
	goft.Error(this.Client.NetworkingV1beta1().Ingresses(ns).
		Delete(c,name,v1.DeleteOptions{}))
	return gin.H{
		"code":20000,
		"data":"OK",
	}
}
func(this *IngressCtl) ListAll(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code":20000,
		 "data":this.IngressService.ListIngress(ns), //暂时 不分页

	}
}

func(this *IngressCtl) GetIngressForYaml(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	name:=c.DefaultQuery("name","")
	ingress,err:=this.Client.NetworkingV1beta1().Ingresses(ns).Get(c,name,v1.GetOptions{})
	goft.Error(err)
	b,err:=yaml.Marshal(ingress)
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":string(b), //暂时 不分页
	}
}
func(this *IngressCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/ingress",this.ListAll)
	goft.Handle("DELETE","/ingress",this.RmIngress)
	goft.Handle("POST","/ingress",this.PostIngress)
	goft.Handle("GET","/ingressyaml",this.GetIngressForYaml)
}
