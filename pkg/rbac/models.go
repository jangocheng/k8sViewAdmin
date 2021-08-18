package rbac

import "k8s.io/api/rbac/v1"

type RoleModel struct {
	Name string
	NameSpace string
	CreateTime string
}
type RoleBindingModel struct {
	Name string
	NameSpace string
	CreateTime string
	RoleRef v1.RoleRef
	Subject []v1.Subject  //包含了 绑定用户 数据
}

//UserAccount 模型
type UAModel struct {
	Name string
	CreateTime string
}
//提交用户时的对象模型
type PostUAModel struct {
	CN string `json:"cn" binding:"required,min=2"`
	O string  `json:"o"`
}




