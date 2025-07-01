package model

type OSRMTransport struct {
	PathVars    map[string]string // e.g. {"coordinates": "..."}
	QueryParams map[string]any    // e.g. {"sources": "...", "destinations": "..."}
}
