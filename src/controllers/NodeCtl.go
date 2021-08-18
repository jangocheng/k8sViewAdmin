package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8sapi/src/models"
	"k8sapi/src/services"
)

//@controller
type NodeCtl struct {
	 NodeService *services.NodeService `inject:"-"`
	 Client *kubernetes.Clientset  `inject:"-"`
}

func NewNodeCtl() *NodeCtl {
	return &NodeCtl{}
}
func(this *NodeCtl) ListAll(c *gin.Context) goft.Json{
	return gin.H{
		"code":20000,
		"data":this.NodeService.ListAllNodes(),
	}
}
//保存node
func(this *NodeCtl) SaveNode(c *gin.Context) goft.Json{
	nodeModel:=&models.PostNodeModel{}
	_=c.ShouldBindJSON(nodeModel)
	node:=this.NodeService.LoadOrginNode(nodeModel.Name) //取出原始node 信息
	if node==nil{
		panic("no such node")
	}
	node.Labels=nodeModel.OrginLabels  //覆盖标签
	node.Spec.Taints=nodeModel.OrginTaints //覆盖 污点
    _,err:=this.Client.CoreV1().Nodes().Update(c,node,v1.UpdateOptions{})
    goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *NodeCtl) LoadDetail(c *gin.Context) goft.Json{
	nodeName:=c.Param("node")
	return gin.H{
		"code":20000,
		"data":this.NodeService.LoadNode(nodeName),
	}

}
func(this *NodeCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/nodes",this.ListAll)
	goft.Handle("GET","/nodes/:node",this.LoadDetail) //加载详细
	goft.Handle("POST","/nodes",this.SaveNode)  //保存
}
func(*NodeCtl) Name() string{
	return "NodeCtl"
}