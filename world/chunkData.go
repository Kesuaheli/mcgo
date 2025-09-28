package world

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"mcgo/types"
)

type XZPosition struct {
	x int32
	z int32
}
type World struct {
	chunks map[XZPosition]*Chunk
}

type Chunk struct {
	dimensionHeight int

	chunkSections []*ChunkSection
}

type ChunkSection struct {
	blocks [16][16][16]*Block
}

type Block struct {
	id     string
	states map[string]any
}
type ChunkSectionFormat struct {
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

// convert internal to ChunkSectionNetProtocol
func (chunksection *ChunkSection) ToBytesOnWire() []byte {

	blockCount := uint16(0)
	// analyse data for paletted container

	bufBlockData := &bytes.Buffer{}
	blockIdToRegistryId := map[string]uint16{
		"minecraft:air":         0,
		"minecraft:stone":       1,
		"minecraft:grass_block": 9,
		"minecraft:dirt":        10,
		"minecraft:cobblestone": 14,
	}

	for y := range 16 {
		for z := range 16 {
			for x := range 4 {
				data := uint64(0)
				for i := range 4 {
					data <<= 15
					block := chunksection.blocks[y][z][x*4+i]
					if block.id != "minecraft:air" && block.id != "minecraft:void_air" && block.id != "minecraft:cave_air" {
						blockCount++
					}
					data |= uint64(blockIdToRegistryId[block.id])
				}
				binary.Write(bufBlockData, binary.BigEndian, data)
			}
		}
	}

	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, blockCount)
	binary.Write(buf, binary.BigEndian, uint8(15))
	// no palete
	// block data
	buf.Write(bufBlockData.Bytes())

	// biome data
	// single value
	binary.Write(buf, binary.BigEndian, uint8(0))
	// value is set as VarInt. 57 -> minecraft:the_void
	types.WriteVarInt(buf, 57)

	return buf.Bytes()

}

func NewChunkSection() *ChunkSection {
	chunkSection := ChunkSection{}
	for y := range 16 {
		for z := range 16 {
			for x := range 16 {
				chunkSection.blocks[y][z][x] = &Block{id: "minecraft:air"}
			}
		}
	}

	return &chunkSection
}

func (chunksection *ChunkSection) SetBlock(x, y, z int, id string) {
	chunksection.blocks[y][z][x].id = id
	chunksection.blocks[y][z][x].states = nil
}

func (chunksection *ChunkSection) FillWithBlocks(xMin, yMin, zMin, xMax, yMax, zMax int, id string) {
	for yOffset := range (yMax + 1) - yMin {
		y := yOffset + yMin
		for zOffset := range (zMax + 1) - zMin {
			z := zOffset + zMin
			for xOffset := range (xMax + 1) - xMin {
				x := xOffset + xMin
				chunksection.blocks[y][z][x].id = id
				chunksection.blocks[y][z][x].states = nil
			}
		}
	}
}

func NewChunk() (chunk *Chunk) {
	chunk = &Chunk{
		dimensionHeight: 24,
	}
	for range chunk.dimensionHeight {
		chunk.chunkSections = append(chunk.chunkSections, NewChunkSection())
	}
	return chunk
}

func (chunk *Chunk) GetChunkSection(index uint) *ChunkSection {
	return chunk.chunkSections[index]
}

func NewWorld() (world *World) {
	world = &World{
		chunks: map[XZPosition]*Chunk{},
	}

	worldSize := int32(5)
	worldSizeInDirection := worldSize / 2
	for x := -worldSizeInDirection; x <= worldSizeInDirection; x++ {
		for z := -worldSizeInDirection; z <= worldSizeInDirection; z++ {
			world.chunks[XZPosition{x, z}] = NewChunk()
		}
	}

	return world
}

func (world *World) GetChunk(x, z int32) (chunk *Chunk) {
	return world.chunks[XZPosition{x, z}]
}

func (chunk *Chunk) ToChunkData() []byte {
	buf := &bytes.Buffer{}

	// heightmaps
	buf.WriteByte(0)

	// data
	chunksectionBuf := &bytes.Buffer{}
	for _, chunksection := range chunk.chunkSections {
		chunksectionBuf.Write(chunksection.ToBytesOnWire())
	}
	types.WriteVarInt(buf, int32(chunksectionBuf.Len()))
	buf.Write(chunksectionBuf.Bytes())

	// block entities
	buf.WriteByte(0)

	return buf.Bytes()
}
