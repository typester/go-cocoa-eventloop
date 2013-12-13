package eventloop

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation

#import <Foundation/Foundation.h>

void Run() {
    @autoreleasepool {
        [[NSRunLoop currentRunLoop] runMode:NSDefaultRunLoopMode
                                 beforeDate:[NSDate dateWithTimeIntervalSinceNow:0.1]];
    }
}
*/
import "C"

import (
	"runtime"
)

var (
	mainfunc = make(chan func())
	stop     = make(chan bool)
)

// https://code.google.com/p/go-wiki/wiki/LockOSThread
func init() {
	runtime.LockOSThread()
}

func Run() {
	for {
		select {
		case f := <-mainfunc:
			f()
		case <-stop:
			return
		default:
			C.Run()
		}
	}
}

func Stop() {
	stop <- true
}

func Do(f func()) {
	done := make(chan bool, 1)
	mainfunc <- func() {
		f()
		done <- true
	}
	<-done
}
