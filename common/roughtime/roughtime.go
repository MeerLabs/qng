// Copyright (c) 2017-2019 The Qitmeer developers
//
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package roughtime

import (
	"log"
	"time"
)

const RecalibrationInterval = time.Hour

var offset time.Duration

func Init() {
	log.Println("Init roughtime... ...")
	recalibrateRoughtime()
	runT := time.NewTimer(RecalibrationInterval)
	go func() {
		for {
			select {
			case <-runT.C:
				recalibrateRoughtime()
			}
		}
	}()
}

func recalibrateRoughtime() {
}

// Since returns the duration since t, based on the roughtime response
func Since(t time.Time) time.Duration {
	return Now().Sub(t)
}

// Until returns the duration until t, based on the roughtime response
func Until(t time.Time) time.Duration {
	return t.Sub(Now())
}

// Now returns the current local time given the roughtime offset.
func Now() time.Time {
	if offset <= 0 {
		return time.Now()
	}
	return time.Now().Add(offset)
}

func Offset() time.Duration {
	return offset
}
