/**
 *  Copyright 2014 Paul Querna
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */

package memory

import (
	"bytes"
	"testing"
)

func TestCache(t *testing.T) {
	c := New()
	v := c.Get("foobar")
	if v != nil {
		t.Fatalf("should't of found key foobar in %v", c)
		return
	}

	v = c.Alloc(1024)
	v.Key = "foo"
	for i := 0; i < 1024; i++ {
		v.Data[0][i] = 'a'
	}
	c.Set(v)

	b := c.Get("foo")

	if b == nil {
		t.Fatalf("missing key 'foo' in %v", c)
		return
	}
	cmp := bytes.Compare(v.Data[0], b.Data[0])

	if cmp != 0 {
		t.Fatalf("buffers didn't compare: %v  %v != %v", cmp, v, b)
	}
	b.Unref()
}

func TestBig(t *testing.T) {
	c := New()
	v := c.Alloc((MAX_BLOCK_SIZE * 5) + 1)

	if len(v.Data) != 6 {
		t.Fatalf("expected len=6, got len=%v", len(v.Data))
	}
	if len(v.Data[5]) != 1 {
		t.Fatalf("expected block at offset 5 to be of length 1, got %v", len(v.Data[5]))
	}
}
