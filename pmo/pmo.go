package pmo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// Configuration of PMO client
type Configuration struct {
	Username      string               `json:"username"`
	Password      string               `json:"password"`
	FilterUsers   []string             `json:"filterUsers"`
	LoginURL      string               `json:"loginUrl"`
	PeopleListURL string               `json:"peopleListUrl"`
	Spreadsheet   EngineersSpreadsheet `json:"Spreadsheet"`
}

// PMO representation
type PMO struct {
	config Configuration
	client *http.Client
}

// NewPMO returns prepared PMO structure
func NewPMO(config Configuration) PMO {
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

	var pmo = PMO{config: config, client: client}
	return pmo
}

// Login to PMO
func (pmo *PMO) Login() {
	config := pmo.config
	resp, err := pmo.client.PostForm(config.LoginURL, url.Values{"j_username": {config.Username}, "j_password": {config.Password}})
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
func (pmo *PMO) request(url string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Error happened %v, req: %v", err, req)
	}
	resp, err := pmo.client.Do(req)
	if err != nil {
		log.Fatalf("error on sending request: %v, resp: %v", err, resp)
	}
	return resp
}

// get list of engineers by sending request to
func (pmo *PMO) engineers() []Person {
	resp := pmo.request(pmo.config.PeopleListURL)
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Fatalf("Error when closing: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("request to %q completed with status %s\n", resp.Request.URL, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error on reading: %v. body: %v", err, body)
	}
	var peopleResponse = new(APIResponse)
	err = json.Unmarshal(body, &peopleResponse)
	if err != nil {
		log.Fatalf("something happened during unmarshall: %v. body: %v\n", err, string(body))
	}
	return peopleResponse.Data
}

// FilterEngineers returns only data for subset of engineers defined in `filter``
func (pmo *PMO) FilterEngineers(filter []string) []Person {
	// initialize return value
	filteredEngineers := make([]Person, 0)

	// initialize temporary map for filtering
	filterMap := make(map[string]bool)
	for _, u := range filter {
		filterMap[strings.Replace(strings.ToLower(u), " ", "", -1)] = true
	}

	for _, val := range pmo.engineers() {
		targetKey := strings.Replace(strings.ToLower(val.Name), " ", "", -1)
		if filterMap[targetKey] {
			filteredEngineers = append(filteredEngineers, val)
		}
	}
	return filteredEngineers
}

// FilterEngineersByConfig using filter defined in config
func (pmo *PMO) FilterEngineersByConfig() []Person {
	return pmo.FilterEngineers(pmo.config.FilterUsers)
}
