package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vistrcm/pmoclient/pmo"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

// EngineersSheet represents metadata about google sheet with engineers data.
type EngineersSheet struct {
	srv           *sheets.Service
	spreadsheetID string
	namesRange    string
	appendRange   string
}

// GetNames return names defined in spreadsheet
func (es *EngineersSheet) GetNames() []string {
	resp, err := es.srv.Spreadsheets.Values.Get(es.spreadsheetID, es.namesRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	var result []string
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			name := row[0].(string)
			result = append(result, name)
		}
	}
	return result
}

func slice2str(data []string) string {
	return strings.Join(pmo.RemoveDuplicates(data), "\n")
}

func sliceI2str(data []int) string {
	var new []string
	for _, v := range data {
		new = append(new, strconv.Itoa(v))
	}
	return strings.Join(pmo.RemoveDuplicates(new), "\n")
}

// AppendPerson to append  Person to the spreadsheet.
func (es *EngineersSheet) AppendPerson(engineer pmo.Person) {

	// engineer.ID                int
	// engineer.EmployeeID        string
	// engineer.Location          string
	// engineer.Manager           string
	// engineer.Grade             string
	// engineer.Specialization    string
	// engineer.WorkProfile       string
	// engineer.Position          string
	// engineer.FullName          string
	// engineer.AssignmentStart   []int
	// engineer.Project           []string
	// engineer.Account           []string
	// engineer.AssignmentFinish  []int
	// engineer.AssignmentComment []string
	// engineer.Involvements      []int
	// engineer.StaffPositionID   int
	// engineer.ProjectID         int
	// engineer.Involvement       int
	// engineer.Status            string
	// engineer.BenchStart        int
	// engineer.DaysOnBench       int
	// engineer.DaysAvailable     int
	// engineer.DaysOnBenchAlt    int
	// engineer.BenchStartAlt     int
	// engineer.TotalInvolvement  string
	// engineer.NewBenchStart     string
	// engineer.CanBeMovedToBench bool

	values := []interface{}{
		engineer.ID,
		engineer.EmployeeID,
		engineer.Location,
		engineer.Manager,
		engineer.Grade,
		engineer.Specialization,
		engineer.WorkProfile,
		engineer.Position,
		engineer.FullName,
		sliceI2str(engineer.AssignmentStart),
		slice2str(engineer.Project),
		slice2str(engineer.Account),
		sliceI2str(engineer.AssignmentFinish),
		slice2str(engineer.AssignmentComment),
		sliceI2str(engineer.Involvements),
		engineer.StaffPositionID,
		engineer.ProjectID,
		engineer.Involvement,
		engineer.Status,
		engineer.BenchStart,
		engineer.DaysOnBench,
		engineer.DaysAvailable,
		engineer.DaysOnBenchAlt,
		engineer.BenchStartAlt,
		engineer.TotalInvolvement,
		engineer.NewBenchStart,
		engineer.CanBeMovedToBench,
	}

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	_, err := es.srv.Spreadsheets.Values.Append(es.spreadsheetID, es.appendRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data from sheet: %v", err)
	}
}

// NewEngineersSheet generates new
func NewEngineersSheet(spreadsheetID string, secretFile string) EngineersSheet {
	client := clientFromFile(secretFile)
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("unable to retrieve Sheets client: %v", err)
	}

	// do some work
	es := EngineersSheet{srv: srv, spreadsheetID: spreadsheetID, namesRange: "list!A2:A35", appendRange: "AutofillFromPMO!A2"}
	return es
}

func clientFromFile(secretFile string) *http.Client {
	b, err := ioutil.ReadFile(secretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved gdoc_client_secret.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)
	return client
}

// getClient retrieves a token, saves the token, then returns the generated client
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	defer checkDefer(f.Close)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Save a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	log.Printf("Saving credentials file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer checkDefer(f.Close)
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("unable to encode token: %v", err)
	}
}

// checkDefer helper to catch errors in deferred functions
func checkDefer(f func() error) {
	if err := f(); err != nil {
		fmt.Println("Received error:", err)
	}
}
