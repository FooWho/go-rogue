// Package sparseset is an implementation of a Sparse Set for the Entity
// Component System.
//
// Set is a generational sparse set that stores components of type T in packed,
// contiguous memory for cache-friendly iteration and O(1) lookup, insertion,
// and removal.
//
// A Set consists of three slices:
//   - dense ([]T): A packed, contiguous slice holding the actual component
//     values at indices 0..[Set.Len]()-1.
//   - denseID ([][Entity]): A packed slice parallel to dense (denseID[i] is the
//     full generational Entity owning dense[i]). Used both when iterating over
//     components and to update sparse during O(1) swap-and-pop removals.
//   - sparse ([]int): Indexed by [Entity.Index]() (the lower 16 bits of an
//     Entity), mapping an entity's index slot to its location in dense, or -1
//     if the entity does not have this component.
//
// Entity handles are 32-bit unsigned integers with two parts. The lower 16-bits
// are the index and the upper 16 bits are the generation. This limits the length
// of the sparse slice to 65,536 elements (rather than 2^32 elements). The
// generational part of the handle means we don't need to worry about one
// entity referencing a stale index to another element. If we have a goblin
// with index 5 in generation 1, and it dies, we don't have to worry
// about the index being released back to the system for reuse. Any handle
// to that original goblin will be index 5, generation 1. When index 5 is reused,
// it will be index 5, generation 2, and so on. We won't run into an index
// collision unless some very long-lived entity holds onto a stale index for
// 64k generations.
//
// Example:
// Suppose three entities exist at indices 1, 2, and 500. Entities 1 and 500
// have Velocity, while entity 2 does not:
//   - dense has length 2: [vel_1, vel_500] (at indices 0 and 1).
//   - denseID has length 2: [entity_1, entity_500] (at indices 0 and 1).
//   - sparse has length >= 501: sparse[1] = 0, sparse[500] = 1, and all other
//     slots (including sparse[0] and sparse[2..499]) are -1 (assuming no other
//     entities with a velocity exist in the world).
package sparseset
