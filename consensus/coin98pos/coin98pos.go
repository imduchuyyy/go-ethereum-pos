package coin98pos

import (
  "math/big"
  "golang.org/x/crypto/sha3"
  "io"
  "fmt"
  "errors"
  "sync"

  "github.com/ethereum/go-ethereum/params"
  "github.com/ethereum/go-ethereum/rpc"
  "github.com/ethereum/go-ethereum/consensus"
  "github.com/ethereum/go-ethereum/consensus/clique"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/common/worker"
  "github.com/ethereum/go-ethereum/trie"
  "github.com/ethereum/go-ethereum/core/state"
  "github.com/ethereum/go-ethereum/rlp"
  lru "github.com/hashicorp/golang-lru"
)

const (
  inmemorySnapshots = 128 // Number of recent vote snapshots to keep in memory
  inmemorySignatures = 4096 // Number of recent block signatures to keep in memory
)

var (
  epochLength = uint64(900)
  extraSeal = crypto.SignatureLength
  coin98posNonce = types.EncodeNonce(0)
)

type ValidatorNode struct {
  address common.Address
  staked *big.Int
}

type Coin98Pos struct {
  chainConfig *params.ChainConfig
  config *params.Coin98PosConfig

  signer common.Address
  signFn clique.SignerFn
  lock sync.RWMutex

  recents *lru.ARCCache
  signatures *lru.ARCCache
  validatorSignatures *lru.ARCCache
  verifiedHeaders *lru.ARCCache
}

var (
  errTooManyUncles = errors.New("too many uncles")
	errInvalidNonce     = errors.New("invalid nonce")
	errInvalidUncleHash = errors.New("invalid uncle hash")

  // ecrecover
  errMissingSignature = errors.New("extra-data 65 byte suffix signature missing")
)

// ecrecover extracts the Ethereum account address from a signed header.
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
	pubkey, err := crypto.Ecrecover(sigHash(header).Bytes(), signature)
	if err != nil {
		return common.Address{}, err
	}
	var signer common.Address
	copy(signer[:], crypto.Keccak256(pubkey[1:])[12:])

	sigcache.Add(hash, signer)
	return signer, nil
}

func New(chainConfig *params.ChainConfig) *Coin98Pos {
  if chainConfig.Coin98Pos.Epoch == 0 {
    chainConfig.Coin98Pos.Epoch = epochLength
  }

  recents, _ := lru.NewARC(inmemorySnapshots)
  signatures, _ := lru.NewARC(inmemorySnapshots)
  validatorSignatures, _ := lru.NewARC(inmemorySnapshots)
  verifiedHeaders, _ := lru.NewARC(inmemorySnapshots)

  c := &Coin98Pos{
    chainConfig: chainConfig,
    config: chainConfig.Coin98Pos,

    recents: recents,
    signatures: signatures,
    validatorSignatures: validatorSignatures,
    verifiedHeaders: verifiedHeaders,
  }

  return c
}

func (c *Coin98Pos) APIs(chain consensus.ChainHeaderReader) []rpc.API {
  return []rpc.API{{
		Namespace: "coin98pos",
		Version:   "1.0",
		Service:   &API{chain: chain, coin98pos: c},
		Public:    false,
	}}
}

func (c *Coin98Pos) Author(header *types.Header) (common.Address, error) {
  return ecrecover(header, c.signatures)
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
func (c *Coin98Pos) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	// Transition isn't triggered yet, use the legacy rules for calculation
  return common.Big0
}

func (c *Coin98Pos) Close() error {
	return nil
}

// Finalize implements consensus.Engine, setting the final state on the header
func (c *Coin98Pos) Finalize(_chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, _txs []*types.Transaction, _uncles []*types.Header) {
	// Finalize is different with Prepare, it can be used in both block generation
	// and verification. So determine the consensus rules by header type.
	// The block reward is no longer handled here. It's done by the
	// external consensus engine.
	header.Root = state.IntermediateRoot(true)
}


// FinalizeAndAssemble implements consensus.Engine, setting the final state and
// assembling the block.
func (c *Coin98Pos) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error) {
	// FinalizeAndAssemble is different with Prepare, it can be used in both block
	// generation and verification. So determine the consensus rules by header type.
	// Finalize and assemble the block
	c.Finalize(chain, header, state, txs, uncles)
	return types.NewBlock(header, txs, uncles, receipts, trie.NewStackTrie(nil)), nil
}

// Prepare implements consensus.Engine, initializing the difficulty field of a
// header to conform to the coin98 protocol. The changes are done inline.
func (c *Coin98Pos) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	// Transition isn't triggered yet, use the legacy rules for preparation.
  header.Coinbase = common.Address{}
  header.Nonce = coin98posNonce

  number := header.Number.Uint64()

  snap, err := c.snapshot(chain, number-1, header.ParentHash, nil)
	return nil
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (c *Coin98Pos) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	// The seal verification is done by the external consensus engine,
	// return directly without pushing any block back. In another word
	// coin98 won't return any result by `results` channel which may
	// blocks the receiver logic forever.
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (c *Coin98Pos) SealHash(header *types.Header) (hash common.Hash) {
  hasher := sha3.NewLegacyKeccak256()

  encodeSigHeader(hasher, header, c.chainConfig.ChainID)
  hasher.Sum(hash[:0])
  return hash
}

func encodeSigHeader(w io.Writer, header *types.Header, chainId *big.Int) {
	err := rlp.Encode(w, []interface{}{
		chainId,
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
		header.Extra[:len(header.Extra)-65], // this will panic if extra is too short, should check before calling encodeSigHeader
		header.MixDigest,
		header.Nonce,
	})
	if err != nil {
		panic("can't encode: " + err.Error())
	}
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum consensus engine.
func (c *Coin98Pos) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header, seal bool) error {
	// Short circuit if the parent is not known
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return c.verifyHeader(chain, header, parent)
}

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Ethereum consensus engine. The difference between the coin98 and classic is
// (a) The following fields are expected to be constants:
//     - difficulty is expected to be 0
// 	   - nonce is expected to be 0
//     - unclehash is expected to be Hash(emptyHeader)
//     to be the desired constants
// (b) the timestamp is not verified anymore
// (c) the extradata is limited to 32 bytes
func (c *Coin98Pos) verifyHeader(chain consensus.ChainHeaderReader, header, parent *types.Header) error {
	// Ensure that the header's extra-data section is of a reasonable size
	if len(header.Extra) > 32 {
		return fmt.Errorf("extra-data longer than 32 bytes (%d)", len(header.Extra))
	}
	// Verify the seal parts. Ensure the nonce and uncle hash are the expected value.
	if header.Nonce != coin98posNonce {
		return errInvalidNonce
	}
	if header.UncleHash != types.EmptyUncleHash {
		return errInvalidUncleHash
	}
	// Verify that the gas limit is <= 2^63-1
	if header.GasLimit > params.MaxGasLimit {
		return fmt.Errorf("invalid gasLimit: have %v, max %v", header.GasLimit, params.MaxGasLimit)
	}
	// Verify that the gasUsed is <= gasLimit
	if header.GasUsed > header.GasLimit {
		return fmt.Errorf("invalid gasUsed: have %d, gasLimit %d", header.GasUsed, header.GasLimit)
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(common.Big1) != 0 {
		return consensus.ErrInvalidNumber
	}

  return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
// VerifyHeaders expect the headers to be ordered and continuous.
func (c *Coin98Pos) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
  abort := make(chan struct{})
  results := make(chan error, len(headers))

  worker.Submit(func() {
		for _, header := range headers {
      parent := chain.GetHeader(header.ParentHash, header.Number.Uint64() - 1)
			err := c.verifyHeader(chain, header, parent)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	})
	return abort, results
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of the Ethereum consensus engine.
func (c *Coin98Pos) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	// Verify that there is no uncle block. It's explicitly disabled in the beacon
	if len(block.Uncles()) > 0 {
		return errTooManyUncles
	}
	return nil
}

// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewLegacyKeccak256()

	rlp.Encode(hasher, []interface{}{
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
		header.Extra[:len(header.Extra)-65], // Yes, this will panic if extra is too short
		header.MixDigest,
		header.Nonce,
	})
	hasher.Sum(hash[:0])
	return hash
}

func SigHash(header *types.Header) (hash common.Hash) {
	return sigHash(header)
}
