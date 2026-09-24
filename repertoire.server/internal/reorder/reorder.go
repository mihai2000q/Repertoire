package reorder

import (
	"errors"
	"reflect"
	"repertoire/server/internal/httperror"

	"github.com/google/uuid"
)

// Config holds optional configuration for reordering
type Config struct {
	IDField               string // Field name for ID (default: "ID")
	OrderField            string // Field name for Order (default: "Order")
	StartOffset           uint   // Starting value for order (default: 0)
	EntityNotFoundMsg     string // Error message when entity not found
	OverEntityNotFoundMsg string // Error message when over entity not found
}

// DefaultConfig returns a config with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		IDField:               "ID",
		OrderField:            "Order",
		StartOffset:           0,
		EntityNotFoundMsg:     "entity not found",
		OverEntityNotFoundMsg: "over entity not found",
	}
}

// MoveEntity updates Order fields when moving an entity within a slice.
// It does NOT reorder the slice – only updates the Order property.
func MoveEntity[T any](entities []T, id, overID uuid.UUID, config *Config) *httperror.ErrorCode {
	cfg := mergeConfig(config)

	if len(entities) == 0 {
		return httperror.ConflictError(errors.New("empty slice"))
	}

	index, overIndex, errCode := getIndexes(entities, id, overID, cfg)
	if errCode != nil {
		return errCode
	}

	return updateOrderFields(entities, index, overIndex, cfg)
}

// MergeConfig merges user config with defaults
func mergeConfig(userConfig *Config) *Config {
	config := DefaultConfig()
	if userConfig == nil {
		return config
	}

	if userConfig.IDField != "" {
		config.IDField = userConfig.IDField
	}
	if userConfig.OrderField != "" {
		config.OrderField = userConfig.OrderField
	}
	if userConfig.StartOffset != 0 {
		config.StartOffset = userConfig.StartOffset
	}
	if userConfig.EntityNotFoundMsg != "" {
		config.EntityNotFoundMsg = userConfig.EntityNotFoundMsg
	}
	if userConfig.OverEntityNotFoundMsg != "" {
		config.OverEntityNotFoundMsg = userConfig.OverEntityNotFoundMsg
	}

	return config
}

// getIndexes finds the indices of the entity and the target entity
func getIndexes[T any](entities []T, id, overID uuid.UUID, cfg *Config) (int, int, *httperror.ErrorCode) {
	var index, overIndex int
	found := false
	overFound := false

	for i := 0; i < len(entities); i++ {
		item := reflect.ValueOf(entities[i])
		idField := item.FieldByName(cfg.IDField)

		if !idField.IsValid() {
			return -1, -1, httperror.InternalServerError(errors.New("ID field '" + cfg.IDField + "' not found on type"))
		}

		// Type check to avoid panic
		idValue, ok := idField.Interface().(uuid.UUID)
		if !ok {
			return -1, -1, httperror.InternalServerError(errors.New("ID field is not of type uuid.UUID"))
		}

		if idValue == id {
			index = i
			found = true
		}
		if idValue == overID {
			overIndex = i
			overFound = true
		}

		if found && overFound {
			break // Early exit when both found
		}
	}

	if !found {
		return -1, -1, httperror.NotFoundError(errors.New(cfg.EntityNotFoundMsg))
	}
	if !overFound {
		return -1, -1, httperror.NotFoundError(errors.New(cfg.OverEntityNotFoundMsg))
	}

	return index, overIndex, nil
}

// updateOrderFields updates the Order fields based on the move operation.
func updateOrderFields[T any](entities []T, index, overIndex int, cfg *Config) *httperror.ErrorCode {
	if index < 0 || index >= len(entities) || overIndex < 0 || overIndex >= len(entities) {
		return httperror.ConflictError(errors.New("index out of range"))
	}

	// Validate Order field exists and is either uint or *uint
	testVal := reflect.ValueOf(entities[0])
	orderFieldVal := testVal.FieldByName(cfg.OrderField)
	if !orderFieldVal.IsValid() {
		return httperror.InternalServerError(errors.New("Order field '" + cfg.OrderField + "' not found on type"))
	}

	kind := orderFieldVal.Kind()
	if kind != reflect.Uint && !(kind == reflect.Ptr && orderFieldVal.Type().Elem().Kind() == reflect.Uint) {
		return httperror.InternalServerError(errors.New("order field must be of type uint or *uint"))
	}

	if index < overIndex {
		// Moving forward: items between index+1 and overIndex shift left (decrease by 1)
		for i := index + 1; i <= overIndex; i++ {
			setOrder(&entities[i], cfg.StartOffset+uint(i-1), cfg.OrderField)
		}
	} else if index > overIndex {
		// Moving backward: items between overIndex and index shift right (increase by 1)
		for i := overIndex; i <= index; i++ {
			setOrder(&entities[i], cfg.StartOffset+uint(i+1), cfg.OrderField)
		}
	} else {
		// Same position, nothing to do
		return nil
	}

	// Set the moved entity's Order to the target position
	setOrder(&entities[index], cfg.StartOffset+uint(overIndex), cfg.OrderField)

	return nil
}

// setOrder sets the Order field using reflection.
// Supports both uint and *uint field types.
func setOrder[T any](entity *T, order uint, fieldName string) {
	v := reflect.ValueOf(entity).Elem()
	field := v.FieldByName(fieldName)

	if !field.IsValid() || !field.CanSet() {
		return
	}

	kind := field.Kind()
	if kind == reflect.Uint {
		field.SetUint(uint64(order))
	} else if kind == reflect.Ptr && field.Type().Elem().Kind() == reflect.Uint {
		// Allocate a new uint and set the pointer
		newVal := reflect.New(field.Type().Elem())
		newVal.Elem().SetUint(uint64(order))
		field.Set(newVal)
	}
}
