package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	//"runtime/trace"
	//"os"

	"github.com/wallix/awless/logger"
)

const url = "http://sevoneURL/api/v2/"

//Device in SevOne API
type Device struct {
	TotalElements int             `json:"totalElements"`
	DeviceContent []DeviceContent `json:"content"`
	PageNumber    int             `json:"pageNumber"`
	PageSize      int             `json:"pageSize"`
	TotalPages    int             `json:"totalPages"`
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
	TotalElements int             `json:"totalElements"`
	ObjectContent []ObjectContent `json:"content"`
	PageNumber    int             `json:"pageNumber"`
	PageSize      int             `json:"pageSize"`
	TotalPages    int             `json:"totalPages"`
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
	TotalElements    int                `json:"totalElements"`
	IndicatorContent []IndicatorContent `json:"content"`
	PageNumber       int                `json:"pageNumber"`
	PageSize         int                `json:"pageSize"`
	TotalPages       int                `json:"totalPages"`
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
type Metric struct {
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
	logger.Info("Obtained token")
	fmt.Println(token)

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
		logger.Error("getDevices1: ", err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("getDevices2: ", err)
	}

	if string(body[:]) == "[]" {
		logger.Error("getDevices3: No data returned.")
	}

	return body
}

//function to get list of objects for each device (url, token, cookie, object)
func getObjects(t string, d int) []byte {
	urlget := fmt.Sprintf("%s%s%d%s", url, "devices/", d, "/objects?size=10000")

	req, _ := http.NewRequest("GET", urlget, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-auth-token", t)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("getObjects1: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

//function to get list of indicators for each object (url, object, device, objectindicator)
func getIndicators(t string, d int, o int) []byte {
	urlget := fmt.Sprintf("%s%s%d%s%d%s", url, "devices/", d, "/objects/", o, "/indicators?size=10000")

	req, _ := http.NewRequest("GET", urlget, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-auth-token", t)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("getIndicators1: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

//function to get metric value(s) for each indicator (url, token, cookie, indicator)
func getIndicatorMetric(t string, d int, o int, i int, startTime int64, endTime int64) []byte {
	urlget := fmt.Sprintf("%s%s%d%s%d%s%d%s%d%s%d", url, "devices/", d, "/objects/", o, "/indicators/", i, "/data?startTime=", startTime, "&endTime=", endTime)

	req, _ := http.NewRequest("GET", urlget, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-auth-token", t)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("getIndicatorMetric1: ", err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func main() {

	//trace.Start(os.Stderr)
	//defer trace.Stop()

	now := time.Now()
	endTime := now.Unix() * 1000
	startTime := endTime - 300000

	//get AuthToken to use for the API GET calls that follow
	token := authToken()

	//get list of Devices
	msg := getDevices(token)

	//Unmarshal "devices" json string into a list of device IDs (int) as "device"
	byt := []byte(msg)
	var d Device
	err := json.Unmarshal(byt, &d)
	if err != nil {
		fmt.Println(err)
		return
	}

	//sem represents semaphore so we can isolate goroutines
	totalDevices := len(d.DeviceContent)
	concurrency := 4
	semDev := make(chan bool, concurrency)

	fmt.Println("{")

	for device := 0; device < totalDevices; device++ {
		semDev <- true

		//create "localDevice" variable for use within goroutine so "device" doesn't change value
		localDevice := device
		go func(current int) {
			msg2 := getObjects(token, d.DeviceContent[localDevice].ID)
			byt2 := []byte(msg2)
			var o Object
			err := json.Unmarshal(byt2, &o)
			if err != nil {
				fmt.Println(err)
				return
			}

			totalObjects := len(o.ObjectContent)
			concurrency2 := 8
			semObj := make(chan bool, concurrency2)

			for object := 0; object < totalObjects; object++ {
				semObj <- true

				localObject := object
				go func(current int) {
					msg3 := getIndicators(token, d.DeviceContent[localDevice].ID, o.ObjectContent[localObject].ID)
					byt3 := []byte(msg3)
					var i Indicator
					err := json.Unmarshal(byt3, &i)
					if err != nil {
						fmt.Println(err)
						return
					}

					totalIndicators := len(i.IndicatorContent)
					concurrency3 := 10
					semInd := make(chan bool, concurrency3)

					for indicator := 0; indicator < totalIndicators; indicator++ {
						semInd <- true

						localIndicator := indicator
						go func(current int) {
							msg4 := getIndicatorMetric(token, d.DeviceContent[localDevice].ID, o.ObjectContent[localObject].ID, i.IndicatorContent[localIndicator].ID, startTime, endTime)
							byt4 := []byte(msg4)
							var m []Metric
							err := json.Unmarshal(byt4, &m)
							if err != nil {
								fmt.Println(err)
								return
							}
							for metric := range m {
								//fmt.Println(fmt.Sprintf("%s.%s.%s: %f", d.DeviceContent[localDevice].Name, o.ObjectContent[localObject].Name, i.IndicatorContent[localIndicator].Name, m[metric].Value))
								fmt.Println(fmt.Sprintf("  \"%s|ST[device:%s,object:%s\": "+"\"%f\",", i.IndicatorContent[localIndicator].Name, d.DeviceContent[localDevice].Name, o.ObjectContent[localObject].Name, m[metric].Value))
							}
							<-semInd
						}(indicator)
					}
					for i := 0; i < cap(semInd); i++ {
						semInd <- true
					}
					<-semObj
				}(object)
			}
			for i := 0; i < cap(semObj); i++ {
				semObj <- true
			}
			<-semDev
		}(device)
	}
	for i := 0; i < cap(semDev); i++ {
		semDev <- true
	}
	fmt.Println("}")
	elapsed := time.Since(now)
	logger.Info("Script finished in ", elapsed)
}
