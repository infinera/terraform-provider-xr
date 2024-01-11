package xrcm_pf

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ns "terraform-provider-xrcm/internal/service/xrns"

	"github.com/google/martian/v3/log"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL       string
	HTTPClient    *http.Client
	Token         string
	Auth          AuthStruct
	Devicemap     map[string]string
	GetTimeout    time.Duration
	DeleteTimeout time.Duration
	UpdateTimeout time.Duration
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	Token         string `json:"acess_token"`
	Id_token      string `json:"id_token"`
	Refresh_token string `json:"refresh_token"`
	Scope         string `json:"scope"`
	Token_type    string `json:"token_type"`
}

// NewClient -
func NewClient(host, username, password *string) (*Client, error) {
	getTimeout, err := strconv.Atoi(os.Getenv("GET_TIMEOUT"))
	if err != nil {
		getTimeout = 0
	}
	updateTimeout, err := strconv.Atoi(os.Getenv("UPDATE_TIMEOUT"))
	if err != nil {
		updateTimeout = 4
	}
	deleteTimeout, err := strconv.Atoi(os.Getenv("DELETE_TIMEOUT"))
	if err != nil {
		deleteTimeout = 5
	}

	log.Debugf("NewClient: getTimeout = %d, updateTimeout = %d, deleteTimeout = %d", getTimeout, updateTimeout, deleteTimeout)

	c := Client{
		HTTPClient: &http.Client{Timeout: time.Duration(getTimeout) * time.Second},
		// Default Hashicups URL
		HostURL: HostURL,
		Auth: AuthStruct{
			Username: *username,
			Password: *password,
		},
		UpdateTimeout: time.Duration(updateTimeout) * time.Second,
		GetTimeout:    time.Duration(getTimeout) * time.Second,
		DeleteTimeout: time.Duration(deleteTimeout) * time.Second,
	}

	if host != nil {
		c.HostURL = *host
	}

	ar, err := c.SignIn()
	if err != nil {
		return nil, err
	}

	c.Token = ar.Token
	//fmt.Println("ar Token:" + ar.Token)
	//fmt.Println("c Token:" + c.Token)

	log.SetLevel(log.Debug)

	return &c, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", c.Token)

	if req.Method == "GET" {
		c.HTTPClient.Timeout = c.GetTimeout
	} else if req.Method != "DELETE" {
		c.HTTPClient.Timeout = c.UpdateTimeout
	} else {
		c.HTTPClient.Timeout = c.DeleteTimeout
	}

	log.Debugf("doRequest: method = %s, Timeout = %v", req.Method, c.HTTPClient.Timeout)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		log.Debugf("doRequest: Send HTTP Request error %v", err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Debugf("doRequest: Can not read Reponse Body. error %v", err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

// executes commands on specified device;
// devicename	in cfg would be mapped to device id, optional attribute
// deviceid		deviceid associated, optional attribute; either devicename or deviceid is used
func (c *Client) ExecuteDeviceHttpCommand(devicename string, command, commanduri string, commandBody []byte) (result []byte, deviceid string, err error) {

	deviceid, found := c.GetDeviceIdFromName(devicename)
	if !found {
		return nil, devicename, errors.New("device not found : " + devicename)
	}
	log.Debugf("ExecuteDeviceHttpCommand:New HTTP Request %s/api/v1/devices/%s/%s/", c.HostURL, deviceid, commanduri)
	// fmt.Println("deviceid:", deviceid, "command body"+string(commandBody))
	req, err := http.NewRequest(command, fmt.Sprintf("%s/api/v1/devices/%s/%s/", c.HostURL, deviceid, commanduri), bytes.NewBuffer(commandBody))
	log.Debugf("ExecuteDeviceHttpCommand: Create HTTP Request %v", req)
	if err != nil {
		log.Errorf("ExecuteDeviceHttpCommand: Device ID = %s, Create New HTTP Request failed error %v", deviceid, err)
		return nil, deviceid, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		log.Errorf("ExecuteDeviceHttpCommand: Send HTTP NewRequest failed for devide  error= %v", err)
		return nil, deviceid, err
	}
	log.Debugf("ExecuteDeviceHttpCommand: Send HTTP Request SUCCESS. Response = %s", string(body))
	return body, deviceid, err
}

func (c *Client) ExecuteDeviceHttpCommandByID(deviceid string, command, commanduri string, commandBody []byte) (result []byte, err error) {

	log.Debugf("ExecuteDeviceHttpCommand:New HTTP Request %s/api/v1/devices/%s/%s/", c.HostURL, deviceid, commanduri)
	// fmt.Println("deviceid:", deviceid, "command body"+string(commandBody))
	req, err := http.NewRequest(command, fmt.Sprintf("%s/api/v1/devices/%s/%s/", c.HostURL, deviceid, commanduri), bytes.NewBuffer(commandBody))
	log.Debugf("ExecuteDeviceHttpCommand: Create HTTP Request %v", req)
	if err != nil {
		log.Errorf("ExecuteDeviceHttpCommand: Device ID = %s, Create New HTTP Request failed error %v", deviceid, err)
		return nil, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		log.Errorf("ExecuteDeviceHttpCommand: Send HTTP NewRequest failed for devide  error= %v", err)
		return nil, err
	}
	log.Debugf("ExecuteDeviceHttpCommand: Send HTTP Request SUCCESS. Response = %s", string(body))
	return body, err
}

func (c *Client) ExecuteHttpCommand(command, commanduri string, commandBody []byte) (result []byte, err error) {
	// TODO: Remove the API base hardcoding
	log.Debugf("ExecuteHttpCommand:New HTTP Request %s/api/v1/%s", c.HostURL, commanduri)
	req, err := http.NewRequest(command, fmt.Sprintf("%s/api/v1/%s", c.HostURL, commanduri), bytes.NewBuffer(commandBody))
	// fmt.Println(req)
	if err != nil {
		log.Debugf("ExecuteHttpCommand: NewRequest error %v", err, c.HostURL, commanduri)
		return nil, err
	}
	log.Debugf("ExecuteHttpCommand: Send HTTP Request %v", req)

	body, err := c.doRequest(req)
	if err != nil {
		log.Errorf("ExecuteHttpCommand: URL= %s, Send HTTP NewRequest failed error %v", c.HostURL, err)
		return nil, err
	}

	log.Debugf("ExecuteHttpCommand: Send HTTP Request SUCCESS. Response = %s", string(body))
	return body, err
}

func getDevices(body []byte) (dev []map[string]interface{}) {

	log.Debugf("getDevices: ")
	var devices []map[string]interface{}
	dec := json.NewDecoder(strings.NewReader(string(body)))
	for {
		var data map[string]interface{}
		if err := dec.Decode(&data); err == io.EOF {
			break
		} else if err != nil {
			log.Errorf("getDevices: Can't parse the data error" + err.Error())
			break
		}
		// fmt.Println("Device Data: ", data)
		devices = append(devices, data["result"].(map[string]interface{}))
	}
	log.Debugf("getDevices: number of devices = %d", len(devices))
	return devices
}

func (c *Client) DiscoverDevices(deviceMap *map[string]string) (err error) {
	// Get Devices
	// Store devices in map
	// fmt.Println("getting devices")
	log.Debugf("DiscoverDevices")
	body, err := c.ExecuteHttpCommand("GET", "devices", nil)
	// fmt.Println("Device list" + string(body))
	if err != nil {
		log.Errorf("DiscoverDevices: Can't get the devices error" + err.Error())
		return
	}

	devices := getDevices(body)
	c.Devicemap = make(map[string]string)

	for _, v := range devices {
		name, dId := getNameAndId(v)
		if (name != "") && (dId != "") {
			c.Devicemap[name] = dId
		}
	}
	//fmt.Println("Devices Found", c.Devicemap)
	log.Debugf("DiscoverDevices: number of devices = %d", len(c.Devicemap))
	return nil
}

func isOnline(d map[string]interface{}) bool {
	metadata := d["metadata"].(map[string]interface{})
	connection := metadata["connection"].(map[string]interface{})
	status := connection["status"].(string)
	return status == "ONLINE"
}

func getNameAndId(d map[string]interface{}) (string, string) {
	if isOnline(d) {
		return d["name"].(string), d["id"].(string)
	}
	return "", ""
}

func (c *Client) GetDeviceIdFromName(devicename string) (dev string, found bool) {
	dId, ok := c.Devicemap[devicename]
	if !ok {
		// Invoke the XR Naming service if it's enabled via the environment variable
		if ep, nsEnabled := os.LookupEnv("XRCM_NAMING_SERVICE"); nsEnabled {
			nsClient := ns.XrnsClient{Endpoint: ep}
			device, err := nsClient.GetDeviceByName(devicename)
			if err != nil {
				log.Errorf("Error: Failed device lookup - %v\n", err)
				return devicename, false
			}
			c.Devicemap[devicename] = device.GetId()
			return device.GetId(), true
		}

		c.DiscoverDevices(&c.Devicemap)
		dId, ok = c.Devicemap[devicename]
	}

	log.Debugf("GetDeviceIdFromName: devicename = %s, ID = %s", devicename, dId)
	if ok {
		return dId, true // found
	}

	return devicename, false // not found
}
