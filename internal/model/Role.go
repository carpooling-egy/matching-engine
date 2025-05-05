package model

import "matching-engine/internal/enums"

type Role interface {
	GetRoleType() enums.RoleType
}
