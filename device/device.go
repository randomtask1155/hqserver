package device 

import(
	"time"
	"fmt"
)

// MaxEvent history defaults to 10
var (
	MaxEvents = 10
	UnHealthyStatus = "unhealthy"
	HealthyStatus = "healthy"
)

type Device struct {
	Controller DeviceController
	Info DeviceStatus
	EventTrack int
}

type DeviceController interface {
	Monitor()
}

type DeviceStatus struct {
	Status string `json:"status"`
	Name string `json:"name"`
	Events []string `json:"events"` 
}

func (d *Device) AddEvent(event string) {
	event = fmt.Sprintf("%s: %s", time.Now(), event)
	if len(d.Info.Events) <= 0 {
		d.Info.Events = make([]string,MaxEvents)
		d.EventTrack = 0
	}
	if d.EventTrack > (MaxEvents-1) { 
		d.EventTrack = 0
	}
	d.Info.Events[d.EventTrack] = event
	d.EventTrack += 1
}

func (d *Device) UpdateStatus(health, event string){
	d.Info.Status = health 
	if event != "" {
		d.AddEvent(event)
	}
}