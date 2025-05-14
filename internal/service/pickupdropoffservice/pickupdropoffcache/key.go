package pickupdropoffcache

type Key struct {
	offerID   string
	requestID string
}

func NewKey(offerID, requestID string) Key {
	return Key{
		offerID:   offerID,
		requestID: requestID,
	}
}

func (k Key) OfferID() string {
	return k.offerID
}

func (k Key) RequestID() string {
	return k.requestID
}

func (k Key) String() string {
	return "CacheKey{OfferID: " + k.offerID + ", RequestID: " + k.requestID + "}"
}
