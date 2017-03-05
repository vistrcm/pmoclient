// Package pmoclient implements a simple client for GD's PMO tool.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/vistrcm/pmoclient/pmo"
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
)

// some constants
const FormatString = "%s\t%s\t%s\t%s\t%s\t%s\t%v\n"
const RelativeConfigFilePath = "/.config/pmoclient.json"

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
	resp, err := client.PostForm(config.LoginUrl, url.Values{"j_username": {config.Username}, "j_password": {config.Password}})
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
	resp := request(config.PeopleListUrl)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
	}
	var peopleResponse = new(pmo.APIResponse)
	err = json.Unmarshal(body, &peopleResponse)
	if err != nil {
		log.Fatalf("something happened during unmarshall: %v. engineers: %v", err, engineers)
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
		if filterMap[strings.ToLower(val.EmployeeId)] {
			filteredEngineers = append(filteredEngineers, val)
		}
	}
	return filteredEngineers
}

// finction to print table representation of engineers
func printTable() {
	// Observe how the b's and the d's, despite appearing in the
	// second cell of each line, belong to different columns.
	//w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	w := tabwriter.NewWriter(os.Stdout, 5, 0, 1, ' ', 0)
	// print header
	fmt.Fprintf(w, FormatString,
		"Name",
		"Grade",
		"Profile",
		"Account",
		"Project",
		"Manager",
		"Status")
	// iterate over engineers and print only required from config
	filtered := pmo.ByLocation(filterEngineers())
	sort.Sort(filtered) // sort by location
	for _, enginer := range filtered {
		fmt.Fprintf(w, FormatString,
			enginer.FullName,
			enginer.Grade,
			enginer.WorkProfile,
			strings.Join(pmo.RemoveDups(enginer.Account), ","),
			strings.Join(pmo.RemoveDups(enginer.Project), ","),
			enginer.Manager,
			strings.Join(pmo.RemoveDups(enginer.AssignmentStatuses()), ","))
	}

	w.Flush()

}

func main() {
	// read config
	config = pmo.ReadConfig(RelativeConfigFilePath)
	login()
	printTable() // print table representation of engineers
}
