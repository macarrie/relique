package rsync

import (
	"crypto/sha256"
	"encoding/json"
	"strconv"

	"github.com/monmohan/xferspdy"
)

type Signature struct {
	Sig *xferspdy.Fingerprint
}

//type Fingerprint struct {
//	Blocksz  uint32
//	BlockMap map[uint32]map[[sha256.Size]byte]Block
//	Source   string
//}
type jsonFingerprint struct {
	Blocksz  uint32
	BlockMap map[string]map[string]xferspdy.Block
	Source   string
}

func NewSignature(fg *xferspdy.Fingerprint) *Signature {
	return &Signature{
		Sig: fg,
	}
}

func (f *Signature) UnmarshalJSON(b []byte) error {
	var jsonFg jsonFingerprint
	if err := json.Unmarshal(b, &jsonFg); err != nil {
		return err
	}

	f.Sig = xferspdy.NewFingerprint("/dev/null", 1024)
	f.Sig.Blocksz = jsonFg.Blocksz
	f.Sig.Source = jsonFg.Source

	// Convert map[string]map[string]Block to map[uint32]map[[sha256.Size]byte]Block
	blocksMap := make(map[uint32]map[[sha256.Size]byte]xferspdy.Block)
	for i, block := range jsonFg.BlockMap {
		index, err := strconv.Atoi(i)
		if err != nil {
			return err
		}
		innerMap := make(map[[sha256.Size]byte]xferspdy.Block)
		for innerIndex, innerBlock := range block {
			var innerIndexConv [sha256.Size]byte
			copy(innerIndexConv[:], innerIndex)
			innerMap[innerIndexConv] = innerBlock
		}
		blocksMap[uint32(index)] = innerMap
	}

	f.Sig.BlockMap = blocksMap

	return nil
}

func (f *Signature) MarshalJSON() ([]byte, error) {
	jsonFg := jsonFingerprint{}

	// Convert map[uint32]map[[sha256.Size]byte]Block to map[string]map[string]Block
	blocksMap := make(map[string]map[string]xferspdy.Block)
	if f.Sig != nil {
		jsonFg.Blocksz = f.Sig.Blocksz
		jsonFg.Source = f.Sig.Source

		for i, blockMap := range f.Sig.BlockMap {
			str := strconv.Itoa(int(i))
			innerMap := make(map[string]xferspdy.Block)
			for bytes, block := range blockMap {
				index := string(bytes[:])
				innerMap[index] = block
			}
			blocksMap[str] = innerMap
		}
	}

	jsonFg.BlockMap = blocksMap

	return json.Marshal(jsonFg)
}
