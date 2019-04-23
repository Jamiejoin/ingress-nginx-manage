package k8sconn

import (
	"k8s.io/client-go/kubernetes"
	"flag"
	"log"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
)

var Clientset *kubernetes.Clientset

func K8sconn(){
	k8sconfig := flag.String("k8sconfig","./k8sconfig","kubernetes config file path")
	flag.Parse()
	config , err := clientcmd.BuildConfigFromFlags("",*k8sconfig)
	if err != nil {
		log.Println(err)
	}
	Clientset , err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("The configuration file is correct")
	}
}
