package sumppump 

import (
	"github.com/randomtask1155/hqserver/device"
	"time"
	"net/http"
	"fmt"
	"os"
	"io/ioutil"
)

var SumpPump *device.Device

func init() {
	SumpPump = &device.Device{
		Controller: &Controller{
			Interval: 900, // seconds
			Address: os.Getenv("SUMP_ADDRESS"),
		},
		Info: device.DeviceStatus{
			Name: "Sump Pump",
		},
	}
	go SumpPump.Controller.Monitor()
}

type Controller struct {
	Interval int64
	Address string
}

func(c *Controller) Monitor() {

	tr := &http.Transport{
		MaxIdleConns:       -1,
		IdleConnTimeout:    1 * time.Second,
		}
	client := &http.Client{Transport: tr,
			Timeout: 3 * time.Second,
		}
	
	SumpPump.UpdateStatus(device.UnHealthyStatus, "initializing")
	for {
		time.Sleep(time.Duration(c.Interval) * time.Second)
		resp, err := client.Get(c.Address)
		if err != nil {
			SumpPump.UpdateStatus(device.UnHealthyStatus, err.Error())
			continue
		}
		defer resp.Body.Close()
		ioutil.ReadAll(resp.Body) // make sure we read response even though we don't care about one to keep golang happy
		if resp.StatusCode != 200 {
			SumpPump.UpdateStatus(device.UnHealthyStatus, fmt.Sprintf("Sump Pump returned response code %d", resp.StatusCode))
		} else if resp.StatusCode == 200 {
			SumpPump.UpdateStatus(device.HealthyStatus, "Sump Pump became Healthy")
		}
	}
}