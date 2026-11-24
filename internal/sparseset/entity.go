package sparseset

import "math"

// Entity is a 32 bit unsigned integer used as a handle to an entity. The
// lower 16 bits are the entity index. The upper 16 bits are the generation.
// The first generation is generation 1. There is no generation 0. This
// preserves 0 as the [NullEntity]. An [Entity] is valid as long as the handle
// is not the [NullEntity]. An [Entity] is alive if we still need to be
// tracking its components. An [Entity] is dead if we no longer need to track
// its components and we want to reclaim the handle. An [Entity] is stale
// when some other [Entity] is holding onto a handle that has either been
// returned to the pool or has been returned to the pool and subsequently
// reassigned.
type Entity uint32

const (
	NullEntity Entity = 0

	indexBits = 16
	indexMask = (1 << indexBits) - 1
)

// NewEntity generates an entity handle with the specified index and generation.
func NewEntity(index, gen uint16) Entity {
	return Entity(uint32(gen)<<indexBits | uint32(index))
}

// Index returns the index (the lower 16 bits) of an entity handle.
func (e Entity) Index() uint16 {
	return uint16(e & indexMask)
}

// Generation returns the generation (the upper 16 bits) of an entity handle.
func (e Entity) Generation() uint16 {
	return uint16(e >> indexBits)
}

// IsValid returns true if the entity handle e is not the [NullEntity]. This
// doesn't mean that this handle refers to a "live" entity, just that the
// handle is not null.
func (e Entity) IsValid() bool {
	return e != NullEntity
}

// EntityPool is the ID allocator and lifecycle manager. It creates unique
// Entity handles when objects are created, recycles those index slots when
// objects are destroyed, and checks whether a handle still refers to a
// living entity.
type EntityPool struct {
	generations []uint16
	freeIndices []uint16
}

// NewEntityPool initializes and returns a pointer to a new EntityPool
// with the specified initial capacity.
func NewEntityPool(initialCap int) *EntityPool {
	return &EntityPool{
		generations: make([]uint16, 0, initialCap),
		freeIndices: make([]uint16, 0, initialCap),
	}
}

// Spawn generates a new [Entity] from the pool. It will panic if the
// pool of entity handles has been exhausted. This will only occur if there
// are already 65,536 living (or simultaneously active) entities and we try
// to [EntityPool.Spawn] a new entity.
func (p *EntityPool) Spawn() Entity {
	if len(p.freeIndices) > 0 {
		idx := p.freeIndices[0]
		copy(p.freeIndices, p.freeIndices[1:])
		p.freeIndices = p.freeIndices[:len(p.freeIndices)-1]

		return NewEntity(idx, p.generations[idx])
	}

	if len(p.generations) > math.MaxUint16 {
		panic("entity pool exhausted")
	}

	idx := uint16(len(p.generations))
	gen := uint16(1)
	p.generations = append(p.generations, gen)
	return NewEntity(idx, gen)
}

// Free is called at the end of the [Entity] life cycle so that the
// [EntityPool] can reclaim the handle to be used again.
func (p *EntityPool) Free(e Entity) {
	if !p.IsAlive(e) {
		return
	}
	idx := e.Index()
	p.generations[idx]++
	if p.generations[idx] == 0 {
		p.generations[idx] = 1
	}
	p.freeIndices = append(p.freeIndices, idx)
}

// IsAlive returns true if the given [Entity] is currently alive and
// false otherwise.
func (p *EntityPool) IsAlive(e Entity) bool {
	idx := int(e.Index())
	return e != NullEntity && idx < len(p.generations) && p.generations[idx] == e.Generation()
}
