package main

const (
	width  = float32(128)
	height = float32(128)
	size   = float32(16)
	wid    = size / width
	hgt    = size / height

	BlockNone         = int32(0)
	BlockGrass        = int32(1)
	BlockStone        = int32(2)
	BlockStoneSlopeXp = int32(3)
)

type Block struct {
	Type  int32
	Light int32
	// SunLight int32
	Color int32
}

func BlockTexture(typ int32) (float32, float32, float32, float32) {
	switch typ {
	case BlockGrass:
		return 0.0, 0.0, 1.0 * wid, 1.0 * hgt
	case BlockStone:
		return 1.0 * wid, 0.0, 2.0 * wid, 1.0 * hgt
	}
	return 0.0, 0.0, 1.0, 1.0
}

func BlockTextureIndex(typ int32) float32 {
	switch typ {
	case BlockGrass:
		return 0.0
	case BlockStone:
		return 1.0
	}
	return 0.0
}

func BlockClosed(typ int32) bool {
	switch typ {
	case BlockStoneSlopeXp:
		return false
	case BlockNone:
		return false
	}
	return true
}

func BlockPointerClosed(o *Block) bool {
	if o == nil {
		return true
	}
	return BlockClosed(o.Type)
}
