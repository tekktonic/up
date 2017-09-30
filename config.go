package main

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
)

type configuration struct {
	Server string `json:"server"`
	Key string `json:"key"`
	Owner string `json:"owner"`
	Max int `json:"max"`
}

var config configuration;

func (c configuration) String() string {
	return "Server: " + c.Server +
	 "\nKey: " + c.Key +
	 "\nOwner: " + c.Owner
}

func readConfig() {
	confstring, err := ioutil.ReadFile("config.json")

	if (err != nil) {
		log.Fatal("Unable to open config file\n" + err.Error())
	}
	
	err = json.Unmarshal(confstring, &config)

	if (err != nil) {
		log.Fatal(err)
	}
	fmt.Println((string)(confstring))
	fmt.Println(config)
	
}
