package exporter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var serviceListUri string = "/nacos/v1/ns/service/list"
var instanceListUri string = "/nacos/v1/ns/instance/list"
var endpoint string = "http://127.0.0.1"

func ServiceList(pageNo int, pageSize int) {
	requestUrl := endpoint + serviceListUri + "?pageNo=" + strconv.Itoa(pageNo) + "&pageSize=" + strconv.Itoa(pageSize)
	fmt.Println("requestUrl: ", requestUrl)

	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	defer resp.Body.Close()

	respBody := resp.Body

	respBodyByte, err := ioutil.ReadAll(respBody)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	fmt.Println(string(respBodyByte))

	//go InstanceList("")
}

func InstanceList(serviceName string) {
	requestUrl := endpoint + instanceListUri + "?serviceName=" + serviceName

	fmt.Println("requestUrl: ", requestUrl)
}
