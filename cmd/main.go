package main

import (
	"flag"
	"fmt"
	"github.com/ExportersTools/nacos_service_exporter/v1/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var endPoint string
var nameSpaceId string

func init() {

	flag.StringVar(&endPoint, "endpoint", "http://127.0.0.1:8848", "http://127.0.0.1:8848")
	flag.StringVar(&nameSpaceId, "nameSpaceId", "", "public")

	flag.Parse()
}

func main() {
	//exporter.Endpoint = "http://10.4.35.12:8848"
	//exporter.NameSpaceId = "b8293f1a-da23-4561-b6ae-920d5d662e5f"

	fmt.Println("#===============================================#")
	fmt.Println("endPoint: ", endPoint)
	fmt.Println("nameSpaceId: ", nameSpaceId)
	fmt.Println("#===============================================#")
	exporter.Endpoint = endPoint
	exporter.NameSpaceId = nameSpaceId
	exporter.Interval = 5
	go exporter.ServiceListProm()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":11111", nil)
}
