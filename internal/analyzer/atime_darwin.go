//go:build darwin

package analyzer

import "syscall"

func atimeFromStat(st *syscall.Stat_t) (sec int64, nsec int64) {
	return int64(st.Atimespec.Sec), int64(st.Atimespec.Nsec)
}
