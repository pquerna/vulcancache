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
)

type basicCache struct{}

func newBasicCache(co *CacheOptions) *basicCache {
	return &basicCache{}
}

func (bc *basicCache) Recv(r request.Request) RecvCmd {
	return RECV_LOOKUP
}

func (bc *basicCache) Hash(r request.Request, hf HashFunc) {

}

func (bc *basicCache) Get(key string, r request.Request) (CacheGetCmd, CachedResponse) {
	return CACHE_GET_MISS, nil
}

func (bc *basicCache) Store(key string, r request.Request, a request.Attempt) CacheStoreCmd {
	return CACHE_STORE_MISS
}

func (bc *basicCache) Miss(r request.Request) MissCmd {
	return MISS_FETCH
}
