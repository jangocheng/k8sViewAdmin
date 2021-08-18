package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8sapi/src/helpers"
	"k8sapi/src/models"
	"k8sapi/src/wscore"
	"log"
)

//@Controller
type WsCtl struct {
	Client *kubernetes.Clientset  `inject:"-"`
	Config *rest.Config  `inject:"-"`
	SysConfig *models.SysConfig   `inject:"-"`
}

func NewWsCtl() *WsCtl {
	return &WsCtl{}
}

func(this *WsCtl) Connect(c *gin.Context) (v goft.Void)    {
	client,err:=wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)  //升级
	if err!=nil{
		log.Println(err)
		return
	}else{
		wscore.ClientMap.Store(client)

		return
	}
}
func(this *WsCtl) PodConnect(c *gin.Context) (v goft.Void){
	ns:=c.Query("ns")
	pod:=c.Query("pod")
	container:=c.Query("c")
	wsClient,err:=wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)
	goft.Error(err)
	shellClient:=wscore.NewWsShellClient(wsClient)
	err=helpers.HandleCommand(ns,pod,container,this.Client,this.Config,[]string{"sh"}).
		Stream(remotecommand.StreamOptions{
			Stdin:shellClient,
			Stdout:shellClient,
			Stderr:shellClient,
			Tty:true,
		})
	return
}

func(this *WsCtl) NodeConnect(c *gin.Context) (v goft.Void){
	nodeName:=c.Query("node")
	nodeConfig:=helpers.GetNodeConfig(this.SysConfig,nodeName) //读取配置文件
	wsClient,err:=wscore.Upgrader.Upgrade(c.Writer,c.Request,nil)
	goft.Error(err)
	shellClient:=wscore.NewWsShellClient(wsClient)
	//session, err := helpers.SSHConnect(helpers.TempSSHUser,  helpers.TempSSHPWD, helpers.TempSSHIP ,22 )
	session, err := helpers.SSHConnect(nodeConfig.User,  nodeConfig.Pass, nodeConfig.Ip ,22 )
	fmt.Println(err)
	goft.Error(err)
	defer session.Close()
	session.Stdout = shellClient
	session.Stderr =shellClient
	session.Stdin=shellClient
    err= session.RequestPty("xterm-256color", 300, 500, helpers.NodeShellModes)
	goft.Error(err)

	err=session.Run("sh")
	goft.Error(err)
	return
}
func(this *WsCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/ws",this.Connect)
	goft.Handle("GET","/podws",this.PodConnect)
	goft.Handle("GET","/nodews",this.NodeConnect)
}
func(this *WsCtl) Name() string{
	return "WsCtl"
}
