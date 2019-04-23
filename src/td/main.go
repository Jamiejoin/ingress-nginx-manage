package main



import (
	"td/k8sconn"
)



func main() {
	k8sconn.K8sconn()
	router := InitRouter()

	router.Run(":8000")
}



