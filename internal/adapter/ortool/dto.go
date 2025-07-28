package ortool

// NewORToolData creates a new ORToolData instance with the provided parameters.
func NewORToolData(
	timeMatrix [][]int,
	timeWindows [][2]int,
	noOfRidersPerRequest []int,
	vehicleCapacity int,
	pickupAndDropoffs [][2]int,
	maxRouteDuration int,
) *ORToolData {
	return &ORToolData{
		TimeMatrix:           timeMatrix,
		TimeWindows:          timeWindows,
		NoOfRidersPerRequest: noOfRidersPerRequest,
		VehicleCapacity:      vehicleCapacity,
		PickupAndDropoffs:    pickupAndDropoffs,
		MaxRouteDuration:     maxRouteDuration,
	}
}

type ORToolData struct {
	TimeMatrix              [][]int  `json:"time_matrix"`
	TimeWindows             [][2]int `json:"time_windows"`
	NoOfRidersPerRequest    []int    `json:"no_of_riders_per_request"`
	VehicleCapacity         int      `json:"vehicle_capacity"`
	PickupAndDropoffs       [][2]int `json:"pickup_and_dropoffs"`
	MaxRouteDuration        int      `json:"max_route_duration,omitempty"`
	Timeout                 int      `json:"timeout,omitempty"`                    // Timeout in milliseconds for the OR-Tools solver
	Method                  string   `json:"method,omitempty"`                     // Method to use for solving, e.g., "parallel_cheapest_insertion" or :path_cheapest_arc" or "global_cheapest_arch" or ...
	EnableGuidedLocalSearch bool     `json:"enable_guided_local_search,omitempty"` // Enable guided local search for the solver
}

type ORToolSolutionStep struct {
	Node        int `json:"node"`
	ArrivalTime int `json:"arrival_time"`
}

type ORToolSolutionResponse struct {
	Success bool                 `json:"success"`
	Route   []ORToolSolutionStep `json:"route"`
}

func (s *ORToolData) SetMethod(method string) {
	s.Method = method
}

func (s *ORToolData) SetTimeout(timeout int) {
	s.Timeout = timeout
}

func (s *ORToolData) SetEnableGuidedLocalSearch(enable bool) {
	s.EnableGuidedLocalSearch = enable
}
