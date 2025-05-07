package routing

type Client[TransReq any, TransRes any] interface {
	Post(endpoint string, req TransReq) (TransRes, error)
}
