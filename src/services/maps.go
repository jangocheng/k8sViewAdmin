package services

import (
	"fmt"
	"k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/networking/v1beta1"
	"k8sapi/src/helpers"
	"k8sapi/src/models"
	"sort"
	"sync"
)


type MapItems []*MapItem
type MapItem struct {
	key string
	value interface{}
}
func(this *MapItem) String() string{
	return this.key
}
//把sync.map  转为 自定义切片
func convertToMapItems(m sync.Map) MapItems{
	items:=make(MapItems,0)
	m.Range(func(key, value interface{}) bool {
		items=append(items,&MapItem{key:key.(string),value:value})
		return true
	})
	return items
}
func(this MapItems) Len() int{
	return len(this)
}
func(this MapItems) Less(i, j int) bool{
	return this[i].key<this[j].key
}
func(this MapItems) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}

type V1Deployment []*v1.Deployment
func(this V1Deployment) Len() int{
	return len(this)
}
func(this V1Deployment) Less(i, j int) bool{
	//根据时间排序    倒排序
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func(this V1Deployment) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}

//对deployments的集合进行定义
type DeploymentMap struct {
	data sync.Map  // [key string] []*v1.Deployment    key=>namespace
}
//添加
func(this *DeploymentMap) Add(dep *v1.Deployment){

	if list,ok:=this.data.Load(dep.Namespace);ok{
		list=append(list.([]*v1.Deployment),dep)
		this.data.Store(dep.Namespace,list)
	}else{
		this.data.Store(dep.Namespace,[]*v1.Deployment{dep})
	}
}
//更新
func(this *DeploymentMap) Update(dep *v1.Deployment) error {
	if list,ok:=this.data.Load(dep.Namespace);ok{
		for i,range_dep:=range list.([]*v1.Deployment){
			if range_dep.Name==dep.Name{
				list.([]*v1.Deployment)[i]=dep
			}
		}
		return nil
	}
	return fmt.Errorf("deployment-%s not found",dep.Name)
}
// 删除
func(this *DeploymentMap) Delete(dep *v1.Deployment){
	if list,ok:=this.data.Load(dep.Namespace);ok{
		for i,range_dep:=range list.([]*v1.Deployment){
			if range_dep.Name==dep.Name{
				newList:= append(list.([]*v1.Deployment)[:i], list.([]*v1.Deployment)[i+1:]...)
				this.data.Store(dep.Namespace,newList)
				break
			}
		}
	}
}
func(this *DeploymentMap) ListByNS(ns string) ([]*v1.Deployment,error){
	if list,ok:=this.data.Load(ns);ok {
		 getList:= V1Deployment(list.([]*v1.Deployment))
		 sort.Sort(getList)
		 return getList,nil
	}
	return nil,fmt.Errorf("record not found")
}
func(this *DeploymentMap) GetDeployment(ns string,depname string) (*v1.Deployment,error){
	if list,ok:=this.data.Load(ns);ok {
		for _,item:=range list.([]*v1.Deployment){
			if item.Name==depname{
				return item,nil
			}
		}
	}
	return nil,fmt.Errorf("record not found")
}


// ReplicaSet 集合
type RsMapStruct struct {
	data sync.Map   // [key string] []*appv1.ReplicaSet    key=>namespace
}
func(this *RsMapStruct) Add(rs *v1.ReplicaSet){
	if list,ok:=this.data.Load(rs.Namespace);ok{
		list=append(list.([]*v1.ReplicaSet),rs)
		this.data.Store(rs.Namespace,list)
	}else{
		this.data.Store(rs.Namespace,[]*v1.ReplicaSet{rs})
	}
}
func(this *RsMapStruct) Update(rs *v1.ReplicaSet) error {
	if list,ok:=this.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*v1.ReplicaSet){
			if range_rs.Name==rs.Name{
				list.([]*v1.ReplicaSet)[i]=rs
			}
		}
		return nil
	}
	return fmt.Errorf("rs-%s not found",rs.Name)
}
func(this *RsMapStruct) Delete(rs *v1.ReplicaSet){
	if list,ok:=this.data.Load(rs.Namespace);ok{
		for i,range_rs:=range list.([]*v1.ReplicaSet){
			if range_rs.Name==rs.Name{
				newList:= append(list.([]*v1.ReplicaSet)[:i], list.([]*v1.ReplicaSet)[i+1:]...)
				this.data.Store(rs.Namespace,newList)
				break
			}
		}
	}
}
//普普通通的函数， 就是根据ns获取 对应的rs列表
func(this *RsMapStruct) ListByNameSpace(ns string) ([]*v1.ReplicaSet,error){
	if list,ok:=this.data.Load(ns);ok {
		return list.([]*v1.ReplicaSet),nil
	}
	return nil,fmt.Errorf("pods not found ")
}


type CoreV1Pods []*corev1.Pod
func(this CoreV1Pods) Len() int{
	return len(this)
}
func(this CoreV1Pods) Less(i, j int) bool{
	//根据时间排序    正排序
	return this[i].CreationTimestamp.Time.Before(this[j].CreationTimestamp.Time)
}
func(this CoreV1Pods) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}
// 保存Pod集合
type PodMapStruct struct {
	data sync.Map  // [key string] []*v1.Pod    key=>namespace
}
//根据节点名称 获取pods数量
func(this *PodMapStruct) GetNum(nodeName string) (num int ){
	 this.data.Range(func(key, value interface{}) bool {
		 list:=value.([]*corev1.Pod)
		 for _,pod:=range list{
		 	if pod.Spec.NodeName==nodeName{
		 		num++
			}
		 }
		 return true
	 })
	return
}
func(this *PodMapStruct) ListByNs(ns string) []*corev1.Pod{
	if list,ok:=this.data.Load(ns);ok{
		 ret:=list.([]*corev1.Pod)
		 sort.Sort(CoreV1Pods(ret))//排序
		 return ret
	}
	return nil
}

func(this *PodMapStruct) Get(ns string,podName string) *corev1.Pod{
	if list,ok:=this.data.Load(ns);ok{
		for _,pod:=range list.([]*corev1.Pod){
			if pod.Name==podName{
				return pod
			}
		}
	}
	return nil
}
func(this *PodMapStruct) Add(pod *corev1.Pod){
	if list,ok:=this.data.Load(pod.Namespace);ok{
		list=append(list.([]*corev1.Pod),pod)
		this.data.Store(pod.Namespace,list)
	}else{
		this.data.Store(pod.Namespace,[]*corev1.Pod{pod})
	}
}
func(this *PodMapStruct) Update(pod *corev1.Pod) error {
	if list,ok:=this.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				list.([]*corev1.Pod)[i]=pod
			}
		}
		return nil
	}
	return fmt.Errorf("Pod-%s not found",pod.Name)
}
func(this *PodMapStruct) Delete(pod *corev1.Pod){
	if list,ok:=this.data.Load(pod.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Pod){
			if range_pod.Name==pod.Name{
				newList:= append(list.([]*corev1.Pod)[:i], list.([]*corev1.Pod)[i+1:]...)
				this.data.Store(pod.Namespace,newList)
				break
			}
		}
	}
}
//根据标签获取 POD列表   本方法做了修改
func(this *PodMapStruct) ListByLabels(ns string,labels []map[string]string) ([]*corev1.Pod,error){
	ret:=make([]*corev1.Pod,0)
	if list,ok:=this.data.Load(ns);ok {
		for _,item:=range list.([]*corev1.Pod){
			for _,label:=range labels{
				i := 0
				for k1, vv1 := range label {
					for k2, v2 := range item.Labels {
						if k1 == k2 && vv1 == v2 {
							i++
						}
					}
				}
				if i == len(label) {
					ret = append(ret, item)
				}
			}
		}
		return ret,nil
	}
	return nil,fmt.Errorf("pods not found ")
}
func(this *PodMapStruct) DEBUG_ListByNS(ns string) ([]*corev1.Pod){
	ret:=make([]*corev1.Pod,0)
	if list,ok:=this.data.Load(ns);ok {
		for _,pod:=range list.([]*corev1.Pod){
			ret=append(ret,pod)
		}

	}
	return ret
}


// namespace相关
type NsMapStruct struct {
	data sync.Map  // [key string] []*corev1.Namespace    key=>namespace的名称
}
func(this *NsMapStruct) Get(ns string) *corev1.Namespace{
	if item,ok:=this.data.Load(ns);ok{
		 return item.(*corev1.Namespace)
	}
	return nil
}
func(this *NsMapStruct) Add(ns *corev1.Namespace){
	this.data.Store(ns.Name,ns)
}
func(this *NsMapStruct) Update(ns *corev1.Namespace) {
	this.data.Store(ns.Name,ns)
}
func(this *NsMapStruct) Delete(ns *corev1.Namespace){
	this.data.Delete(ns.Name)
}
//显示所有的 namespace
func(this *NsMapStruct) ListAll()[]*models.NsModel{

	 //this.data.Range(func(key, value interface{}) bool {
		//ret=append(ret,&models.NsModel{Name:key.(string)})
	 //	return true
	 //})
	 items:=convertToMapItems(this.data)
	 sort.Sort(items)
	  ret:=make([]*models.NsModel,len(items))
	 for index,item:=range items{
		 ret[index]=&models.NsModel{Name:item.key}
	 }

	return ret
}

// event 事件map 相关
// EventSet 集合 用来保存事件, 只保存最新的一条
type EventMapStruct struct {
	data sync.Map   // [key string] *v1.Event
	// key=>namespace+"_"+kind+"_"+name 这里的name 不一定是pod ,这样确保唯一
}
func(this *EventMapStruct) GetMessage(ns string,kind string,name string) string{
	key:=fmt.Sprintf("%s_%s_%s",ns,kind,name)
	if v,ok:=this.data.Load(key);ok{
		return v.(*corev1.Event).Message
	}
	return ""
}


type V1Beta1Ingress []*v1beta1.Ingress
func(this V1Beta1Ingress) Len() int{
	return len(this)
}
func(this V1Beta1Ingress) Less(i, j int) bool{
	//根据时间排序    倒排序
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func(this V1Beta1Ingress) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}
type IngressMapStruct struct {
	data sync.Map   // [ns string] []*v1beta1.Ingress
}
//获取单个Ingress
func(this *IngressMapStruct) Get(ns string,name string) *v1beta1.Ingress{
	if items,ok:=this.data.Load(ns);ok{
		for _,item:=range items.([]*v1beta1.Ingress){
			if item.Name==name{

				return item
			}
		}
	}
	return nil
}
func(this *IngressMapStruct) Add(ingress *v1beta1.Ingress){
	if list,ok:=this.data.Load(ingress.Namespace);ok{
		list=append(list.([]*v1beta1.Ingress),ingress)
		this.data.Store(ingress.Namespace,list)
	}else{
		this.data.Store(ingress.Namespace,[]*v1beta1.Ingress{ingress})
	}
}
func(this *IngressMapStruct) Update(ingress *v1beta1.Ingress) error {
	if list,ok:=this.data.Load(ingress.Namespace);ok{
		for i,range_pod:=range list.([]*v1beta1.Ingress){
			if range_pod.Name==ingress.Name{
				list.([]*v1beta1.Ingress)[i]=ingress
			}
		}
	 return nil
	}
	return fmt.Errorf("ingress-%s not found",ingress.Name)
}
func(this *IngressMapStruct) Delete(ingress *v1beta1.Ingress){
	if list,ok:=this.data.Load(ingress.Namespace);ok{
		for i,range_ingress:=range list.([]*v1beta1.Ingress){
			if range_ingress.Name==ingress.Name{
				newList:= append(list.([]*v1beta1.Ingress)[:i], list.([]*v1beta1.Ingress)[i+1:]...)
				this.data.Store(ingress.Namespace,newList)
				break
			}
		}
	}
}
func(this *IngressMapStruct) ListAll(ns string)[]*v1beta1.Ingress{
	     if list,ok:=this.data.Load(ns);ok{
	     	 newList:=list.([]*v1beta1.Ingress)
			 sort.Sort(V1Beta1Ingress(newList))//  按时间倒排序
			 return newList
	     	 //之前获取列表代码是写在这的，现在移动到了 IngressService中
		 }
	     return []*v1beta1.Ingress{} //返回空列表
}



type ServiceMapStruct struct {
	data sync.Map   // [ns string] []*v1.Service
}

type CoreV1Service []*corev1.Service
func(this CoreV1Service) Len() int{
	return len(this)
}
func(this CoreV1Service) Less(i, j int) bool{
	//根据时间排序    倒排序
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func(this CoreV1Service) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}
//获取单个Service
func(this *ServiceMapStruct) Get(ns string,name string) *corev1.Service{
	if items,ok:=this.data.Load(ns);ok{
		for _,item:=range items.([]*corev1.Service){
			if item.Name==name{
				return item
			}
		}
	}
	return nil
}
func(this *ServiceMapStruct) Add(svc *corev1.Service){
	if list,ok:=this.data.Load(svc.Namespace);ok{
		list=append(list.([]*corev1.Service),svc)
		this.data.Store(svc.Namespace,list)
	}else{
		this.data.Store(svc.Namespace,[]*corev1.Service{svc})
	}
}
func(this *ServiceMapStruct) Update(svc *corev1.Service) error {
	if list,ok:=this.data.Load(svc.Namespace);ok{
		for i,range_pod:=range list.([]*corev1.Service){
			if range_pod.Name==svc.Name{
				list.([]*corev1.Service)[i]=svc
			}
		}
		return nil
	}
	return fmt.Errorf("service-%s not found",svc.Name)
}
func(this *ServiceMapStruct) Delete(svc *corev1.Service){
	if list,ok:=this.data.Load(svc.Namespace);ok{
		for i,range_svc:=range list.([]*corev1.Service){
			if range_svc.Name==svc.Name{
				newList:= append(list.([]*corev1.Service)[:i], list.([]*corev1.Service)[i+1:]...)
				this.data.Store(svc.Namespace,newList)
				break
			}
		}
	}
}
func(this *ServiceMapStruct) ListAll(ns string)[]*models.ServiceModel{
	if list,ok:=this.data.Load(ns);ok{
		newList:=list.([]*corev1.Service)
		sort.Sort(CoreV1Service(newList))//  按时间倒排序
		ret:=make([]*models.ServiceModel,len(newList))
		for i,item:=range newList{
			ret[i]=&models.ServiceModel{
				Name:item.Name,
				CreateTime:item.CreationTimestamp.Format("2006-01-02 15:04:05"),
				NameSpace:item.Namespace,
			}
		}
		return ret
	}
	return []*models.ServiceModel{} //返回空列表
}



type CoreV1Secret []*corev1.Secret
func(this CoreV1Secret) Len() int{
	return len(this)
}
func(this CoreV1Secret) Less(i, j int) bool{
	//根据时间排序    倒排序
	return this[i].CreationTimestamp.Time.After(this[j].CreationTimestamp.Time)
}
func(this CoreV1Secret) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}
//SecretMap
type SecretMapStruct struct {
	data sync.Map   // [ns string] []*v1.Secret
}
func(this *SecretMapStruct) Get(ns string,name string) *corev1.Secret{
	if items,ok:=this.data.Load(ns);ok{
		for _,item:=range items.([]*corev1.Secret){
			if item.Name==name{
				return item
			}
		}
	}
	return nil
}
func(this *SecretMapStruct) Add(item *corev1.Secret){
	if list,ok:=this.data.Load(item.Namespace);ok{
		list=append(list.([]*corev1.Secret),item)
		this.data.Store(item.Namespace,list)
	}else{
		this.data.Store(item.Namespace,[]*corev1.Secret{item})
	}
}
func(this *SecretMapStruct) Update(item *corev1.Secret) error {
	if list,ok:=this.data.Load(item.Namespace);ok{
		for i,range_item:=range list.([]*corev1.Secret){
			if range_item.Name==item.Name{
				list.([]*corev1.Secret)[i]=item
			}
		}
		return nil
	}
	return fmt.Errorf("Secret-%s not found",item.Name)
}
func(this *SecretMapStruct) Delete(svc *corev1.Secret){
	if list,ok:=this.data.Load(svc.Namespace);ok{
		for i,range_item:=range list.([]*corev1.Secret){
			if range_item.Name==svc.Name{
				newList:= append(list.([]*corev1.Secret)[:i], list.([]*corev1.Secret)[i+1:]...)
				this.data.Store(svc.Namespace,newList)
				break
			}
		}
	}
}
func(this *SecretMapStruct) ListAll(ns string)[]*corev1.Secret{
	if list,ok:=this.data.Load(ns);ok{
		newList:=list.([]*corev1.Secret)
		sort.Sort(CoreV1Secret(newList))//  按时间倒排序

		return newList
	}
	return []*corev1.Secret{} //返回空列表
}


//ConfigMapMap

type CoreV1ConfigMap []*cm
func(this CoreV1ConfigMap) Len() int{
	return len(this)
}
func(this CoreV1ConfigMap) Less(i, j int) bool{
	//根据时间排序    倒排序
	return this[i].cmdata.CreationTimestamp.Time.After(this[j].cmdata.CreationTimestamp.Time)
}
func(this CoreV1ConfigMap) Swap(i, j int){
	this[i],this[j]=this[j],this[i]
}
//给configmap的特殊struct
type cm struct {
	cmdata *corev1.ConfigMap
	md5 string  //对cm的data进行md5存储，防止过度更新
}

func newcm(c *corev1.ConfigMap) *cm  {
	return &cm{
		cmdata:c,//原始对象
		md5:helpers.Md5Data(c.Data),
	}
}
type ConfigMapStruct struct {
	data sync.Map   // [ns string] []*cm
}
func(this *ConfigMapStruct) Get(ns string,name string) *corev1.ConfigMap{
	if items,ok:=this.data.Load(ns);ok{
		for _,item:=range items.([]*cm){
			if item.cmdata.Name==name{
				return item.cmdata
			}
		}
	}
	return nil
}
func(this *ConfigMapStruct) Add(item *corev1.ConfigMap){
	if list,ok:=this.data.Load(item.Namespace);ok{
		list=append(list.([]*cm),newcm(item))
		this.data.Store(item.Namespace,list)
	}else{
		this.data.Store(item.Namespace,[]*cm{newcm(item)})
	}
}
//返回值 是true 或false . true代表有值更新了， 否则返回false
func(this *ConfigMapStruct) Update(item *corev1.ConfigMap) bool {
	if list,ok:=this.data.Load(item.Namespace);ok{
		for i,range_item:=range list.([]*cm){
			//这里做判断，如果没变化就不做 更新
			if range_item.cmdata.Name==item.Name && !helpers.CmIsEq(range_item.cmdata.Data,item.Data){
				list.([]*cm)[i]=newcm(item)
				return true //代表有值更新了
			}
		}
	}
	return 	  false
}
func(this *ConfigMapStruct) Delete(svc *corev1.ConfigMap){
	if list,ok:=this.data.Load(svc.Namespace);ok{
		for i,range_item:=range list.([]*cm){
			if range_item.cmdata.Name==svc.Name{
				newList:= append(list.([]*cm)[:i], list.([]*cm)[i+1:]...)
				this.data.Store(svc.Namespace,newList)
				break
			}
		}
	}
}
func(this *ConfigMapStruct) ListAll(ns string)[]*corev1.ConfigMap{
	ret:=[]*corev1.ConfigMap{}
	if list,ok:=this.data.Load(ns);ok{
		newList:=list.([]*cm)
		sort.Sort(CoreV1ConfigMap(newList))//  按时间倒排序
		for _,cm:=range newList{
			ret=append(ret,cm.cmdata)
		}
	}
	return ret//返回空列表
}


//node map
type NodeMapStruct struct {
	data sync.Map   // [nodename string] *v1.Node   注意里面不是切片
}
func(this *NodeMapStruct) Get(name string) *corev1.Node{
	if node,ok:=this.data.Load(name);ok{
		 return node.(*corev1.Node)
	}
	return nil
}
func(this *NodeMapStruct) Add(item *corev1.Node){
	 //直接覆盖
		this.data.Store(item.Name,item)
}

func(this *NodeMapStruct) Update(item *corev1.Node) bool {
	this.data.Store(item.Name,item)
	return true
}
func(this *NodeMapStruct) Delete(node *corev1.Node){
	 this.data.Delete(node.Name)
}
func(this *NodeMapStruct) ListAll()[]*corev1.Node{
	 ret:=[]*corev1.Node{}
	 this.data.Range(func(key, value interface{}) bool {
		 ret=append(ret,value.(*corev1.Node))
		 return true
	 })
	return ret//返回空列表
}