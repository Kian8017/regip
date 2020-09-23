/* Copyright (c) 2020 Kian Musser.
 * This file is part of regip.
 *
 * regip is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * regip is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with regip.  If not, see <https://www.gnu.org/licenses/>.
 */

package regip

import (
	"sync"
)

type Filter func(*Flow) *Flow
type FilterFn func(Resource) bool

func NewFlowFilter(nt FilterFn) Filter {
	return func(inputF *Flow) *Flow {
		output := NewFlow()
		go func() {
			for {
				inp, ok := inputF.Get()
				if !ok {
					output.Stop()
					break
				}
				keep := nt(inp)
				if keep {
					ok = output.Push(inp)
					if !ok {
						inputF.Stop()
						break
					}
				}
			}
		}()
		return output
	}
}

type Flow struct {
	resources chan Resource
	isErr     bool
	errChan   chan struct{}
	mut       sync.Mutex
}

func NewFlow() *Flow {
	var f Flow
	var m sync.Mutex
	f.resources = make(chan Resource)
	f.errChan = make(chan struct{})
	f.mut = m
	return &f
}

func (f *Flow) Get() (Resource, bool) {
	if f.err() {
		return nil, false
	}
	select {
	case <-f.errChan:
		return nil, false
	case res := <-f.resources:
		return res, true
	}
}

func (f *Flow) Push(res Resource) bool {
	if f.err() {
		return false
	}
	select {
	case <-f.errChan:
		return false
	case f.resources <- res:
		return true
	}
}

func (f *Flow) err() bool {
	f.mut.Lock()
	defer f.mut.Unlock()
	return f.isErr
}

func (f *Flow) Stop() {
	if f.err() {
		return // Someone already stopped
	}
	f.mut.Lock()
	f.isErr = true
	f.mut.Unlock()
	close(f.errChan)
}
