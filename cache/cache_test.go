// Copyright 2019 asana Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"os"
	"sync"
	"testing"
	"time"
)

func TestCacheIncr(t *testing.T) {
	bm, err := NewCache("memory", `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	//timeoutDuration := 10 * time.Second

	_ = bm.Put("edwardhey", 0, time.Second*20)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			bm.Incr("edwardhey")
		}()
	}
	wg.Wait()
	if bm.Get("edwardhey").(int) != 10 {
		t.Error("Incr err")
	}
}

func TestCache(t *testing.T) {
	bm, err := NewCache(MemoryProvider, `{"interval":20}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("asana") {
		t.Error("check err")
	}

	if err = bm.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("asana"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("asana"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}
	_ = bm.Delete("asana")
	if bm.IsExist("asana") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = bm.Put("asana", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}
	if v := bm.Get("asana"); v.(string) != "author" {
		t.Error("get err")
	}

	if err = bm.Put("asana1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"asana", "asana1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
}

func TestGCache(t *testing.T) {
	bm, err := NewCache(GCacheProvider, `{"size":20,"type":"arc"}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if bm.IsExist("asana") {
		t.Error("check err")
	}

	if err = bm.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	if err = bm.Incr("asana"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("asana"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}
	_ = bm.Delete("asana")
	if bm.IsExist("asana") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = bm.Put("asana", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}
	if v := bm.Get("asana"); v.(string) != "author" {
		t.Error("get err")
	}

	if err = bm.Put("asana1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"asana", "asana1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}
}

func TestSyncronizerCache(t *testing.T) {
	bm, err := NewCache(MemoryProvider, `{"interval":20}`)
	gc, err2 := NewCache(GCacheProvider, `{"size":20,"type":"arc"}`)

	sync := NewSynchronizer(bm, gc)

	if err != nil {
		t.Error("init err")
	}

	if err2 != nil {
		t.Error("init err")
	}

	timeoutDuration := 10 * time.Second
	if err = sync.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !sync.IsExist("asana") {
		t.Error("check err")
	}

	if v := sync.Get("asana", timeoutDuration); v.(int) != 1 {
		t.Error("get err")
	}

	time.Sleep(30 * time.Second)

	if sync.IsExist("asana") {
		t.Error("check err")
	}

	if err = sync.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}

	sync.Delete("asana")
	if sync.IsExist("asana") {
		t.Error("delete err")
	}

	//test GetMulti
	if err = sync.Put("asana", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !sync.IsExist("asana") {
		t.Error("check err")
	}
	if v := sync.Get("asana", 0); v.(string) != "author" {
		t.Error("get err")
	}

	if err = sync.Put("asana1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !sync.IsExist("asana1") {
		t.Error("check err")
	}
}

func TestFileCache(t *testing.T) {
	bm, err := NewCache(FileProvider, `{"CachePath":"cache","FileSuffix":".bin","DirectoryLevel":"2","EmbedExpiry":"0"}`)
	if err != nil {
		t.Error("init err")
	}
	timeoutDuration := 10 * time.Second
	if err = bm.Put("asana", 1, timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}

	if err = bm.Incr("asana"); err != nil {
		t.Error("Incr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 2 {
		t.Error("get err")
	}

	if err = bm.Decr("asana"); err != nil {
		t.Error("Decr Error", err)
	}

	if v := bm.Get("asana"); v.(int) != 1 {
		t.Error("get err")
	}
	_ = bm.Delete("asana")
	if bm.IsExist("asana") {
		t.Error("delete err")
	}

	//test string
	if err = bm.Put("asana", "author", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana") {
		t.Error("check err")
	}
	if v := bm.Get("asana"); v.(string) != "author" {
		t.Error("get err")
	}

	//test GetMulti
	if err = bm.Put("asana1", "author1", timeoutDuration); err != nil {
		t.Error("set Error", err)
	}
	if !bm.IsExist("asana1") {
		t.Error("check err")
	}

	vv := bm.GetMulti([]string{"asana", "asana1"})
	if len(vv) != 2 {
		t.Error("GetMulti ERROR")
	}
	if vv[0].(string) != "author" {
		t.Error("GetMulti ERROR")
	}
	if vv[1].(string) != "author1" {
		t.Error("GetMulti ERROR")
	}

	_ = os.RemoveAll("cache")
}
