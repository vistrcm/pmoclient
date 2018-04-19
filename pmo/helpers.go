package pmo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sort"
	"strings"
	"text/tabwriter"
)

// ReadConfig gets information from config file and creates structure `Configuration`
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

// RemoveDuplicates helper function to remove duplicates
func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] {
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

// PrintTable prints table representation of engineers
func PrintTable(engineers []Person, formatString string) {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, ' ', 0)
	// print header
	_, err := fmt.Fprintf(w, formatString,
		"Name",
		"Grade",
		"Profile",
		"Account",
		"Project",
		"Manager",
		"Status")
	if err != nil {
		log.Fatalf("[ERR] railed to print formatString %q. E: %v", formatString, err)
	}
	// iterate over engineers and print only required from config
	filtered := ByLocation(engineers)
	sort.Sort(filtered) // sort by location
	for _, engineer := range filtered {
		_, err := fmt.Fprintf(w, formatString,
			engineer.FullName,
			engineer.Grade,
			engineer.WorkProfile,
			strings.Join(RemoveDuplicates(engineer.Account), ","),
			strings.Join(RemoveDuplicates(engineer.Project), ","),
			engineer.Manager,
			strings.Join(RemoveDuplicates(engineer.AssignmentStatuses()), ","))
		if err != nil {
			log.Fatalf("[ERR] failed on writing %v to tabwriter. E: %v", engineer, err)
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatalln("can not flush tabwriter")
	}

}
