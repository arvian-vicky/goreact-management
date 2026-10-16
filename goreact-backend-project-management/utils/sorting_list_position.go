package utils

import (
	"arfian-nurdavicky/project-management/models"

	"github.com/google/uuid"
)

func SortListsByPosition(lists []models.List, order []uuid.UUID) []models.List {
	ordered := make([]models.List, 0, len(order))

	// mapping biar gampang lookup
	listMap := make(map[uuid.UUID]models.List)
	for _, l := range lists {
		listMap[l.PublicID] = l
	}

	// urutkan sesuai array order
	for _, id := range order {
		if l, ok := listMap[id]; ok {
			ordered = append(ordered, l)
		}
	}

	return ordered
}
