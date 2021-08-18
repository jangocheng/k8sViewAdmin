package rbac

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8sapi/src/helpers"
	"k8sapi/src/models"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type RBACCtl struct {
	RoleService *RoleService `inject:"-"`
	SaService *SaService `inject:"-"`
	Client *kubernetes.Clientset `inject:"-"`
	SysConfig *models.SysConfig `inject:"-"`
}

func NewRBACCtl() *RBACCtl {
	return &RBACCtl{}
}
func(this *RBACCtl) Roles(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code":20000,
		"data":this.RoleService.ListRoles(ns),
	}
}
func(this *RBACCtl) ClusterRoles(c *gin.Context) goft.Json{
	return gin.H{
		"code":20000,
		"data":this.RoleService.ListClusterRoles(),
	}
}


//获取角色详细
func(this *RBACCtl) RolesDetail(c *gin.Context) goft.Json{
	ns:=c.Param("ns")
	rname:=c.Param("rolename")
	return gin.H{
		"code":20000,
		"data":this.RoleService.GetRole(ns,rname),
	}
}
////获取集群角色详细
func(this *RBACCtl) ClusterRolesDetail(c *gin.Context) goft.Json{

	rname:=c.Param("cname") //集群角色名
	return gin.H{
		"code":20000,
		"data":this.RoleService.GetClusterRole(rname),
	}
}
func(this *RBACCtl) CreateRole(c *gin.Context) goft.Json{
	role:=rbacv1.Role{} //原生的k8s role 对象
	goft.Error(c.ShouldBindJSON(&role))
	role.APIVersion="rbac.authorization.k8s.io/v1"
	role.Kind="Role"
	_,err:=this.Client.RbacV1().Roles(role.Namespace).Create(c,&role,metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
//创建集群角色
func(this *RBACCtl) CreateClusterRole(c *gin.Context) goft.Json{
	clusterRole:=rbacv1.ClusterRole{} //原生的k8s role 对象
	goft.Error(c.ShouldBindJSON(&clusterRole))
	clusterRole.APIVersion="rbac.authorization.k8s.io/v1"
	clusterRole.Kind="ClusterRole"
	_,err:=this.Client.RbacV1().ClusterRoles().Create(c,&clusterRole,metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
//更新集群角色
func(this *RBACCtl) UpdateClusterRolesDetail(c *gin.Context) goft.Json{
	cname:=c.Param("cname") //集群角色名
	clusterRole:=this.RoleService.GetClusterRole(cname)
	postRole:=rbacv1.ClusterRole{}
	goft.Error(c.ShouldBindJSON(&postRole))  //获取提交过来的对象

	clusterRole.Rules=postRole.Rules //目前修改只允许修改 rules，其他不允许。大家可以自行扩展，如标签也允许修改
	_,err:=this.Client.RbacV1().ClusterRoles().Update(c,clusterRole,metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}


//更新角色
func(this *RBACCtl) UpdateRolesDetail(c *gin.Context) goft.Json{
	ns:=c.Param("ns")
	rname:=c.Param("rolename")
	role:=this.RoleService.GetRole(ns,rname)
	postRole:=rbacv1.Role{}
	goft.Error(c.ShouldBindJSON(&postRole))  //获取提交过来的对象

	role.Rules=postRole.Rules //目前修改只允许修改 rules，其他不允许。大家可以自行扩展，如标签也允许修改
	_,err:=this.Client.RbacV1().Roles(role.Namespace).Update(c,role,metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) RoleBindingList(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code":20000,
		"data":this.RoleService.ListRoleBindings(ns),
	}
}
func(this *RBACCtl) ClusterRoleBindingList(c *gin.Context) goft.Json{

	return gin.H{
		"code":20000,
		"data":this.RoleService.ListClusterRoleBindings(),
	}
}

func(this *RBACCtl) CreateRoleBinding(c *gin.Context) goft.Json{
	rb:=&rbacv1.RoleBinding{}
	goft.Error(c.ShouldBindJSON(rb))
	_,err:=this.Client.RbacV1().RoleBindings(rb.Namespace).Create(c,rb,metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) CreateClusterRoleBinding(c *gin.Context) goft.Json{
	rb:=&rbacv1.ClusterRoleBinding{}
	goft.Error(c.ShouldBindJSON(rb))
	_,err:=this.Client.RbacV1().ClusterRoleBindings().Create(c,rb,metav1.CreateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
//从rolebinding中 增加或删除用户
func(this *RBACCtl) AddUserToRoleBinding(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	name:=c.DefaultQuery("name","")//rolebinding 名称
	t:=c.DefaultQuery("type","") //如果没传值就是增加，传值（不管什么代表删除)
	subject:=rbacv1.Subject{}// 传过来
	goft.Error(c.ShouldBindJSON(&subject))
	if subject.Kind=="ServiceAccount"{
		subject.APIGroup=""
	}
	rb:=this.RoleService.GetRoleBinding(ns,name) //通过名称获取 rolebinding对象
	if t!=""{//代表删除
		for i,sub:=range rb.Subjects{
			if sub.Kind==subject.Kind && sub.Name==subject.Name{
				rb.Subjects=append(rb.Subjects[:i], rb.Subjects[i+1:]...)
				break //确保只删一个（哪怕有同名同kind用户)
			}
		}
		fmt.Println(rb.Subjects)
	}else{
		rb.Subjects=append(rb.Subjects,subject)
	}
	_,err:=this.Client.RbacV1().RoleBindings(ns).Update(c,rb,metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) AddUserToClusterRoleBinding(c *gin.Context) goft.Json{

	name:=c.DefaultQuery("name","")//clusterrolebinding 名称
	t:=c.DefaultQuery("type","") //如果没传值就是增加，传值（不管什么代表删除)
	subject:=rbacv1.Subject{}// 传过来
	goft.Error(c.ShouldBindJSON(&subject))
	if subject.Kind=="ServiceAccount"{
		subject.APIGroup=""
	}
	rb:=this.RoleService.GetClusterRoleBinding(name) //通过名称获取 clusterrolebinding对象
	if t!=""{//代表删除
		for i,sub:=range rb.Subjects{
			if sub.Kind==subject.Kind && sub.Name==subject.Name{
				rb.Subjects=append(rb.Subjects[:i], rb.Subjects[i+1:]...)
				break //确保只删一个（哪怕有同名同kind用户)
			}
		}
	}else{
		rb.Subjects=append(rb.Subjects,subject)
	}
	_,err:=this.Client.RbacV1().ClusterRoleBindings().Update(c,rb,metav1.UpdateOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}

func(this *RBACCtl) DeleteRole(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	name:=c.DefaultQuery("name","")
	err:=this.Client.RbacV1().Roles(ns).Delete(c,name,metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) DeleteClusterRole(c *gin.Context) goft.Json{
	name:=c.DefaultQuery("name","")
	err:=this.Client.RbacV1().ClusterRoles().Delete(c,name,metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) DeleteRoleBinding(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	name:=c.DefaultQuery("name","")
	err:=this.Client.RbacV1().RoleBindings(ns).Delete(c,name,metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) DeleteClusterRoleBinding(c *gin.Context) goft.Json{

	name:=c.DefaultQuery("name","")
	err:=this.Client.RbacV1().ClusterRoleBindings().Delete(c,name,metav1.DeleteOptions{})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
func(this *RBACCtl) SaList(c *gin.Context) goft.Json{
	ns:=c.DefaultQuery("ns","default")
	return gin.H{
		"code":20000,
		"data":this.SaService.ListSa(ns),
	}
}

//UserAccount 列表
func(this *RBACCtl) UaList(c *gin.Context) goft.Json{
		uaPath:="./k8susers"  //写死的路径存储证书
	keyReg:=regexp.MustCompile(".*_key.pem")
    users:=[]*UAModel{}
    suffix:=".pem"
	err:= filepath.Walk(uaPath, func(p string, f os.FileInfo, err error) error {
		if f.IsDir() {return nil}
		if path.Ext(f.Name())==suffix{
			if !keyReg.MatchString(f.Name()){
				users=append(users,&UAModel{
					Name:strings.Replace(f.Name(),suffix,"",-1),
					CreateTime:f.ModTime().Format("2006-01-02 15:04:05"),
				})
			}
		}
		return nil
	})
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":users,
	}

}
func(this *RBACCtl) PostUa(c *gin.Context) goft.Json{
	postModel:=&PostUAModel{}
	goft.Error(c.ShouldBindJSON(postModel))
	helpers.GenK8sUser(postModel.CN,postModel.O)
	return gin.H{
		"code":20000,
		"data":"success",
	}

}
func(this *RBACCtl) DeleteUa(c *gin.Context) goft.Json{
	postModel:=&PostUAModel{}
	goft.Error(c.ShouldBindJSON(postModel))
	helpers.DeleteK8sUser(postModel.CN)
	return gin.H{
		"code":20000,
		"data":"success",
	}

}
//生成 config文件
func(this *RBACCtl) Clientconfig(c *gin.Context) goft.Json{
	user:=c.DefaultQuery("user","")
	if user==""{
		panic("no such user")
	}
	cfg:=api.NewConfig()
	clusterName:="kubernetes"
	cfg.Clusters[clusterName]=&api.Cluster{
		Server:this.SysConfig.K8s.ClusterInfo.EndPoint,
		CertificateAuthorityData:helpers.CertData(this.SysConfig.K8s.ClusterInfo.CaFile),
	}
	contextName:=fmt.Sprintf("%s@kubernetes",user)
	cfg.Contexts[contextName]=&api.Context{
		AuthInfo:user,
		Cluster:clusterName,
	}
	cfg.CurrentContext=contextName
	userCertFile:=fmt.Sprintf("%s/%s.pem",this.SysConfig.K8s.ClusterInfo.UserCert,user)
	userCertKeyFile:=fmt.Sprintf("%s/%s_key.pem",this.SysConfig.K8s.ClusterInfo.UserCert,user)
	cfg.AuthInfos[user]=&api.AuthInfo{
		ClientKeyData:helpers.CertData(userCertKeyFile),
		ClientCertificateData:helpers.CertData(userCertFile),
	}
	fileContent,err:=clientcmd.Write(*cfg)
	goft.Error(err)

	return gin.H{
		"code":20000,
		"data":string(fileContent),
	}
}

func(*RBACCtl) Name() string{
	return "RBACCtl"
}
func(this *RBACCtl)  Build(goft *goft.Goft){
	goft.Handle("GET","/clusterroles",this.ClusterRoles)
	goft.Handle("DELETE","/clusterroles",this.DeleteClusterRole)

	goft.Handle("GET","/roles",this.Roles)
	goft.Handle("GET","/clusterroles/:cname",this.ClusterRolesDetail)
	goft.Handle("GET","/roles/:ns/:rolename",this.RolesDetail)
	goft.Handle("POST","/roles/:ns/:rolename",this.UpdateRolesDetail) //修改角色
	goft.Handle("POST","/clusterroles/:cname",this.UpdateClusterRolesDetail)
	goft.Handle("GET","/rolebindings",this.RoleBindingList)
	goft.Handle("POST","/rolebindings",this.CreateRoleBinding)
	goft.Handle("PUT","/rolebindings",this.AddUserToRoleBinding) //添加用户到binding
	goft.Handle("DELETE","/rolebindings",this.DeleteRoleBinding)
	goft.Handle("POST","/roles",this.CreateRole)
	goft.Handle("POST","/clusterroles",this.CreateClusterRole) //创建集群角色
	goft.Handle("DELETE","/roles",this.DeleteRole)



	goft.Handle("GET","/clusterrolebindings",this.ClusterRoleBindingList)
	goft.Handle("POST","/clusterrolebindings",this.CreateClusterRoleBinding)
	goft.Handle("PUT","/clusterrolebindings",this.AddUserToClusterRoleBinding)
	goft.Handle("DELETE","/clusterrolebindings",this.DeleteClusterRoleBinding)


	goft.Handle("GET","/sa",this.SaList)

	goft.Handle("GET","/ua",this.UaList)
	goft.Handle("POST","/ua",this.PostUa)
	goft.Handle("DELETE","/ua",this.DeleteUa)

	goft.Handle("GET","/clientconfig",this.Clientconfig)

}