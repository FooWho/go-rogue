package sparseset

import "math"

type Entity uint32

const (
	NullEntity Entity = 0

	indexBits = 16
	indexMask = (1 << indexBits) - 1
)

func NewEntity(index, gen uint16) Entity {
	return Entity(uint32(gen)<<indexBits | uint32(index))
}

func (e Entity) Index() uint16 {
	return uint16(e & indexMask)
}

func (e Entity) Generation() uint16 {
	return uint16(e >> indexBits)
}

func (e Entity) IsValid() bool {
	return e != NullEntity
}

type EntityPool struct {
	generations []uint16
	freeIndices []uint16
}

func NewEntityPool(initialCap int) *EntityPool {
	return &EntityPool{
		generations: make([]uint16, 0, initialCap),
		freeIndices: make([]uint16, 0, initialCap),
	}
}

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

func (p *EntityPool) IsAlive(e Entity) bool {
	idx := int(e.Index())
	return e != NullEntity && idx < len(p.generations) && p.generations[idx] == e.Generation()
}
