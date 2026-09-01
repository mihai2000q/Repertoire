package deduplicate

import "github.com/google/uuid"

func Deduplicate(ids []uuid.UUID) map[uuid.UUID]bool {
	idsSet := make(map[uuid.UUID]bool, len(ids))
	for _, id := range ids {
		idsSet[id] = true
	}
	return idsSet
}
