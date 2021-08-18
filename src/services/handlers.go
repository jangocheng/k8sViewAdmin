package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	"k8sapi/src/wscore"
	"log"
)

//处理deployment 回调的handler
type DepHandler struct {
	DepMap *DeploymentMap `inject:"-"`
	DepService *DeploymentService `inject:"-"`
}
func(this *DepHandler) OnAdd(obj interface{}){
	this.DepMap.Add(obj.(*v1.Deployment))
	ns:=obj.(*v1.Deployment).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"deployments",
			"result":gin.H{"ns":ns,"data":this.DepService.ListAll(ns)},
		},
	)
  }
func(this *DepHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.DepMap.Update(newObj.(*v1.Deployment))
	if err!=nil{
		log.Println(err)
	}else{
		ns:=newObj.(*v1.Deployment).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"deployments",
				"result":gin.H{"ns":ns,"data":this.DepService.ListAll(ns)},
			},
		)
 	}
}
func(this *DepHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*v1.Deployment);ok{
		this.DepMap.Delete(d)
		ns:=obj.(*v1.Deployment).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"deployments",
				"result":gin.H{"ns":ns,"data":this.DepService.ListAll(ns)},
			},
		)
	}
}

type RsHandler struct {
	RsMap *RsMapStruct `inject:"-"`
}
func(this *RsHandler) OnAdd(obj interface{}){
	this.RsMap.Add(obj.(*v1.ReplicaSet))
}
func(this *RsHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.RsMap.Update(newObj.(*v1.ReplicaSet))
	if err!=nil{
		log.Println(err)
	}
}
func(this *RsHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*v1.ReplicaSet);ok{
		this.RsMap.Delete(d)
	}
}

// pod相关的回调handler
type PodHandler struct {
	PodMap *PodMapStruct `inject:"-"`
	PodService *PodService `inject:"-"`
}
func(this *PodHandler) OnAdd(obj interface{}){
	this.PodMap.Add(obj.(*corev1.Pod))
    ns:=obj.(*corev1.Pod).Namespace
    wscore.ClientMap.SendAll(
		gin.H{
			"type":"pods",
			"result":gin.H{"ns":ns,
				"data":this.PodService.PagePods(ns,1,5)},
		},
	)

}
func(this *PodHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.PodMap.Update(newObj.(*corev1.Pod))
	if err!=nil{
		log.Println(err)
	}else{
		ns:=newObj.(*corev1.Pod).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"pods",
				"result":gin.H{"ns":ns,
					"data":this.PodService.PagePods(ns,1,5)},
			},
		)
	}
}
func(this *PodHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*corev1.Pod);ok {
		this.PodMap.Delete(d)
		ns:=obj.(*corev1.Pod).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"pods",
				"result":gin.H{"ns":ns,
					"data":this.PodService.PagePods(ns,1,5)},
			},
		)
	}
}


// namespace 相关的回调handler
type NsHandler struct {
	NsMap *NsMapStruct `inject:"-"`
}
func(this *NsHandler) OnAdd(obj interface{}){
	this.NsMap.Add(obj.(*corev1.Namespace))
}
func(this *NsHandler) OnUpdate(oldObj, newObj interface{}){
	 this.NsMap.Update(newObj.(*corev1.Namespace))

}
func(this *NsHandler)	OnDelete(obj interface{}){
	if d,ok:=obj.(*corev1.Namespace);ok{
		this.NsMap.Delete(d)
	}
}

// event 事件相关的handler
type EventHandler struct {
	EventMap *EventMapStruct  `inject:"-"`
}
func(this *EventHandler) storeData(obj interface{},isdelete bool){
	if event,ok:=obj.(*corev1.Event);ok{
		key:=fmt.Sprintf("%s_%s_%s",event.Namespace,event.InvolvedObject.Kind,event.InvolvedObject.Name)
		if !isdelete{
			this.EventMap.data.Store(key,event)
		}else{
			this.EventMap.data.Delete(key)
		}
	}
}
func(this *EventHandler) OnAdd(obj interface{}){
	this.storeData(obj,false)
}
func(this *EventHandler) OnUpdate(oldObj, newObj interface{}){
	this.storeData(newObj,false)
}
func(this *EventHandler) OnDelete(obj interface{}){
	this.storeData(obj,true)
}


// ingress相关handler
type IngressHandler struct {
	IngressMap *IngressMapStruct  `inject:"-"`
	IngressService *IngressService `inject:"-"`
}
func(this *IngressHandler) OnAdd(obj interface{}){
	this.IngressMap.Add(obj.(*v1beta1.Ingress))
	ns:=obj.(*v1beta1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns":ns,
				"data":this.IngressService.ListIngress(ns)},
		},
	)
}
func(this *IngressHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.IngressMap.Update(newObj.(*v1beta1.Ingress))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*v1beta1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns":ns,
				"data":this.IngressService.ListIngress(ns)},
		},
	)

}
func(this *IngressHandler) OnDelete(obj interface{}){
	this.IngressMap.Delete(obj.(*v1beta1.Ingress))
	ns:=obj.(*v1beta1.Ingress).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"ingress",
			"result":gin.H{"ns":ns,
				"data":this.IngressService.ListIngress(ns)},
		},
	)
}


// Service 相关handler
type ServiceHandler struct {
	SvcMap *ServiceMapStruct  `inject:"-"`
}
func(this *ServiceHandler) OnAdd(obj interface{}){
	this.SvcMap.Add(obj.(*corev1.Service))
	ns:=obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"service",
			"result":gin.H{"ns":ns,
				"data":this.SvcMap.ListAll(ns)},
		},
	)
}
func(this *ServiceHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.SvcMap.Update(newObj.(*corev1.Service))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"service",
			"result":gin.H{"ns":ns,
				"data":this.SvcMap.ListAll(ns)},
		},
	)
}
func(this *ServiceHandler) OnDelete(obj interface{}){
	this.SvcMap.Delete(obj.(*corev1.Service))
	ns:=obj.(*corev1.Service).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"service",
			"result":gin.H{"ns":ns,
				"data":this.SvcMap.ListAll(ns)},
		},
	)
}


//Secret相关的handler
type SecretHandler struct {
	SecretMap *SecretMapStruct  `inject:"-"`
	SecretService *SecretService  `inject:"-"`
}
func(this *SecretHandler) OnAdd(obj interface{}){
	this.SecretMap.Add(obj.(*corev1.Secret))
	ns:=obj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns":ns,
				"data":this.SecretService.ListSecret(ns)},
		},
	)
}
func(this *SecretHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.SecretMap.Update(newObj.(*corev1.Secret))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns":ns,
				"data":this.SecretService.ListSecret(ns)},
		},
	)
}
func(this *SecretHandler) OnDelete(obj interface{}){
	this.SecretMap.Delete(obj.(*corev1.Secret))
	ns:=obj.(*corev1.Secret).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"secret",
			"result":gin.H{"ns":ns,
				"data":this.SecretService.ListSecret(ns)},
		},
	)
}


//CofigMap相关的handler
type ConfigMapHandler struct {
	ConfigMap *ConfigMapStruct  `inject:"-"`
	ConfigMapService *ConfigMapService  `inject:"-"`
}
func(this *ConfigMapHandler) OnAdd(obj interface{}){
	this.ConfigMap.Add(obj.(*corev1.ConfigMap))
	ns:=obj.(*corev1.ConfigMap).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cm",
			"result":gin.H{"ns":ns,
				"data":this.ConfigMapService.ListConfigMap(ns)},
		},
	)
}
func(this *ConfigMapHandler) OnUpdate(oldObj, newObj interface{}){
	//重点： 只要update返回true 才会发送 。否则不发送
	if this.ConfigMap.Update(newObj.(*corev1.ConfigMap)){
		ns:=newObj.(*corev1.ConfigMap).Namespace
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"cm",
				"result":gin.H{"ns":ns,
					"data":this.ConfigMapService.ListConfigMap(ns)},
			},
		)
	}
}
func(this *ConfigMapHandler) OnDelete(obj interface{}){
	this.ConfigMap.Delete(obj.(*corev1.ConfigMap))
	ns:=obj.(*corev1.ConfigMap).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"cm",
			"result":gin.H{"ns":ns,
				"data":this.ConfigMapService.ListConfigMap(ns)},
		},
	)
}

//Node相关的handler
type NodeMapHandler struct {
	NodeMap *NodeMapStruct  `inject:"-"`
	NodeService *NodeService `inject:"-"`
}
func(this *NodeMapHandler) OnAdd(obj interface{}){
	this.NodeMap.Add(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type":"node",
			"result":gin.H{"ns":"node",
				"data":this.NodeService.ListAllNodes()},
		},
	)
}
func(this *NodeMapHandler) OnUpdate(oldObj, newObj interface{}){
	//重点： 只要update返回true 才会发送 。否则不发送
	if this.NodeMap.Update(newObj.(*corev1.Node)){
		wscore.ClientMap.SendAll(
			gin.H{
				"type":"node",
				"result":gin.H{"ns":"node",
					"data":this.NodeService.ListAllNodes()},
			},
		)
	}
}
func(this *NodeMapHandler) OnDelete(obj interface{}){
	this.NodeMap.Delete(obj.(*corev1.Node))

	wscore.ClientMap.SendAll(
		gin.H{
			"type":"node",
			"result":gin.H{"ns":"node",
				"data":this.NodeService.ListAllNodes()},
		},
	)
}