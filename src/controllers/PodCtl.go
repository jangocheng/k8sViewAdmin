package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8sapi/src/services"
)

type PodCtl struct {
	PodService *services.PodService `inject:"-"`
	Helper *services.Helper `inject:"-"`

}

func NewPodCtl() *PodCtl {
	return &PodCtl{}
}
//获取 容器
func(this *PodCtl) Containers(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	podname:=c.DefaultQuery("name","")
	return gin.H{
		"code":20000,
		"data":this.PodService.GetPodContainer(ns,podname),
	}

}
func(this *PodCtl) GetAll(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	page:=c.DefaultQuery("current","1") //当前页
	size:=c.DefaultQuery("size","5")
	return gin.H{
		"code":20000,
		"data":this.PodService.PagePods(ns,this.Helper.StrToInt(page,1),
			this.Helper.StrToInt(size,5)),
	}

}
func(this *PodCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/pods",this.GetAll)
	goft.Handle("GET","/pods/containers",this.Containers)
}
func(*PodCtl) Name() string{
	return "PodCtl"
}