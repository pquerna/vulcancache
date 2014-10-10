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

func (hc *HttpCacher) doMiss(r request.Request) (*http.Response, error) {
	// OK, missed
	mcmd := hc.miss.Miss(r)
	switch mcmd {
	case MISS_FETCH:
		// TOOD: should there be a doFetch?
		return nil, nil
	case MISS_PASS:
		return hc.doPass(r)
	}
	panic("not reached")
}

func (hc *HttpCacher) doLookup(r request.Request) (*http.Response, error) {
	// TODO: hashme
	// until hashing done, assume miss.
	//	getCmd := hc.Hash(r, hf)
	getCmd := CACHE_GET_PASS
	switch getCmd {
	case CACHE_GET_HIT:
		return hc.doLookup(r)
	case CACHE_GET_MISS:
		return hc.doMiss(r)
	case CACHE_GET_PASS:
		return hc.doPass(r)
	}
	panic("not reached")
}

func (hc *HttpCacher) doPass(r request.Request) (*http.Response, error) {
	r.SetUserData(USERDATA_CACHE_PASS, true)
	return nil, nil
}

func (hc *HttpCacher) doReject(r request.Request) (*http.Response, error) {
	// TODO: better reject
	return nil, errors.New("rejected")
}

const (
	USERDATA_CACHE_PASS = "cache:pass"
)

func (hc *HttpCacher) ProcessRequest(r request.Request) (*http.Response, error) {
	recvCmd := hc.recv.Recv(r)
	switch recvCmd {
	case RECV_LOOKUP:
		return hc.doLookup(r)
	case RECV_PASS:
		return hc.doPass(r)
	case RECV_REJECT:
		return hc.doReject(r)
	}
	panic("not reached")
}

func (hc *HttpCacher) ProcessResponse(r request.Request, a request.Attempt) {
	// we were told to pass, so ignore the request.Attempt
	_, ok := r.GetUserData(USERDATA_CACHE_PASS)
	if ok {
		return
	}
	// TODO: look at Cache-Control headers if we should cache
}
