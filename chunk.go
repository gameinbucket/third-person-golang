package main

import (
	"math/rand"
	"unsafe"

	g "./graphics"
)

const (
	chunkShift      = 4
	chunkSize       = 1 << chunkShift
	chunkShiftSlice = chunkShift + chunkShift
	chunkSlice      = chunkSize * chunkSize
	chunkAll        = chunkSlice * chunkSize
	PosX            = 0
	NegX            = 1
	PosY            = 2
	NegY            = 3
	PosZ            = 4
	NegZ            = 5
	Sides           = 6
)

type Chunk struct {
	Blocks     [chunkAll]Block
	Mesh       *g.Buffer
	beginSide  [Sides]unsafe.Pointer
	countSide  [Sides]int32
	visibility uint32
	x          int
	y          int
	z          int
}

func (o *Chunk) Init(x, y, z int) {
	o.x = x
	o.y = y
	o.z = z
	bx := 0
	by := 0
	bz := 0
	for i := 0; i < chunkAll; i++ {
		if y == 1 && x != 0 && z != 0 {
			o.Blocks[i].Type = BlockNone
			continue
		}
		if by == 0 {
			o.Blocks[i].Type = BlockGrass
		} else if bx >= 2 && bx <= chunkSize-2 && bz >= 2 && bz <= chunkSize-2 {
			o.Blocks[i].Type = rand.Int31n(3)
		} else {
			o.Blocks[i].Type = BlockNone
		}
		bx++
		if bx == chunkSize {
			bx = 0
			by++
			if by == chunkSize {
				by = 0
				bz++
			}
		}
	}

	/*o.Blocks[(chunkSize>>1)+1+((1)<<chunkShift)+((chunkSize>>1)<<chunkShiftSlice)].Type = BlockStone
	for i := 1; i < 6; i++ {
		o.Blocks[(chunkSize>>1)+2+((i)<<chunkShift)+((chunkSize>>1)<<chunkShiftSlice)].Type = BlockStone
	}

	o.Blocks[(chunkSize>>1)+((1)<<chunkShift)+((chunkSize>>1)<<chunkShiftSlice)].Light = TorchLight*/

	o.Blocks[(1)+((1)<<chunkShift)+((1)<<chunkShiftSlice)].Light = packRgb(255, 230, 200)
}

func (o *Chunk) GetBlockUnsafe(x, y, z int) int32 {
	return o.Blocks[x+(y<<chunkShift)+(z<<chunkShiftSlice)].Type
}

func (o *Chunk) GetColorUnsafe(x, y, z int) int32 {
	return o.Blocks[x+(y<<chunkShift)+(z<<chunkShiftSlice)].Color
}

func (o *Chunk) GetBlockPointerUnsafe(x, y, z int) *Block {
	return &o.Blocks[x+(y<<chunkShift)+(z<<chunkShiftSlice)]
}
