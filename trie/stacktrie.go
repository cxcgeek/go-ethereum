// Copyright 2020 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// ReStackTrie is a reimplementation of the Stacktrie, that fixes
// bugs in the previous implementation, and which also implements
// its own hashing mechanism which is more specific and hopefully
// more efficient that the default hasher.
type ReStackTrie struct {
	nodeType  uint8            // node type (as in branch, ext, leaf)
	val       []byte           // value contained by this node if it's a leaf
	key       []byte           // key chunk covered by this (full|ext) node
	keyOffset int              // offset of the key chunk inside a full key
	children  [16]*ReStackTrie // list of children (for fullnodes and exts)
}

// NewReStackTrie allocates and initializes an empty trie.
func NewReStackTrie() *ReStackTrie {
	return &ReStackTrie{
		nodeType: 3,
	}
}

// List all values that ReStackTrie#nodeType can hold
const (
	branchNode = iota
	extNode
	leafNode
	emptyNode
	hashedNode
)

func (st *ReStackTrie) TryUpdate(key, value []byte) error {
	k := keybytesToHex(key)
	if len(value) == 0 {
		panic("deletion not supported")
	}
	st.insert(k[:len(k)-1], value)
	return nil
}

// Helper function that, given a full key, determines the index
// at which the chunk pointed by st.keyOffset is different from
// the same chunk in the full key.
func (st *ReStackTrie) getDiffIndex(key []byte) int {
	diffindex := 0
	for ; diffindex < len(st.key) && st.key[diffindex] == key[st.keyOffset+diffindex]; diffindex++ {
	}
	return diffindex
}

// Helper function to that inserts a (key, value) pair into
// the trie.
func (st *ReStackTrie) insert(key, value []byte) {
	switch st.nodeType {
	case branchNode: /* Branch */
		idx := int(key[st.keyOffset])
		if st.children[idx] == nil {
			st.children[idx] = NewReStackTrie()
			st.children[idx].keyOffset = st.keyOffset + 1
		}
		for i := idx - 1; i >= 0; i-- {
			if st.children[i] != nil {
				if st.children[i].nodeType != hashedNode {
					st.children[i].val = st.children[i].hash()
					st.children[i].key = nil
					st.children[i].nodeType = hashedNode
				}

				break
			}

		}
		st.children[idx].insert(key, value)
	case extNode: /* Ext */
		// Compare both key chunks and see where they differ
		diffidx := st.getDiffIndex(key)

		// Check if chunks are identical. If so, recurse into
		// the child node. Otherwise, the key has to be split
		// into 1) an optional common prefix, 2) the fullnode
		// representing the two differing path, and 3) a leaf
		// for each of the differentiated subtrees.
		if diffidx == len(st.key) {
			// Ext key and key segment are identical, recurse into
			// the child node.
			st.children[0].insert(key, value)
			return
		}
		// Save the original part. Depending if the break is
		// at the extension's last byte or not, create an
		// intermediate extension or use the extension's child
		// node directly.
		var n *ReStackTrie
		if diffidx < len(st.key)-1 {
			n = NewReStackTrie()
			n.key = st.key[diffidx+1:]
			n.children[0] = st.children[0]
			n.nodeType = extNode
		} else {
			// Break on the last byte, no need to insert
			// an extension node: reuse the current node
			n = st.children[0]
		}
		n.keyOffset = st.keyOffset + diffidx + 1

		var p *ReStackTrie
		if diffidx == 0 {
			// the break is on the first byte, so
			// the current node is converted into
			// a branch node.
			st.children[0] = nil
			p = st
			st.nodeType = branchNode
		} else {
			// the common prefix is at least one byte
			// long, insert a new intermediate branch
			// node.
			st.children[0] = NewReStackTrie()
			st.children[0].nodeType = branchNode
			st.children[0].keyOffset = st.keyOffset + diffidx
			p = st.children[0]
		}

		n.val = n.hash()
		n.nodeType = hashedNode
		n.key = nil

		// Create a leaf for the inserted part
		o := NewReStackTrie()
		o.keyOffset = st.keyOffset + diffidx + 1
		o.key = key[o.keyOffset:]
		o.val = value
		o.nodeType = leafNode

		// Insert both child leaves where they belong:
		origIdx := st.key[diffidx]
		newIdx := key[diffidx+st.keyOffset]
		p.children[origIdx] = n
		p.children[newIdx] = o
		st.key = st.key[:diffidx]

	case leafNode: /* Leaf */
		// Compare both key chunks and see where they differ
		diffidx := st.getDiffIndex(key)

		// Overwriting a key isn't supported, which means that
		// the current leaf is expected to be split into 1) an
		// optional extension for the common prefix of these 2
		// keys, 2) a fullnode selecting the path on which the
		// keys differ, and 3) one leaf for the differentiated
		// component of each key.
		if diffidx >= len(st.key) {
			panic("Trying to insert into existing key")
		}

		// Check if the split occurs at the first nibble of the
		// chunk. In that case, no prefix extnode is necessary.
		// Otherwise, create that
		var p *ReStackTrie
		if diffidx == 0 {
			// Convert current leaf into a branch
			st.nodeType = branchNode
			p = st
			st.children[0] = nil
		} else {
			// Convert current node into an ext,
			// and insert a child branch node.
			st.nodeType = extNode
			st.children[0] = NewReStackTrie()
			st.children[0].nodeType = branchNode
			st.children[0].keyOffset = st.keyOffset + diffidx
			p = st.children[0]
		}

		// Create the two child leaves: the one containing the
		// original value and the one containing the new value
		// The child leave will be hashed directly in order to
		// free up some memory.
		origIdx := st.key[diffidx]
		p.children[origIdx] = NewReStackTrie()
		p.children[origIdx].nodeType = leafNode
		p.children[origIdx].key = st.key[diffidx+1:]
		p.children[origIdx].val = st.val
		p.children[origIdx].keyOffset = p.keyOffset + 1

		p.children[origIdx].val = p.children[origIdx].hash()
		p.children[origIdx].nodeType = hashedNode
		p.children[origIdx].key = nil

		newIdx := key[diffidx+st.keyOffset]
		p.children[newIdx] = NewReStackTrie()
		p.children[newIdx].nodeType = leafNode
		p.children[newIdx].key = key[p.keyOffset+1:]
		p.children[newIdx].val = value
		p.children[newIdx].keyOffset = p.keyOffset + 1

		st.key = st.key[:diffidx]
	case emptyNode: /* Empty */
		st.nodeType = leafNode
		st.key = key[st.keyOffset:]
		st.val = value
	case hashedNode:
		panic("trying to insert into hash")
	default:
		panic("invalid type")
	}
}

// rawHPRLP is called when the length of the RLP of a node is
// less than 32. It will return the un-hashed payload.
func rawHPRLP(key, val []byte, leaf bool) []byte {
	payload := [32]byte{}

	// payload size - none of the components are larger
	// than 56 since the whole size is smaller than 32
	val_header_len := 1
	if len(val) == 1 && val[0] < 128 {
		val_header_len = 0
	}
	var key_header_len int // Length taken by the size header
	var hplen int          // Length taken by the hex prefix (0 or 1)
	switch len(key) {
	case 0:
		// Key is 0 byte long, so the whole byte will
		// be less than 128, no header needed. No HP
		// will be present either.
		key_header_len = 0
		hplen = 0
	case 1:
		// Key is 1 byte long, so the whole byte will
		// be less than 128, no header needed. The HP
		// will be present, however.
		key_header_len = 0
		hplen = 1
	default:
		// Key is longer than 1 byte, and including
		// the hex prefix that's more than 2 bytes.
		// And of course it's less than 56 bytes.
		key_header_len = 1
		hplen = 1
	}
	payload[0] = 192 + byte(hplen+key_header_len+len(key)/2+val_header_len+len(val))

	pos := 1
	// Add key, if present
	if len(key) > 0 {
		// add length prefix if needed
		if key_header_len == 1 {
			payload[1] = byte(1 + len(key)/2)
			pos++
		}

		// hex prefix
		payload[pos] = byte(16 * (len(key) % 2))
		if leaf {
			payload[pos] |= 32
		}
		if len(key)%2 == 0 {
			// Advance to next byte iff the key has an odd length
			// of nibbles.
			pos++
		}

		// copy key data
		for i, nibble := range key {
			offset := 1 - uint((len(key)+i)%2)
			payload[pos] |= byte(int(nibble) << (4 * offset))
			if offset == 0 {
				pos++
			}
		}
	}

	// copy value data. If the payload isn't a single byte
	// lower than 128, also add the header.
	if len(val) > 1 || val[0] >= 128 {
		payload[pos] = byte(len(val))
		if len(val) > 1 {
			payload[pos] += 128
		}
		pos += 1

	}
	copy(payload[pos:], val)
	pos += len(val)

	// In case the payload is only one byte,
	// no header is needed.
	if pos == 2 {
		return payload[1:pos]
	}

	return payload[:pos]
}

// writeEvenHP writes a key with its hex prefix into a writer (presumably, the
// input of a hasher) and then writes the value. The value can be a maximum of
// 256 bytes, as it is only concerned with writing account leaves and optimize
// for this use case.
func writeHPRLP(writer io.Writer, key, val []byte, leaf bool) {
	// DEBUG don't remove yet
	//var writer bytes.Buffer

	// Determine the _t_ part of the hex prefix
	hp := byte(0)
	if leaf {
		hp = 32
	}

	const maxHeaderSize = 1 /* key byte list header */ +
		1 /* list header for key + value */ +
		1 /* potential size byte if total size > 56 */ +
		1 /* hex prefix if key is even-length*/
	header := [maxHeaderSize]byte{}
	keyOffset := 0
	headerPos := maxHeaderSize - 1

	// Add the hex prefix to its own byte if the key length is even, and
	// as the most significant nibble of the key if it's odd.
	// In the latter case, the first nibble of the key will be part of
	// the header and it will be skipped later when it's added to the
	// hasher sponge.
	if len(key)%2 == 0 {
		header[headerPos] = hp
	} else {
		header[headerPos] = hp | key[0] | 16
		keyOffset = 1
	}
	headerPos--

	// Add the key byte header, the key is 32 bytes max so it's always
	// under 56 bytes - no extra byte needed.
	keyByteSize := byte(len(key) / 2)
	if len(key) > 1 || header[len(header)-1] > 128 {
		header[headerPos] = 0x80 + keyByteSize + 1 /* HP */
		headerPos--
	}

	// If this is a leaf being inserted, the header length for the
	// value part will be two bytes as the leaf is more than 56 bytes
	// long.
	valHeaderLen := 1
	if len(val) == 1 && val[0] < 128 {
		// Don't reserve space for the header if this
		// is an integer < 128
		valHeaderLen = 0
	}
	if len(val) > 56 {
		valHeaderLen = 2
	}

	// Add the global header, with optional length, and specify at
	// which byte the header is starting.
	payloadSize := int(keyByteSize) + (len(header) - headerPos - 1) +
		valHeaderLen + len(val) /* value + rlp header */
	var start int
	if payloadSize > 56 {
		header[headerPos] = byte(payloadSize)
		headerPos--
		header[headerPos] = 0xf8
		start = headerPos
	} else {
		header[headerPos] = 0xc0 + byte(payloadSize)
		start = headerPos
	}

	// Write the header into the sponge
	writer.Write(header[start:])

	// Write the key into the sponge
	var m byte
	for i, nibble := range key {
		// Skip the first byte if the key has an odd-length, since
		// it has already been written with the header.
		if i >= keyOffset {
			if (i-keyOffset)%2 == 0 {
				m = nibble
			} else {
				writer.Write([]byte{m*16 + nibble})
			}
		}
	}

	// Write the RLP prefix to the value if needed
	if len(val) > 56 {
		writer.Write([]byte{0xb8, byte(len(val))})
	} else if len(val) > 1 || val[0] >= 128 {
		writer.Write([]byte{0x80 + byte(len(val))})
	}
	writer.Write(val)

	// DEBUG don't remove yet
	//if leaf {
	//fmt.Println("leaf rlp ", writer)
	//} else {
	//fmt.Println("ext rlp ", writer)
	//}
	//io.Copy(w, &writer)
}

func (st *ReStackTrie) hash() []byte {
	/* Shortcut if node is already hashed */
	if st.nodeType == hashedNode {
		return st.val
	}

	d := sha3.NewLegacyKeccak256()
	switch st.nodeType {
	case branchNode:
		payload := [544]byte{}
		pos := 3 // maximum header length given what we know
		for i, v := range st.children {
			if v != nil {
				// Write a 32 byte list to the sponge
				childhash := v.hash()
				if len(childhash) == 1 && childhash[0] < 128 {
					payload[pos] = childhash[0]
				} else {
					payload[pos] = 128 + byte(len(childhash))
					pos++
				}
				copy(payload[pos:pos+len(childhash)], childhash)
				pos += len(childhash)
				st.children[i] = nil // Reclaim mem from subtree
			} else {
				// Write an empty list to the sponge
				payload[pos] = 0x80
				pos++
			}
		}
		// Add an empty 17th value
		payload[pos] = 0x80
		pos++

		// Compute the header, length size is either 0, 1 or 2 bytes since
		// there are at least 17 empty list headers, and at most 16 hashes
		// plus an empty header for the value.
		var start int
		if pos-3 < 56 {
			payload[2] = 0xc0 + byte(pos-3)
			start = 2
		} else if pos-3 < 256 {
			payload[2] = byte(pos - 3)
			payload[1] = 0xf8
			start = 1
		} else {
			payload[2] = byte(pos - 3)
			payload[1] = byte((pos - 3) >> 8)
			payload[0] = 0xf9
			start = 0
		}

		// Do not hash if the payload length is less than 32 bytes
		if pos-start < 32 {
			return payload[start:pos]
		}
		d.Write(payload[start:pos])
	case extNode:
		ch := st.children[0].hash()
		if (len(st.key)/2)+1+len(ch) < 29 {
			rlp := [32]byte{}

			// Copy key data, including hex prefix
			if len(st.key)%2 == 0 {
				rlp[2] = 0
			} else {
				rlp[2] = 16 + st.key[0]
			}
			rlp[1] = byte(128 + 1 + len(st.key)/2)
			for i := 0; i < len(st.key); i++ {
				rlp[3-len(st.key)%2+i/2] = st.key[i] << uint(4*((i+1+len(st.key))%2))
			}

			copy(rlp[3+len(st.key)/2:], ch)
			rlp[0] = byte(192 + 2 + len(st.key)/2 + len(ch))

			return rlp[:3+len(st.key)/2+len(ch)]
		}
		writeHPRLP(d, st.key, ch, false)
		st.children[0] = nil // Reclaim mem from subtree
	case leafNode:
		if (len(st.key)/2)+1+len(st.val) < 29 {
			return rawHPRLP(st.key, st.val, true)
		}
		writeHPRLP(d, st.key, st.val, true)
	case emptyNode:
		return emptyRoot[:]
	default:
		panic("Invalid node type")
	}
	return d.Sum(nil)
}

func (st *ReStackTrie) Hash() (h common.Hash) {
	return common.BytesToHash(st.hash())
}
