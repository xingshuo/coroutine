package coroutine

import (
	"errors"
	"fmt"
)

type CoFunc func(co *Coroutine, a ...interface{}) []interface{}

// 基于goroutine的并发特性考虑, 移除了lua coroutine的"normal"状态
type CoState uint32

const (
	Suspended CoState = iota // yield or not start
	Running
	Dead // finished or stopped with an error
)

func (s CoState) String() string {
	switch s {
	case Suspended:
		return "Suspended"
	case Running:
		return "Running"
	case Dead:
		return "Dead"
	default:
		return "Unknown"
	}
}

type Coroutine struct {
	waitOut chan []interface{}
	waitIn  chan []interface{}
	fn      CoFunc
	// 捕获fn执行异常
	fnWrapper func(a ...interface{})
	state     CoState
	err       error
}

func (co *Coroutine) Yield(out ...interface{}) []interface{} {
	co.state = Suspended
	co.waitOut <- out
	return <-co.waitIn
}

func (co *Coroutine) Resume(in ...interface{}) (error, []interface{}) {
	if co.state == Dead {
		return errors.New("cannot resume dead coroutine"), nil
	}
	if co.state != Suspended {
		return fmt.Errorf("resume not suspended coroutine %s", co.state), nil
	}
	co.state = Running
	if co.fnWrapper == nil { // not start
		co.fnWrapper = func(a ...interface{}) {
			defer func() {
				co.state = Dead
				if e := recover(); e != nil {
					co.err = fmt.Errorf("%w", e)
					co.waitOut <- nil
				}
			}()
			co.waitOut <- co.fn(co, a...)
		}
		go co.fnWrapper(in...)
	} else { // wake up
		co.waitIn <- in
	}
	out := <-co.waitOut
	if co.err != nil {
		return co.err, nil
	}
	return nil, out
}

func (co *Coroutine) Status() CoState {
	return co.state
}

func Create(f CoFunc) *Coroutine {
	co := &Coroutine{
		waitIn:  make(chan []interface{}),
		waitOut: make(chan []interface{}),
		fn:      f,
		state:   Suspended,
	}
	return co
}
