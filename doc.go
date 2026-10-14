// Package mapper provides type-safe struct mapping with custom transformation functions.
//
// Mapper is a high-performance, thread-safe library for mapping between Go structs
// using custom transformation functions. It uses generics for compile-time type safety
// and a registry pattern for flexible mapping registration.
//
// # Basic Usage
//
//	// Define mapping function
//	func MapUserToDTO(user User) (UserDTO, error) {
//		return UserDTO{
//			Name: user.FirstName + " " + user.LastName,
//			Age:  user.Age,
//		}, nil
//	}
//
//	// Register the mapping
//	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)
//
//	// Use the mapping
//	user := User{FirstName: "John", LastName: "Doe", Age: 30}
//	dto, err := mapper.Map[User, UserDTO](user)
//
// # Dynamic Mapping Management
//
// The library supports dynamic registration and deregistration of mappings:
//
//	// Register a mapping
//	mapper.RegisterMapper[User, UserDTO](MapUserToDTO)
//
//	// Check if mapping exists
//	if mapper.HasMapping[User, UserDTO]() {
//		// mapping is available
//	}
//
//	// Unregister a specific mapping
//	removed := mapper.UnregisterMapper[User, UserDTO]()
//
//	// Clear all mappings (useful for testing)
//	mapper.ClearRegistry()
//
// # Thread Safety
//
// All mapping operations are thread-safe. The registry uses sync.RWMutex to allow
// concurrent read operations while ensuring safe writes during registration and deregistration.
//
// # Performance
//
// Mapper is optimized for performance:
//   - Your mapping function is called directly; reflection is only used to
//     look up the type pair in the registry, never to copy fields
//   - No allocations in Map beyond those of your mapping function
//   - List mapping with a single registry lookup and a preallocated result
//
// # Error Handling
//
// Mapping functions should return errors for any validation failures.
// Map returns these errors unchanged. MapList wraps them with the index of
// the failing element, so use errors.Is or errors.As to inspect them.
//
// When no mapping function is registered for a type pair, Map and MapList
// return an error that wraps ErrMapperNotFound.
//
// # Type Safety
//
// The generic type parameters ensure compile-time type safety.
// Type mismatches are caught at compile time, not runtime. Whether a mapping
// is registered for a type pair is only known at runtime.
package mapper
