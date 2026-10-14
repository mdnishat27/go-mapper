package mapper

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ErrMapperNotFound is returned (wrapped) by Map and MapList when no mapping
// function is registered for the requested type pair
var ErrMapperNotFound = errors.New("no mapper registered")

type mapKey struct {
	Source reflect.Type
	Dest   reflect.Type
}

// MappingFunction defines the signature for custom mapping functions
type MappingFunction[Source any, Destination any] func(Source) (Destination, error)

type mappingRegistry struct {
	mu       sync.RWMutex
	mappings map[mapKey]any
}

var globalRegistry = &mappingRegistry{
	mappings: make(map[mapKey]any),
}

// typeOf returns the static type of T. Unlike reflect.TypeOf(value) it also
// works when T is an interface type.
func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func getKey[Source any, Destination any]() mapKey {
	return mapKey{Source: typeOf[Source](), Dest: typeOf[Destination]()}
}

// lookup returns the mapping function registered for the given type pair
func lookup[Source any, Destination any]() (MappingFunction[Source, Destination], error) {
	key := getKey[Source, Destination]()

	globalRegistry.mu.RLock()
	mappingFn, exists := globalRegistry.mappings[key]
	globalRegistry.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w for %s to %s", ErrMapperNotFound, key.Source, key.Dest)
	}

	return mappingFn.(MappingFunction[Source, Destination]), nil
}

// RegisterMapper registers a custom mapping function for the given type pair.
// A previously registered function for the same pair is replaced.
// It panics if fn is nil.
func RegisterMapper[Source any, Destination any](fn MappingFunction[Source, Destination]) {
	key := getKey[Source, Destination]()

	if fn == nil {
		panic(fmt.Sprintf("mapper: nil mapping function for %s to %s", key.Source, key.Dest))
	}

	globalRegistry.mu.Lock()
	globalRegistry.mappings[key] = fn
	globalRegistry.mu.Unlock()
}

// Map performs type-safe mapping from source to destination using registered mapping function
func Map[Source any, Destination any](source Source) (Destination, error) {
	fn, err := lookup[Source, Destination]()
	if err != nil {
		var zeroDest Destination
		return zeroDest, err
	}

	return fn(source)
}

// MapList maps a slice of source structs to a slice of destination structs.
// It stops at the first element that fails and returns a nil slice together
// with an error that wraps the mapping function's error and names the index.
func MapList[Source any, Destination any](sources []Source) ([]Destination, error) {
	fn, err := lookup[Source, Destination]()
	if err != nil {
		return nil, err
	}

	if sources == nil {
		return nil, nil
	}

	destList := make([]Destination, 0, len(sources))
	for i, item := range sources {
		result, err := fn(item)
		if err != nil {
			return nil, fmt.Errorf("mapping element %d: %w", i, err)
		}

		destList = append(destList, result)
	}

	return destList, nil
}

// ClearRegistry removes all registered mappings (useful for testing)
func ClearRegistry() {
	globalRegistry.mu.Lock()
	globalRegistry.mappings = make(map[mapKey]any)
	globalRegistry.mu.Unlock()
}

// UnregisterMapper removes a specific mapping function for the given type pair
func UnregisterMapper[Source any, Destination any]() bool {
	key := getKey[Source, Destination]()

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	_, exists := globalRegistry.mappings[key]
	if exists {
		delete(globalRegistry.mappings, key)
	}

	return exists
}

// HasMapping checks if a mapping function is registered for the given type pair
func HasMapping[Source any, Destination any]() bool {
	key := getKey[Source, Destination]()

	globalRegistry.mu.RLock()
	_, exists := globalRegistry.mappings[key]
	globalRegistry.mu.RUnlock()

	return exists
}
