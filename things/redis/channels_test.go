//
// Copyright (c) 2018
// Mainflux
//
// SPDX-License-Identifier: Apache-2.0
//

package redis_test

import (
	"fmt"
	"testing"

	"github.com/mainflux/mainflux/things/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	channelCache := redis.NewChannelCache(cacheClient)

	cid := uint64(123)
	tid := uint64(321)

	cases := []struct {
		desc string
		cid  uint64
		tid  uint64
	}{
		{
			desc: "connect thing to channel",
			cid:  cid,
			tid:  tid,
		},
		{
			desc: "connect already connected thing to channel",
			cid:  cid,
			tid:  tid,
		},
	}
	for _, tc := range cases {
		err := channelCache.Connect(cid, tid)
		assert.Nil(t, err, fmt.Sprintf("%s: fail to connect due to: %s\n", tc.desc, err))
	}
}

func TestHasThing(t *testing.T) {
	channelCache := redis.NewChannelCache(cacheClient)

	cid := uint64(123)
	tid := uint64(321)

	err := channelCache.Connect(cid, tid)
	require.Nil(t, err, fmt.Sprintf("connect thing to channel: fail to connect due to: %s\n", err))

	cases := map[string]struct {
		cid       uint64
		tid       uint64
		hasAccess bool
	}{
		"access check for thing that has access": {
			cid:       cid,
			tid:       tid,
			hasAccess: true,
		},
		"access check for thing without access": {
			cid:       cid,
			tid:       cid,
			hasAccess: false,
		},
		"access check for non-existing channel": {
			cid:       tid,
			tid:       tid,
			hasAccess: false,
		},
	}

	for desc, tc := range cases {
		hasAccess := channelCache.HasThing(tc.cid, tc.tid)
		assert.Equal(t, tc.hasAccess, hasAccess, fmt.Sprintf("%s: expected %t got %t\n", desc, tc.hasAccess, hasAccess))
	}
}
func TestDisconnect(t *testing.T) {
	channelCache := redis.NewChannelCache(cacheClient)

	cid := uint64(123)
	tid := uint64(321)
	tid2 := uint64(322)

	err := channelCache.Connect(cid, tid)
	require.Nil(t, err, fmt.Sprintf("connect thing to channel: fail to connect due to: %s\n", err))

	cases := []struct {
		desc      string
		cid       uint64
		tid       uint64
		hasAccess bool
	}{
		{
			desc:      "disconnecting connected thing",
			cid:       cid,
			tid:       tid,
			hasAccess: false,
		},
		{
			desc:      "disconnecting non-connected thing",
			cid:       cid,
			tid:       tid2,
			hasAccess: false,
		},
	}
	for _, tc := range cases {
		err := channelCache.Disconnect(tc.cid, tc.tid)
		assert.Nil(t, err, fmt.Sprintf("%s: fail due to: %s\n", tc.desc, err))

		hasAccess := channelCache.HasThing(tc.cid, tc.tid)
		assert.Equal(t, tc.hasAccess, hasAccess, fmt.Sprintf("access check after %s: expected %t got %t\n", tc.desc, tc.hasAccess, hasAccess))
	}
}

func TestRemove(t *testing.T) {
	channelCache := redis.NewChannelCache(cacheClient)

	cid := uint64(123)
	cid2 := uint64(124)
	tid := uint64(321)

	err := channelCache.Connect(cid, tid)
	require.Nil(t, err, fmt.Sprintf("connect thing to channel: fail to connect due to: %s\n", err))

	cases := []struct {
		desc      string
		cid       uint64
		tid       uint64
		err       error
		hasAccess bool
	}{
		{
			desc:      "Remove channel from cache",
			cid:       cid,
			tid:       tid,
			err:       nil,
			hasAccess: false,
		},
		{
			desc:      "Remove non-cached channel from cache",
			cid:       cid2,
			tid:       tid,
			err:       nil,
			hasAccess: false,
		},
	}

	for _, tc := range cases {
		err := channelCache.Remove(tc.cid)
		assert.Nil(t, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		hasAcces := channelCache.HasThing(tc.cid, tc.tid)
		assert.Equal(t, tc.hasAccess, hasAcces, "%s - check access after removing channel: expected %t got %t\n", tc.desc, tc.hasAccess, hasAcces)
	}
}
