// Copyright 2016 The go-ethereum Authors
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

// Code generated by "stringer -type=nodeEvent"; DO NOT EDIT

package discover

import "fmt"

const (
	_nodeEvent_name_0 = "invalidEventpingPacketpongPacketfindnodePacketneighborsPacket"
	_nodeEvent_name_1 = "pongTimeoutpingTimeoutneighboursTimeout"
)

var (
	_nodeEvent_index_0 = [...]uint8{0, 12, 22, 32, 46, 61}
	_nodeEvent_index_1 = [...]uint8{0, 11, 22, 39}
)

func (i nodeEvent) String() string {
	switch {
	case 0 <= i && i <= 4:
		return _nodeEvent_name_0[_nodeEvent_index_0[i]:_nodeEvent_index_0[i+1]]
	case 261 <= i && i <= 263:
		i -= 261
		return _nodeEvent_name_1[_nodeEvent_index_1[i]:_nodeEvent_index_1[i+1]]
	default:
		return fmt.Sprintf("nodeEvent(%d)", i)
	}
}
