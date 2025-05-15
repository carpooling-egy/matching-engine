package errors

const (
	ErrNilOfferNode            = "offer node is nil"
	ErrNilOffer                = "offer node has nil offer"
	ErrEmptyUserID             = "user ID is empty"
	ErrEmptyOfferID            = "offer ID is empty"
	ErrEmptyOfferIDORRequestID = "offer ID or request ID is empty"
	ErrNilPath                 = "offer node has nil path"
	ErrEmptyPath               = "offer node has empty path"
	ErrNilMatchedRequests      = "offer node has nil newly assigned matched requests"
	ErrEmptyMatchedRequests    = "offer node has empty newly assigned matched requests"
	ErrNoOffersOrRequests      = "no offers or requests provided"
)
