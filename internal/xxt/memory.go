// Copyright (c) 2022, Geert JM Vanderkelen

package xxt

import (
	"fmt"
	"runtime"
)

type MemoryUse struct {
	start *runtime.MemStats
	end   *runtime.MemStats
}

func NewMemoryUse() *MemoryUse {
	m := &MemoryUse{}
	m.start = &runtime.MemStats{}
	m.end = &runtime.MemStats{}

	runtime.GC()
	runtime.ReadMemStats(m.start)
	return m
}

func (m *MemoryUse) Stop() {
	runtime.ReadMemStats(m.end)
}

func (m *MemoryUse) DiffAlloc() uint64 {
	return m.end.Alloc - m.start.Alloc
}

func (m MemoryUse) String() string {
	return fmt.Sprintf("DiffTotalAlloc = % 10d\tNumGC = % 5d\n",
		m.end.Alloc-m.start.Alloc, m.end.NumGC)
}
