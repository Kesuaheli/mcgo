package world

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"mcgo/types"
)

type ChunkSection struct {
	BlockCount   int16
	BitsPerBlock uint8
	Data         []uint64
	Palette      []int32
}

type SingleValuePalette uint32

type Palette interface {
	GetFormat() PaletteFormat
	Values() []int32
}

type IndirectPalette struct {
	PaletteLength int32
	Palette       []int32
}

type DirectPalette struct{}

type PalettedContainerType int
type PaletteFormat int

const (
	PalettedContainerTypeBlocks PalettedContainerType = iota // Representing all the blocks in the chunk section.
	PalettedContainerTypeBiomes                              // Representing 4×4×4 biome regions in the chunk section.

	PaletteFormatSingleValue PaletteFormat = iota // When this palette is used, the Data Array sent/received is empty, since entries can be inferred from the palette's single value.
	PaletteFormatIndirect                         // This is an actual palette which lists the entries used. Values in the Data Array are indices into the palette, which in turn gives a proper registry ID.
	PaletteFormatDirect                           // Registry IDs are stored directly as entries in the Data Array.
)

func (p PalettedContainerType) GetEntries() int {
	switch p {
	case PalettedContainerTypeBlocks:
		return 4096
	case PalettedContainerTypeBiomes:
		return 64
	default:
		panic("Unknown container type")
	}
}

func (p PaletteFormat) String() string {
	switch p {
	case PaletteFormatSingleValue:
		return "SingleValuePalette"
	case PaletteFormatIndirect:
		return "IndirectPalette"
	case PaletteFormatDirect:
		return "DirectPalette"
	}
	panic("Unknown palette format")
}

func (p SingleValuePalette) GetFormat() PaletteFormat {
	return PaletteFormatSingleValue
}
func (p SingleValuePalette) Values() []int32 {
	return []int32{int32(p)}
}

func (p IndirectPalette) GetFormat() PaletteFormat {
	return PaletteFormatIndirect
}
func (p IndirectPalette) Values() []int32 {
	return p.Palette
}

func (p DirectPalette) GetFormat() PaletteFormat {
	return PaletteFormatDirect
}
func (p DirectPalette) Values() []int32 {
	return []int32{}
}

func ParseChunkData(buf *[]byte) {
	const NUM_OF_SECTIONS = 24
	r := bytes.NewReader(*buf)
	for i := range NUM_OF_SECTIONS {
		blockCounts, _ := types.Read(2, r)
		blockCount := binary.BigEndian.Uint16(blockCounts)
		fmt.Printf("Non-Air block count in section %d/%d: %d\n", i+1, NUM_OF_SECTIONS, blockCount)

		parsePalettedContainer(r, PalettedContainerTypeBlocks)
		parsePalettedContainer(r, PalettedContainerTypeBiomes)
	}
}

func parsePalettedContainer(r io.Reader, containerType PalettedContainerType) {
	bitsPerEntry, _ := types.ReadOne(r)
	paletteFormat := containerType.GetFormat(&bitsPerEntry)

	fmt.Printf("Bits per entry: %d → %s\n", bitsPerEntry, paletteFormat)

	var palette Palette

	switch paletteFormat {
	case PaletteFormatSingleValue:
		value, _ := types.ReadVarInt(r)
		palette = SingleValuePalette(value)

	case PaletteFormatIndirect:
		p := IndirectPalette{}
		p.PaletteLength, _ = types.ReadVarInt(r)
		p.Palette = make([]int32, 0, p.PaletteLength)
		for range p.PaletteLength {
			paletteEntry, _ := types.ReadVarInt(r)
			p.Palette = append(p.Palette, paletteEntry)
		}
		palette = p

	case PaletteFormatDirect:
		palette = DirectPalette{}
	default:
		panic("Palette format " + fmt.Sprint(paletteFormat) + " not implemented")
	}

	fmt.Printf("Got Palette: %+v\n", palette.Values())

	if paletteFormat == PaletteFormatSingleValue {
		// No block data to read
		return
	}

	var blocksPerLong = 64 / int(bitsPerEntry)
	var numOfLongs = int(math.Ceil(float64(containerType.GetEntries()) / float64(blocksPerLong)))

	blockData := make([]uint64, numOfLongs)
	err := binary.Read(r, binary.BigEndian, blockData)
	if err != nil {
		fmt.Printf("binary.Read failed: %v\n", err)
		return
	}

	entries := make([]int, 0, containerType.GetEntries())
	for _, long := range blockData {
		for range blocksPerLong {
			mask := (2 ^ uint64(bitsPerEntry)) - 1
			entry := int(long & mask)
			entries = append(entries, entry)
			long >>= bitsPerEntry
		}
	}
	fmt.Printf("Got Entries (x%d): %+v\n", len(entries), entries)
}

func (container PalettedContainerType) GetFormat(bitsPerEntry *uint8) PaletteFormat {
	switch container {
	case PalettedContainerTypeBlocks:
		switch {
		case *bitsPerEntry == 0:
			return PaletteFormatSingleValue
		case *bitsPerEntry < 4:
			fmt.Printf("Invalid bits per entry: %d. Rounded up to 4\n", *bitsPerEntry)
			*bitsPerEntry = 4
			fallthrough
		case *bitsPerEntry <= 8:
			return PaletteFormatIndirect
		case *bitsPerEntry < 15 || *bitsPerEntry > 15:
			fmt.Printf("Invalid bits per entry: %d. Set to 15\n", *bitsPerEntry)
			*bitsPerEntry = 15
			fallthrough
		default:
			return PaletteFormatDirect
		}
	case PalettedContainerTypeBiomes:
		switch {
		case *bitsPerEntry == 0:
			return PaletteFormatSingleValue
		case *bitsPerEntry <= 3:
			return PaletteFormatIndirect
		case *bitsPerEntry < 7 || *bitsPerEntry > 7:
			fmt.Printf("Invalid bits per entry: %d. Set to 7\n", *bitsPerEntry)
			*bitsPerEntry = 7
			fallthrough
		default:
			return PaletteFormatDirect
		}
	default:
		panic("Unknown container type")
	}
}
