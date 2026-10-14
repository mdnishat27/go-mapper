package mapper

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

type SourceStruct struct {
	Field1 string
	Field2 int
}

type DestinationStruct struct {
	FieldA string
	FieldB int
}

type SourceStruct2 struct {
	Field1 string
	Field2 int
}

type DestinationStruct2 struct {
	FieldA string
	FieldB int
	FieldC string
}

func MapSourceStructToDestinationStructNoPointer(sourceStruct SourceStruct) (DestinationStruct, error) {
	dest := DestinationStruct{}
	dest.FieldA = sourceStruct.Field1
	dest.FieldB = sourceStruct.Field2 + 2
	return dest, nil
}

func MapSourceStructToDestinationStruct(source SourceStruct) (*DestinationStruct, error) {
	dest := &DestinationStruct{}
	dest.FieldA = source.Field1
	dest.FieldB = source.Field2 + 2
	return dest, nil
}

func MapDestinationStructToSourceStruct(source DestinationStruct) (SourceStruct, error) {
	dest := SourceStruct{}
	dest.Field1 = source.FieldA
	dest.Field2 = source.FieldB + 2
	return dest, nil
}

func MapSourceStructToDestinationStruct2(source SourceStruct2) (*DestinationStruct2, error) {
	dest := &DestinationStruct2{}
	dest.FieldA = source.Field1
	dest.FieldB = source.Field2 + 2
	dest.FieldC = "Hello " + source.Field1
	return dest, nil
}

func TestMapper(t *testing.T) {
	tests := []struct {
		name  string
		setup func()
		test  func(t *testing.T)
	}{
		{
			name: "Single struct to pointer mapping",
			setup: func() {
				RegisterMapper[SourceStruct, *DestinationStruct](MapSourceStructToDestinationStruct)
			},
			test: func(t *testing.T) {
				source := SourceStruct{Field1: "Hello", Field2: 42}

				dest, err := Map[SourceStruct, *DestinationStruct](source)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if dest == nil {
					t.Fatal("Expected destination to not be nil")
				}
				if dest.FieldA != "Hello" {
					t.Errorf("Expected FieldA='Hello', got '%s'", dest.FieldA)
				}
				if dest.FieldB != 44 { // 42 + 2
					t.Errorf("Expected FieldB=44, got %d", dest.FieldB)
				}
			},
		},
		{
			name: "Single struct to struct mapping",
			setup: func() {
				RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
			},
			test: func(t *testing.T) {
				source := SourceStruct{Field1: "World", Field2: 100}

				dest, err := Map[SourceStruct, DestinationStruct](source)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if dest.FieldA != "World" {
					t.Errorf("Expected FieldA='World', got '%s'", dest.FieldA)
				}
				if dest.FieldB != 102 { // 100 + 2
					t.Errorf("Expected FieldB=102, got %d", dest.FieldB)
				}
			},
		},
		{
			name: "Reverse mapping (destination to source)",
			setup: func() {
				RegisterMapper[DestinationStruct, SourceStruct](MapDestinationStructToSourceStruct)
			},
			test: func(t *testing.T) {
				source := DestinationStruct{FieldA: "Test", FieldB: 50}

				dest, err := Map[DestinationStruct, SourceStruct](source)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if dest.Field1 != "Test" {
					t.Errorf("Expected Field1='Test', got '%s'", dest.Field1)
				}
				if dest.Field2 != 52 { // 50 + 2
					t.Errorf("Expected Field2=52, got %d", dest.Field2)
				}
			},
		},
		{
			name: "Complex struct with additional field",
			setup: func() {
				RegisterMapper[SourceStruct2, *DestinationStruct2](MapSourceStructToDestinationStruct2)
			},
			test: func(t *testing.T) {
				source := SourceStruct2{Field1: "Complex", Field2: 25}

				dest, err := Map[SourceStruct2, *DestinationStruct2](source)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if dest == nil {
					t.Fatal("Expected destination to not be nil")
				}
				if dest.FieldA != "Complex" {
					t.Errorf("Expected FieldA='Complex', got '%s'", dest.FieldA)
				}
				if dest.FieldB != 27 { // 25 + 2
					t.Errorf("Expected FieldB=27, got %d", dest.FieldB)
				}
				if dest.FieldC != "Hello Complex" {
					t.Errorf("Expected FieldC='Hello Complex', got '%s'", dest.FieldC)
				}
			},
		},
		{
			name: "Nil slice mapping",
			setup: func() {
				RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
			},
			test: func(t *testing.T) {
				var sources []SourceStruct

				destinations, err := MapList[SourceStruct, DestinationStruct](sources)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if len(destinations) != 0 {
					t.Errorf("Expected empty slice, got length %d", len(destinations))
				}
			},
		},
		{
			name: "Multiple items slice mapping",
			setup: func() {
				RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
			},
			test: func(t *testing.T) {
				sources := []SourceStruct{
					{Field1: "First", Field2: 10},
					{Field1: "Second", Field2: 20},
					{Field1: "Third", Field2: 30},
				}

				destinations, err := MapList[SourceStruct, DestinationStruct](sources)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if len(destinations) != 3 {
					t.Fatalf("Expected 3 destinations, got %d", len(destinations))
				}

				// Verify each mapped item
				if destinations[0].FieldA != "First" {
					t.Errorf("Expected first FieldA='First', got '%s'", destinations[0].FieldA)
				}
				if destinations[0].FieldB != 12 {
					t.Errorf("Expected first FieldB=12, got %d", destinations[0].FieldB)
				}

				if destinations[1].FieldA != "Second" {
					t.Errorf("Expected second FieldA='Second', got '%s'", destinations[1].FieldA)
				}
				if destinations[1].FieldB != 22 {
					t.Errorf("Expected second FieldB=22, got %d", destinations[1].FieldB)
				}

				if destinations[2].FieldA != "Third" {
					t.Errorf("Expected third FieldA='Third', got '%s'", destinations[2].FieldA)
				}
				if destinations[2].FieldB != 32 {
					t.Errorf("Expected third FieldB=32, got %d", destinations[2].FieldB)
				}
			},
		},
		{
			name: "Pointer slice mapping",
			setup: func() {
				RegisterMapper[SourceStruct, *DestinationStruct](MapSourceStructToDestinationStruct)
			},
			test: func(t *testing.T) {
				sources := []SourceStruct{
					{Field1: "Ptr1", Field2: 5},
					{Field1: "Ptr2", Field2: 15},
				}

				destinations, err := MapList[SourceStruct, *DestinationStruct](sources)

				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if len(destinations) != 2 {
					t.Fatalf("Expected 2 destinations, got %d", len(destinations))
				}

				// Verify pointers are not nil
				if destinations[0] == nil {
					t.Fatal("Expected first destination to not be nil")
				}
				if destinations[1] == nil {
					t.Fatal("Expected second destination to not be nil")
				}

				// Verify mapped values
				if destinations[0].FieldA != "Ptr1" {
					t.Errorf("Expected first FieldA='Ptr1', got '%s'", destinations[0].FieldA)
				}
				if destinations[0].FieldB != 7 {
					t.Errorf("Expected first FieldB=7, got %d", destinations[0].FieldB)
				}

				if destinations[1].FieldA != "Ptr2" {
					t.Errorf("Expected second FieldA='Ptr2', got '%s'", destinations[1].FieldA)
				}
				if destinations[1].FieldB != 17 {
					t.Errorf("Expected second FieldB=17, got %d", destinations[1].FieldB)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear registry for clean test
			ClearRegistry()

			// Setup test
			tt.setup()

			// Run test
			tt.test(t)
		})
	}
}

func TestHasMapping(t *testing.T) {
	ClearRegistry()

	// Test before registration
	if HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected no mapping to be registered")
	}

	// Register and test
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

	if !HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected mapping to be registered")
	}
}

func TestMapUnregistered(t *testing.T) {
	ClearRegistry()

	source := SourceStruct{Field1: "test", Field2: 123}
	_, err := Map[SourceStruct, DestinationStruct](source)

	if err == nil {
		t.Fatal("Expected error for unregistered mapping")
	}
	if !errors.Is(err, ErrMapperNotFound) {
		t.Errorf("Expected error to wrap ErrMapperNotFound, got '%v'", err)
	}

	expectedError := "no mapper registered for mapper.SourceStruct to mapper.DestinationStruct"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestMapListUnregistered(t *testing.T) {
	ClearRegistry()

	destinations, err := MapList[SourceStruct, DestinationStruct]([]SourceStruct{{Field1: "test", Field2: 123}})

	if !errors.Is(err, ErrMapperNotFound) {
		t.Errorf("Expected error to wrap ErrMapperNotFound, got '%v'", err)
	}
	if destinations != nil {
		t.Errorf("Expected nil slice, got %v", destinations)
	}
}

func TestMapPropagatesError(t *testing.T) {
	ClearRegistry()

	mappingErr := errors.New("mapping failed")
	RegisterMapper[SourceStruct, DestinationStruct](func(SourceStruct) (DestinationStruct, error) {
		return DestinationStruct{}, mappingErr
	})

	_, err := Map[SourceStruct, DestinationStruct](SourceStruct{})

	// Map must return the mapping function's error unchanged
	if !errors.Is(err, mappingErr) {
		t.Errorf("Expected error '%v', got '%v'", mappingErr, err)
	}
}

func TestMapListPropagatesError(t *testing.T) {
	ClearRegistry()

	mappingErr := errors.New("mapping failed")
	calls := 0
	RegisterMapper[SourceStruct, DestinationStruct](func(source SourceStruct) (DestinationStruct, error) {
		calls++
		if source.Field2 < 0 {
			return DestinationStruct{}, mappingErr
		}
		return DestinationStruct{FieldA: source.Field1, FieldB: source.Field2}, nil
	})

	sources := []SourceStruct{
		{Field1: "ok", Field2: 1},
		{Field1: "bad", Field2: -1},
		{Field1: "ok", Field2: 3},
	}

	destinations, err := MapList[SourceStruct, DestinationStruct](sources)

	if !errors.Is(err, mappingErr) {
		t.Fatalf("Expected error to wrap '%v', got '%v'", mappingErr, err)
	}
	if !strings.Contains(err.Error(), "element 1") {
		t.Errorf("Expected error to name the failing index, got '%v'", err)
	}
	if destinations != nil {
		t.Errorf("Expected nil slice on error, got %v", destinations)
	}
	// MapList must stop at the first failing element
	if calls != 2 {
		t.Errorf("Expected 2 calls to the mapping function, got %d", calls)
	}
}

func TestMapListNilAndEmpty(t *testing.T) {
	ClearRegistry()
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

	// A nil slice maps to a nil slice
	destinations, err := MapList[SourceStruct, DestinationStruct](nil)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if destinations != nil {
		t.Errorf("Expected nil slice, got %v", destinations)
	}

	// An empty slice maps to an empty, non-nil slice
	destinations, err = MapList[SourceStruct, DestinationStruct]([]SourceStruct{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if destinations == nil || len(destinations) != 0 {
		t.Errorf("Expected empty non-nil slice, got %#v", destinations)
	}
}

func TestPointerSource(t *testing.T) {
	ClearRegistry()

	RegisterMapper[*SourceStruct, DestinationStruct](func(source *SourceStruct) (DestinationStruct, error) {
		return DestinationStruct{FieldA: source.Field1, FieldB: source.Field2}, nil
	})

	// Pointer and value sources are different type pairs
	if HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected no mapping for the value source type")
	}

	dest, err := Map[*SourceStruct, DestinationStruct](&SourceStruct{Field1: "Ptr", Field2: 7})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if dest.FieldA != "Ptr" || dest.FieldB != 7 {
		t.Errorf("Expected {Ptr 7}, got %+v", dest)
	}
}

type stringerStruct struct {
	Value string
}

func (s stringerStruct) String() string {
	return s.Value
}

func TestInterfaceTypes(t *testing.T) {
	ClearRegistry()

	RegisterMapper[fmt.Stringer, string](func(source fmt.Stringer) (string, error) {
		return source.String(), nil
	})

	if !HasMapping[fmt.Stringer, string]() {
		t.Fatal("Expected mapping to be registered")
	}

	dest, err := Map[fmt.Stringer, string](stringerStruct{Value: "Hello"})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if dest != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", dest)
	}

	// The mapping is keyed by the interface type, not the concrete type
	if HasMapping[stringerStruct, string]() {
		t.Error("Expected no mapping for the concrete type")
	}
}

func registerFirstLocalType() func() (int, error) {
	type local struct{}

	RegisterMapper[local, int](func(local) (int, error) { return 1, nil })

	return func() (int, error) { return Map[local, int](local{}) }
}

func registerSecondLocalType() func() (int, error) {
	type local struct{}

	RegisterMapper[local, int](func(local) (int, error) { return 2, nil })

	return func() (int, error) { return Map[local, int](local{}) }
}

// Distinct types that share a name must not overwrite each other's mapping.
func TestSameNameTypesDoNotCollide(t *testing.T) {
	ClearRegistry()

	mapFirst := registerFirstLocalType()
	mapSecond := registerSecondLocalType()

	first, err := mapFirst()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	second, err := mapSecond()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if first != 1 || second != 2 {
		t.Errorf("Expected 1 and 2, got %d and %d", first, second)
	}
}

func TestRegisterMapperReplaces(t *testing.T) {
	ClearRegistry()

	RegisterMapper[SourceStruct, int](func(SourceStruct) (int, error) { return 1, nil })
	RegisterMapper[SourceStruct, int](func(SourceStruct) (int, error) { return 2, nil })

	dest, err := Map[SourceStruct, int](SourceStruct{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if dest != 2 {
		t.Errorf("Expected the second registration to win, got %d", dest)
	}
}

func TestRegisterMapperNilPanics(t *testing.T) {
	ClearRegistry()

	defer func() {
		if recover() == nil {
			t.Error("Expected RegisterMapper to panic for a nil mapping function")
		}
		if HasMapping[SourceStruct, DestinationStruct]() {
			t.Error("Expected no mapping to be registered")
		}
	}()

	RegisterMapper[SourceStruct, DestinationStruct](nil)
}

// Run with -race to verify the registry is safe for concurrent use.
func TestConcurrentAccess(t *testing.T) {
	ClearRegistry()

	const goroutines = 8
	const iterations = 500

	sources := []SourceStruct{{Field1: "a", Field2: 1}, {Field1: "b", Field2: 2}}

	var wg sync.WaitGroup
	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for i := 0; i < iterations; i++ {
				RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

				// Another goroutine may have unregistered the mapping in between
				if _, err := Map[SourceStruct, DestinationStruct](sources[0]); err != nil && !errors.Is(err, ErrMapperNotFound) {
					t.Errorf("Unexpected error: %v", err)
				}
				if _, err := MapList[SourceStruct, DestinationStruct](sources); err != nil && !errors.Is(err, ErrMapperNotFound) {
					t.Errorf("Unexpected error: %v", err)
				}

				HasMapping[SourceStruct, DestinationStruct]()
				UnregisterMapper[SourceStruct, DestinationStruct]()
			}
		}()
	}
	wg.Wait()
}

func TestUnregisterMapper(t *testing.T) {
	ClearRegistry()

	// Register a mapping first
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
	RegisterMapper[SourceStruct, *DestinationStruct](MapSourceStructToDestinationStruct)

	// Verify mapping exists
	if !HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected mapping to be registered")
	}

	// Test unregistration
	removed := UnregisterMapper[SourceStruct, DestinationStruct]()
	if !removed {
		t.Error("Expected unregistration to return true")
	}

	// Verify mapping no longer exists
	if HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected mapping to be unregistered")
	}

	// Test mapping after unregistration fails
	source := SourceStruct{Field1: "test", Field2: 123}
	_, err := Map[SourceStruct, DestinationStruct](source)
	if !errors.Is(err, ErrMapperNotFound) {
		t.Errorf("Expected error to wrap ErrMapperNotFound, got '%v'", err)
	}

	// Other type pairs are left untouched
	if !HasMapping[SourceStruct, *DestinationStruct]() {
		t.Error("Expected the pointer destination mapping to still be registered")
	}

	// Test unregistering non-existent mapping returns false
	removedAgain := UnregisterMapper[SourceStruct, DestinationStruct]()
	if removedAgain {
		t.Error("Expected second unregistration to return false")
	}
}

func TestClearRegistry(t *testing.T) {
	ClearRegistry()

	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
	RegisterMapper[DestinationStruct, SourceStruct](MapDestinationStructToSourceStruct)

	ClearRegistry()

	if HasMapping[SourceStruct, DestinationStruct]() || HasMapping[DestinationStruct, SourceStruct]() {
		t.Error("Expected all mappings to be removed")
	}

	// The registry stays usable after being cleared
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)
	if !HasMapping[SourceStruct, DestinationStruct]() {
		t.Error("Expected mapping to be registered after clearing")
	}
}

var (
	sinkDest     DestinationStruct
	sinkDestList []DestinationStruct
)

func BenchmarkMap(b *testing.B) {
	ClearRegistry()
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

	source := SourceStruct{Field1: "benchmark", Field2: 42}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkDest, _ = Map[SourceStruct, DestinationStruct](source)
	}
}

func BenchmarkMapParallel(b *testing.B) {
	ClearRegistry()
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

	source := SourceStruct{Field1: "benchmark", Field2: 42}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var dest DestinationStruct
		for pb.Next() {
			dest, _ = Map[SourceStruct, DestinationStruct](source)
		}
		_ = dest
	})
}

func BenchmarkMapList(b *testing.B) {
	ClearRegistry()
	RegisterMapper[SourceStruct, DestinationStruct](MapSourceStructToDestinationStructNoPointer)

	sourceList := make([]SourceStruct, 1000)
	for i := range sourceList {
		sourceList[i] = SourceStruct{Field1: "benchmark", Field2: i}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sinkDestList, _ = MapList[SourceStruct, DestinationStruct](sourceList)
	}
}
