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

package nvm

import "C"

import (
	"regexp"
	"unsafe"

	"github.com/nebulasio/go-nebulas/common/trie"

	log "github.com/sirupsen/logrus"
)

var (
	keyPattern = regexp.MustCompile("^@([a-zA-Z_].*?)\\[(.+?)\\]$")
)

// hashStorageKey return the key hash.
// There are two kinds of key, the one is ItemKey, the other is Map-ItemKey.
// ItemKey in SmartContract is used for object storage.
// For example, the ItemKey for the statement "token.totalSupply = 1000" is "totalSupply".
// Map-ItemKey in SmartContrat is used for Map storage.
// For example, the Map-ItemKey for the statement "token.balances.set('addr1', 100)" is "@balances[addr1]".
func hashStorageKey(key string) []byte {
	var domainKey, itemKey string

	matches := keyPattern.FindAllStringSubmatch(key, -1)
	if matches == nil {
		domainKey = ""
		itemKey = key
	} else {
		domainKey = matches[0][1]
		itemKey = matches[0][2]
	}

	return trie.HashDomains(domainKey, itemKey)
}

//export StorageGetFunc
func StorageGetFunc(handler unsafe.Pointer, key *C.char) *C.char {
	_, storage := getEngineAndStorage(uint64(uintptr(handler)))
	if storage == nil {
		return nil
	}

	val, err := storage.Get([]byte(hashStorageKey(C.GoString(key))))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "nvm.StorageGetFunc",
			"handler": uint64(uintptr(handler)),
			"key":     C.GoString(key),
		}).Warn("StorageGetFunc get key failed.")
		return nil
	}
	return C.CString(string(val))
}

//export StoragePutFunc
func StoragePutFunc(handler unsafe.Pointer, key *C.char, value *C.char) int {
	_, storage := getEngineAndStorage(uint64(uintptr(handler)))
	if storage == nil {
		return 1
	}

	// log.Errorf("[--------------] StoragePutFunc, storage = %v; {%v: %v}", storage, C.GoString(key), C.GoString(value))

	_, err := storage.Put([]byte(hashStorageKey(C.GoString(key))), []byte(C.GoString(value)))
	if err != nil {
		log.WithFields(log.Fields{
			"func":    "nvm.StoragePutFunc",
			"handler": uint64(uintptr(handler)),
			"key":     C.GoString(key),
			"err":     err,
		}).Error("StoragePutFunc put key failed.")
		return 1
	}
	return 0
}

//export StorageDelFunc
func StorageDelFunc(handler unsafe.Pointer, key *C.char) int {
	_, storage := getEngineAndStorage(uint64(uintptr(handler)))
	if storage == nil {
		return 1
	}

	_, err := storage.Del([]byte(hashStorageKey(C.GoString(key))))

	if err != nil {
		log.WithFields(log.Fields{
			"func":    "nvm.StorageDelFunc",
			"handler": uint64(uintptr(handler)),
			"key":     C.GoString(key),
			"err":     err,
		}).Warn("StorageDelFunc del key failed.")
		return 1
	}

	return 0
}
