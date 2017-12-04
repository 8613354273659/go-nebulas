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
	"reflect"
	"testing"
	"time"

	pb "github.com/gogo/protobuf/proto"
	"github.com/nebulasio/go-nebulas/core/pb"
	"github.com/nebulasio/go-nebulas/crypto"
	"github.com/nebulasio/go-nebulas/crypto/keystore"
	"github.com/nebulasio/go-nebulas/crypto/keystore/secp256k1"
	"github.com/nebulasio/go-nebulas/storage"
	"github.com/nebulasio/go-nebulas/util"
	"github.com/nebulasio/go-nebulas/util/byteutils"
	"github.com/stretchr/testify/assert"
)

func TestBlockHeader(t *testing.T) {
	type fields struct {
		hash       byteutils.Hash
		parentHash byteutils.Hash
		stateRoot  byteutils.Hash
		nonce      uint64
		coinbase   *Address
		timestamp  int64
		chainID    uint32
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"full struct",
			fields{
				[]byte("124546"),
				[]byte("344543"),
				[]byte("43656"),
				3546456,
				&Address{[]byte("hello")},
				time.Now().Unix(),
				1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BlockHeader{
				hash:       tt.fields.hash,
				parentHash: tt.fields.parentHash,
				stateRoot:  tt.fields.stateRoot,
				nonce:      tt.fields.nonce,
				coinbase:   tt.fields.coinbase,
				timestamp:  tt.fields.timestamp,
				chainID:    tt.fields.chainID,
			}
			proto, _ := b.ToProto()
			ir, _ := pb.Marshal(proto)
			nb := new(BlockHeader)
			pb.Unmarshal(ir, proto)
			nb.FromProto(proto)
			b.timestamp = nb.timestamp
			if !reflect.DeepEqual(*b, *nb) {
				t.Errorf("Transaction.Serialize() = %v, want %v", *b, *nb)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	type fields struct {
		header       *BlockHeader
		transactions Transactions
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"full struct",
			fields{
				&BlockHeader{
					hash:       []byte("124546"),
					parentHash: []byte("344543"),
					stateRoot:  []byte("43656"),
					txsRoot:    []byte("43656"),
					dposContext: &corepb.DposContext{
						DynastyRoot:     []byte("43656"),
						NextDynastyRoot: []byte("43656"),
						DelegateRoot:    []byte("43656"),
					},
					nonce:     3546456,
					coinbase:  &Address{[]byte("hello")},
					timestamp: time.Now().Unix(),
					chainID:   1,
				},
				Transactions{
					&Transaction{
						[]byte("123455"),
						&Address{[]byte("1235")},
						&Address{[]byte("1245")},
						util.NewUint128(),
						456,
						time.Now().Unix(),
						&corepb.Data{Type: TxPayloadBinaryType, Payload: []byte("hello")},
						1,
						util.NewUint128(),
						util.NewUint128(),
						uint8(keystore.SECP256K1),
						nil,
					},
					&Transaction{
						[]byte("123455"),
						&Address{[]byte("1235")},
						&Address{[]byte("1245")},
						util.NewUint128(),
						456,
						time.Now().Unix(),
						&corepb.Data{Type: TxPayloadBinaryType, Payload: []byte("hello")},
						1,
						util.NewUint128(),
						util.NewUint128(),
						uint8(keystore.SECP256K1),
						nil,
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				header:       tt.fields.header,
				transactions: tt.fields.transactions,
			}
			proto, _ := b.ToProto()
			ir, _ := pb.Marshal(proto)
			nb := new(Block)
			pb.Unmarshal(ir, proto)
			nb.FromProto(proto)
			b.header.timestamp = nb.header.timestamp
			b.transactions[0].timestamp = nb.transactions[0].timestamp
			b.transactions[1].timestamp = nb.transactions[1].timestamp
			if !reflect.DeepEqual(*b.header, *nb.header) {
				t.Errorf("Transaction.Serialize() = %v, want %v", *b.header, *nb.header)
			}
			if !reflect.DeepEqual(*b.transactions[0], *nb.transactions[0]) {
				t.Errorf("Transaction.Serialize() = %v, want %v", *b.transactions[0], *nb.transactions[0])
			}
			if !reflect.DeepEqual(*b.transactions[1], *nb.transactions[1]) {
				t.Errorf("Transaction.Serialize() = %v, want %v", *b.transactions[1], *nb.transactions[1])
			}
		})
	}
}

func TestBlock_LinkParentBlock(t *testing.T) {
	storage, _ := storage.NewMemoryStorage()
	bc, _ := NewBlockChain(0, storage)
	genesis := bc.genesisBlock
	assert.Equal(t, genesis.Height(), uint64(1))
	block1 := &Block{
		header: &BlockHeader{
			hash:       []byte("124546"),
			parentHash: GenesisHash,
			stateRoot:  []byte("43656"),
			txsRoot:    []byte("43656"),
			dposContext: &corepb.DposContext{
				DynastyRoot:     []byte("43656"),
				NextDynastyRoot: []byte("43656"),
				DelegateRoot:    []byte("43656"),
			},
			nonce:     3546456,
			coinbase:  &Address{[]byte("hello")},
			timestamp: time.Now().Unix(),
			chainID:   1,
		},
		transactions: []*Transaction{},
	}
	assert.Equal(t, block1.Height(), uint64(0))
	assert.Equal(t, block1.LinkParentBlock(genesis), true)
	assert.Equal(t, block1.Height(), uint64(2))
	assert.Equal(t, block1.ParentHash(), genesis.Hash())
	block2 := &Block{
		header: &BlockHeader{
			hash:       []byte("124546"),
			parentHash: []byte("344543"),
			stateRoot:  []byte("43656"),
			txsRoot:    []byte("43656"),
			dposContext: &corepb.DposContext{
				DynastyRoot:     []byte("43656"),
				NextDynastyRoot: []byte("43656"),
				DelegateRoot:    []byte("43656"),
			},
			nonce:     3546456,
			coinbase:  &Address{[]byte("hello")},
			timestamp: time.Now().Unix(),
			chainID:   1,
		},
		transactions: []*Transaction{},
	}
	assert.Equal(t, block2.LinkParentBlock(genesis), false)
	assert.Equal(t, block2.Height(), uint64(0))
}

func TestBlock_CollectTransactions(t *testing.T) {
	storage, _ := storage.NewMemoryStorage()
	bc, _ := NewBlockChain(0, storage)
	tail := bc.tailBlock
	assert.Panics(t, func() { tail.CollectTransactions(1) })

	ks := keystore.DefaultKS
	priv := secp256k1.GeneratePrivateKey()
	pubdata, _ := priv.PublicKey().Encoded()
	from, _ := NewAddressFromPublicKey(pubdata)
	ks.SetKey(from.ToHex(), priv, []byte("passphrase"))
	ks.Unlock(from.ToHex(), []byte("passphrase"), time.Second*60*60*24*365)

	key, _ := ks.GetUnlocked(from.ToHex())
	signature, _ := crypto.NewSignature(keystore.SECP256K1)
	signature.InitSign(key.(keystore.PrivateKey))

	priv1 := secp256k1.GeneratePrivateKey()
	pubdata1, _ := priv1.PublicKey().Encoded()
	to, _ := NewAddressFromPublicKey(pubdata1)
	priv2 := secp256k1.GeneratePrivateKey()
	pubdata2, _ := priv2.PublicKey().Encoded()
	coinbase, _ := NewAddressFromPublicKey(pubdata2)

	block0 := NewBlock(0, from, tail)
	block0.Seal()
	//bc.BlockPool().push(block0)
	bc.SetTailBlock(block0)

	block := NewBlock(0, coinbase, block0)

	tx1 := NewTransaction(0, from, to, util.NewUint128FromInt(1), 1, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx1.Sign(signature)
	tx2 := NewTransaction(0, from, to, util.NewUint128FromInt(1), 2, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx2.Sign(signature)
	tx3 := NewTransaction(0, from, to, util.NewUint128FromInt(1), 0, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx3.Sign(signature)
	tx4 := NewTransaction(0, from, to, util.NewUint128FromInt(1), 4, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx4.Sign(signature)
	tx5 := NewTransaction(0, from, to, util.NewUint128FromInt(1), 3, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx5.Sign(signature)
	tx6 := NewTransaction(1, from, to, util.NewUint128FromInt(1), 1, TxPayloadBinaryType, []byte("nas"), TransactionGasPrice, util.NewUint128FromInt(200000))
	tx6.Sign(signature)

	bc.txPool.Push(tx1)
	bc.txPool.Push(tx2)
	bc.txPool.Push(tx3)
	bc.txPool.Push(tx4)
	bc.txPool.Push(tx5)
	bc.txPool.Push(tx6)

	assert.Equal(t, len(block.transactions), 0)
	assert.Equal(t, bc.txPool.cache.Len(), 5)
	block.CollectTransactions(bc.txPool.cache.Len())
	assert.Equal(t, len(block.transactions), 4)
	assert.Equal(t, block.txPool.cache.Len(), 0)

	assert.Equal(t, block.Sealed(), false)
	balance := block.GetBalance(block.header.coinbase.address)
	assert.Equal(t, balance.Cmp(util.NewUint128().Int), 0)
	block.Seal()
	assert.Equal(t, block.Sealed(), true)
	assert.Equal(t, block.transactions[0], tx1)
	assert.Equal(t, block.transactions[1], tx2)
	assert.Equal(t, block.StateRoot().Equals(block.accState.RootHash()), true)
	assert.Equal(t, block.TxsRoot().Equals(block.txsTrie.RootHash()), true)
	balance = block.GetBalance(block.header.coinbase.address)
	assert.Equal(t, balance.Cmp(BlockReward.Int), 0)
	// mock net message
	proto, _ := block.ToProto()
	ir, _ := pb.Marshal(proto)
	nb := new(Block)
	pb.Unmarshal(ir, proto)
	nb.FromProto(proto)
	nb.LinkParentBlock(bc.tailBlock)
	assert.Nil(t, nb.Verify(0))
}
