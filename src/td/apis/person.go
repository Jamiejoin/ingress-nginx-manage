package apis

import (
	"github.com/gin-gonic/gin"
	"net/http"
  . "td/searchingress"
    "td/monitor"
)

//func IndexApi(c *gin.Context) {
//	c.String(http.StatusOK, "It works")
//}

func SearchIngressApi(c *gin.Context) {
	productline := c.Query("productline")
	servicename := c.Query("servicename")

	if productline  == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters productline "})
		return
	}

	if servicename == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters Servicename "})
		return
	}

	 res := Search(servicename,productline)

	 c.JSON(200,res)

}

func NamespaceApi(c *gin.Context) {

	res := Listnamespace()

	c.JSON(200,res)

}


func ListIngressApi(c *gin.Context) {
	productline := c.Query("productline")

    if productline  ==""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters Productline"})
		return
	}

	res := Listingress(productline)

	c.JSON(200,res)
}


func ListMonitor(c *gin.Context){
	productline := c.Query("productline")
	servicename := c.Query("servicename")

	if productline  == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters productline "})
		return
	}

	if servicename == ""{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters Servicename "})
		return
	}

	res := monitor.SearchMonitor(servicename,productline)

	c.JSON(200,res)
}


func MonitorNamespace(c *gin.Context){


	res := monitor.MonitorNamespace(Listnamespace())

	c.JSON(200,res)



}
