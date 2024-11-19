package audio

import (
	"github.com/Dmytro-Kucherenko/users-sam/internal/common/types"
)

type GetAppsFilters struct {
	Indices types.Optional[[]uint32]
	Names   types.Optional[[]string]
}
