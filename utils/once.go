/*
Copyright 2014 The Camlistore Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file_utils except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"sync"
	"sync/atomic"
)

// A Once will perform a successful action exactly once.
//
// Unlike a sync.Once, this Once's func returns an error
// and is re-armed on failure.
type Once struct {
	m    sync.Mutex
	done uint32
}

// Do calls the function f if and only if Do has not been invoked
// without error for this instance of Once.  In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will
// invoke f, even if f has a different value in each invocation unless
// f returns an error.  A new instance of Once is required for each
// function to execute.
//
// Do is intended for initialization that must be run exactly once.  Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	err := config.once.Do(func() error { return config.init(filename) })
//
// sync.Once.Do(f func())是一个挺有趣的东西,能保证once只执行一次，
// 无论你是否更换once.Do(xx)这里的方法,这个sync.Once块只会执行一次。
func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		// 1表示已经成功执行过，保证只成功执行一次
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = f()
		// 成功设置为1
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}
