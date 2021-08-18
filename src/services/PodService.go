package services

import (
	"github.com/gin-gonic/gin"
	"k8sapi/src/models"
)

//@Service
type PodService struct {
	PodMap *PodMapStruct `inject:"-"`
	Common *CommonService `inject:"-"`
	EventMap *EventMapStruct `inject:"-"`
	Helper *Helper `inject:"-"`
}

func NewPodService() *PodService {
	return &PodService{}
}
//获取Pod容器列表
func(this *PodService) GetPodContainer(ns,podname string) []*models.ContainerModel{
	ret:=make([]*models.ContainerModel,0)
	pod:=this.PodMap.Get(ns,podname)
	if pod!=nil{
		for _,c:=range pod.Spec.Containers{
			ret=append(ret,&models.ContainerModel{
				Name: c.Name,
			})
		}
	}
	return ret
}

//分页PODS的输出
func(this *PodService) PagePods(ns string,page,size int ) *ItemsPage{
	pods:=this.ListByNs(ns).([]*models.Pod)
	readyCount:=0 //就绪的pod数量
	allCount:=0 //总数量
	ipods:=make([]interface{},len(pods))
	for i,pod:=range pods{
		allCount++
		ipods[i]=pod
		if pod.IsReady{
			readyCount++
		}
	}
	return this.Helper.PageResource(
		 page,
		 size,
		ipods).SetExt(gin.H{"ReadyNum":readyCount,"AllNum":allCount})
}
//不分页， 显示全部
func(this *PodService) ListByNs(ns string ) interface{}{
	podList:=this.PodMap.ListByNs(ns)
	 ret:=make([]*models.Pod,0)
	for _,pod:=range podList{
		ret=append(ret,&models.Pod{
			Name:pod.Name,
			NameSpace:pod.Namespace,
			Images:this.Common.GetImagesByPod(pod.Spec.Containers),
			NodeName:pod.Spec.NodeName,
			Phase:string(pod.Status.Phase),// 阶段
			IsReady:this.Common.PosIsReady(pod), //是否就绪
			IP:[]string{pod.Status.PodIP,pod.Status.HostIP},
			Message:this.EventMap.GetMessage(pod.Namespace,"Pod",pod.Name),
			CreateTime:pod.CreationTimestamp.Format("2006-01-02 15:04:05"),
		})
	}
	return ret
}