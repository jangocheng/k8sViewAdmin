package models

import "k8s.io/api/core/v1"

//使用量 包含cpu 内存 pods
type NodeUsage struct {
	Pods int
	Cpu float64
	Memory float64
}
func NewNodeUsage(pods int, cpu float64, memory float64) *NodeUsage {
	return &NodeUsage{Pods: pods, Cpu: cpu, Memory: memory}
}


//func NewNodeUsage(pods int) *NodeUsage {
//	return &NodeUsage{Pods: pods}
//}

//容量
type NodeCapacity struct {
	Cpu int64
	Memory int64
	Pods int64
}

func NewNodeCapacity(cpu int64, memory int64, pods int64) *NodeCapacity {
	return &NodeCapacity{Cpu: cpu, Memory: memory, Pods: pods}
}
//保存用
type PostNodeModel struct {
	Name string
	OrginLabels map[string]string  //原始标签 ---->前端 是一个对象
	OrginTaints []v1.Taint //原始污点
}
//节点模型
type NodeModel struct {
	Name string
	IP string
	HostName string

	OrginLabels map[string]string  //原始标签
	OrginTaints []v1.Taint //原始污点
	Lables []string //标签  列表展现
	Taints []string //污点   列表展现
	Capacity *NodeCapacity  //容量 包含了cpu 内存和pods数量
	Usage *NodeUsage //资源 使用情况
	CreateTime string

}
