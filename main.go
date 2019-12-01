package main 

import (
	"github.com/randomtask1155/hqserver/device"
	"github.com/randomtask1155/hqserver/monitors/sumppump"
	"github.com/randomtask1155/hqserver/monitors/roku"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"log"
)

var (
	logger         *log.Logger
	devices []*device.Device
)


func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hq/status", statusHandler).Methods("GET")
	return r
}

func main() {
	logger = log.New(os.Stdout, "logger: ", log.Ldate|log.Ltime|log.Lshortfile)
	devices = append(devices, sumppump.SumpPump)
	devices = append(devices, roku.RokuBoxes)

	r := newRouter()
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), r)
	if err != nil {
		panic(err)
	}	
}