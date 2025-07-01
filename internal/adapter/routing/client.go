package routing

type HTTPMethod int

const (
	MethodPost HTTPMethod = iota
	MethodGet
)

func (m HTTPMethod) String() string {
	switch m {
	case MethodPost:
		return "POST"
	case MethodGet:
		return "GET"
	default:
		return "UNKNOWN"
	}
}

type Client[TransReq any, TransRes any] interface {
	Post(endpoint string, req TransReq) (TransRes, error)
	Get(endpoint string, params TransReq) (TransRes, error)
}
