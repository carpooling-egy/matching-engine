package routing

type OperationMapper[
	DomainReq any,
	DomainRes any,
	TransReq any,
	TransRes any,
] interface {
	ToTransport(DomainReq) (TransReq, error)
	FromTransport(TransRes) (DomainRes, error)
}
