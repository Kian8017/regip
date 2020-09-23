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
