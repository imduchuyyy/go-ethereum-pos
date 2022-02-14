package Pos

import (
  "time"
	"math/big"
  "io"
  "errors"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/log"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/core/state"
  "github.com/ethereum/go-ethereum/trie"
  "github.com/ethereum/go-ethereum/params"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/ethdb"
  "github.com/ethereum/go-ethereum/rpc"
  "github.com/ethereum/go-ethereum/rlp"
  "github.com/ethereum/go-ethereum/crypto"
  lru "github.com/hashicorp/golang-lru"
  "golang.org/x/crypto/sha3"
)


const (
	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	inmemorySignatures = 4096 // Number of recent block signatures to keep in memory

	wiggleTime = 500 * time.Millisecond // Random delay (per signer) to allow concurrent signers
)

var (
  extraSeal = crypto.SignatureLength
) 

var (
  errMissingSignature = errors.New("extra-data 65 byte signature suffix missing")
)

type Pos struct {
  config *params.CliqueConfig
  db ethdb.Database
  recents    *lru.ARCCache

  signatures *lru.ARCCache
}

func New(config *params.CliqueConfig, db ethdb.Database) *Pos {
	// Set any missing consensus parameters to their defaults
	conf := *config
	if conf.Epoch == 0 {
		conf.Epoch = 0
	}
	// Allocate the snapshot caches and create the engine

  recents, _ := lru.NewARC(inmemorySnapshots)
	signatures, _ := lru.NewARC(inmemorySignatures)

	return &Pos{
		config: &conf,
		db: db,
    recents: recents,
    signatures: signatures,
	}
}

func (Pos *Pos) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
  log.Info("will verfiyHeader")
  return nil
}

func (Pos *Pos) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
  log.Info("will verfiyHeaders")
  abort := make(chan struct{})
  results := make(chan error, len(headers))
  
  go func() {
    for _, header := range headers {
      err := Pos.VerifyHeader(chain, header, false)

      select {
        case <- abort: return
        case results <- err: return
      }
    }
  }()
  return abort, results
}

func (Pos *Pos) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
  log.Info("will verfiy uncles")
  return nil
}
func (Pos *Pos) VerifySeal(chain consensus.ChainReader, header *types.Header) error{
  log.Info("will verfiy VerifySeal")
  return nil

}
func (Pos *Pos) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error{
  log.Info("will prepare the block")
  parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
  if parent == nil {
    return consensus.ErrUnknownAncestor

  }
  header.Difficulty = Pos.CalcDifficulty(chain, header.Time, parent)
  return nil

}
func (Pos *Pos) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
  return new(big.Int) 
}

func (pos *Pos) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header,
	) {
	// No block rewards in PoA, so the state remains as is and uncles are dropped
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	header.UncleHash = types.CalcUncleHash(nil)
	return
}

// FinalizeAndAssemble implements consensus.Engine, ensuring no uncles are set,
// nor block rewards given, and returns the final block.
func (pos *Pos) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB,
	txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// No block rewards in PoA, so the state remains as is and uncles are dropped
	var err error
	header.Root = state.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	if err != nil {
		return nil, err
	}
	header.UncleHash = types.CalcUncleHash(nil)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, nil, receipts, trie.NewStackTrie(nil)), nil
}

func (Pos *Pos) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) (error){
  log.Info("will Seal the block")
  //time.Sleep(15 * time.Second)
  return nil
}

func ecrecover(header *types.Header, sigcache *lru.ARCCache) (common.Address, error) {
	// If the signature's already cached, return that
	hash := header.Hash()
	if address, known := sigcache.Get(hash); known {
		return address.(common.Address), nil
	}
	// Retrieve the signature from the header extra-data
	if len(header.Extra) < extraSeal {
		return common.Address{}, errMissingSignature
	}
	signature := header.Extra[len(header.Extra)-extraSeal:]

	// Recover the public key and the Ethereum address
	pubkey, err := crypto.Ecrecover(SealHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

func (pos *Pos) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header, pos.signatures)
}

func (pos *Pos) SealHash(header *types.Header) common.Hash {
	return SealHash(header)
}

func SealHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()
	encodeSigHeader(hasher, header)
	hasher.Sum(hash[:0])
	return hash
}

func (pos *Pos) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return []rpc.API{{
		Namespace: "Pos",
		Version:   "1.0",
		Service:   &API{chain: chain, pos: pos},
		Public:    false,
	}}
}

func encodeSigHeader(w io.Writer, header *types.Header) {
	err := rlp.Encode(w, []interface{}{
		header.ParentHash,
		header.UncleHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Difficulty,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		header.Extra[:len(header.Extra)-crypto.SignatureLength], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

func (pos *Pos) Close() error {
	return nil
}

