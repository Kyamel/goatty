//go:build darwin

package hinters

import "golang.org/x/sys/unix"

func getUptime() int64 {
	var uptime unix.Timespec
	if err := unix.ClockGettime(unix.CLOCK_UPTIME_RAW, &uptime); err != nil {
		return 0
	}
	return uptime.Sec
}
