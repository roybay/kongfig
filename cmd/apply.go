package cmd

import (
	"encoding/json"
	"log"
	"os"

	"github.com/pagerinc/kongfig/kong"
)

func Apply() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c := kong.Config{}
	json.NewDecoder(file).Decode(&c)

	for _, s := range c.Services {
		c.UpdateService(s)
		c.CreateRoutes(s)
		c.GetRoutes(s)
	}
}
