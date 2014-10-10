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
	"github.com/cznic/bufs"

	"container/list"
	"sync"
	"sync/atomic"
)

type Value struct {
	Key  string
	Data []byte
	ref  int32
}

func (v *Value) Unref() {
	atomic.AddInt32(&v.ref, -1)

	c := atomic.LoadInt32(&v.ref)
	if c == 0 {
		bufs.GCache.Put(v.Data)
	}
}

func (v *Value) Ref() {
	atomic.AddInt32(&v.ref, 1)
}

type Cache struct {
	mtx  sync.Mutex
	kv   map[string]*list.Element
	list *list.List
}

func New() *Cache {
	return &Cache{
		kv:   make(map[string]*list.Element),
		list: list.New(),
	}
}

func (c *Cache) unlink(elem *list.Element) {
	// LOCK HELD BY CALLER
	v := c.list.Remove(elem).(*Value)
	v.Unref()
}

func (c *Cache) Set(v *Value) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	oldelem, ok := c.kv[v.Key]

	if ok {
		c.unlink(oldelem)
	}

	elem := c.list.PushFront(v)
	c.kv[v.Key] = elem
}

func (c *Cache) Get(key string) *Value {
	// Returns a Value or nil. If Value, the caller
	// must call v.Unref() when they are done using Value.
	c.mtx.Lock()
	defer c.mtx.Unlock()

	elem, ok := c.kv[key]

	if !ok {
		return nil
	}

	c.list.MoveToFront(elem)
	v := elem.Value.(*Value)
	v.Ref()

	return v
}
