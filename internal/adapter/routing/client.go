package routing

type RoutingClient[TransReq any, TransRes any] interface {
	Post(endpoint string, req TransReq) (TransRes, error)
}
