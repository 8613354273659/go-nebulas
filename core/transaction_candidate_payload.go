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
	"encoding/json"

	"github.com/nebulasio/go-nebulas/util"
	log "github.com/sirupsen/logrus"
)

// Candidate Action
const (
	LoginAction  = "login"
	LogoutAction = "logout"
)

// CandidatePayload carry candidate application
type CandidatePayload struct {
	Action string
}

// LoadCandidatePayload from bytes
func LoadCandidatePayload(bytes []byte) (*CandidatePayload, error) {
	payload := &CandidatePayload{}
	if err := json.Unmarshal(bytes, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

// NewCandidatePayload with comments
func NewCandidatePayload(action string) *CandidatePayload {
	return &CandidatePayload{
		Action: action,
	}
}

// ToBytes serialize payload
func (payload *CandidatePayload) ToBytes() ([]byte, error) {
	return json.Marshal(payload)
}

// Execute the candidate payload in tx
func (payload *CandidatePayload) Execute(tx *Transaction, block *Block) (*util.Uint128, error) {
	candidate := tx.from.Bytes()
	counter := util.NewUint128()
	switch payload.Action {
	case LoginAction:
		counter.Add(counter.Int, one.Int)
		if _, err := block.dposContext.candidateTrie.Put(candidate, candidate); err != nil {
			return counter, err
		}
		log.WithFields(log.Fields{
			"func":      "Payload.Candidate",
			"block":     block,
			"tx":        tx,
			"candidate": tx.from.String(),
		}).Info("Candidate Login.")
	case LogoutAction:
		counter.Add(counter.Int, one.Int)
		if err := block.kickoutCandidate(candidate); err != nil {
			return counter, err
		}
		log.WithFields(log.Fields{
			"func":      "Payload.Candidate",
			"block":     block,
			"tx":        tx,
			"candidate": tx.from.String(),
		}).Info("Candidate Logout.")
	default:
		return counter, ErrInvalidCandidatePayloadAction
	}
	return counter, nil
}
