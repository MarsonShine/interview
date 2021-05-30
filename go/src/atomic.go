package interview

import "sync/atomic"

var value int32

func cas(newValue int32) {
	for {
		v := value
		if atomic.CompareAndSwapInt32(&value, v, (v + newValue)) {
			break
		}
	}
}
