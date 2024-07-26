// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	len  int64
	want int
}{
	{5, 5},
	{10, 10},
	{100, 100},
}

func TestString(t *testing.T) {
	for index, test := range tests {
		t.Run(fmt.Sprintf("Test Name %d", index), func(t *testing.T) {
			got := String(test.len)
			assert.Equal(t, test.len, int64(len(got)))
		})
	}
}

func TestAlpha(t *testing.T) {
	for index, test := range tests {
		t.Run(fmt.Sprintf("Test Alpha %d", index), func(t *testing.T) {
			got := Alpha(test.len)
			assert.Equal(t, test.len, int64(len(got)))
		})
	}
}

func TestAlphaNum(t *testing.T) {
	for index, test := range tests {
		t.Run(fmt.Sprintf("Test AlphaNum %d", index), func(t *testing.T) {
			got := AlphaNum(test.len)
			assert.Equal(t, test.len, int64(len(got)))
		})
	}
}

func TestSeq(t *testing.T) {
	got := Seq(20)
	assert.Len(t, got, 20)
}
