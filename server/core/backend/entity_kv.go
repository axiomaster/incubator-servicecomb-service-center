/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package backend

import (
	"github.com/apache/incubator-servicecomb-service-center/pkg/util"
	"github.com/apache/incubator-servicecomb-service-center/server/core"
	"sync"
)

var (
	defaultKvEntity *KvEntity
	newKvEntityOnce sync.Once
)

type KvEntity struct {
	Cacher
	Indexer
}

func (se *KvEntity) Run() {
	if r, ok := se.Cacher.(Runnable); ok {
		r.Run()
	}
}

func (se *KvEntity) Stop() {
	if r, ok := se.Cacher.(Runnable); ok {
		r.Stop()
	}
}

func (se *KvEntity) Ready() <-chan struct{} {
	if r, ok := se.Cacher.(Runnable); ok {
		return r.Ready()
	}
	return closedCh
}

func NewKvEntity(name string, cfg *Config) *KvEntity {
	var entity KvEntity
	switch {
	case core.ServerInfo.Config.EnableCache && cfg.InitSize > 0:
		entity.Cacher = NewKvCacher(name, cfg)
		entity.Indexer = NewCacheIndexer(cfg, entity.Cache())
	default:
		util.Logger().Infof(
			"core will not cache '%s' and ignore all events of it, cache enabled: %v, init size: %d",
			name, core.ServerInfo.Config.EnableCache, cfg.InitSize)
		entity.Indexer = NewCommonIndexer(cfg.Key, cfg.Parser)
	}
	return &entity
}

func DefaultKvEntity() *KvEntity {
	newKvEntityOnce.Do(func() {
		defaultKvEntity = &KvEntity{
			Indexer: NewCommonIndexer(Configure().Key, BytesParser),
		}
	})
	return defaultKvEntity
}
