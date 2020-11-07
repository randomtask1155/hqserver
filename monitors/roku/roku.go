package roku


import (
	"sync"
	"os"
	//"net/http"
	"encoding/xml"
	"encoding/json"
	"time"
	"fmt"
	"github.com/randomtask1155/hqserver/device"
	roku "github.com/randomtask1155/rokuremote"
)

var(
	RokuBoxes *device.Device
)

// RokuStatus used to store and inform relavent status information about the player
type RokuStatus struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    string `json:"port"`
	Status  string `json:"status"`
}

// RokuRequest used for servicing roku requests
type RokuRequest struct {
	Name string `json:"name"`
}

// RokuPlayer represents a single roku device
type RokuPlayer struct {
	Player roku.Player
	Mutex  sync.Mutex
	Error  error      // if player has some error save it here so we can monitor and recover
	Status RokuStatus // current state of player
}

// Controller main controller for the roku device
type Controller struct {
	Players []RokuPlayer
}

// ActiveApp response struct for roku query
type ActiveApp struct {
	App string
}

// ActiveResponse xml response from roku /query/active-app
/*
<?xml version="1.0" encoding="UTF-8" ?>
<active-app>
	<app id="837" subtype="ndka" type="appl" version="1.0.80000286">YouTube</app>
</active-app>
*/
type ActiveResponse struct {
	App  xml.Name `xml:"active-app"`
	Name string   `xml:"app"`
	ID   string   `xml:"id,attr"`
}

func init() {
	RokuBoxes = &device.Device{
		Controller: &Controller{
			Players: SetupNewRokuPlayers(),
		},
		Info: device.DeviceStatus{
			Name: "Roku Players",
		},
	}
	go RokuBoxes.Controller.Monitor()
}



func(c *Controller) Monitor() {

	RokuBoxes.Info.Events = make([]string,len(c.Players))
	RokuBoxes.UpdateStatus(device.HealthyStatus, "status ok")
	for {		
		for i := range c.Players{
			err := c.Players[i].UpdateStatus()
			if err != nil {
				RokuBoxes.Info.Events[i] = fmt.Sprintf("%s:%s:%s", c.Players[i].Player.NickName, c.Players[i].Player.Address, err)
			} else {
				RokuBoxes.Info.Events[i] = fmt.Sprintf("%s:%s:%s", c.Players[i].Player.NickName, c.Players[i].Player.Address, c.Players[i].Status.Status)
			}
		}
		time.Sleep(60 * time.Second)
	}
}



// SetupNewRokuPlayers updates the global RokuPlayer list
func SetupNewRokuPlayers() []RokuPlayer {
	RokuPlayers := make([]RokuPlayer, 0)
	status := make([]RokuStatus, 0)
	json.Unmarshal([]byte(os.Getenv("ROKU_PLAYERS")), &status)

	for i := range status {
		var err error
		RokuPlayers = append(RokuPlayers, RokuPlayer{})
		RokuPlayers[i].Player, err = roku.ConnectName(status[i].Address, status[i].Name)

		if err != nil {
			fmt.Printf("connecting to roku %s failed: %s\n", status[i].Name, err)
			RokuPlayers[i].Error = err
			RokuPlayers[i].Player = roku.Player{NickName: status[i].Name, Address: status[i].Address}
		}
		RokuPlayers[i].Status = RokuStatus{status[i].Name, status[i].Address, RokuPlayers[i].Player.Port, ""}
	}
	return RokuPlayers
}

// rokuGetHandler returns the roku player device information
/*
func rokuGetHandler(w http.ResponseWriter, r *http.Request) {

	rps := make([]RokuStatus, 0)
	for i := range RokuPlayers {
		err := RokuPlayers[i].UpdateStatus()
		if err != nil {
			RokuPlayers[i].Status.Status = fmt.Sprintf("%s", err)
			logger.Println(err)
		}
		rps = append(rps, RokuPlayers[i].Status)
	}
	b, err := json.Marshal(&rps)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Printf("Failed to marshal RokuPlayers: %s\n", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func rokuHomeHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to read request body: %s", err)))
		logger.Printf("Failed to read request body: %s\n", err)
		return
	}

	rr := RokuRequest{}
	err = json.Unmarshal(b, &rr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Failed to unmarshal request body: %s", err)))
		logger.Printf("Failed to unmarshal request body: %s\n", err)
		return
	}

	for i := range RokuPlayers {
		if RokuPlayers[i].Status.Name == rr.Name {
			RokuPlayers[i].Player.Home()
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte(fmt.Sprintf("{\"status\": \"success\"}")))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(fmt.Sprintf("Could not find roku player %s", rr.Name)))
}*/

// UpdateStatus sets the status for the given roku player
func (rp *RokuPlayer) UpdateStatus() error {
	b, err := rp.Player.Get("/query/active-app")
	if err != nil {
		return fmt.Errorf("failed to fetch roku status %s: %s", rp.Status.Name, err)
	}
	ar := ActiveResponse{}
	err = xml.Unmarshal(b, &ar)
	if err != nil {
		return fmt.Errorf("failed to unmarshal active-app response: %s", err)
	}
	rp.Status.Status = ar.Name
	return nil
}