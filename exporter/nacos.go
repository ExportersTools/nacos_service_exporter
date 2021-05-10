package exporter

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	// 服务实例数量 Metrics
	ServiceInstanceCountMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "serviceInstanceCount",
		Help: "Nacos Service Instance Count",
	},
		[]string{"service", "nameSpaceId"},
	)
	// 是否正常访问 nacos 接口 Metrics
	UpMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Nacos Service Exporter Up",
	},
	)
	// 所有服务数量 Metrics
	AllServiceInstanceCountMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "allServiceInstanceCount",
		Help:        "All Nacos Service Instance Count",
		ConstLabels: map[string]string{"nameSpaceId": NameSpaceId},
	},
	)
	// 所有服务数量 Metrics
	AllServiceCountMetric = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "allServiceCount",
		Help:        "All Nacos Service Count",
		ConstLabels: map[string]string{"nameSpaceId": NameSpaceId},
	},
	)
	//
	EndpointServiceCountMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "endpointServiceCount",
		Help: "Nacos endpoint Service Count",
	},
		[]string{"endpoint", "nameSpaceId"},
	)

	// 所有服务实例数量
	allServiceInstanceCount = 0
	// 间隔多久收集一次数据
	Interval = 30
	// nacos serviceList API
	serviceListUri string = "/nacos/v1/ns/service/list"
	// nacos instanceList API
	instanceListUri string = "/nacos/v1/ns/instance/list"
	// nacos endpoint
	Endpoint string = "http://127.0.0.1"
	// nacos nameSpaceId
	NameSpaceId string

	EndpointServiceCount = map[string]int{}
)

type ServiceListResponse struct {
	Count int      `json:"count"`
	Doms  []string `json:"doms"`
}

type InstanceListResponse struct {
	Dom             string                      `json:"dom"`
	CacheMillis     int64                       `json:"cacheMillis"`
	UseSpecifiedURL bool                        `json:"useSpecifiedURL"`
	Checksum        string                      `json:"checksum"`
	LastRefTime     int64                       `json:"lastRefTime"`
	Env             string                      `json:"env"`
	Clusters        string                      `json:"clusters"`
	Hosts           []InstanceListResponseHosts `json:"hosts"`
}
type InstanceListResponseHosts struct {
	Valid      bool    `json:"valid"`
	Marked     bool    `json:"marked"`
	InstanceId string  `json:"instanceId"`
	Port       int64   `json:"port"`
	Ip         string  `json:"ip"`
	Weight     float32 `json:"weight"`
	Healthy    bool    `json:"healthy"`
}

// nacos 是否启动
func setPromUp(upStatus int) {
	UpMetric.Set(float64(upStatus))
}

// nacos service count
func setPromServiceInstanceCount(serviceName string, serviceCount int) {
	ServiceInstanceCountMetric.With(prometheus.Labels{"service": serviceName, "nameSpaceId": NameSpaceId}).Set(float64(serviceCount))
}

func ServiceListProm() {
	//
	ticker := time.NewTicker(time.Duration(Interval) * time.Second)

	for {
		ServiceList()
		fmt.Println("Now: ", time.Now().Unix())
		//fmt.Println(EndpointServiceCount)
		initVar()
		<-ticker.C
	}
}

func initVar() {
	allServiceInstanceCount = 0
	for k, _ := range EndpointServiceCount {
		delete(EndpointServiceCount, k)
	}
}

// 列出 Service 列表
func ServiceList() {
	var pageNo int = 1
	var pageSize int = 20
	var pageCount int = 0

	for {
		var serviceListResponse ServiceListResponse
		requestUrl := Endpoint + serviceListUri + "?pageNo=" + strconv.Itoa(pageNo) + "&pageSize=" + strconv.Itoa(pageSize) + "&namespaceId=" + NameSpaceId

		resp, err := http.Get(requestUrl)
		if err != nil {
			fmt.Println("http request failed, err: ", err.Error())
			setPromUp(0)
			break
		}

		respBody := resp.Body
		respBodyByte, err := ioutil.ReadAll(respBody)
		resp.Body.Close()

		if err != nil {
			fmt.Println("read respBody failed, err: ", err.Error())
			setPromUp(0)
			break
		}

		err = json.Unmarshal(respBodyByte, &serviceListResponse)
		if err != nil {
			fmt.Println("json unmarshal failed, err: ", err.Error())
			setPromUp(0)
			break
		}

		setPromUp(1)
		AllServiceCountMetric.Set(float64(serviceListResponse.Count))

		for _, service := range serviceListResponse.Doms {
			InstanceList(service)
		}

		pageCount += len(serviceListResponse.Doms)
		if pageCount >= serviceListResponse.Count {
			break
		}
		pageNo += 1
		//break
	}

	AllServiceInstanceCountMetric.Set(float64(allServiceInstanceCount))

	for endpoint, count := range EndpointServiceCount {
		EndpointServiceCountMetric.With(prometheus.Labels{"endpoint": endpoint, "nameSpaceId": NameSpaceId}).Set(float64(count))
	}
}

// 列出 Service 中 Instance 列表
func InstanceList(serviceName string) {
	var instanceListResponse InstanceListResponse
	requestUrl := Endpoint + instanceListUri + "?serviceName=" + serviceName + "&namespaceId=" + NameSpaceId

	resp, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("request service instance failed, err: ", err.Error())
		setPromServiceInstanceCount(serviceName, 0)
		return
	}

	defer resp.Body.Close()

	respBody := resp.Body

	respBodyByte, err := ioutil.ReadAll(respBody)
	if err != nil {
		fmt.Println("read service instance responseBody failed, err: ", err.Error())
		setPromServiceInstanceCount(serviceName, 0)
		return
	}

	err = json.Unmarshal(respBodyByte, &instanceListResponse)
	if err != nil {
		fmt.Println("json unmarshal service instance failed, err: ", err.Error())
		setPromServiceInstanceCount(serviceName, 0)
		return
	}

	for _, h := range instanceListResponse.Hosts {
		EndpointServiceCount[h.Ip] += 1
	}

	//serviceCount := len(instanceListResponse.Hosts)
	var serviceCount int
	for _, hosts := range instanceListResponse.Hosts {
		if hosts.Valid {
			serviceCount += 1
		}
	}
	setPromServiceInstanceCount(serviceName, serviceCount)

	allServiceInstanceCount += serviceCount
}
