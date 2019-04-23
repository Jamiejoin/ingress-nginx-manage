package monitor

import (
	"k8s.io/api/core/v1"
	"encoding/json"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/prometheus/common/log"
	. "td/k8sconn"
	dto "github.com/prometheus/client_model/go"
	"k8s.io/apimachinery/pkg/labels"
	"github.com/prometheus/prom2json"
	"strings"
	"fmt"
	"reflect"
)

//读取监控数据
func SearchMonitor(Servicename string,Productline string)(res map[string]interface{}){

	//找到pod ip
	res = make(map[string]interface{})
	ingress_config,err := Clientset.ExtensionsV1beta1().Ingresses(Productline).Get(Servicename,metav1.GetOptions{ResourceVersion:"0"})
	if err != nil{
		return res
	}
	namespace := ingress_config.Namespace
	annotations := ingress_config.Annotations["kubectl.kubernetes.io/last-applied-configuration"]
	anmap := make(map[string]interface{})
	err = json.Unmarshal([]byte(annotations), &anmap)
	if err != nil{
		return res
	}
	ascription_nginx := anmap["metadata"].(map[string]interface{})["annotations"].(map[string]interface{})["kubernetes.io/ingress.class"]
	if ascription_nginx == nil{
		return res
	}

	// 通过inress class查找所属nginx
	deploymentList, err := Clientset.AppsV1beta1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil{
		return res
	}
	// 查找对应ip地址
	result := []*prom2json.Family{}
	for _,item := range deploymentList.Items {
		var ingress_pods *v1.PodList
		args := item.Spec.Template.Spec.Containers[0].Args
		deploy_lab := labels.Set(item.Spec.Selector.MatchLabels).AsSelector().String()
		if (len(args) != 0) {
			for _, arg := range args {
				if strings.Contains(arg, ascription_nginx.(string)) {
					ingress_pods, err = Clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: deploy_lab})
					if err != nil {
						return res
					}
				}
			}
			if ingress_pods !=nil{

				for _, pod := range ingress_pods.Items {
				if len(pod.Status.HostIP) != 0 && len(pod.Status.PodIP) != 0 {
					//转换owl数据到json
					url := "http://" + pod.Status.PodIP + ":10254/metrics"
					mfChan := make(chan *dto.MetricFamily, 1024)

					go func() {
						err := prom2json.FetchMetricFamilies(url, mfChan)
						if err != nil {
							log.Errorln(err)
							defer close(mfChan)

						}
					}()


					for mf := range mfChan {
						result = append(result, prom2json.NewFamily(mf))
						}

					//jsonText, err := json.Marshal(result)
					//if err != nil {
					//	log.Fatalln("error marshaling JSON:", err)
					//}
					//fmt.Println(string(jsonText))
					}
					//循环读取切片


							}
						}

					}

				}
				fmt.Println(result[30])
			return res
		}





func  MonitorNamespace(namespace_list []string)(res interface{}){
	//监控Namespace剩余资源
	type Tags struct {
		Namespace string `json:"namespace"`

	}
	type Data struct {
		Metric string `json:"metric"`
		Data_type string `json:"data_type"`
		Value int64 `json:"value"`
		Tags Tags `json:"tags"`

	}

	var NamespaceUsed []Data
	for _,v := range namespace_list {
		fmt.Println(v)

		resource, err := Clientset.CoreV1().ResourceQuotas(v).Get(v, metav1.GetOptions{ResourceVersion: "0"})
		if err != nil {
			continue
		}
		//获取cpu
		cpu := resource.Status.Used["limits.cpu"]
		usedcpu := reflect.ValueOf(&cpu).Elem().FieldByName("i").FieldByName("value").Int()
		var cpudata Data
		cpudata.Metric = "k8s.namespace_cpu_used"
		cpudata.Data_type = "gauge"
		cpudata.Value = usedcpu
		cpudata.Tags.Namespace = v

		//获取内存
		memory := resource.Status.Used["limits.memory"]
		usedmemory := reflect.ValueOf(&memory).Elem().FieldByName("i").FieldByName("value").Int()
		var memorydata Data
		memorydata.Metric = "k8s.namespace_memory_used"
		memorydata.Data_type = "gauge"
		memorydata.Value = usedmemory
		memorydata.Tags.Namespace = v

		//定义结构体数组转json
		NamespaceUsed = append(NamespaceUsed, cpudata)
		NamespaceUsed = append(NamespaceUsed, memorydata)
		//

	}
	return NamespaceUsed

}

