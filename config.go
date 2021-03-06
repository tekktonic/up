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
	TimelineSize int `json:"timelinesize"`
	Max int `json:"max"`
	DbFile string `json:"dbfile"`
	Port string `json:"port"`
}

var config configuration;

func (c configuration) String() string {
	return "Server: " + c.Server +
	 "\nKey: " + c.Key +
	 "\nOwner: " + c.Owner
}

func readConfig(file string) {
	confstring, err := ioutil.ReadFile(file)

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
