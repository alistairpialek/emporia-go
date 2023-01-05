package emporia

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// EnergyInput is required to get energy values.
type EnergyInput struct {
	DeviceGID int
	Channel   int
	Start     string
	End       string
	Scale     string
	Unit      string
}

// EnergyOutput models energy values returned.
type EnergyOutput struct {
	UsageList []float64 `json:"usageList"`
}

type CircuitEnergy struct {
	Name   string
	Energy float64
}

func (e *Emporia) energyEndpoint(en EnergyInput) string {
	return fmt.Sprintf("AppAPI?apiMethod=getChartUsage&deviceGid=%d&channel=%d&start=%s&end=%s&scale=%s&energyUnit=%s", en.DeviceGID, en.Channel, en.Start, en.End, en.Scale, en.Unit)
}

func (e *Emporia) GetDayEnergy(token *string) (ce []CircuitEnergy, err error) {
	if e.Circuits == nil {
		log.Print("No circuits configured to collect energy for.")
		return nil, nil
	}

	for _, c := range e.Circuits {
		input := EnergyInput{
			DeviceGID: e.DeviceGID,
			Channel:   c.Channel,
			Start:     time.Now().UTC().Format(time.RFC3339),
			End:       time.Now().UTC().Format(time.RFC3339),
			Scale:     "1D",
			Unit:      "KilowattHours",
		}

		endpoint := e.energyEndpoint(input)
		log.Printf("URL: %s", endpoint)

		resp, err := e.getRequest(token, endpoint)
		if err != nil {
			return nil, err
		}

		log.Printf("Response: %s", *resp)

		var usage EnergyOutput
		err = json.Unmarshal([]byte(*resp), &usage)
		if err != nil {
			return nil, err
		}

		//niceEnergy := fmt.Sprintf("%.2f", usage.UsageList[0])
		circuit := CircuitEnergy{
			Name:   c.Name,
			Energy: usage.UsageList[0],
		}

		log.Printf("Circuit Name: %s", circuit.Name)
		log.Printf("Circuit Energy: %f", circuit.Energy)

		ce = append(ce, circuit)

	}

	return ce, nil
}
