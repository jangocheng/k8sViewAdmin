package rbac

import (
	corev1 "k8s.io/api/core/v1"
	"github.com/gin-gonic/gin"
	"k8s.io/api/rbac/v1"
	"k8sapi/src/wscore"
	"log"
)

type RoleHander struct {
	RoleMap *RoleMapStruct  `inject:"-"`
	RoleService *RoleService  `inject:"-"`
}
func(this *RoleHander) OnAdd(obj interface{}){
	this.RoleMap.Add(obj.(*v1.Role))
	ns:=obj.(*v1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"role",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoles(ns)},
		},
	)
}
func(this *RoleHander) OnUpdate(oldObj, newObj interface{}){
	err:=this.RoleMap.Update(newObj.(*v1.Role))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*v1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"role",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoles(ns)},
		},
	)
}
func(this *RoleHander) OnDelete(obj interface{}){
	this.RoleMap.Delete(obj.(*v1.Role))
	ns:=obj.(*v1.Role).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"role",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoles(ns)},
		},
	)
}


type RoleBindingHander struct {
	RoleBindingMap *RoleBindingMapStruct  `inject:"-"`
	RoleService *RoleService  `inject:"-"`
}
func(this *RoleBindingHander) OnAdd(obj interface{}){
	this.RoleBindingMap.Add(obj.(*v1.RoleBinding))
	ns:=obj.(*v1.RoleBinding).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"rolebinding",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoleBindings(ns)},
		},
	)
}
func(this *RoleBindingHander) OnUpdate(oldObj, newObj interface{}){
	err:=this.RoleBindingMap.Update(newObj.(*v1.RoleBinding))
	if err!=nil{
		log.Println(err)
		return
	}
	ns:=newObj.(*v1.RoleBinding).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"rolebinding",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoleBindings(ns)},
		},
	)
}
func(this *RoleBindingHander) OnDelete(obj interface{}){
	this.RoleBindingMap.Delete(obj.(*v1.RoleBinding))
	ns:=obj.(*v1.RoleBinding).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"rolebinding",
			"result":gin.H{"ns":ns,
				"data":this.RoleService.ListRoleBindings(ns)},
		},
	)
}

type ClusterRoleBindingHander struct {
	ClusterRoleBindingMap *ClusterRoleBindingMapStruct  `inject:"-"`
	RoleService *RoleService  `inject:"-"`
}
func(this *ClusterRoleBindingHander) OnAdd(obj interface{}){
	this.ClusterRoleBindingMap.Add(obj.(*v1.ClusterRoleBinding))
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrolebinding",
			"result":gin.H{"ns":"clusterrolebinding",
				"data":this.RoleService.ListClusterRoleBindings()},
		},
	)
}
func(this *ClusterRoleBindingHander) OnUpdate(oldObj, newObj interface{}){
	err:=this.ClusterRoleBindingMap.Update(newObj.(*v1.ClusterRoleBinding))
	if err!=nil{
		log.Println(err)
		return
	}
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrolebinding",
			"result":gin.H{"ns":"clusterrolebinding",
				"data":this.RoleService.ListClusterRoleBindings()},
		},
	)
}
func(this *ClusterRoleBindingHander) OnDelete(obj interface{}){
	this.ClusterRoleBindingMap.Delete(obj.(*v1.ClusterRoleBinding))
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrolebinding",
			"result":gin.H{"ns":"clusterrolebinding",
				"data":this.RoleService.ListClusterRoleBindings()},
		},
	)
}

type SaHander struct {
	SaMap *SaMapStruct  `inject:"-"`
	SaService *SaService  `inject:"-"`
}
func(this *SaHander) OnAdd(obj interface{}){
	this.SaMap.Add(obj.(*corev1.ServiceAccount))
	ns:=obj.(*corev1.ServiceAccount).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"sa",
			"result":gin.H{"ns":ns,
				"data":this.SaService.ListSa(ns)},
		},
	)
}
func(this *SaHander) OnUpdate(oldObj, newObj interface{}){
	err:=this.SaMap.Update(newObj.(*corev1.ServiceAccount))
	if err!=nil{
		return
	}
	ns:=newObj.(*corev1.ServiceAccount).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"sa",
			"result":gin.H{"ns":ns,
				"data":this.SaService.ListSa(ns)},
		},
	)
}
func(this *SaHander) OnDelete(obj interface{}){
	this.SaMap.Delete(obj.(*corev1.ServiceAccount))
	ns:=obj.(*corev1.ServiceAccount).Namespace
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"sa",
			"result":gin.H{"ns":ns,
				"data":this.SaService.ListSa(ns)},
		},
	)
}


type ClusterRoleHandler struct {
	ClusterRoleMap *ClusterRoleMapStruct  `inject:"-"`
 	RoleService *RoleService `inject:"-"`
}
func(this *ClusterRoleHandler) OnAdd(obj interface{}){
	this.ClusterRoleMap.Add(obj.(*v1.ClusterRole))
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrole",
			"result":gin.H{"ns":"clusterrole",
				"data":this.RoleService.ListClusterRoles()},
		},
	)
}
func(this *ClusterRoleHandler) OnUpdate(oldObj, newObj interface{}){
	err:=this.ClusterRoleMap.Update(newObj.(*v1.ClusterRole))
	if err!=nil{
		return
	}
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrole",
			"result":gin.H{"ns":"clusterrole",
				"data":this.RoleService.ListClusterRoles()},
		},
	)
}
func(this *ClusterRoleHandler) OnDelete(obj interface{}){
	this.ClusterRoleMap.Delete(obj.(*v1.ClusterRole))
	wscore.ClientMap.SendAll(
		gin.H{
			"type":"clusterrole",
			"result":gin.H{"ns":"clusterrole",
				"data":this.RoleService.ListClusterRoles()},
		},
	)
}