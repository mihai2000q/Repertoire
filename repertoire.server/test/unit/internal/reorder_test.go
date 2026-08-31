package internal

import (
	"net/http"
	"repertoire/server/internal/reorder"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ------ Helper Types ------

type sectionEntity struct {
	ID    uuid.UUID
	Order uint
}

type trackEntity struct {
	ID      uuid.UUID
	TrackNo *uint // pointer to uint
}

// Helper to create uint pointer
func uintPtr(v uint) *uint {
	return &v
}

// ------ Test Cases ------

func TestMoveEntity_WhenEntitiesEmpty_ShouldReturnConflictError(t *testing.T) {
	// given
	var entities []sectionEntity
	id := uuid.New()
	overID := uuid.New()

	// when
	errCode := reorder.MoveEntity(entities, id, overID, nil)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusConflict, errCode.Code)
	assert.Equal(t, "empty slice", errCode.Error.Error())
}

func TestMoveEntity_WhenEntityNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	id1, id2 := uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id1, Order: 0},
		{ID: id2, Order: 1},
	}
	missingID := uuid.New()

	// when
	errCode := reorder.MoveEntity(entities, missingID, id1, nil)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "entity not found", errCode.Error.Error())
}

func TestMoveEntity_WhenOverEntityNotFound_ShouldReturnNotFoundError(t *testing.T) {
	// given
	id1, id2 := uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id1, Order: 0},
		{ID: id2, Order: 1},
	}
	missingID := uuid.New()

	// when
	errCode := reorder.MoveEntity(entities, id1, missingID, nil)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "over entity not found", errCode.Error.Error())
}

func TestMoveEntity_WhenMovingForward_ShouldUpdateOrdersCorrectly(t *testing.T) {
	// given
	id0, id1, id2, id3 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id0, Order: 0},
		{ID: id1, Order: 1},
		{ID: id2, Order: 2},
		{ID: id3, Order: 3},
	}
	// Move entity at index 0 to position of entity at index 2 (overIndex=2)
	// Expected Order updates (no slice reorder):
	// - Entity at index 1 (id1) gets Order = StartOffset + (i-1) = 0 + (1-1) = 0
	// - Entity at index 2 (id2) gets Order = 0 + (2-1) = 1
	// - Entity at index 0 (id0) gets Order = StartOffset + overIndex = 0 + 2 = 2
	// - Entity at index 3 (id3) stays 3

	// when
	errCode := reorder.MoveEntity(entities, id0, id2, nil)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, uint(2), entities[0].Order)
	assert.Equal(t, uint(0), entities[1].Order)
	assert.Equal(t, uint(1), entities[2].Order)
	assert.Equal(t, uint(3), entities[3].Order)
}

func TestMoveEntity_WhenMovingBackward_ShouldUpdateOrdersCorrectly(t *testing.T) {
	// given
	id0, id1, id2, id3 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id0, Order: 0},
		{ID: id1, Order: 1},
		{ID: id2, Order: 2},
		{ID: id3, Order: 3},
	}
	// Move entity at index 3 to position of entity at index 1 (overIndex=1)
	// Expected Order updates:
	// - Entity at index 1 (id1) gets Order = StartOffset + (i+1) = 0 + (1+1) = 2
	// - Entity at index 2 (id2) gets Order = 0 + (2+1) = 3
	// - Entity at index 3 (id3) gets Order = StartOffset + overIndex = 0 + 1 = 1
	// - Entity at index 0 (id0) stays 0

	// when
	errCode := reorder.MoveEntity(entities, id3, id1, nil)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, uint(0), entities[0].Order)
	assert.Equal(t, uint(2), entities[1].Order)
	assert.Equal(t, uint(3), entities[2].Order)
	assert.Equal(t, uint(1), entities[3].Order)
}

func TestMoveEntity_WhenSamePosition_ShouldDoNothing(t *testing.T) {
	// given
	id0, id1 := uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id0, Order: 0},
		{ID: id1, Order: 1},
	}
	originalOrders := []uint{0, 1}

	// when
	errCode := reorder.MoveEntity(entities, id0, id0, nil)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, originalOrders[0], entities[0].Order)
	assert.Equal(t, originalOrders[1], entities[1].Order)
}

func TestMoveEntity_WithStartOffsetOne_ShouldStartFromOne(t *testing.T) {
	// given
	id0, id1, id2, id3 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	entities := []sectionEntity{
		{ID: id0, Order: 1},
		{ID: id1, Order: 2},
		{ID: id2, Order: 3},
		{ID: id3, Order: 4},
	}
	config := &reorder.Config{
		StartOffset: 1,
	}
	// Move index 0 to overIndex 2 (forward)
	// Expected: index1 -> 1, index2 -> 2, index0 -> 3, index3 -> 4

	// when
	errCode := reorder.MoveEntity(entities, id0, id2, config)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, uint(3), entities[0].Order) // moved gets overIndex+1 = 3
	assert.Equal(t, uint(1), entities[1].Order) // i-1+1 = 1
	assert.Equal(t, uint(2), entities[2].Order) // i-1+1 = 2
	assert.Equal(t, uint(4), entities[3].Order) // unchanged
}

func TestMoveEntity_WithCustomFieldNames_ShouldUseCustomFields(t *testing.T) {
	// given
	id0, id1, id2 := uuid.New(), uuid.New(), uuid.New()
	type customEntity struct {
		CustomID    uuid.UUID
		CustomOrder uint
	}
	entities := []customEntity{
		{CustomID: id0, CustomOrder: 0},
		{CustomID: id1, CustomOrder: 1},
		{CustomID: id2, CustomOrder: 2},
	}
	config := &reorder.Config{
		IDField:    "CustomID",
		OrderField: "CustomOrder",
	}

	// when - move first to third position (index 0 -> overIndex 2)
	errCode := reorder.MoveEntity(entities, id0, id2, config)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, uint(2), entities[0].CustomOrder)
	assert.Equal(t, uint(0), entities[1].CustomOrder)
	assert.Equal(t, uint(1), entities[2].CustomOrder)
}

func TestMoveEntity_WithPointerUint_ShouldUpdatePointerValues(t *testing.T) {
	// given
	id0, id1, id2 := uuid.New(), uuid.New(), uuid.New()
	entities := []trackEntity{
		{ID: id0, TrackNo: uintPtr(1)},
		{ID: id1, TrackNo: uintPtr(2)},
		{ID: id2, TrackNo: uintPtr(3)},
	}
	config := &reorder.Config{
		OrderField:  "TrackNo",
		StartOffset: 1,
	}
	// Move index 0 to overIndex 2 (forward)
	// Expected: index1 gets 1, index2 gets 2, index0 gets 3

	// when
	errCode := reorder.MoveEntity(entities, id0, id2, config)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, entities[0].TrackNo)
	assert.NotNil(t, entities[1].TrackNo)
	assert.NotNil(t, entities[2].TrackNo)
	assert.Equal(t, uint(3), *entities[0].TrackNo) // moved
	assert.Equal(t, uint(1), *entities[1].TrackNo) // shifted left
	assert.Equal(t, uint(2), *entities[2].TrackNo) // shifted left
}

func TestMoveEntity_WithPointerUint_Backward_ShouldUpdateCorrectly(t *testing.T) {
	// given
	id0, id1, id2, id3 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	entities := []trackEntity{
		{ID: id0, TrackNo: uintPtr(1)},
		{ID: id1, TrackNo: uintPtr(2)},
		{ID: id2, TrackNo: uintPtr(3)},
		{ID: id3, TrackNo: uintPtr(4)},
	}
	config := &reorder.Config{
		OrderField:  "TrackNo",
		StartOffset: 1,
	}
	// Move index 3 to overIndex 1 (backward)
	// Expected: index1 gets 3, index2 gets 4, index3 gets 2, index0 stays 1

	// when
	errCode := reorder.MoveEntity(entities, id3, id1, config)

	// then
	assert.Nil(t, errCode)
	assert.Equal(t, uint(1), *entities[0].TrackNo)
	assert.Equal(t, uint(3), *entities[1].TrackNo)
	assert.Equal(t, uint(4), *entities[2].TrackNo)
	assert.Equal(t, uint(2), *entities[3].TrackNo)
}

func TestMoveEntity_WithPointerUint_NilField_ShouldAllocate(t *testing.T) {
	// given
	id0, id1 := uuid.New(), uuid.New()
	entities := []trackEntity{
		{ID: id0, TrackNo: nil},
		{ID: id1, TrackNo: nil},
	}
	config := &reorder.Config{
		OrderField:  "TrackNo",
		StartOffset: 1,
	}
	// Move index 0 to overIndex 1 (forward)
	// Expected: index0 gets 2 (overIndex+1), index1 gets 1 (i-1+1)

	// when
	errCode := reorder.MoveEntity(entities, id0, id1, config)

	// then
	assert.Nil(t, errCode)
	assert.NotNil(t, entities[0].TrackNo)
	assert.NotNil(t, entities[1].TrackNo)
	assert.Equal(t, uint(2), *entities[0].TrackNo) // moved
	assert.Equal(t, uint(1), *entities[1].TrackNo) // shifted
}

func TestMoveEntity_WithInvalidField_ShouldReturnInternalServerError(t *testing.T) {
	// given
	type invalidEntity struct {
		ID    uuid.UUID
		Wrong string
	}
	entities := []invalidEntity{
		{ID: uuid.New(), Wrong: "not numeric"},
	}
	config := &reorder.Config{
		OrderField: "Wrong",
	}

	// when
	errCode := reorder.MoveEntity(entities, entities[0].ID, entities[0].ID, config)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusInternalServerError, errCode.Code)
	assert.Contains(t, errCode.Error.Error(), "order field must be of type uint or *uint")
}

func TestMoveEntity_WithCustomErrorMessages_ShouldUseThem(t *testing.T) {
	// given
	entities := []sectionEntity{
		{ID: uuid.New(), Order: 0},
	}
	missingID := uuid.New()
	config := &reorder.Config{
		EntityNotFoundMsg:     "custom entity missing",
		OverEntityNotFoundMsg: "custom over missing",
	}

	// when
	errCode := reorder.MoveEntity(entities, missingID, entities[0].ID, config)

	// then
	require.NotNil(t, errCode)
	assert.Equal(t, http.StatusNotFound, errCode.Code)
	assert.Equal(t, "custom entity missing", errCode.Error.Error())

	// also test over not found
	errCode = reorder.MoveEntity(entities, entities[0].ID, missingID, config)
	require.NotNil(t, errCode)
	assert.Equal(t, "custom over missing", errCode.Error.Error())
}
