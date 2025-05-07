package mappers

import (
	"encoding/json"
	"matching-engine/internal/adapter/messaging/natsjetstream/mappers/converters"
	"matching-engine/internal/model"
)

// JsonMapper implements the Mapper interface for JSON serialization
type JsonMapper struct {
	resultConverter *converters.ResultConverter
}

// Marshal serializes a matching result to JSON bytes
func (mapper *JsonMapper) Marshal(result *model.MatchingResult) ([]byte, error) {
	// Convert the domain model to DTO using the converter
	resultDTO := mapper.resultConverter.ToDTO(result)

	// Marshal the DTO to JSON
	return json.Marshal(resultDTO)
}

// NewJsonMapper creates a new JsonMapper
func NewJsonMapper() Mapper {
	return &JsonMapper{
		resultConverter: converters.NewResultConverter(),
	}
}
