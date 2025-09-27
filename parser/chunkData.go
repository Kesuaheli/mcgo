package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"mcgo/types"
)

type ChunkSection struct {
	BlockCount   int16
	BitsPerBlock uint8
	Data         []uint64
	Palette      []int32
}

type IndirectPalette struct {
	PaletteLength int32
	Palette       []int32
}

func ParseChunkData(buf *[]byte) map[string]any {
	// sections := make([]ChunkSection, 0)

	section := ChunkSection{}
	section.BlockCount, _ = types.PopShort(buf)
	section.BitsPerBlock, _ = types.PopUByte(buf)

	const BlockEntries = 4096
	var blocksPerLong = 64 / int(section.BitsPerBlock)
	var numOfLongs = int(math.Ceil(float64(BlockEntries) / float64(blocksPerLong)))

	palette := IndirectPalette{}
	palette.PaletteLength, _ = types.PopVarInt(buf)
	palette.Palette = make([]int32, palette.PaletteLength)
	for i := int32(0); i < palette.PaletteLength; i++ {
		palette.Palette[i], _ = types.PopVarInt(buf)
	}
	fmt.Printf("Got Palette: Palette Length: %d, Palette: %+v\n", palette.PaletteLength, palette.Palette)

	blockData := make([]uint64, numOfLongs)
	r := bytes.NewBuffer(*buf)
	err := binary.Read(r, binary.BigEndian, blockData)
	if err != nil {
		fmt.Printf("binary.Read failed: %v\n", err)
		return nil
	}
	fmt.Printf("Got Block Data (x%d): %+v\n", len(blockData), blockData)

	blocks := make([]int, 0, BlockEntries)
	for _, long := range blockData {
		fmt.Printf("Long: %064b\n", long)
		for range blocksPerLong {
			mask := (2 ^ uint64(section.BitsPerBlock)) - 1
			block := int(long & mask)
			blocks = append(blocks, block)
			long >>= section.BitsPerBlock
		}
	}
	fmt.Printf("Got Blocks (x%d): %+v\n", len(blocks), blocks)

	return nil
}
