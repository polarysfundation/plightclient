package astrum

import (
	"encoding/binary"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
)

type AstrumPLightClient struct{}

type InputData struct {
	Height uint64
	Hash   common.Hash
}

type AstrumHeader struct {
	ParentHash common.Hash
	UncleHash  common.Hash
	Coinbase   common.Address
	Difficulty *big.Int
	Number     *big.Int
	GasLimit   uint64
	GasUsed    uint64
	Timestamp  uint64
	Nonce      types.BlockNonce
}

func (a *AstrumPLightClient) protocol(input []byte, chain consensus.ChainHeaderReader) (b []byte, err error) {
	data, err := a.decodeInput(input)
	if err != nil {
		return b, err
	}

	b, err = a.encodeHeader(chain, data.Height)
	if err != nil {
		return b, err
	}

	return b, nil

}

func (a *AstrumPLightClient) encodeHeader(chain consensus.ChainHeaderReader, height uint64) (result []byte, err error) {
	header := chain.GetHeaderByNumber(height)

	astrumHeader := &AstrumHeader{
		ParentHash: header.ParentHash,
		UncleHash:  header.UncleHash,
		Coinbase:   header.Coinbase,
		Difficulty: header.Difficulty,
		Number:     header.Number,
		GasLimit:   header.GasLimit,
		GasUsed:    header.GasUsed,
		Timestamp:  header.Time,
		Nonce:      header.Nonce,
	}
	copy(result[0:32], astrumHeader.ParentHash[:])
	copy(result[32:64], astrumHeader.UncleHash[:])
	copy(result[64:96], astrumHeader.Coinbase[12:])
	binary.BigEndian.PutUint64(result[96:128], astrumHeader.Difficulty.Uint64())
	binary.BigEndian.PutUint64(result[128:160], astrumHeader.Number.Uint64())
	binary.BigEndian.PutUint64(result[160:192], astrumHeader.GasLimit)
	binary.BigEndian.PutUint64(result[192:224], astrumHeader.GasUsed)
	binary.BigEndian.PutUint64(result[224:256], astrumHeader.Timestamp)
	copy(result[256:288], astrumHeader.Nonce[:])

	return result, nil
}

func (a *AstrumPLightClient) decodeInput(input []byte) (result InputData, err error) {

	if len(input) != 64 {
		return result, errors.New("invalid input length")
	}

	result.Height = binary.BigEndian.Uint64(input[0:32])
	copy(result.Hash[:], input[32:64])

	return result, nil
}
