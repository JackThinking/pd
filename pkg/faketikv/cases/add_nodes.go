// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package cases

import (
	"github.com/pingcap/kvproto/pkg/metapb"
	"github.com/pingcap/pd/server/core"
	log "github.com/sirupsen/logrus"
)

func newAddNodes() *Conf {
	var conf Conf

	for i := 1; i <= 8; i++ {
		conf.Stores = append(conf.Stores, Store{
			ID:        uint64(i),
			Status:    metapb.StoreState_Up,
			Capacity:  10 * gb,
			Available: 9 * gb,
		})
	}

	var id idAllocator
	id.setMaxID(20)
	for i := 0; i < 1000; i++ {
		peers := []*metapb.Peer{
			{Id: id.nextID(), StoreId: uint64(i)%4 + 1},
			{Id: id.nextID(), StoreId: uint64(i+1)%4 + 1},
			{Id: id.nextID(), StoreId: uint64(i+2)%4 + 1},
		}
		conf.Regions = append(conf.Regions, Region{
			ID:     id.nextID(),
			Peers:  peers,
			Leader: peers[0],
			Size:   96 * mb,
		})
	}
	conf.MaxID = id.maxID

	conf.Checker = func(regions *core.RegionsInfo) bool {
		res := true
		counts := make([]int, 0, 8)
		for i := 1; i <= 8; i++ {
			count := regions.GetStoreLeaderCount(uint64(i))
			counts = append(counts, count)
			if count > 130 || count < 120 {
				res = false
			}
		}
		log.Infof("leader counts: %v", counts)
		return res
	}
	return &conf
}