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

package vulcancache

import (
	"github.com/mailgun/vulcand/Godeps/_workspace/src/github.com/mailgun/vulcan/request"

	"errors"
	"net/http"
)

type HashFunc func(v string)

type Hash interface {
	// opaque hashing function, call for each string you
	// want included to make a request unique.
	Hash(r request.Request, hf HashFunc)
}

type RecvCmd int

const (
	RECV_LOOKUP RecvCmd = iota
	RECV_PASS   RecvCmd = iota
	RECV_REJECT RecvCmd = iota
)

type Recv interface {
	Recv(r request.Request) RecvCmd
}

type CacheGetCmd int

const (
	CACHE_GET_HIT  CacheGetCmd = iota
	CACHE_GET_MISS CacheGetCmd = iota
	CACHE_GET_PASS CacheGetCmd = iota
)

type CacheStoreCmd int

const (
	CACHE_STORE_STORED CacheStoreCmd = iota
	CACHE_STORE_MISS   CacheStoreCmd = iota
)

type CachedResponse interface {
	// TODO: more.
}

type Cache interface {
	Get(key string, r request.Request) (CacheGetCmd, CachedResponse)
	Store(key string, r request.Request, a request.Attempt) CacheStoreCmd
}

type MissCmd int

const (
	MISS_FETCH MissCmd = iota
	MISS_PASS  MissCmd = iota
)

type Miss interface {
	Miss(r request.Request) MissCmd
}

type HttpCacher struct {
	recv  Recv
	hash  Hash
	cache Cache
	miss  Miss
}

type CacheOptions struct {
	Recv  Recv
	Hash  Hash
	Cache Cache
	Miss  Miss
	// TODO: other options
}

func NewHttpCacher(co *CacheOptions) *HttpCacher {
	bc := newBasicCache(co)

	cache := &HttpCacher{
		recv:  bc,
		hash:  bc,
		cache: bc,
		miss:  bc,
	}

	if co.Recv != nil {
		cache.recv = co.Recv
	}

	if co.Hash != nil {
		cache.hash = co.Hash
	}

	if co.Cache != nil {
		cache.cache = co.Cache
	}

	if co.Miss != nil {
		cache.miss = co.Miss
	}

	return cache
}

func (hc *HttpCacher) lookup(r request.Request) (*http.Response, error) {
	// TODO: hashme
	return nil, nil
}

func (hc *HttpCacher) pass(r request.Request) (*http.Response, error) {
	// TODO: pass
	return nil, nil
}

func (hc *HttpCacher) reject(r request.Request) (*http.Response, error) {
	// TODO: better reject
	return nil, errors.New("rejected")
}

func (hc *HttpCacher) ProcessRequest(r request.Request) (*http.Response, error) {
	recvCmd := hc.recv.Recv(r)
	switch recvCmd {
	case RECV_LOOKUP:
		return hc.lookup(r)
	case RECV_PASS:
		return hc.pass(r)
	case RECV_REJECT:
		return hc.reject(r)
	}
	panic("not reached")
}

func (hc *HttpCacher) ProcessResponse(r request.Request, a request.Attempt) {

}
