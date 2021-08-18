package rbac

import "k8s.io/api/rbac/v1"

//@Service
type RoleService struct {
   RoleMap *RoleMapStruct  `inject:"-"`
   ClusterRoleMap *ClusterRoleMapStruct  `inject:"-"`
   RoleBindingMap *RoleBindingMapStruct  `inject:"-"`
	ClusterRoleBindingMap *ClusterRoleBindingMapStruct  `inject:"-"`
}

func NewRoleService() *RoleService {
	return &RoleService{}
}
func(this *RoleService) ListRoles(ns string) []*RoleModel {
	list:=this.RoleMap.ListAll(ns)
	ret:=make([]*RoleModel,len(list))
	for i,item:=range list{
		ret[i]=&RoleModel{
			Name:item.Name,
			CreateTime:item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace:item.Namespace,
		}
	}
	return ret
}
func(this *RoleService) ListClusterRoles() []*v1.ClusterRole {
	 return this.ClusterRoleMap.ListAll()
}
func(this *RoleService) ListClusterRoleBindings() []*v1.ClusterRoleBinding {
	return this.ClusterRoleBindingMap.ListAll()
}
func(this *RoleService) ListRoleBindings(ns string) []*RoleBindingModel {
	list:=this.RoleBindingMap.ListAll(ns)
	ret:=make([]*RoleBindingModel,len(list))
	for i,item:=range list{
		ret[i]=&RoleBindingModel{
			Name:item.Name,
			CreateTime:item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace:item.Namespace,
			Subject:item.Subjects,
			RoleRef:item.RoleRef,
		}
	}
	return ret
}
func(this *RoleService) GetRole(ns ,name string) *v1.Role{
	rb:= this.RoleMap.Get(ns,name)
	if rb==nil{
		panic("no such role")
	}
	return rb
}
func(this *RoleService) GetClusterRole(name string) *v1.ClusterRole{
	rb:= this.ClusterRoleMap.Get(name)
	if rb==nil{
		panic("no such cluster-role")
	}
	return rb
}
func(this *RoleService) GetRoleBinding(ns ,name string) *v1.RoleBinding{
	rb:= this.RoleBindingMap.Get(ns,name)
	if rb==nil{
		panic("no such rolebinding")
	}
	return rb
}
func(this *RoleService) GetClusterRoleBinding(name string) *v1.ClusterRoleBinding{
	crb:= this.ClusterRoleBindingMap.Get(name)
	if crb==nil{
		panic("no such clusterrolebinding")
	}
	return crb
}
