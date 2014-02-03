package gosk

import (
	"errors"
	"reflect"
	"sync"
	"sync/atomic"
)

type Dataflow struct {
	isSet int32
	cond  *sync.Cond
	val   reflect.Value
}

func NewDataflow() *Dataflow {
	return &Dataflow{
		isSet: 0,
		cond:  sync.NewCond(new(sync.Mutex)),
	}
}

func NewDataflowWith(val reflect.Value) *Dataflow {
	return &Dataflow{
		isSet: 1,
		cond:  sync.NewCond(new(sync.Mutex)),
		val:   val,
	}
}

func (d *Dataflow) Get() reflect.Value {
	if atomic.LoadInt32(&d.isSet) == 1 {
		return d.val
	}

	d.cond.L.Lock()
	for atomic.LoadInt32(&d.isSet) == 0 {
		d.cond.Wait()
	}
	d.cond.L.Unlock()
	return d.val
}

func (d *Dataflow) Set(val reflect.Value) error {
	if atomic.LoadInt32(&d.isSet) == 1 {
		return errors.New("can't assign twice")
	}

	d.cond.L.Lock()
	d.val = val
	atomic.StoreInt32(&d.isSet, 1)
	d.cond.L.Unlock()
	d.cond.Broadcast()

	return nil
}
