package pipeline

import (
	"sync"
	"testing"
	"time"

	"github.com/dchaykin/mygolib/log"
	"github.com/stretchr/testify/require"
)

func TestPipelineRunsAndStops(t *testing.T) {
	var callCount int

	f := CreateFunc("TestFunc", 300, func(mutexState *sync.Mutex) {
		mutexState.Lock()
		callCount++
		log.Info("TestFunc called %d times", callCount)
		mutexState.Unlock()
	})

	p := Pipeline{
		Fn: []PipelineFunc{f},
	}

	p.Start()
	time.Sleep(2 * time.Second)
	p.Stop()

	log.Info("TestFunc called %d times", callCount)

	require.Less(t, 0, callCount, "Expected function to be called at least once")
}

func TestPanicRecoveryAndRestart(t *testing.T) {
	var callCount int

	f := CreateFunc("MayPanic", 100, func(mutexState *sync.Mutex) {
		callCount++
		log.Info("MayPanic called %d times", callCount)
		panic("simulated panic")
	})

	p := Pipeline{
		Fn: []PipelineFunc{f},
	}

	p.Start()
	time.Sleep(20 * time.Second)
	p.Stop()

	require.LessOrEqual(t, 2, callCount, "Expected function to be called more than 2 times")
}

func TestMultipleFunctionsRun(t *testing.T) {
	called := map[string]int{
		"A": 0,
		"B": 0,
	}

	f1 := CreateFunc("FuncA", 200, func(mutexState *sync.Mutex) {
		mutexState.Lock()
		called["A"]++
		log.Info("FuncA called %d times", called["A"])
		mutexState.Unlock()
	})

	f2 := CreateFunc("FuncB", 300, func(mutexState *sync.Mutex) {
		mutexState.Lock()
		called["B"]++
		log.Info("FuncB called %d times", called["B"])
		mutexState.Unlock()
	})

	p := Pipeline{
		Fn: []PipelineFunc{f1, f2},
	}

	p.Start()
	time.Sleep(2 * time.Second)
	p.Stop()

	require.Len(t, called, 2, "Expected both functions to be called")
	require.LessOrEqual(t, 3, called["A"], "Expected function A to be called")
	require.LessOrEqual(t, 2, called["B"], "Expected function B to be called")
}
