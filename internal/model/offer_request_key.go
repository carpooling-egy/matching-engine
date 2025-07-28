package model

type OfferRequestKey struct {
	offerID   string
	requestID string
}

func NewOfferRequestKey(offerID, requestID string) OfferRequestKey {
	return OfferRequestKey{
		offerID:   offerID,
		requestID: requestID,
	}
}

func (k OfferRequestKey) OfferID() string {
	return k.offerID
}

func (k OfferRequestKey) RequestID() string {
	return k.requestID
}

func (k OfferRequestKey) String() string {
	return "OfferRequestKey{OfferID: " + k.offerID + ", RequestID: " + k.requestID + "}"
}
