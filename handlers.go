package main 

import (
	"github.com/randomtask1155/hqserver/device"
	"net/http"
	"encoding/json"
)

func reportMonitors() []device.DeviceStatus {

	ds := make([]device.DeviceStatus,0)
	for i := range devices {
		ds = append(ds, devices[i].Info)
	}
	return ds
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	d := make([]device.DeviceStatus,len(devices))
	for i := range devices {
		d[i] = devices[i].Info
	}
	b, err := json.Marshal(d)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}