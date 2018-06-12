package pmo

import (
	"strings"
	"fmt"
)

// employee data
type employee struct {
	EmployeeID int    `json:"employeeId"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
}

// engineerManagers represents Engineering Managers in for of employ
type engineerManagers struct {
	Employee   employee `json:"employee"`
	Discipline string   `json:"discipline"`
}

// assignment defines engineer assignment
type assignment struct {
	ID          int    `json:"id"`
	EmployeeID  int    `json:"employeeId"`
	Account     string `json:"account"`
	Project     string `json:"project"`
	Start       string `json:"start"`
	Finish      string `json:"finish"`
	StartDate   string `json:"startDate"`
	FinishDate  string `json:"finishDate"`
	Involvement int    `json:"involvement"`
	Status      string `json:"status"`
	Comment     string `json:"comment"`
}

// Person contains person-related information presented in PMO
type Person struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	Username         string             `json:"username"`
	Grade            string             `json:"grade"`
	Specialization   string             `json:"specialization"`
	Profile          string             `json:"profile"`
	Position         string             `json:"position"`
	ServiceLine      string             `json:"serviceLine"`
	Location         string             `json:"location"`
	Manager          string             `json:"manager"`
	AvailableDays    int                `json:"availableDays"`
	DaysOnBench      int                `json:"daysOnBench"`
	Assignments      []assignment       `json:"assignments"`
	EngineerManagers []engineerManagers `json:"engineerManagers"`
	InBusinessTrip   bool               `json:"inBusinessTrip"`
}

// GetAssignments returns assignments in form `account-project-involvement`
func (p *Person) GetAssignmentsString() []string {
	var assignments []string
	for _, assignment := range p.Assignments {
		assignments = append(assignments, fmt.Sprintf("%q-%q-%d", assignment.Account, assignment.Project, assignment.Involvement))
	}
	return assignments
}

// GetEngineerManagers return list of managers
func (p *Person)GetEngineerManagers() []string {
	var managers []string
	for _, manager := range p.EngineerManagers {
		managers = append(managers, manager.Employee.Username)
	}
	return managers
}

// GetAccounts returns list of accounts this engineer is working on
func (p *Person) GetAccounts() []string {
	var accounts []string
	for _, assignment := range p.Assignments {
		accounts = append(accounts, assignment.Account)
	}
	return RemoveDuplicates(accounts)
}

// GetProjects returns list of projects this engineer is assigned to
func (p *Person) GetProjects() []string {
	var projects []string
	for _, assignment := range p.Assignments {
		projects = append(projects, assignment.Project)
	}
	return RemoveDuplicates(projects)
}

// AssignmentStatuses returns list of assignment statuses
func (p *Person) AssignmentStatuses() []string {
	result := make([]string, 0)
	for _, element := range p.Assignments {
		result = append(result, element.Status)
	}
	return result
}

// GetAccountsString return list of accounts as a string
func (p *Person)GetAccountsString() string {
	return strings.Join(p.GetAccounts(), ",")
}

// GetProjectsString return list of projects as a string
func (p *Person)GetProjectsString() string {
	return strings.Join(p.GetProjects(), ",")
}

// ByLocation implements sort.Interface for []Person based on
// the Location field.
type ByLocation []Person

func (slice ByLocation) Len() int {
	return len(slice)
}

func (slice ByLocation) Less(i, j int) bool {
	return slice[i].Location < slice[j].Location
}

func (slice ByLocation) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// APIResponse represent fields in PMO api response
type APIResponse struct {
	Data     []Person `json:"data"`
	Messages []string `json:"messages"`
}

// EngineersSpreadsheet is google spreadsheet with engineering data
type EngineersSpreadsheet struct {
	SpreadsheetID string `json:"SpreadsheetID"`
	SecretFile    string `json:"SecretFile"`
}
