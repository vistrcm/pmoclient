package pmo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
)

func ReadConfig(relativeConfigFilePath string) Configuration {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configFileName := usr.HomeDir + relativeConfigFilePath
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatalf("error %q happened reading config file %v. raw: %v", err, configFileName, raw)
	}

	config := Configuration{}
	err = json.Unmarshal(raw, &config)
	if err != nil {
		log.Fatalf("something happened during unmarshall config: %q. engineers: %v", err, config)
	}

	return config

}

// helper function to remove duplicates
func RemoveDups(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
