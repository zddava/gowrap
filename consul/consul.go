package consul

import (
	"bytes"
	"context"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/zddava/gowrap/json"
)

type (
	ConsulClient struct {
		Host string
		Port int
	}

	ServiceInstance struct {
		ID                string            `json:"ID"`
		Service           string            `json:"Service,omitempty"`
		Name              string            `json:"Name"`
		Tags              []string          `json:"Tags,omitempty"`
		Address           string            `json:"Address"`
		Port              int               `json:"Port"`
		Meta              map[string]string `json:"Meta,omitempty"`
		EnableTagOverride bool              `json:"EnableTagOverride"`
		Check             `json:"Check,omitempty"`
		Weights           `json:"Weights,omitempty"`
		CurWeight         int `json:"CurWeights,omitempty"`
	}

	Check struct {
		DeregisterCriticalServiceAfter string   `json:"DeregisterCriticalServiceAfter"`
		Args                           []string `json:"Args,omitempty"`
		HTTP                           string   `json:"HTTP"`
		Interval                       string   `json:"Interval,omitempty"`
		TTL                            string   `json:"TTL,omitempty"`
	}

	Weights struct {
		Passing int `json:"Passing"`
		Warning int `json:"Warning"`
	}
)

func NewClient(host string, port int) *ConsulClient {
	return &ConsulClient{
		Host: host,
		Port: port,
	}
}

func (consulClient *ConsulClient) Register(ctx context.Context, serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, weights *Weights) error {

	instance := &ServiceInstance{
		ID:                instanceId,
		Name:              serviceName,
		Address:           instanceHost,
		Port:              instancePort,
		Meta:              meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	if weights != nil {
		instance.Weights = *weights
	} else {
		instance.Weights = Weights{
			Passing: 10,
			Warning: 1,
		}
	}

	byteData, err := json.Marshal(instance)
	if err != nil {
		log.Printf("json format err: %s", err)
		return err
	}

	req, err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/register",
		bytes.NewReader(byteData))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	client.Timeout = time.Second * 2
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("register service err : %s", err)
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("register service http request errCode : %v", resp.StatusCode)
		return fmt.Errorf("register service http request errCode : %v", resp.StatusCode)
	}

	log.Println("register service success")
	return nil
}

func (consulClient *ConsulClient) Deregister(ctx context.Context, instanceId string) error {
	req, err := http.NewRequest("PUT",
		"http://"+consulClient.Host+":"+strconv.Itoa(consulClient.Port)+"/v1/agent/service/deregister/"+instanceId, nil)

	if err != nil {
		log.Printf("req format err: %s", err)
		return err
	}

	client := http.Client{}
	client.Timeout = time.Second * 2

	resp, err := client.Do(req)

	if err != nil {
		log.Printf("deregister service err : %s", err)
		return err
	}

	resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("deresigister service http request errCode : %v", resp.StatusCode)
		return fmt.Errorf("deresigister service http request errCode : %v", resp.StatusCode)
	}

	log.Println("deregister service success")

	return nil
}
