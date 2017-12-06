// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	"errors"
)

// Payload Types
const (
	TxPayloadBinaryType = "binary"
	TxPayloadDeployType = "deploy"
	TxPayloadCallType   = "call"
	TxPayloadVoteType   = "vote"
)

// Error Types
var (
	ErrInvalidTxPayloadType        = errors.New("invalid transaction data payload type")
	ErrInvalidContractAddress      = errors.New("invalid contract address")
	ErrInsufficientBalance         = errors.New("insufficient balance")
	ErrInvalidSignature            = errors.New("invalid transaction signature")
	ErrInvalidTransactionHash      = errors.New("invalid transaction hash")
	ErrMissingParentBlock          = errors.New("cannot find a on-chain block's parent block in storage")
	ErrTooFewCandidates            = errors.New("too few candidates in consensus")
	ErrNotBlockForgTime            = errors.New("now is not time to forg block")
	ErrInvalidBlockHash            = errors.New("invalid block hash")
	ErrInvalidBlockStateRoot       = errors.New("invalid block state root hash")
	ErrInvalidBlockTxsRoot         = errors.New("invalid block txs root hash")
	ErrInvalidBlockDposContextRoot = errors.New("invalid block dpos context root hash")
	ErrInvalidChainID              = errors.New("invalid transaction chainID")
	ErrDuplicatedTransaction       = errors.New("duplicated transaction")
	ErrSmallTransactionNonce       = errors.New("cannot accept a transaction with smaller nonce")
	ErrLargeTransactionNonce       = errors.New("cannot accept a transaction with too bigger nonce")
	ErrDuplicatedBlock             = errors.New("duplicated block")
	ErrInvalidAddress              = errors.New("address: invalid address")
	ErrInvalidAddressDataLength    = errors.New("address: invalid address data length")
	ErrOutofGasLimit               = errors.New("out of gas limit")
	ErrDoubleSealBlock             = errors.New("cannot seal a block twice")
)

// TxPayload stored in tx
type TxPayload interface {
	ToBytes() ([]byte, error)
	Execute(tx *Transaction, block *Block) error
}

// MessageType
const (
	MessageTypeNewBlock = "newblock"
	MessageTypeNewTx    = "newtx"
)

// Consensus interface
type Consensus interface {
	VerifyBlock(block *Block, parent *Block) error
	FastVerifyBlock(block *Block) error
}
