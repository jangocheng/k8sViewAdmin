package rbac

import (
	corev1 "k8s.io/api/core/v1"
)

//@Service
type SaService struct {
	SaMap *SaMapStruct  `inject:"-"`
}
func NewSaService() *SaService {
	return &SaService{}
}
func(this *SaService) ListSa(ns string) []*corev1.ServiceAccount {
	 //sa:=corev1.ServiceAccount{}

	return this.SaMap.ListAll(ns)
}