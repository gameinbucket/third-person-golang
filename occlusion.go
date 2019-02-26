package main

import (
	"fmt"
)

const (
	printOcclusion = false
)

var (
	visionSliceA [3]int
	visionSliceB [3]int
)

func (obj *Chunk) calculateOcclusion() {
	for sideA := uint(0); sideA < Sides; sideA++ {
		ax := sliceX[sideA]
		ay := sliceY[sideA]
		az := sliceZ[sideA]
		for sideB := sideA + 1; sideB < Sides; sideB++ {
			bx := sliceX[sideB]
			by := sliceY[sideB]
			bz := sliceZ[sideB]
			if sliceTowards[sideA] > 0 {
				visionSliceA[2] = chunkSize - 1
			} else {
				visionSliceA[2] = 0
			}
			if sliceTowards[sideB] > 0 {
				visionSliceB[2] = chunkSize - 1
			} else {
				visionSliceB[2] = 0
			}
		sliceLabel:
			for visionSliceA[1] = 0; visionSliceA[1] < chunkSize; visionSliceA[1]++ {
				for visionSliceA[0] = 0; visionSliceA[0] < chunkSize; visionSliceA[0]++ {
					for visionSliceB[1] = 0; visionSliceB[1] < chunkSize; visionSliceB[1]++ {
						for visionSliceB[0] = 0; visionSliceB[0] < chunkSize; visionSliceB[0]++ {
							fromX := float32(visionSliceA[ax]) + 0.5
							fromY := float32(visionSliceA[ay]) + 0.5
							fromZ := float32(visionSliceA[az]) + 0.5
							toX := float32(visionSliceB[bx]) + 0.5
							toY := float32(visionSliceB[by]) + 0.5
							toZ := float32(visionSliceB[bz]) + 0.5
							if printOcclusion {
								fmt.Println(fromX, fromY, fromZ, toX, toY, toZ)
							}
							if raycast(obj, fromX, fromY, fromZ, toX, toY, toZ) {
								obj.visibility |= 1 << (sideA*6 + sideB)
								if printOcclusion {
									fmt.Println("true")
								}
								break sliceLabel
							}
						}
					}
				}
			}
		}
	}
	if false {
		for sideA := uint(0); sideA < Sides; sideA++ {
			for sideB := sideA + 1; sideB < Sides; sideB++ {
				if obj.visibility&(1<<(sideA*6+sideB)) > 0 {
					fmt.Println(sideA, "<->", sideB)
				}
			}
		}
	}
}
