package models

type ClusterInfo struct {
	EndPoint string  `yaml:"endpoint"`
	CaFile string  `yaml:"cafile"`
	UserCert string `yaml:"usercert"`  //用户证书存放位置
}
type NodesConfig struct {
	Name string
	Ip string
	User string
	Pass string
}
type K8sConfig struct {
	Nodes []*NodesConfig
	ClusterInfo *ClusterInfo `yaml:"cluster-info"`
}
type SysConfig struct {
	K8s *K8sConfig
}
