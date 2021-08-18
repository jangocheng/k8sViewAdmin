package deployoments

import (
	"github.com/gin-gonic/gin"
	"github.com/shenyisyn/goft-gin/goft"
	"k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8sapi/src/services"
)

type DeploymentCtlV2 struct {
	K8sClient *kubernetes.Clientset `inject:"-"`
	DeployMap *services.DeploymentMap `inject:"-"`
	RsMap *services.RsMapStruct  `inject:"-"`
	PodMap *services.PodMapStruct `inject:"-"`
}

func NewDeploymentCtlV2() *DeploymentCtlV2 {
	return &DeploymentCtlV2{}
}
//快捷创建时  需要 初始化一些 标签
func(this *DeploymentCtlV2) initLabel(deploy *v1.Deployment){
		if deploy.Spec.Selector==nil{
			deploy.Spec.Selector=&metav1.LabelSelector{MatchLabels: map[string]string{"jtapp":deploy.Name,}}
		}
		if deploy.Spec.Selector.MatchLabels==nil{
			deploy.Spec.Selector.MatchLabels=map[string]string{"jtapp":deploy.Name}
		}
		if deploy.Spec.Template.ObjectMeta.Labels==nil{
			deploy.Spec.Template.ObjectMeta.Labels=map[string]string{"jtapp":deploy.Name}
		}
		deploy.Spec.Selector.MatchLabels["jtapp"]=deploy.Name
		deploy.Spec.Template.ObjectMeta.Labels["jtapp"]=deploy.Name
}
func(this *DeploymentCtlV2) SaveDeployment(c *gin.Context) goft.Json{
	dep:=&v1.Deployment{}
	goft.Error(c.ShouldBindJSON(dep))
	if(c.Query("fast")!=""){ //代表是快捷创建 。 要预定义一些值
		 this.initLabel(dep)
	}
	update:=c.Query("update") //代表是更新
	if update!=""{
		_,err:=this.K8sClient.AppsV1().Deployments(dep.Namespace).Update(c,dep,metav1.UpdateOptions{})
		goft.Error(err)
	}else{
		_,err:=this.K8sClient.AppsV1().Deployments(dep.Namespace).Create(c,dep,metav1.CreateOptions{})
		goft.Error(err)
	}

	return gin.H{
		"code":20000,
		"data":"success",
	}
}

func(this *DeploymentCtlV2) RmDeployment(c *gin.Context) goft.Json{
	ns:=c.Param("ns")
	name:=c.Param("name")

	 err:=this.K8sClient.AppsV1().Deployments(ns).Delete(c,name,metav1.DeleteOptions{})
	 goft.Error(err)
	return gin.H{
		"code":20000,
		"data":"success",
	}
}
//加载deploy 详细
func(this *DeploymentCtlV2) LoadDeploy(c *gin.Context) goft.Json{
	ns:=c.Param("ns")
	name:=c.Param("name")
	dep,err:=this.DeployMap.GetDeployment(ns,name)// 原生
	goft.Error(err)
	return gin.H{
		"code":20000,
		"data":dep,
	}
}
func(this *DeploymentCtlV2) isRsFromDep(dep *v1.Deployment,set v1.ReplicaSet) bool{
	for _,ref:=range set.OwnerReferences{
		if ref.Kind=="Deployment" && ref.Name==dep.Name{
			return true
		}
	}
	return false
}
//获取deployment下的 ReplicaSet的 标签集合
func(this *DeploymentCtlV2) getLablesByDep(dep *v1.Deployment,ns string ) ([]map[string]string,error){
	rslist,err:= this.RsMap.ListByNameSpace(ns)  // 根据namespace 取到 所有rs
	goft.Error(err)
	ret:=make([]map[string]string,0)
	for _,item:=range rslist{
		if this.isRsFromDep(dep,*item){
			s,err:= metav1.LabelSelectorAsMap(item.Spec.Selector)
			if err!=nil{
				return nil,err
			}
			ret=append(ret,s)
		}
	}
	return ret,nil
}
//加载deploymeny的pods列表
func(this *DeploymentCtlV2) LoadDeployPods(c *gin.Context) goft.Json{
	ns:=c.Param("ns")
	name:=c.Param("name")
	dep,err:=this.DeployMap.GetDeployment(ns,name)// 原生
	goft.Error(err)

	labels,err:=this.getLablesByDep(dep,ns) //根据deployment过滤出 rs，然后直接获取标签
	goft.Error(err)
	podList,err :=this.PodMap.ListByLabels(ns,labels)
	goft.Error(err)


	return gin.H{
		"code":20000,
		"data":podList,
	}
}


func(this *DeploymentCtlV2)  Build(goft *goft.Goft){
	//路由
	goft.Handle("GET","/deployments/:ns/:name",this.LoadDeploy)
	goft.Handle("POST","/deployments",this.SaveDeployment)

	//根据Deployment获取PODS
	goft.Handle("GET","/deployments-pods/:ns/:name",this.LoadDeployPods)
	//删除deploy
	goft.Handle("DELETE","/deployments/:ns/:name",this.RmDeployment)
}
func(*DeploymentCtlV2) Name() string{
	return "DeploymentCtlV2"
}