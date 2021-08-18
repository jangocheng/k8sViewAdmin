package configs

import (
	"k8sapi/pkg/rbac"
	"k8sapi/src/services"
)

//注入 回调handler
type K8sHandler struct {}

func NewK8sHandler() *K8sHandler {
	return &K8sHandler{}
}
// deployment handler
func(this *K8sHandler) DepHandlers() *services.DepHandler{
	return &services.DepHandler{}
}
// pod handler
func(this *K8sHandler) PodHandlers() *services.PodHandler{
	return &services.PodHandler{}
}
//ns handler
func(this *K8sHandler) NsHandlers() *services.NsHandler{
	return &services.NsHandler{}
}

// event handler
func(this *K8sHandler) EventHandlers() *services.EventHandler{
	return &services.EventHandler{}
}
// IngressHandler
func(this *K8sHandler) IngressHandler() *services.IngressHandler{
	return &services.IngressHandler{}
}

// ServiceHandler
func(this *K8sHandler) ServiceHandler() *services.ServiceHandler{
	return &services.ServiceHandler{}
}

// SecretHandler
func(this *K8sHandler) SecretHandler() *services.SecretHandler{
	return &services.SecretHandler{}
}

// ConfigMapHandler
func(this *K8sHandler) ConfigMapHandler() *services.ConfigMapHandler{
	return &services.ConfigMapHandler{}
}

// ConfigMapHandler
func(this *K8sHandler) ConfigNodeHandler() *services.NodeMapHandler{
	return &services.NodeMapHandler{}
}

// RoleHandler
func(this *K8sHandler) ConfigRoleHandler() *rbac.RoleHander{
	return &rbac.RoleHander{}
}
// RoleBindingHandler
func(this *K8sHandler) ConfigRoleBindingHandler() *rbac.RoleBindingHander{
	return &rbac.RoleBindingHander{}
}
// RoleBindingHandler
func(this *K8sHandler) ConfigSaHandler() *rbac.SaHander{
	return &rbac.SaHander{}
}
// ClusterRoleHandler
func(this *K8sHandler) ConfigClusterRoleHandler() *rbac.ClusterRoleHandler{
	return &rbac.ClusterRoleHandler{}
}
// ClusterRoleBindingHandler
func(this *K8sHandler) ConfigClusterRoleBindingHandler() *rbac.ClusterRoleBindingHander{
	return &rbac.ClusterRoleBindingHander{}
}
// RShandler
func(this *K8sHandler) ConfigRsHandler() *services.RsHandler{
	return &services.RsHandler{}
}