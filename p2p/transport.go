// Copyright 2015 The go-ethereum Authors
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

package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	// total timeout for encryption handshake and protocol
	// handshake in both directions.
	handshakeTimeout = 5 * time.Second

	// This is the timeout for sending the disconnect reason.
	// This is shorter than the usual timeout because we don't want
	// to wait if the connection is known to be bad anyway.
	discWriteTimeout = 1 * time.Second
)

// rlpxTransport is the transport used by actual (non-test) connections.
// It wraps an RLPx connection with locks and read/write deadlines.
type rlpxTransport struct {
	rmu, wmu sync.Mutex
	conn     *rlpx.Conn
}

func newRLPX(conn net.Conn, dialDest *ecdsa.PublicKey) transport {
	conn.SetDeadline(time.Now().Add(handshakeTimeout))
	return &rlpxTransport{conn: rlpx.NewConn(conn, dialDest)}
}

func (t *rlpxTransport) ReadMsg() (Msg, error) {
	t.rmu.Lock()
	defer t.rmu.Unlock()

	t.conn.SetReadDeadline(time.Now().Add(frameReadTimeout))

	var (
		msg Msg
		err error
	)

	msg.Code, msg.Size, msg.Payload, err = t.conn.ReadMsg()
	msg.meterSize = msg.Size
	msg.ReceivedAt = time.Now()

	return msg, err
}

func (t *rlpxTransport) WriteMsg(msg Msg) error {
	t.wmu.Lock()
	defer t.wmu.Unlock()

	t.conn.SetWriteDeadline(time.Now().Add(frameWriteTimeout))
	// write message
	size, err := t.conn.WriteMsg(msg.Code, msg.Size, msg.Payload)
	if err != nil {
		return err
	}
	// set metrics
	msg.meterSize = size
	if metrics.Enabled && msg.meterCap.Name != "" { // don't meter non-subprotocol messages
		m := fmt.Sprintf("%s/%s/%d/%#02x", egressMeterName, msg.meterCap.Name, msg.meterCap.Version, msg.meterCode)
		metrics.GetOrRegisterMeter(m, nil).Mark(int64(msg.meterSize))
		metrics.GetOrRegisterMeter(m+"/packets", nil).Mark(1)
	}
	return nil
}

func (t *rlpxTransport) close(err error) {
	t.wmu.Lock()
	defer t.wmu.Unlock()

	// Tell the remote end why we're disconnecting if possible.
	if t.conn != nil {
		if r, ok := err.(DiscReason); ok && r != DiscNetworkError {
			// frameRW tries to send DiscReason to disconnected peer
			// if the connection is net.Pipe (in-memory simulation)
			// it hangs forever, since net.Pipe does not implement
			// a write deadline. Because of this only try to send
			// the disconnect reason message if there is no error.
			deadline := time.Now().Add(discWriteTimeout)
			if err := t.conn.SetWriteDeadline(deadline); err == nil {
				size, data, _ := rlp.EncodeToReader([]interface{}{r})
				t.conn.WriteMsg(discMsg, uint32(size), data)
			}
		}
	}
	t.conn.Close()
}

func (t *rlpxTransport) doEncHandshake(prv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	return t.conn.Handshake(prv)
}

func (t *rlpxTransport) doProtoHandshake(our *protoHandshake) (their *protoHandshake, err error) {
	// Writing our handshake happens concurrently, we prefer
	// returning the handshake read error. If the remote side
	// disconnects us early with a valid reason, we should return it
	// as the error so it can be tracked elsewhere.
	werr := make(chan error, 1)
	go func() { werr <- Send(t, handshakeMsg, our) }()
	if their, err = readProtocolHandshake(t); err != nil {
		<-werr // make sure the write terminates too
		return nil, err
	}
	if err := <-werr; err != nil {
		return nil, fmt.Errorf("write error: %v", err)
	}
	// If the protocol version supports Snappy encoding, upgrade immediately
	t.conn.SetSnappy(their.Version >= snappyProtocolVersion)

	return their, nil
}

func readProtocolHandshake(rw MsgReader) (*protoHandshake, error) {
	msg, err := rw.ReadMsg()
	if err != nil {
		return nil, err
	}
	if msg.Size > baseProtocolMaxMsgSize {
		return nil, fmt.Errorf("message too big")
	}
	if msg.Code == discMsg {
		// Disconnect before protocol handshake is valid according to the
		// spec and we send it ourself if the post-handshake checks fail.
		// We can't return the reason directly, though, because it is echoed
		// back otherwise. Wrap it in a string instead.
		var reason [1]DiscReason
		rlp.Decode(msg.Payload, &reason)
		return nil, reason[0]
	}
	if msg.Code != handshakeMsg {
		return nil, fmt.Errorf("expected handshake, got %x", msg.Code)
	}
	var hs protoHandshake
	if err := msg.Decode(&hs); err != nil {
		return nil, err
	}
	if len(hs.ID) != 64 || !bitutil.TestBytes(hs.ID) {
		return nil, DiscInvalidIdentity
	}
	return &hs, nil
}
