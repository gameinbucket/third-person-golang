package main

import (
	g "./graphics"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"sync"
	"unsafe"
)

const (
	worldShiftX     = 2
	worldShiftY     = 1
	worldShiftZ     = 2
	worldShiftSlice = worldShiftX + worldShiftY
	worldChunksX    = 1 << worldShiftX
	worldChunksY    = 1 << worldShiftY
	worldChunksZ    = 1 << worldShiftZ
	worldChunkSlice = worldChunksX * worldChunksY
	worldChunkAll   = worldChunkSlice * worldChunksZ

	worldBlocksX = worldChunksX * chunkSize
	worldBlocksY = worldChunksY * chunkSize
	worldBlocksZ = worldChunksZ * chunkSize

	meshColorSize  = chunkSize + 1
	meshColorSlice = meshColorSize * meshColorSize

	printDebug           = false
	printWorldRenderTime = false
)

var (
	worldChunks [worldChunkAll]Chunk

	sliceX         = [Sides]int{2, 2, 1, 1, 0, 0}
	sliceY         = [Sides]int{0, 0, 2, 2, 1, 1}
	sliceZ         = [Sides]int{1, 1, 0, 0, 2, 2}
	sliceTowards   = [Sides]int{1, -1, 1, -1, 1, -1}
	mask           [chunkSlice]bool
	meshColor      [meshColorSize * meshColorSlice]int32
	meshAmbient    [chunkAll][Sides][4]byte
	slice          [3]int
	sliceTemp      [3]int
	sliceTemp2     [3]int
	meshBuffer     = g.NewBufferUtil(chunkAll*Sides*4*(3+3+3), chunkAll*Sides*6)
	meshAttributes = []g.Attribute{
		{Typ: gl.FLOAT, Size: 3},
		{Typ: gl.FLOAT, Size: 3},
		{Typ: gl.FLOAT, Size: 3}}
)

func MakeWorld() {
	cx := 0
	cy := 0
	cz := 0
	for i := 0; i < worldChunkAll; i++ {
		worldChunks[i].Init(cx, cy, cz)
		cx++
		if cx == worldChunksX {
			cx = 0
			cy++
			if cy == worldChunksY {
				cy = 0
				cz++
			}
		}
	}
}

func WorldGetBlockUnsafe(cx, cy, cz, bx, by, bz int) int32 {
	obj := worldChunks[cx+(cy<<worldShiftX)+(cz<<worldShiftSlice)]
	return obj.Blocks[bx+(by<<chunkShift)+(bz<<chunkShiftSlice)].Type
}

func WorldGetChunk(x, y, z int) *Chunk {
	if x < 0 || x >= worldChunksX {
		return nil
	}
	if y < 0 || y >= worldChunksY {
		return nil
	}
	if z < 0 || z >= worldChunksZ {
		return nil
	}
	return &worldChunks[x+(y<<worldShiftX)+(z<<worldShiftSlice)]
}

func WorldGetBlock(cx, cy, cz, bx, by, bz int) int32 {
	for bx < 0 {
		bx += chunkSize
		cx--
	}
	for bx >= chunkSize {
		bx -= chunkSize
		cx++
	}
	for by < 0 {
		by += chunkSize
		cy--
	}
	for by >= chunkSize {
		by -= chunkSize
		cy++
	}
	for bz < 0 {
		bz += chunkSize
		cz--
	}
	for bz >= chunkSize {
		bz -= chunkSize
		cz++
	}
	obj := WorldGetChunk(cx, cy, cz)
	if obj == nil {
		return BlockNone
	}
	return obj.GetBlockUnsafe(bx, by, bz)
}

func WorldGetColor(cx, cy, cz, bx, by, bz int) int32 {
	for bx < 0 {
		bx += chunkSize
		cx--
	}
	for bx >= chunkSize {
		bx -= chunkSize
		cx++
	}
	for by < 0 {
		by += chunkSize
		cy--
	}
	for by >= chunkSize {
		by -= chunkSize
		cy++
	}
	for bz < 0 {
		bz += chunkSize
		cz--
	}
	for bz >= chunkSize {
		bz -= chunkSize
		cz++
	}
	obj := WorldGetChunk(cx, cy, cz)
	if obj == nil {
		return 0
	}
	return obj.GetColorUnsafe(bx, by, bz)
}

func WorldGetBlockPointer(cx, cy, cz, bx, by, bz int) *Block {
	for bx < 0 {
		bx += chunkSize
		cx--
	}
	for bx >= chunkSize {
		bx -= chunkSize
		cx++
	}
	for by < 0 {
		by += chunkSize
		cy--
	}
	for by >= chunkSize {
		by -= chunkSize
		cy++
	}
	for bz < 0 {
		bz += chunkSize
		cz--
	}
	for bz >= chunkSize {
		bz -= chunkSize
		cz++
	}
	obj := WorldGetChunk(cx, cy, cz)
	if obj == nil {
		return nil
	}
	return obj.GetBlockPointerUnsafe(bx, by, bz)
}

func BuildWorld() {

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for i := 0; i < worldChunkAll; i++ {
			worldChunks[i].calculateOcclusion()
		}
		fmt.Println("occlusion done")
		wg.Done()
	}()

	GiveSunlight()
	for i := 0; i < worldChunkAll; i++ {
		worldChunks[i].processLights()
	}
	fmt.Println("lighting finished")
	for i := 0; i < worldChunkAll; i++ {
		worldChunks[i].buildMesh()
	}
	fmt.Println("meshing done")

	wg.Wait()
}

func (obj *Chunk) buildMesh() {
	obj.colorMesh()
	meshBuffer.Zero()
	for side := 0; side < Sides; side++ {
		firstIndice := meshBuffer.IndiceIndex
		ptrx := sliceX[side]
		ptry := sliceY[side]
		ptrz := sliceZ[side]
		toward := sliceTowards[side]
		for slice[2] = 0; slice[2] < chunkSize; slice[2]++ {
			for i := 0; i < chunkSlice; i++ {
				mask[i] = false
			}
			for slice[1] = 0; slice[1] < chunkSize; slice[1]++ {
				for slice[0] = 0; slice[0] < chunkSize; slice[0]++ {
					if mask[slice[0]+(slice[1]<<chunkShift)] {
						continue
					}
					typ := obj.GetBlockUnsafe(slice[ptrx], slice[ptry], slice[ptrz])
					if typ == BlockNone {
						continue
					}
					sliceTemp[0] = slice[0]
					sliceTemp[1] = slice[1]
					sliceTemp[2] = slice[2] + toward
					if BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, sliceTemp[ptrx], sliceTemp[ptry], sliceTemp[ptrz])) {
						continue
					}
					maskWidth := 1
					maskHeight := 1
					xs := slice[ptrx]
					ys := slice[ptry]
					zs := slice[ptrz]
					ambient := meshAmbient[xs+(ys<<chunkShift)+(zs<<chunkShiftSlice)][side]
					light0, light1, light2, light3 := lightOfSide(xs, ys, zs, side)
					if ambient[0] == ambient[1] &&
						ambient[0] == ambient[2] &&
						ambient[0] == ambient[3] &&
						light0 == light1 &&
						light0 == light2 &&
						light0 == light3 {
						limitWidth := chunkSize - slice[0]
						limitHeight := chunkSize - slice[1]
						for maskWidth < limitWidth {
							sliceTemp[0] = slice[0] + maskWidth
							sliceTemp[1] = slice[1]
							if mask[sliceTemp[0]+(sliceTemp[1]<<chunkShift)] {
								break
							}
							sliceTemp[2] = slice[2]
							if obj.GetBlockUnsafe(sliceTemp[ptrx], sliceTemp[ptry], sliceTemp[ptrz]) != typ {
								break
							}
							xsb := sliceTemp[ptrx]
							ysb := sliceTemp[ptry]
							zsb := sliceTemp[ptrz]
							ambientB := meshAmbient[xsb+(ysb<<chunkShift)+(zsb<<chunkShiftSlice)][side]
							if ambientB[0] != ambientB[1] || ambientB[0] != ambientB[2] || ambientB[0] != ambientB[3] || ambient[0] != ambientB[0] {
								break
							}
							light0b, light1b, light2b, light3b := lightOfSide(xsb, ysb, zsb, side)
							if light0b != light1b || light0b != light2b || light0b != light3b || light0 != light0b {
								break
							}
							sliceTemp2[0] = sliceTemp[0]
							sliceTemp2[1] = sliceTemp[1]
							sliceTemp2[2] = sliceTemp[2] + toward
							if BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, sliceTemp2[ptrx], sliceTemp2[ptry], sliceTemp2[ptrz])) {
								break
							}
							maskWidth++
						}
					labelRectangle:
						for maskHeight < limitHeight {
							for inner := 0; inner < maskWidth; inner++ {
								sliceTemp[0] = slice[0] + inner
								sliceTemp[1] = slice[1] + maskHeight
								if mask[sliceTemp[0]+(sliceTemp[1]<<chunkShift)] {
									break labelRectangle
								}
								sliceTemp[2] = slice[2]
								if obj.GetBlockUnsafe(sliceTemp[ptrx], sliceTemp[ptry], sliceTemp[ptrz]) != typ {
									break labelRectangle
								}
								xsb := sliceTemp[ptrx]
								ysb := sliceTemp[ptry]
								zsb := sliceTemp[ptrz]
								ambientB := meshAmbient[xsb+(ysb<<chunkShift)+(zsb<<chunkShiftSlice)][side]
								if ambientB[0] != ambientB[1] || ambientB[0] != ambientB[2] || ambientB[0] != ambientB[3] || ambient[0] != ambientB[0] {
									break labelRectangle
								}
								light0b, light1b, light2b, light3b := lightOfSide(xsb, ysb, zsb, side)
								if light0b != light1b || light0b != light2b || light0b != light3b || light0 != light0b {
									break labelRectangle
								}
								sliceTemp2[0] = sliceTemp[0]
								sliceTemp2[1] = sliceTemp[1]
								sliceTemp2[2] = sliceTemp[2] + toward
								if BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, sliceTemp2[ptrx], sliceTemp2[ptry], sliceTemp2[ptrz])) {
									break labelRectangle
								}
							}
							maskHeight++
						}
						for w := 0; w < maskWidth; w++ {
							for h := 0; h < maskHeight; h++ {
								mask[slice[0]+w+((slice[1]+h)<<chunkShift)] = true
							}
						}
					}
					red0, green0, blue0 := colorizeVertex(light0, float32(ambient[0])/255.0)
					red1, green1, blue1 := colorizeVertex(light1, float32(ambient[1])/255.0)
					red2, green2, blue2 := colorizeVertex(light2, float32(ambient[2])/255.0)
					red3, green3, blue3 := colorizeVertex(light3, float32(ambient[3])/255.0)
					// rgb0 := packRgb(int32(redf), int32(greenf), int32(bluef)) // send only this to shader int32
					meshBuffer.StaticRenderCube(
						side,
						float32(xs+(obj.x<<chunkShift)), float32(ys+(obj.y<<chunkShift)), float32(zs+(obj.z<<chunkShift)),
						float32(maskWidth), float32(maskHeight), BlockTextureIndex(typ),
						red0, green0, blue0,
						red1, green1, blue1,
						red2, green2, blue2,
						red3, green3, blue3)
					slice[0] += maskWidth - 1
				}
			}
		}
		obj.countSide[side] = int32(meshBuffer.IndiceIndex - firstIndice)
		obj.beginSide[side] = unsafe.Pointer(uintptr(firstIndice) * g.Bytes4)
	}
	obj.Mesh = g.NewBufferCopied(meshAttributes, meshBuffer)
}

func lightOfSide(xs, ys, zs, side int) (int32, int32, int32, int32) {
	switch side {
	case PosX:
		return meshColor[xs+1+ys*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+1+ys*meshColorSize+(zs+1)*meshColorSlice]
	case NegX:
		return meshColor[xs+ys*meshColorSize+zs*meshColorSlice],
			meshColor[xs+ys*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+(ys+1)*meshColorSize+zs*meshColorSlice]
	case PosY:
		return meshColor[xs+(ys+1)*meshColorSize+zs*meshColorSlice],
			meshColor[xs+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+zs*meshColorSlice]
	case NegY:
		return meshColor[xs+ys*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+ys*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+ys*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+ys*meshColorSize+(zs+1)*meshColorSlice]
	case PosZ:
		return meshColor[xs+1+ys*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+(ys+1)*meshColorSize+(zs+1)*meshColorSlice],
			meshColor[xs+ys*meshColorSize+(zs+1)*meshColorSlice]
	default:
		return meshColor[xs+ys*meshColorSize+zs*meshColorSlice],
			meshColor[xs+(ys+1)*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+(ys+1)*meshColorSize+zs*meshColorSlice],
			meshColor[xs+1+ys*meshColorSize+zs*meshColorSlice]
	}
}

func (obj *Chunk) colorMesh() {
	for bz := 0; bz < chunkSize; bz++ {
		for by := 0; by < chunkSize; by++ {
			for bx := 0; bx < chunkSize; bx++ {
				index := bx + (by << chunkShift) + (bz << chunkShiftSlice)
				if obj.Blocks[index].Type == BlockNone {
					continue
				}
				aoMMZ := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by-1, bz))
				aoMMM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by-1, bz-1))
				aoMMP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by-1, bz+1))
				aoMZP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by, bz+1))
				aoMZM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by, bz-1))
				aoMPZ := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by+1, bz))
				aoMPP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by+1, bz+1))
				aoMPM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx-1, by+1, bz-1))
				aoZPP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx, by+1, bz+1))
				aoZMP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx, by-1, bz+1))
				aoZPM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx, by+1, bz-1))
				aoZMM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx, by-1, bz-1))
				aoPPZ := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by+1, bz))
				aoPMZ := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by-1, bz))
				aoPZP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by, bz+1))
				aoPZM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by, bz-1))
				aoPMM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by-1, bz-1))
				aoPPM := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by+1, bz-1))
				aoPPP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by+1, bz+1))
				aoPMP := BlockClosed(WorldGetBlock(obj.x, obj.y, obj.z, bx+1, by-1, bz+1))
				/* pos x */
				meshAmbient[index][PosX][0] = ambientOcclusion(aoPMZ, aoPZM, aoPMM)
				meshAmbient[index][PosX][1] = ambientOcclusion(aoPPZ, aoPZM, aoPPM)
				meshAmbient[index][PosX][2] = ambientOcclusion(aoPPZ, aoPZP, aoPPP)
				meshAmbient[index][PosX][3] = ambientOcclusion(aoPMZ, aoPZP, aoPMP)
				/* neg x */
				meshAmbient[index][NegX][0] = ambientOcclusion(aoMMZ, aoMZM, aoMMM)
				meshAmbient[index][NegX][1] = ambientOcclusion(aoMMZ, aoMZP, aoMMP)
				meshAmbient[index][NegX][2] = ambientOcclusion(aoMPZ, aoMZP, aoMPP)
				meshAmbient[index][NegX][3] = ambientOcclusion(aoMPZ, aoMZM, aoMPM)
				/* pos y */
				meshAmbient[index][PosY][0] = ambientOcclusion(aoMPZ, aoZPM, aoMPM)
				meshAmbient[index][PosY][1] = ambientOcclusion(aoMPZ, aoZPP, aoMPP)
				meshAmbient[index][PosY][2] = ambientOcclusion(aoPPZ, aoZPP, aoPPP)
				meshAmbient[index][PosY][3] = ambientOcclusion(aoPPZ, aoZPM, aoPPM)
				/* neg y */
				meshAmbient[index][NegY][0] = ambientOcclusion(aoMMZ, aoZMM, aoMMM)
				meshAmbient[index][NegY][1] = ambientOcclusion(aoPMZ, aoZMM, aoPMM)
				meshAmbient[index][NegY][2] = ambientOcclusion(aoPMZ, aoZMP, aoPMP)
				meshAmbient[index][NegY][3] = ambientOcclusion(aoMMZ, aoZMP, aoMMP)
				/* pos z */
				meshAmbient[index][PosZ][0] = ambientOcclusion(aoPZP, aoZMP, aoPMP)
				meshAmbient[index][PosZ][1] = ambientOcclusion(aoPZP, aoZPP, aoPPP)
				meshAmbient[index][PosZ][2] = ambientOcclusion(aoMZP, aoZPP, aoMPP)
				meshAmbient[index][PosZ][3] = ambientOcclusion(aoMZP, aoZMP, aoMMP)
				/* neg z */
				meshAmbient[index][NegZ][0] = ambientOcclusion(aoMZM, aoZMM, aoMMM)
				meshAmbient[index][NegZ][1] = ambientOcclusion(aoMZM, aoZPM, aoMPM)
				meshAmbient[index][NegZ][2] = ambientOcclusion(aoPZM, aoZPM, aoPPM)
				meshAmbient[index][NegZ][3] = ambientOcclusion(aoPZM, aoZMM, aoPMM)
			}
		}
	}
	for bz := 0; bz < meshColorSize; bz++ {
		for by := 0; by < meshColorSize; by++ {
			for bx := 0; bx < meshColorSize; bx++ {
				red := int32(0)
				green := int32(0)
				blue := int32(0)
				count := int32(0)
				blockZZZ := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx, by, bz)
				blockMZZ := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx-1, by, bz)
				blockMZM := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx-1, by, bz-1)
				blockZZM := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx, by, bz-1)
				blockZMZ := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx, by-1, bz)
				blockMMZ := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx-1, by-1, bz)
				blockMMM := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx-1, by-1, bz-1)
				blockZMM := WorldGetBlockPointer(obj.x, obj.y, obj.z, bx, by-1, bz-1)
				if BlockPointerClosed(blockZZZ) {
					determineLight(blockMZZ, &red, &green, &blue, &count)
					determineLight(blockZMZ, &red, &green, &blue, &count)
					determineLight(blockZZM, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockMZZ) {
					determineLight(blockZZZ, &red, &green, &blue, &count)
					determineLight(blockZMZ, &red, &green, &blue, &count)
					determineLight(blockZZM, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockMZM) {
					determineLight(blockMZZ, &red, &green, &blue, &count)
					determineLight(blockZZM, &red, &green, &blue, &count)
					determineLight(blockMMM, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockZZM) {
					determineLight(blockZZZ, &red, &green, &blue, &count)
					determineLight(blockMZM, &red, &green, &blue, &count)
					determineLight(blockZMM, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockZMZ) {
					determineLight(blockZZZ, &red, &green, &blue, &count)
					determineLight(blockMMZ, &red, &green, &blue, &count)
					determineLight(blockZMM, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockMMZ) {
					determineLight(blockMZZ, &red, &green, &blue, &count)
					determineLight(blockMMM, &red, &green, &blue, &count)
					determineLight(blockZMZ, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockMMM) {
					determineLight(blockMZM, &red, &green, &blue, &count)
					determineLight(blockZMM, &red, &green, &blue, &count)
					determineLight(blockMMZ, &red, &green, &blue, &count)
				}
				if BlockPointerClosed(blockZMM) {
					determineLight(blockZZM, &red, &green, &blue, &count)
					determineLight(blockZMZ, &red, &green, &blue, &count)
					determineLight(blockMMM, &red, &green, &blue, &count)
				}
				index := bx + by*meshColorSize + bz*meshColorSlice
				if count > 0 {
					meshColor[index] = packRgb(red/count, green/count, blue/count)
				} else {
					meshColor[index] = 0
				}
			}
		}
	}
}

func determineLight(o *Block, red *int32, green *int32, blue *int32, count *int32) {
	if o == nil {
		return
	}
	if BlockClosed(o.Type) == false {
		r, g, b := unpackRgb(o.Color)
		*red += r
		*green += g
		*blue += b
		*count++
	}
}

func ambientOcclusion(side1, side2, corner bool) byte {
	if side1 && side2 {
		return 175
	}
	if side1 || side2 || corner {
		return 215
	}
	return 255
}

func RenderWorld(lx, ly, lz int) {
	if printDebug {
		fmt.Println("rendering", viewableNum, "chunks")
	}
	timeA := glfw.GetTime()
	for i := 0; i < viewableNum; i++ {
		obj := viewable[i]
		mesh := obj.Mesh
		if mesh.VertexLimit == 0 {
			continue
		}
		mesh.BindVao()
		if lx == obj.x {
			g.DrawRange(obj.beginSide[PosX], obj.countSide[PosX])
			g.DrawRange(obj.beginSide[NegX], obj.countSide[NegX])
		} else if lx > obj.x {
			g.DrawRange(obj.beginSide[PosX], obj.countSide[PosX])
		} else {
			g.DrawRange(obj.beginSide[NegX], obj.countSide[NegX])
		}
		if ly == obj.y {
			g.DrawRange(obj.beginSide[PosY], obj.countSide[PosY])
			g.DrawRange(obj.beginSide[NegY], obj.countSide[NegY])
		} else if ly > obj.y {
			g.DrawRange(obj.beginSide[PosY], obj.countSide[PosY])
		} else {
			g.DrawRange(obj.beginSide[NegY], obj.countSide[NegY])
		}
		if lz == obj.z {
			g.DrawRange(obj.beginSide[PosZ], obj.countSide[PosZ])
			g.DrawRange(obj.beginSide[NegZ], obj.countSide[NegZ])
		} else if lz > obj.z {
			g.DrawRange(obj.beginSide[PosZ], obj.countSide[PosZ])
		} else {
			g.DrawRange(obj.beginSide[NegZ], obj.countSide[NegZ])
		}
	}
	if printWorldRenderTime {
		timeB := glfw.GetTime()
		fmt.Println("world", timeB-timeA)
	}
}
