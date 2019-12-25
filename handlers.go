package main 

import (
	"github.com/randomtask1155/hqserver/device"
	"net/http"
	"encoding/json"
	"net/url"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(r *http.Request) bool {
	val, ok := r.URL.Query()["token"]
	if !ok || len(val[0]) < 1 {
		logger.Println("token not found")
		return false
	}
	tokenEnc := val[0]
	if tokenEnc == "" {
		logger.Println("token not found")
		return false
	}
	token, err := url.QueryUnescape(tokenEnc)
	if err != nil {
		logger.Println(err)
		return false
	}
	password, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		logger.Println(err)
		return false
	}
	err = bcrypt.CompareHashAndPassword([]byte(accessToken), password)
	if err != nil {
		logger.Println(err)
		return false
	}
	return true
}

func reportMonitors() []device.DeviceStatus {

	ds := make([]device.DeviceStatus,0)
	for i := range devices {
		ds = append(ds, devices[i].Info)
	}
	return ds
}

func statusHandler(w http.ResponseWriter, r *http.Request) {

	if !checkPassword(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
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