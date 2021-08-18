package helpers

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func HandleCommand(ns,pod,container string, client *kubernetes.Clientset,
							config *rest.Config,command []string) remotecommand.Executor{
	option := &v1.PodExecOptions{
		Container:container,
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:true,
	}
	//"nginx-7875f55f56-5qqbr"
	req:=client.CoreV1().RESTClient().Post().Resource("pods").
		Namespace(ns).
		Name(pod).
		SubResource("exec").
		Param("color","false").
		VersionedParams(
			option,
			scheme.ParameterCodec,
		)

	exec,err:=remotecommand.NewSPDYExecutor(config,"POST",
		req.URL())
	if err!=nil{
		panic(err)
	}
	return exec
}
