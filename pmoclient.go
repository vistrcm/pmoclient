// Package pmoclient implements a simple client for GD's PMO tool.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/vistrcm/pmoclient/pmo"
)

// formatString for print results
const formatString = "%s\t%s\t%s\t%s\t%s\t%s\t%v\n"
const relativeConfigFilePath = "/.config/pmoclient.json"

// initialize http client
var cookieJar, _ = cookiejar.New(nil)
var client = &http.Client{
	Timeout: time.Minute * 1,
	Jar:     cookieJar,

	// do not follow redirects
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

var config pmo.Configuration

// login to PMO
func login() {
	resp, err := client.PostForm(config.LoginURL, url.Values{"j_username": {config.Username}, "j_password": {config.Password}})
	if err != nil {
		log.Fatalf("error during login: %v, resp: %v", err, resp)
	}

	// print output here for debug
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
	}
}

//send request to url. Handle errors somehow.
func request(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error happened %v, req: %v", err, req)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error on sending request: %v, resp: %v", err, resp)
	}
	return resp
}

// get list of engineers by sending request to
func engineers() []pmo.Person {
	resp := request(config.PeopleListURL)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Error when closing: %v\n", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
	}
	var peopleResponse = new(pmo.APIResponse)
	err = json.Unmarshal(body, &peopleResponse)
	if err != nil {
		log.Fatalf("something happened during unmarshall: %v. body: %v\n", err, body)
	}
	return peopleResponse.Rows
}

func filterEngineers() []pmo.Person {
	// initialize return value
	filteredEngineers := make([]pmo.Person, 0)

	// initialize temporary map for filtering
	filterMap := make(map[string]bool)
	for _, u := range config.FilterUsers {
		filterMap[strings.ToLower(u)] = true
	}

	for _, val := range engineers() {
		if filterMap[strings.ToLower(val.EmployeeID)] {
			filteredEngineers = append(filteredEngineers, val)
		}
	}
	return filteredEngineers
}

// function to print table representation of engineers
func printTable() {
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
	filtered := pmo.ByLocation(filterEngineers())
	sort.Sort(filtered) // sort by location
	for _, engineer := range filtered {
		_, err := fmt.Fprintf(w, formatString,
			engineer.FullName,
			engineer.Grade,
			engineer.WorkProfile,
			strings.Join(pmo.RemoveDuplicates(engineer.Account), ","),
			strings.Join(pmo.RemoveDuplicates(engineer.Project), ","),
			engineer.Manager,
			strings.Join(pmo.RemoveDuplicates(engineer.AssignmentStatuses()), ","))
		if err != nil {
			log.Fatalf("[ERR] failed on writing %v to tabwriter. E: %v", engineer, err)
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatalln("can not flush tabwriter")
	}

}

func main() {
	// read config
	config = pmo.ReadConfig(relativeConfigFilePath)
	login()
	printTable() // print table representation of engineers
}
