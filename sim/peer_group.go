// Copyright (c) 2018 Aidos Developer

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// This is a rewrite of https://github.com/ripple/rippled/src/test/csf
// covered by:
//------------------------------------------------------------------------------
/*
   This file is part of rippled: https://github.com/ripple/rippled
   Copyright (c) 2012-2017 Ripple Labs Inc.

   Permission to use, copy, modify, and/or distribute this software for any
   purpose  with  or without fee is hereby granted, provided that the above
   copyright notice and this permission notice appear in all copies.

   THE  SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
   WITH  REGARD  TO  THIS  SOFTWARE  INCLUDING  ALL  IMPLIED  WARRANTIES  OF
   MERCHANTABILITY  AND  FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
   ANY  SPECIAL ,  DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
   WHATSOEVER  RESULTING  FROM  LOSS  OF USE, DATA OR PROFITS, WHETHER IN AN
   ACTION  OF  CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
   OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
*/
//==============================================================================

package sim

import (
	"bytes"
	"sort"
	"time"
)

type peerGroup struct {
	peers []*peer
}

func (pg *peerGroup) sort() {
	sort.Slice(pg.peers, func(i, j int) bool {
		return bytes.Compare(pg.peers[i].id[:], pg.peers[j].id[:]) < 0
	})
}

func (pg *peerGroup) trust(o *peerGroup) {
	for _, p := range pg.peers {
		for _, target := range o.peers {
			p.trustGraph.trust(p, target)
		}
	}
}

/** Establish network connection
  Establish outbound connections from all peers in this group to all peers in
  o. If a connection already exists, no new connection is established.
  @param o The group of peers to connect to (will get inbound connections)
  @param delay The fixed messaging delay for all established connections
*/
func (pg *peerGroup) connect(o *peerGroup, delay time.Duration) {
	for _, p := range pg.peers {
		for _, target := range o.peers {
			if p != target {
				p.net.connect(p, target, delay)
			}
		}
	}
}
func (pg *peerGroup) disconnect(o *peerGroup) {
	for _, p := range pg.peers {
		for _, target := range o.peers {
			if p != target {
				p.net.disconnect(p, target)
			}
		}
	}
}

/** Establish trust and network connection

  Establish trust and create a network connection with fixed delay
  from all peers in this group to all peers in o

  @param o The group of peers to trust and connect to
  @param delay The fixed messaging delay for all established connections
*/
func (pg *peerGroup) trustAndConnect(o *peerGroup, delay time.Duration) {
	pg.trust(o)
	pg.connect(o, delay)
}

func (pg *peerGroup) append(o *peerGroup) *peerGroup {
	peers := make([]*peer, len(pg.peers)+len(o.peers))
	copy(peers, pg.peers)
	copy(peers[len(pg.peers):], o.peers)
	pg2 := &peerGroup{
		peers: peers,
	}
	pg2.sort()
	return pg2
}
