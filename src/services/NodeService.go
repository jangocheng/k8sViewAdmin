package services

import (
	"k8s.io/api/core/v1"
	"k8s.io/metrics/pkg/client/clientset/versioned"
	"k8sapi/src/helpers"
	"k8sapi/src/models"
)

//@service
type NodeService struct {
	NodeMap *NodeMapStruct `inject:"-"`
	PodMap *PodMapStruct `inject:"-"`
    Metric *versioned.Clientset `inject:"-"`
}

func NewNodeService() *NodeService {
	return &NodeService{}
}
//保存时用的
func(this *NodeService) LoadOrginNode(nodeName string ) *v1.Node{
	 return this.NodeMap.Get(nodeName)
}
//加载node信息， 给编辑用的
func(this *NodeService) LoadNode(nodeName string ) *models.NodeModel{
	node:= this.NodeMap.Get(nodeName)
	nodeUsage:=helpers.GetNodeUsage(this.Metric,node)
	return &models.NodeModel{
		Name:node.Name,
		IP:node.Status.Addresses[0].Address,
		HostName:node.Status.Addresses[1].Address,
		OrginLabels:node.Labels,
		OrginTaints:node.Spec.Taints,
		Capacity:models.NewNodeCapacity(node.Status.Capacity.Cpu().Value(),
			node.Status.Capacity.Memory().Value(),node.Status.Capacity.Pods().Value()),
		Usage:models.NewNodeUsage(this.PodMap.GetNum(node.Name),nodeUsage[0],nodeUsage[1]),
		CreateTime:node.CreationTimestamp.Format("2006-01-02 15:04:05"),
	}
}
//显示所有节点
func(this *NodeService) ListAllNodes() []*models.NodeModel{
    list:=this.NodeMap.ListAll()
    ret:=make([]*models.NodeModel,len(list))
    for i,node:=range list{
		nodeUsage:=helpers.GetNodeUsage(this.Metric,node)
    	ret[i]=&models.NodeModel{
    		Name:node.Name,
			IP:node.Status.Addresses[0].Address,
			HostName:node.Status.Addresses[1].Address,
			Lables:helpers.FilterLables(node.Labels),
			Taints:helpers.FilterTaints(node.Spec.Taints),
			Capacity:models.NewNodeCapacity(node.Status.Capacity.Cpu().Value(),
				node.Status.Capacity.Memory().Value(),node.Status.Capacity.Pods().Value()),
				Usage:models.NewNodeUsage(this.PodMap.GetNum(node.Name),nodeUsage[0],nodeUsage[1]),
			CreateTime:node.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
	}
    return ret
}
