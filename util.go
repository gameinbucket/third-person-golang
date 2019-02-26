package main

func indexTo3d(i int) (int, int, int) {
	z := i / worldChunkSlice
	i -= z * worldChunkSlice
	y := i / worldChunksX
	x := i % worldChunksX
	return x, y, z
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func packRgb(red, green, blue int32) int32 {
	return (red << 16) | (green << 8) | blue
}

func unpackRgb(rgb int32) (int32, int32, int32) {
	red := (rgb >> 16) & 255
	green := (rgb >> 8) & 255
	blue := rgb & 255
	return red, green, blue
}
