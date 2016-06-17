package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func importConfiguration(configPath string) (configuration, error) {
	fmt.Printf("Importing the configuration information from %v\n", configPath)

	f, err := ioutil.ReadFile(configPath)
	var c configuration

	if err != nil {
		return c, err
	}

	json.Unmarshal(f, &c)

	fmt.Printf("Done. Configuration data: \n %+v \n", c)

	return c, err
}
