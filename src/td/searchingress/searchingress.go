package searchingress

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"encoding/json"
	"strings"
	. "td/k8sconn"
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
)

//获取namespace
func Listnamespace()(namespace_list []string){
	namespace_rs,err := Clientset.CoreV1().Namespaces().List(metav1.ListOptions{ResourceVersion:"0"})

	if err != nil{
		return namespace_list
	}
	for _,item := range namespace_rs.Items{
		name := item.Name
		if name != "kube-system" && name != "kube-public" && name !="spark"{
			namespace_list =append(namespace_list, name)
		}
	}

	return  namespace_list

}


//获取namespace下的所有ingress
func Listingress(productline string)(res map[string]interface{}){
	res = make(map[string]interface{})
	ingress_list, err:= Clientset.ExtensionsV1beta1().Ingresses(productline).List(metav1.ListOptions{})
	if err != nil{
		res["error"] = "Error getting insress list"
		return res
	}

	for _,item := range ingress_list.Items{
		res[item.Name] = item.Spec.Rules[0].Host
	}
	return res

}


//查找具体某个ingress class信息
func Search(Servicename string,Productline string)(res map[string]interface{}){

	res = make(map[string]interface{})
	ingress_config,err := Clientset.ExtensionsV1beta1().Ingresses(Productline).Get(Servicename,metav1.GetOptions{ResourceVersion:"0"})
	if err != nil{
		res["error"] = "Service  not found " + Servicename
		return res
	}
	namespace := ingress_config.Namespace
	annotations := ingress_config.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
	anmap := make(map[string]interface{})
	err = json.Unmarshal([]byte(annotations), &anmap)
	if err != nil{
		res["error"] = "annotations not found ngress_config.Annotations "
		return res
	}
	ascription_nginx := anmap["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubernetes.io/ingress.class"]
	if ascription_nginx == nil{
		res["error"] = "annotations not found ascription_nginx"
		//return res
	}

    // 通过inress class查找所属nginx
	deploymentList, err := Clientset.AppsV1beta1().Deployments(namespace).List(metav1.ListOptions{ResourceVersion:"0"})
	if err != nil{
		res["error"] = "deployment not found"
		return res
	}
	for _,item := range deploymentList.Items{
		args :=item.Spec.Template.Spec.Containers[0].Args
		// 获取deployment的label
		deploy_lab := labels.Set(item.Spec.Selector.MatchLabels).AsSelector().String()
		if (len(args) != 0){
			for _,arg := range args{
				if strings.Contains(arg,ascription_nginx.(string)) {
					//获取nginx数量
					ingress_number := make(map[string]int32)
					ingress_number["avaliable"] = item.Status.AvailableReplicas
					ingress_number["total"] = *item.Spec.Replicas

					//获取具体ingress pod名称和IP地址
					//根据label获取pod
					ingress_pods, err := Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector:deploy_lab,ResourceVersion: "0"})
					if err != nil {
						res["ingress_podinfo"] = ""
					}
					//比对ingress的Pod名称和deployments名称是否相同，如果一致获取pod和node的IP
					podinfo := make(map[string]map[string]string)
					for _, pod := range ingress_pods.Items {
							if len(pod.Status.HostIP) != 0 && len(pod.Status.PodIP) != 0 {
								podinfo[pod.Name] = map[string]string{"PodIp":pod.Status.PodIP,"NodeIp":pod.Status.HostIP}
							}

					}
					res["ingress_podinfo"] = podinfo
					res["ingress_number"] = ingress_number

					//获取backendservice pod 信息
					backend_servicename := ingress_config.Spec.Rules[0].HTTP.Paths[0].Backend.ServiceName
					res["backend_servicename"] = backend_servicename
					service, err := Clientset.CoreV1().Services(namespace).Get(backend_servicename, metav1.GetOptions{ResourceVersion: "0"})
					// labels.Parser
					if err != nil{
						res["backend_podinfo"] = ""
					}
					set := labels.Set(service.Spec.Selector)
					setlab := set.AsSelector().String()

					if pods, err := Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: setlab,ResourceVersion: "0"},); err != nil {
						fmt.Println("no")
					} else {
						backend_podinfo := make(map[string]map[string]string)
						//去重
						for _, v := range pods.Items {
							if _, ok := podinfo[v.Name]; !ok{
								if len(v.Status.HostIP) !=0 && len(v.Status.PodIP) !=0 {
									backend_podinfo[v.Name] =  map[string]string {"PodIp":v.Status.PodIP,"NodeIp":v.Status.HostIP}
									}
								}
							}
							res["backend_podinfo"] = backend_podinfo
						}
					}

				 if strings.Contains(arg,"--configmap=") {
					 //通过ingress class 获取对应的configmap
					 kv := strings.Split(arg, "/")
					 if len(kv) < 1 {
						 res["nginx-config"] = ""
					 } else {
						 configmap_name := kv[1]
						 ingress_configmap, err := Clientset.CoreV1().ConfigMaps(namespace).Get(configmap_name, metav1.GetOptions{})
						 if err != nil {
							 res["nginx-config"] = ""
						 }
						 res["nginx-config"] = ingress_configmap.Data

					 }
				 }
			}
		}
	}

	//获取backendservice pod 信息

	domain := ingress_config.Spec.Rules[0].Host
	res["domain"] = domain
	res["servicename"] = Servicename


	return res
}
