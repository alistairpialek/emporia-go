package main

import (
	"log"

	"github.com/alistairpialek/emporia-go/emporia"
)

func main() {
	a := emporia.Emporia{
		RootTempDir: "/tmp",
		Timezone:    "Country/City",
		Username:    "me@email.com",
		Password:    "",
		ClientID:    "",
		UserPoolID:  "",
		DeviceGID:   123,
		Circuits: []emporia.Circuit{
			{
				Name:    "Circuit Name",
				Channel: 1,
			},
		},
	}

	log.Print("Starting Emporia Report")

	token, err := a.GetLogin()
	if err != nil {
		log.Panicf("login: %s", err)
	}

	ce, err := a.GetDayEnergy(token)
	if err != nil {
		log.Panicf("get energy: %s", err)
	}

	for _, circuit := range ce {
		log.Printf("Name: %s", circuit.Name)
		log.Printf("Energy: %f", circuit.Energy)
	}

	cp, err := a.GetPower(token)
	if err != nil {
		log.Printf("get power: %s", err)
	}

	for _, circuit := range cp {
		log.Printf("Name: %s", circuit.Name)
		log.Printf("Power: %f", circuit.Power)
	}

	err = a.GetCustomerDevices(token)
	if err != nil {
		log.Printf("get customer devices: %s", err)
	}
}
