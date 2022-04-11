package main

import (
	"sync"
	"testing"
)

// Протестируйте производительность множества действительных чисел,
// безопасность, которого обеспечивается sync.Mutex и sync.RWMutex для
// разных вариантов использования: 10% запись, 90% чтение; 50%
// запись, 50% чтение; 90% запись, 10% чтение

type realNumber struct {
	rwm sync.RWMutex
	m   sync.Mutex
}

func (rn *realNumber) RWMRead() {
	rn.rwm.RLock()
	rn.rwm.RUnlock()
}

func (rn *realNumber) RWMWrite() {
	rn.rwm.Lock()
	rn.rwm.Unlock()
}

func (rn *realNumber) MRead() {
	rn.m.Lock()
	rn.m.Unlock()
}

func (rn *realNumber) MWrite() {
	rn.m.Lock()
	rn.m.Unlock()
}

func BenchmarkLesson5_3_RWMutex_10Read_90Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%9 == 0 {
					t.RWMWrite()
				} else {
					t.RWMRead()
				}
				i++
			}
		})
	})

}

func BenchmarkLesson5_3_RWMutex_50Read_50Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%2 == 0 {
					t.RWMRead()
				} else {
					t.RWMWrite()
				}
				i++
			}
		})
	})
}

func BenchmarkLesson5_3_RWMutex_90Read_10Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%9 == 0 {
					t.RWMRead()
				} else {
					t.RWMWrite()
				}
				i++
			}
		})
	})
}

// --------- Mutex ---------------------
func BenchmarkLesson5_3_Mutex_10Read_90Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%9 == 0 {
					t.MRead()
				} else {
					t.MWrite()
				}
				i++
			}
		})
	})
}

func BenchmarkLesson5_3_Mutex_50Read_50Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%2 == 0 {
					t.MRead()
				} else {
					t.MWrite()
				}
				i++
			}
		})
	})
}

func BenchmarkLesson5_3_Mutex_90Read_10Write(b *testing.B) {
	var (
		t realNumber
		i int
	)

	b.Run("", func(b *testing.B) {
		b.SetParallelism(1000)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if i%9 == 0 {
					t.MWrite()
				} else {
					t.MRead()
				}
				i++
			}
		})
	})
}
