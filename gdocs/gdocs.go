package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
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
	cleanRange    string
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

// Clear spreadsheet defined in spreadsheetID
func (es *EngineersSheet) Clear() {
	var vr sheets.ClearValuesRequest

	_, err := es.srv.Spreadsheets.Values.Clear(es.spreadsheetID, es.cleanRange, &vr).Do()
	if err != nil {
		log.Fatalf("Unable to clear data from sheet: %v\n", err)
	}
}

// AppendEngineers from engineers slice to the spreadsheet
func (es *EngineersSheet) AppendEngineers(engineers []pmo.Person) {
	// add header
	values := []interface{}{
		"FullName",
		"Location",
		"Grade",
		"Account",
		"Project",
		"Manager",
		"Profile",
		"Specialization",
		"EmployeeID",
	}

	var vr sheets.ValueRange
	vr.Values = append(vr.Values, values)

	_, err := es.srv.Spreadsheets.Values.Append(es.spreadsheetID, es.appendRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data from sheet: %v", err)
	}

	for _, engineer := range engineers {
		es.appendEngineer(engineer)
	}
}

// appendEngineer to append  Person to the spreadsheet.
func (es *EngineersSheet) appendEngineer(engineer pmo.Person) {
	values := []interface{}{
		engineer.Name,
		engineer.Location,
		engineer.Grade,
		slice2str(engineer.GetAccounts()),
		slice2str(engineer.GetProjects()),
		engineer.Manager,
		engineer.Profile,
		engineer.Specialization,
		engineer.ID,
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
	es := EngineersSheet{
		srv:           srv,
		spreadsheetID: spreadsheetID,
		namesRange:    "list!A2:A100",
		appendRange:   "AutofillFromPMO!A1",
		cleanRange:    "AutofillFromPMO!A1:ZZ1000",
	}
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
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	tokFile := usr.HomeDir + "/.config/pmoclient_gdoc_token.json"
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
