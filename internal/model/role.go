package model

type Role interface {
	AsOffer() (*Offer, bool)
	AsRequest() (*Request, bool)
}
