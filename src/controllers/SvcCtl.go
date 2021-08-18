package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8sapi/src/services"
)

type SvcCtl struct{
	SvcMap *services.ServiceMapStruct `inject:"-"`
}
func NewSvcCtl() *SvcCtl{
  return &SvcCtl{}
}
func(this *SvcCtl) ListAll(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code":20000,
		"data":this.SvcMap.ListAll(ns),
	}
}
func(this *SvcCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/svc",this.ListAll)
}
func(*SvcCtl)  Name() string{
	 return "SvcCtl"
}
