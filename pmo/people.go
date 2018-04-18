package pmo

type assignStatus struct {
	Color string
	Name  string
}

// Person contains person-related information presented in PMO
type Person struct {
	ID                int
	EmployeeID        string
	Location          string
	Manager           string
	Grade             string
	Specialization    string
	WorkProfile       string
	Position          string
	FullName          string
	AssignmentStart   []int
	Project           []string
	Account           []string
	AssignmentFinish  []int
	AssignmentComment []string
	Involvements      []int
	AssignmentStatus  []assignStatus
	StaffPositionID   int
	ProjectID         int
	Involvement       int
	Status            string
	BenchStart        int
	DaysOnBench       int
	DaysAvailable     int
	DaysOnBenchAlt    int
	BenchStartAlt     int
	TotalInvolvement  string
	NewBenchStart     string
	CanBeMovedToBench bool
}

// AssignmentStatuses returns list of assignment statuses
func (p *Person) AssignmentStatuses() []string {
	result := make([]string, 0)
	for _, element := range p.AssignmentStatus {
		result = append(result, element.Name)
	}
	return result
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
	Page    int      `json:"page"`
	Total   int      `json:"total"`
	Records int      `json:"records"`
	Rows    []Person `json:"rows"`
}

// Configuration of PMO client
type Configuration struct {
	Username      string   `json:"username"`
	Password      string   `json:"password"`
	FilterUsers   []string `json:"filterUsers"`
	LoginURL      string   `json:"loginUrl"`
	PeopleListURL string   `json:"peopleListUrl"`
}
