package main

import (
	g "./graphics"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"math"
)

const (
	fully     = 2
	partially = 1
	nothing   = 0

	printOccluderTime = false
	printVision       = false
)

var (
	frustum [6][4]float32

	viewable    [worldChunkAll]*Chunk
	viewableNum int

	visionGoto      [worldChunkAll]bool
	visionQueue     [worldChunkAll]*Chunk
	visionQueueFrom [worldChunkAll]int
	visionQueuePos  int
	visionQueueNum  int
)

func ReadyFrustum() {
	// left
	frustum[0][0] = g.Mvp[3] + g.Mvp[0]
	frustum[0][1] = g.Mvp[7] + g.Mvp[4]
	frustum[0][2] = g.Mvp[11] + g.Mvp[8]
	frustum[0][3] = g.Mvp[15] + g.Mvp[12]
	normalizePlane(0)

	// right
	frustum[1][0] = g.Mvp[3] - g.Mvp[0]
	frustum[1][1] = g.Mvp[7] - g.Mvp[4]
	frustum[1][2] = g.Mvp[11] - g.Mvp[8]
	frustum[1][3] = g.Mvp[15] - g.Mvp[12]
	normalizePlane(1)

	// top
	frustum[2][0] = g.Mvp[3] - g.Mvp[1]
	frustum[2][1] = g.Mvp[7] - g.Mvp[5]
	frustum[2][2] = g.Mvp[11] - g.Mvp[9]
	frustum[2][3] = g.Mvp[15] - g.Mvp[13]
	normalizePlane(2)

	// bottom
	frustum[3][0] = g.Mvp[3] + g.Mvp[1]
	frustum[3][1] = g.Mvp[7] + g.Mvp[5]
	frustum[3][2] = g.Mvp[11] + g.Mvp[9]
	frustum[3][3] = g.Mvp[15] + g.Mvp[13]
	normalizePlane(3)

	// near
	frustum[4][0] = g.Mvp[3] + g.Mvp[2]
	frustum[4][1] = g.Mvp[7] + g.Mvp[6]
	frustum[4][2] = g.Mvp[11] + g.Mvp[10]
	frustum[4][3] = g.Mvp[15] + g.Mvp[14]
	normalizePlane(4)

	// far
	frustum[5][0] = g.Mvp[3] - g.Mvp[2]
	frustum[5][1] = g.Mvp[7] - g.Mvp[6]
	frustum[5][2] = g.Mvp[11] - g.Mvp[10]
	frustum[5][3] = g.Mvp[15] - g.Mvp[14]
	normalizePlane(5)
}

func normalizePlane(i int) {
	n := float32(math.Sqrt(float64(frustum[i][0]*frustum[i][0] + frustum[i][1]*frustum[i][1] + frustum[i][2]*frustum[i][2])))
	frustum[i][0] /= n
	frustum[i][1] /= n
	frustum[i][2] /= n
	frustum[i][3] /= n
}

func Occlude(lx, ly, lz int) {

	timeA := glfw.GetTime()

	viewableNum = 0

	index := lx + (ly << worldShiftX) + (lz << worldShiftSlice)
	if toggleMode || index < 0 || index >= worldChunkAll {
		for viewableNum < worldChunkAll {
			viewable[viewableNum] = &worldChunks[viewableNum]
			viewableNum++
		}
		return
	}

	visionQueuePos = 0
	visionQueue[0] = &worldChunks[index]
	visionQueueFrom[0] = -1
	visionQueueNum = 1

	for i := 0; i < worldChunkAll; i++ {
		visionGoto[i] = true
	}
	for visionQueueNum > 0 {
		B := visionQueue[visionQueuePos]
		from := visionQueueFrom[visionQueuePos]

		viewable[viewableNum] = B
		viewableNum++

		visionQueuePos++
		if visionQueuePos == worldChunkAll {
			visionQueuePos = 0
		}
		visionQueueNum--

		if from != NegX {
			occluderVisit(from, B, PosX)
		}
		if from != PosX {
			occluderVisit(from, B, NegX)
		}
		if from != NegY {
			occluderVisit(from, B, PosY)
		}
		if from != PosY {
			occluderVisit(from, B, NegY)
		}
		if from != NegZ {
			occluderVisit(from, B, PosZ)
		}
		if from != PosZ {
			occluderVisit(from, B, NegZ)
		}
	}

	if printOccluderTime {
		timeB := glfw.GetTime()
		fmt.Println("occluder", timeB-timeA)
	}
}

func occluderVisit(from int, B *Chunk, to int) {
	x := B.x
	y := B.y
	z := B.z
	switch to {
	case PosX:
		x++
		if x == worldChunksX {
			return
		}
	case NegX:
		x--
		if x == -1 {
			return
		}
	case PosY:
		y++
		if y == worldChunksY {
			return
		}
	case NegY:
		y--
		if y == -1 {
			return
		}
	case PosZ:
		z++
		if z == worldChunksZ {
			return
		}
	case NegZ:
		z--
		if z == -1 {
			return
		}
	}
	if printVision {
		fmt.Println("looking at", x, y, z, "from", from, "to", to)
	}
	index := x + (y << worldShiftX) + (z << worldShiftSlice)
	if visionGoto[index] == false {
		return
	}
	if from >= 0 {
		switch from {
		case PosX:
			from = NegX
		case NegX:
			from = PosX
		case PosY:
			from = NegY
		case NegY:
			from = PosY
		case PosZ:
			from = NegZ
		case NegZ:
			from = PosZ
		}
		var sideA, sideB uint
		if from < to {
			sideA = uint(from)
			sideB = uint(to)
		} else {
			sideA = uint(to)
			sideB = uint(from)
		}
		if B.visibility&(1<<(sideA*6+sideB)) == 0 {
			return
		}
	}
	visionGoto[index] = false
	C := &worldChunks[index]
	box := inBox(
		float32((C.x+1)<<chunkShift),
		float32((C.y+1)<<chunkShift),
		float32((C.z+1)<<chunkShift),
		float32(C.x<<chunkShift),
		float32(C.y<<chunkShift),
		float32(C.z<<chunkShift))
	if box == nothing {
		return
	}
	queue := visionQueuePos + visionQueueNum
	if queue >= worldChunkAll {
		queue -= worldChunkAll
	}
	visionQueue[queue] = C
	visionQueueFrom[queue] = to
	visionQueueNum++
}

func inBox(posx, posy, posz, negx, negy, negz float32) int {
	var pvx, pvy, pvz float32
	var nvx, nvy, nvz float32
	result := fully
	for i := 0; i < 6; i++ {
		plane := frustum[i]
		if plane[0] > 0 {
			pvx = posx
			nvx = negx
		} else {
			pvx = negx
			nvx = posx
		}
		if plane[1] > 0 {
			pvy = posy
			nvy = negy
		} else {
			pvy = negy
			nvy = posy
		}
		if plane[2] > 0 {
			pvz = posz
			nvz = negz
		} else {
			pvz = negz
			nvz = posz
		}
		if pvx*plane[0]+pvy*plane[1]+pvz*plane[2]+plane[3] < 0 {
			return nothing
		}
		if nvx*plane[0]+nvy*plane[1]+nvz*plane[2]+plane[3] < 0 {
			result = partially
		}
	}
	return result
}
