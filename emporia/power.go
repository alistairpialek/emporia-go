package emporia

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// PowerInput is required to get power values.
type PowerInput struct {
	DeviceGID int
	Instant   string
	Scale     string
	Unit      string
}

// PowerOutput models power values returned.
type PowerOutput struct {
	DeviceListUsage struct {
		Devices []struct {
			DeviceGID     int `json:"deviceGid"`
			ChannelUsages []struct {
				Name  string  `json:"name"`
				Usage float64 `json:"usage"`
			} `json:"channelUsages"`
		} `json:"devices"`
	} `json:"deviceListUsages"`
}

type CircuitPower struct {
	Name  string
	Power float64
}

func (e *Emporia) powerEndpoint(p PowerInput) string {
	return fmt.Sprintf("AppAPI?apiMethod=getDeviceListUsages&deviceGids=%d&instant=%s&scale=%s&energyUnit=%s", p.DeviceGID, p.Instant, p.Scale, p.Unit)
}

func (e *Emporia) getPower(token *string) (cp []CircuitPower, err error) {
	if e.Circuits == nil {
		log.Print("No circuits configured to collect power for.")
		return
	}

	input := PowerInput{
		DeviceGID: deviceGID,
		// From 5s ago. Without this, occasionally the API will return back no data.
		// I can only assume this is due to the data not being saved yet, so we give the
		// Emporia API some breathing room.
		Instant: time.Now().UTC().Add(time.Duration(-5) * time.Second).Format(time.RFC3339),
		Scale:   "1S",
		Unit:    "KilowattHours",
	}

	endpoint := e.powerEndpoint(input)
	log.Printf("URL: %s", endpoint)

	resp, err := e.getRequest(token, endpoint)
	if err != nil {
		return nil, err
	}

	log.Printf("Response: %s", *resp)

	var usage PowerOutput
	err = json.Unmarshal([]byte(*resp), &usage)
	if err != nil {
		return nil, err
	}

	for _, v := range usage.DeviceListUsage.Devices[0].ChannelUsages {
		for _, c := range *e.Circuits {
			if v.Name == c.Name {
				circuit := CircuitPower{
					Name:  v.Name,
					Power: v.Usage * 3600,
				}

				log.Printf("Circuit Name: %s", circuit.Name)
				log.Printf("Circuit Power: %.3f", circuit.Power)

				cp = append(cp, circuit)
			}
		}
	}

	return cp, nil
}
