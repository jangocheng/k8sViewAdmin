package configs

import (
	"k8s.io/client-go/kubernetes"
	"k8sapi/pkg/rbac"
	"k8sapi/src/services"
)

type K8sMaps struct {
	K8sClient *kubernetes.Clientset `inject:"-"`
}
func NewK8sMaps() *K8sMaps {
	return &K8sMaps{}
}
//初始化 deploymentmap
func(this *K8sMaps) InitDepMap() *services.DeploymentMap{
	return &services.DeploymentMap{}
}
//初始化 podmap
func(this *K8sMaps) InitPodMap() *services.PodMapStruct{
	return &services.PodMapStruct{}
}
//初始化 nsmap
func(this *K8sMaps) InitNsMap() *services.NsMapStruct{
	return &services.NsMapStruct{}
}
//初始化 eventmap
func(this *K8sMaps) InitEventMap() *services.EventMapStruct{
	return &services.EventMapStruct{}
}
//初始化 ingress map
func(this *K8sMaps) InitIngressMap() *services.IngressMapStruct{
	return &services.IngressMapStruct{}
}

//初始化 Service map
func(this *K8sMaps) InitServiceMap() *services.ServiceMapStruct{
	return &services.ServiceMapStruct{}
}

//初始化 Secret map
func(this *K8sMaps) InitSecretMap() *services.SecretMapStruct{
	return &services.SecretMapStruct{}
}

//初始化 ConfigMap map
func(this *K8sMaps) InitConfigMap() *services.ConfigMapStruct{
	return &services.ConfigMapStruct{}
}

//初始化NodeMap
func(this *K8sMaps) InitNodeMap() *services.NodeMapStruct{
	return &services.NodeMapStruct{}
}
//初始化RoleMap
func(this *K8sMaps) InitRoleMap() *rbac.RoleMapStruct{
	return &rbac.RoleMapStruct{}
}
//初始化RoleBindingMap
func(this *K8sMaps) InitRoleBindingMap() *rbac.RoleBindingMapStruct{
	return &rbac.RoleBindingMapStruct{}
}

//初始化RoleBindingMap
func(this *K8sMaps) InitSaMap() *rbac.SaMapStruct{
	return &rbac.SaMapStruct{}
}

//初始化ClusterRole
func(this *K8sMaps) InitClusterRoleMap() *rbac.ClusterRoleMapStruct{
	return &rbac.ClusterRoleMapStruct{}
}
//初始化ClusterRoleBinding
func(this *K8sMaps) InitClusterRoleBindingMap() *rbac.ClusterRoleBindingMapStruct{
	return &rbac.ClusterRoleBindingMapStruct{}
}
//初始化 ReplicaSetMap
func(this *K8sMaps) InitRsMap() *services.RsMapStruct{
	return &services.RsMapStruct{}
}
