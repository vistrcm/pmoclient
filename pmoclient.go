// Package pmoclient implements a simple client for GD's PMO tool.

package main

import (
	"flag"

	"github.com/vistrcm/pmoclient/gdocs"
	"github.com/vistrcm/pmoclient/pmo"
)

// formatString for print results
const formatString = "%s\t%s\t%s\t%s\t%s\t%s\t%v\n"
const relativeConfigFilePath = "/.config/pmoclient.json"

func main() {
	var config pmo.Configuration
	var useSpreadSheet = flag.Bool("spreadsheet", false, "use spreadsheet to get names and update spreadsheet at the end")

	flag.Parse()
	// read config
	config = pmo.ReadConfig(relativeConfigFilePath)

	p := pmo.NewPMO(config)
	p.Login()

	engineersCh := make(chan []pmo.Person)
	doneCh := make(chan bool)

	go printTable(engineersCh, doneCh)

	var engineers []pmo.Person
	if !*useSpreadSheet {
		engineers = p.FilterEngineersByConfig()
		engineersCh <- engineers
	} else {
		es := gdocs.NewEngineersSheet(config.Spreadsheet.SpreadsheetID, config.Spreadsheet.SecretFile)
		filter := es.GetNames()
		engineers = p.FilterEngineers(filter)
		engineersCh <- engineers
		es.Clear()
		es.AppendEngineers(engineers)
	}
	<-doneCh
}

func printTable(engineersCh chan []pmo.Person, doneCh chan bool) {
	engineers := <-engineersCh
	pmo.PrintTable(engineers, formatString) // print table representation of engineers
	doneCh <- true
}
