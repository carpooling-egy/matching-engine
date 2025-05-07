package mappers

import "matching-engine/internal/model"

type Mapper interface {
	Marshal(result *model.MatchingResult) ([]byte, error)
}
