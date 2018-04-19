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

var config pmo.Configuration
var useSpreadSheet = flag.Bool("spreadsheet", false, "use spreadsheet to get names and update spreadsheet at the end")

func main() {
	flag.Parse()

	// read config
	config = pmo.ReadConfig(relativeConfigFilePath)

	p := pmo.NewPMO(config)
	p.Login()

	var engineers []pmo.Person
	if !*useSpreadSheet {
		engineers = p.FilterEngineersByConfig()
	}

	pmo.PrintTable(engineers, formatString) // print table representation of engineers

	if *useSpreadSheet {
		processSpreadsheet(engineers)
	}

}

func processSpreadsheet(engineers []pmo.Person) {
	es := gdocs.NewEngineersSheet(config.Spreadsheet.SpreadsheetID, config.Spreadsheet.SecretFile)
	es.Clear()
	es.AppendEngineers(engineers)
}
