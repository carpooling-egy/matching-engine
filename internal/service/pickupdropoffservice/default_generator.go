package pickupdropoffservice

import "matching-engine/internal/model"

type DefaultGenerator struct {
}

func NewDefaultSelector() *DefaultGenerator {
	return &DefaultGenerator{}
}
func (ds *DefaultGenerator) Generate(request *model.Request, offer *model.Offer) (pickup, dropoff *model.PathPoint, err error) {
	return nil, nil, nil
}
