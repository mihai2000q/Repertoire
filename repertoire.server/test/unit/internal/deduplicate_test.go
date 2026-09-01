package internal

import (
	"repertoire/server/internal/deduplicate"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeduplicate_WhenThereAreDuplicateIDs_ShouldRemoveThemAndReturnUniqueSet(t *testing.T) {
	tests := []struct {
		name     string
		input    []uuid.UUID
		expected []uuid.UUID
	}{
		{
			name:     "empty slice",
			input:    []uuid.UUID{},
			expected: []uuid.UUID{},
		},
		{
			name:     "single element",
			input:    []uuid.UUID{uuid.MustParse("00000000-0000-0000-0000-000000000001")},
			expected: []uuid.UUID{uuid.MustParse("00000000-0000-0000-0000-000000000001")},
		},
		{
			name: "multiple unique elements",
			input: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
			expected: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "with duplicates",
			input: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
			expected: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000002"),
				uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			},
		},
		{
			name: "all duplicates",
			input: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
			expected: []uuid.UUID{
				uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := deduplicate.Deduplicate(tt.input)

			// Verify length
			assert.Len(t, result, len(tt.expected))

			// Verify each expected ID exists in the map
			for _, id := range tt.expected {
				assert.True(t, result[id], "expected ID %v not in result", id)
			}
		})
	}
}
