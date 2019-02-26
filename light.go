package main

import (
	"fmt"
)

const (
	lightQueueLimit = 30 * 30 * 30
)

type Point3d struct {
	X int
	Y int
	Z int
}

var (
	lightQueue    [lightQueueLimit]Point3d
	lightQueuePos int
	lightQueueNum int
)

func colorizeVertex(rgb int32, ambient float32) (float32, float32, float32) {
	r, g, b := unpackRgb(rgb)
	red := (float32(r) / 255.0) * ambient
	green := (float32(g) / 255.0) * ambient
	blue := (float32(b) / 255.0) * ambient
	return red, green, blue
}

func lightVisit(cx, cy, cz, bx, by, bz int, red, green, blue int32) {
	b := WorldGetBlockPointer(cx, cy, cz, bx, by, bz)
	if b == nil || BlockClosed(b.Type) {
		return
	}
	cr, cg, cb := unpackRgb(b.Color)
	if cr >= red || cg >= green || cb >= blue {
		return
	}

	if false {
		fmt.Println(cr, cg, cb, red, green, blue)
	}
	/*if cr > red {
		red = cr
	}
	if cg > green {
		green = cg
	}
	if cb > blue {
		blue = cb
	}*/
	b.Color = packRgb(red, green, blue)

	queue := lightQueuePos + lightQueueNum
	if queue >= lightQueueLimit {
		queue -= lightQueueLimit
	}
	lightQueue[queue].X = bx
	lightQueue[queue].Y = by
	lightQueue[queue].Z = bz
	lightQueueNum++
}

func (obj *Chunk) processLights() {
	for bz := 0; bz < chunkSize; bz++ {
		for by := 0; by < chunkSize; by++ {
			for bx := 0; bx < chunkSize; bx++ {
				node := obj.GetBlockPointerUnsafe(bx, by, bz)
				if node.Light > 0 {
					node.Color = node.Light
					lightQueuePos = 0
					lightQueue[0].X = bx
					lightQueue[0].Y = by
					lightQueue[0].Z = bz
					lightQueueNum = 1
					for lightQueueNum > 0 {
						x := lightQueue[lightQueuePos].X
						y := lightQueue[lightQueuePos].Y
						z := lightQueue[lightQueuePos].Z

						lightQueuePos++
						if lightQueuePos == lightQueueLimit {
							lightQueuePos = 0
						}
						lightQueueNum--

						node := WorldGetBlockPointer(obj.x, obj.y, obj.z, x, y, z)
						if node == nil {
							continue
						}

						r, g, b := unpackRgb(node.Color)
						r = int32(float32(r) * 0.8)
						g = int32(float32(g) * 0.8)
						b = int32(float32(b) * 0.8)
						lightVisit(obj.x, obj.y, obj.z, x-1, y, z, r, g, b)
						lightVisit(obj.x, obj.y, obj.z, x+1, y, z, r, g, b)
						lightVisit(obj.x, obj.y, obj.z, x, y-1, z, r, g, b)
						lightVisit(obj.x, obj.y, obj.z, x, y+1, z, r, g, b)
						lightVisit(obj.x, obj.y, obj.z, x, y, z-1, r, g, b)
						lightVisit(obj.x, obj.y, obj.z, x, y, z+1, r, g, b)
					}
				}
			}
		}
	}
}

func GiveSunlight() {
	cy := worldChunksY - 1
	for cz := 0; cz < worldChunksZ; cz++ {
		for cx := 0; cx < worldChunksX; cx++ {
			for bz := 0; bz < chunkSize; bz++ {
				for bx := 0; bx < chunkSize; bx++ {
					by := chunkSize - 1
					for true {
						node := WorldGetBlockPointer(cx, cy, cz, bx, by, bz)
						if node == nil || BlockClosed(node.Type) {
							break
						}
						node.Light = packRgb(255, 255, 255)
						by--
					}
				}
			}
		}
	}
}
