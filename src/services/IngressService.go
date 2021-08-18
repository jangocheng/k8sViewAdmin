package services

import (
	"context"
	"fmt"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8sapi/src/models"
	"strconv"
	"strings"
)

const (
	OPTION_CROS=iota
	OPTION_LIMIT
	OPTION_REWRITE
)
const (
	OPTOINS_CROS_TAG="nginx.ingress.kubernetes.io/enable-cors"
	OPTIONS_REWRITE_TAG="nginx.ingress.kubernetes.io/rewrite-enable"
)
//@service
type IngressService struct {
     Client *kubernetes.Clientset `inject:"-"`
     IngressMap *IngressMapStruct `inject:"-"`
}
func NewIngressService() *IngressService {
	return &IngressService{}
}
//解析标签
func(this *IngressService) parseAnnotations(annos string) map[string]string{
	replace:=[]string{"\t"," ","\n","\r\n"}
	for _,r:=range replace{
		annos=strings.ReplaceAll(annos,r,"")
	}
	ret:=make(map[string]string)
	list:=strings.Split(annos,";")
	for _,item:=range list{
		annos:=strings.Split(item,":")
		if len(annos)==2{
			ret[annos[0]]=annos[1]
		}
	}
	fmt.Println(ret)
	return ret

}

func(this *IngressService) PostIngress(post *models.IngressPost) error{
	className := "nginx"
	ingressRules := []v1beta1.IngressRule{}
	// 凑 Rule对象
	for _, r := range post.Rules {
		httpRuleValue := &v1beta1.HTTPIngressRuleValue{}
		rulePaths := make([]v1beta1.HTTPIngressPath, 0)
		for _, pathCfg := range r.Paths {
			port,err:=strconv.Atoi(pathCfg.Port)
			if err!=nil{
				return err
			}
			rulePaths = append(rulePaths, v1beta1.HTTPIngressPath{
				Path: pathCfg.Path,
				Backend: v1beta1.IngressBackend{
					ServiceName: pathCfg.SvcName,
					ServicePort: intstr.FromInt(port), //这里需要FromInt
				},
			})
		}
		httpRuleValue.Paths = rulePaths
		rule := v1beta1.IngressRule{
			Host: r.Host,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: httpRuleValue,
			},
		}
		ingressRules = append(ingressRules, rule)
	}

	// 凑 Ingress对象
	ingress := &v1beta1.Ingress{
		TypeMeta:v1.TypeMeta{
			Kind:"Ingress",
			APIVersion:"networking.k8s.io/v1beta1",
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      post.Name,
			Namespace: post.Namespace,
			Annotations:this.parseAnnotations(post.Annotations),
		},
		Spec: v1beta1.IngressSpec{
			IngressClassName: &className,
			Rules:            ingressRules,
		},
	}

	_,err:= this.Client.NetworkingV1beta1().Ingresses(post.Namespace).
		Create(context.Background(),ingress,v1.CreateOptions{})
	return err

}


func(this *IngressService) getIngressOptions(t int,item *v1beta1.Ingress) bool{
	switch t {
	case OPTION_CROS:
		if _,ok:=item.Annotations[OPTOINS_CROS_TAG];ok{
			return true
		}
	case OPTION_REWRITE:
		if _,ok:=item.Annotations[OPTIONS_REWRITE_TAG];ok{
			return true
		}
	}
	return false
}
//前台用于显示Ingress列表
func(this *IngressService) ListIngress(ns string) []*models.IngressModel {
    list:=this.IngressMap.ListAll(ns)
	ret:=make([]*models.IngressModel,len(list))
	for i,item:=range list{
		ret[i]=&models.IngressModel{
			Name:item.Name,
			CreateTime:item.CreationTimestamp.Format("2006-01-02 15:04:05"),
			NameSpace:item.Namespace,
			Host:item.Spec.Rules[0].Host,
			Options:models.IngressOptions{
				IsCros:this.getIngressOptions(OPTION_CROS,item),
				IsRewrite:this.getIngressOptions(OPTION_REWRITE,item),
			},
		}
	}
	return ret
}