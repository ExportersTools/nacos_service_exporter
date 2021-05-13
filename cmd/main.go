package main

import (
	"flag"
	"fmt"
	"github.com/ExportersTools/nacos_service_exporter/v1/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
)

var endPoint string
var nameSpaceId string

func init() {

	flag.StringVar(&endPoint, "endpoint", "http://127.0.0.1:8848", "http://127.0.0.1:8848")
	flag.StringVar(&nameSpaceId, "nameSpaceId", "", "public")

	flag.Parse()
}

func main() {

	endPointEnv := os.Getenv("endPoint")
	if endPointEnv != "" {
		endPoint = endPointEnv
	}

	nameSpaceIdEnv := os.Getenv("nameSpaceId")
	if nameSpaceIdEnv != "" {
		nameSpaceId = nameSpaceIdEnv
	}

	fmt.Println("#===============================================#")
	fmt.Println("endPoint: ", endPoint)
	fmt.Println("nameSpaceId: ", nameSpaceId)
	fmt.Println("#===============================================#")

	exporter.Endpoint = endPoint
	exporter.NameSpaceId = nameSpaceId
	exporter.Interval = 30
	go exporter.ServiceListProm()

	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":11111", nil); err != nil {
		fmt.Println("Listen Service Failed, err: ", err.Error())
	}
}
