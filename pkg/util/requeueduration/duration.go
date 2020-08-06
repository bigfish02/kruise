/*
Copyright 2019 The Kruise Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package requeueduration

import (
	"sync"
	"time"
)

// DurationStore can store a duration map for multiple workloads
type DurationStore struct {
	store sync.Map
}

func (dm *DurationStore) Push(key string, newDuration time.Duration) {
	value, _ := dm.store.LoadOrStore(key, &Duration{})
	requeueDuration, ok := value.(*Duration)
	if !ok {
		dm.store.Delete(key)
		return
	}
	requeueDuration.Update(newDuration)
}

func (dm *DurationStore) Pop(key string) time.Duration {
	value, ok := dm.store.Load(key)
	if !ok {
		return 0
	}
	defer dm.store.Delete(key)
	requeueDuration, ok := value.(*Duration)
	if !ok {
		return 0
	}
	return requeueDuration.Get()
}

// Duration helps calculate the shortest non-zore duration to requeue
type Duration struct {
	sync.Mutex
	duration time.Duration
}

func (rd *Duration) Update(newDuration time.Duration) {
	rd.Lock()
	defer rd.Unlock()
	if newDuration > 0 {
		if rd.duration <= 0 || newDuration < rd.duration {
			rd.duration = newDuration
		}
	}
}

func (rd *Duration) Get() time.Duration {
	rd.Lock()
	defer rd.Unlock()
	return rd.duration
}
