package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"sync"
	"os"
	"runtime/trace"

	"github.com/wallix/awless/logger"
)

const url = "http://sevoneApiUrl/api/v2/"
var now = time.Now()
var endTime = now.Unix() * 1000
var startTime = endTime - 300000
var c = make(chan string)
var wg sync.WaitGroup
var wg2 sync.WaitGroup
var wg3 sync.WaitGroup
//var wg4 sync.WaitGroup

//Device in SevOne API
type Device struct {
	TotalElements int `json:"totalElements"`
	DeviceContent []DeviceContent `json:"content"`
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

//DeviceContent in SevOne API
type DeviceContent struct {
	ID                       int         `json:"id"`
	IsDeleted                bool        `json:"isDeleted"`
	IsNew                    bool        `json:"isNew"`
	Name                     string      `json:"name"`
	AlternateName            string      `json:"alternateName"`
	Description              string      `json:"description"`
	IPAddress                string      `json:"ipAddress"`
	ManualIP                 bool        `json:"manualIP"`
	PeerID                   int         `json:"peerId"`
	PollFrequency            int         `json:"pollFrequency"`
	DateAdded                int64       `json:"dateAdded"`
	LastDiscovery            int64       `json:"lastDiscovery"`
	AllowDelete              bool        `json:"allowDelete"`
	DisablePolling           bool        `json:"disablePolling"`
	DisableConcurrentPolling bool        `json:"disableConcurrentPolling"`
	DisableThresholding      bool        `json:"disableThresholding"`
	Timezone                 string      `json:"timezone"`
	WorkhoursGroupID         int         `json:"workhoursGroupId"`
	NumElements              int         `json:"numElements"`
	PluginInfo               interface{} `json:"pluginInfo"`
	Objects                  interface{} `json:"objects"`
	PluginManagerID          interface{} `json:"pluginManagerId"`
} 

//Object in SevOne API
type Object struct {
	TotalElements int `json:"totalElements"`
	ObjectContent []ObjectContent `json:"content"`
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

//ObjectContent in SevOne API
type ObjectContent struct {
	ID                 int         `json:"id"`
	DeviceID           int         `json:"deviceId"`
	PluginID           int         `json:"pluginId"`
	PluginObjectTypeID int         `json:"pluginObjectTypeId"`
	SubtypeID          int         `json:"subtypeId"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Enabled            string      `json:"enabled"`
	IsEnabled          bool        `json:"isEnabled"`
	IsVisible          bool        `json:"isVisible"`
	IsDeleted          bool        `json:"isDeleted"`
	DateAdded          int64       `json:"dateAdded"`
	Indicators         interface{} `json:"indicators"`
	ExtendedInfo       struct {
		PacketInterval string `json:"packetInterval"`
		Custom         string `json:"custom"`
		IP             string `json:"ip"`
		PacketNumber   string `json:"packetNumber"`
		PacketSize     string `json:"packetSize"`
		DeviceID       string `json:"deviceId"`
		ObjectID       string `json:"objectId"`
	} `json:"extendedInfo"` 
}

//Indicator in SevOne API
type Indicator struct {
	TotalElements int `json:"totalElements"`
	IndicatorContent []IndicatorContent `json:"content"`
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalPages int `json:"totalPages"`
}

//IndicatorContent in SevOne API
type IndicatorContent struct {
	ID                    int         `json:"id"`
	DeviceID              int         `json:"deviceId"`
	ObjectID              int         `json:"objectId"`
	PluginID              int         `json:"pluginId"`
	PluginIndicatorTypeID int         `json:"pluginIndicatorTypeId"`
	Name                  string      `json:"name"`
	Description           string      `json:"description"`
	DataUnits             string      `json:"dataUnits"`
	DisplayUnits          string      `json:"displayUnits"`
	IsEnabled             bool        `json:"isEnabled"`
	IsBaselining          bool        `json:"isBaselining"`
	IsDeleted             bool        `json:"isDeleted"`
	MaxValue              float64     `json:"maxValue"`
	Format                string      `json:"format"`
	LastInvalidationTime  int         `json:"lastInvalidationTime"`
	SyntheticExpression   interface{} `json:"syntheticExpression"`
	EvaluationOrder       int         `json:"evaluationOrder"`
	ExtendedInfo          struct {
		IndicatorID string `json:"indicatorId"`
		OidHigh     string `json:"oidHigh"`
		Oid         string `json:"oid"`
		DeviceID    string `json:"deviceId"`
	} `json:"extendedInfo"`
}

//Metric in SevOne API
type Metric []struct {
	Value float64 `json:"value"`
	Time  int64   `json:"time"`
	Focus int     `json:"focus"`
}

func authToken() string {
	urlget := url + "authentication/signin?nmsLogin=false"

	payload := strings.NewReader("{\n\t\"name\": \"username\",\n\t\"password\": \"password\"\n}")

	req, _ := http.NewRequest("POST", urlget, payload)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error130: ", err)
	}

	if string(body[:]) == "[]" {
		logger.Error("No data returned.")
	}

	var jsonb map[string]string
	err = json.Unmarshal(body, &jsonb)
	if err != nil {
		logger.Error("Error143: ", err)
	}
	token := jsonb["token"]
	//logger.Info("Obtained token")

	return token
}

//function to get list of devices (url, token, cookie, devices)
func getDevices(t string) []byte {
	urlget := url + "devices?size=10000"

	req, _ := http.NewRequest("GET", urlget, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-auth-token", t)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("Error130: ", err)
	}

	if string(body[:]) == "[]" {
		logger.Error("No data returned.")
	}

	return body
}

//function to get list of objects for each device (url, token, cookie, object)
func getObjects(t string, d DeviceContent, c chan string) {
		urlget := fmt.Sprintf("%s%s%d%s", url, "devices/", d.ID, "/objects")

		req, _ := http.NewRequest("GET", urlget, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("x-auth-token", t)

		res, err1 := http.DefaultClient.Do(req)
		if err1 != nil {
			logger.Error(err1)
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		
		byt2 := []byte(body)
		var o Object
		err2 := json.Unmarshal(byt2, &o)
		if err2 != nil {
			fmt.Println("214: ", err2)
			return
		}

		for object := range o.ObjectContent {
			wg2.Add(1)
			go getIndicators(t, d.Name, d.ID, o.ObjectContent[object], c)
			//defer wg2.Done()
			//time.Sleep(1 * time.Second)		
		}
		//wg2.Wait()
		//wg.Done()
}

//function to get list of indicators for each object (url, object, device, objectindicator)
func getIndicators(t string, dn string, d int, o ObjectContent, c chan string) {
		urlget := fmt.Sprintf("%s%s%d%s%d%s", url, "devices/", d, "/objects/", o.ID, "/indicators")

		req, _ := http.NewRequest("GET", urlget, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("x-auth-token", t)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Error(err)
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		byt3 := []byte(body)
		var i Indicator
		err3 := json.Unmarshal(byt3, &i)
		if err3 != nil {
			fmt.Println("244: ", err)
			return
		}


		for indicator := range i.IndicatorContent {
			wg3.Add(1)
			go getIndicatorMetric(t, dn, d, o.Name, o.ID, i.IndicatorContent[indicator], c, startTime, endTime)	
			//defer wg3.Done()
			//time.Sleep(1 * time.Second)
		}	
		wg2.Done()
}

//function to get metric value(s) for each indicator (url, token, cookie, indicator)
func getIndicatorMetric(t string, dn string, d int, on string, o int, i IndicatorContent, c chan string, startTime int64, endTime int64) {
		urlget := fmt.Sprintf("%s%s%d%s%d%s%d%s%d%s%d", url, "devices/", d, "/objects/", o, "/indicators/", i.ID, "/data?startTime=", startTime, "&endTime=", endTime)
		//fmt.Println(urlget)
		req, _ := http.NewRequest("GET", urlget, nil)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("x-auth-token", t)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			logger.Error(err)
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		byt4 := []byte(body)
		var m Metric
		//for range m {
		err4 := json.Unmarshal(byt4, &m)
		if err4 != nil {
			fmt.Println("274: ", err)
			return
		}
		
		
		//fmt.Println(dn, on, i.Name)
		for metric := range m {
			//wg4.Add(1)
			//fmt.Println(fmt.Sprintf("%s.%s.%s: %f", dn, on, i.Name, m[metric].Value))
			//c <- fmt.Sprintf("%s.%s.%s: %f", dn, on, i.Name, m[metric].Value))
			c <- fmt.Sprintf("%s.%s: %f", dn, i.Name, m[metric].Value)
		}
		wg3.Done()
}

func main() {

	trace.Start(os.Stderr)
	defer trace.Stop()
	token := authToken()
	
	msg := getDevices(token)

	//Unmarshal "devices" json string into a list of device IDs (int) as "device"
	byt := []byte(msg)
	var d Device
	err := json.Unmarshal(byt, &d)
	if err != nil {
		fmt.Println("295: ", err)
		return
	}

	//might be good to use goroutines and channels below - research
	//var c chan []byte = make(chan []byte)
	
	defer close(c)

	for device := range d.DeviceContent {
		logger.Info("Pulling device", d.DeviceContent[device].Name)
		wg.Add(1)
		go func() {
			getObjects(token, d.DeviceContent[device], c)
		}()
		Print(c, &wg)
		//defer wg.Done()
	}

	//wg.Wait()
	
	//time.Sleep(90 * time.Second)
	
	//logger.Info("Script finished running.")
	elapsed := time.Since(now)
	fmt.Println(elapsed)
}

//Print prints the contents of channel c
func Print(c <-chan string, wg *sync.WaitGroup) {
	for n := range c { // reads from channel until it's closed
        fmt.Println(n)
    }            
    defer wg.Done()
}
