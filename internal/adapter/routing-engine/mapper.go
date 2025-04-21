package routing_engine

type OperationMapper[
	DomainReq any,
	DomainRes any,
	TransReq any,
	TransRes any,
] interface {
	// ToTransport turns your domain request into whatever your client needs.
	ToTransport(DomainReq) TransReq

	// FromTransport turns the raw client response into your domain model.
	FromTransport(TransRes) DomainRes
}
