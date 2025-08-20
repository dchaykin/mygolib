package pipeline

import (
	"runtime/debug"
	"sync"
	"time"

	"github.com/dchaykin/mygolib/log"
)

type PipelineFunc struct {
	name                     string
	pauseAfterInMilliseconds int
	fn                       func(mutexState *sync.Mutex)
}

func CreateFunc(name string, pauseAfterInMilliseconds int, fn func(mutexState *sync.Mutex)) PipelineFunc {
	return PipelineFunc{
		name:                     name,
		pauseAfterInMilliseconds: pauseAfterInMilliseconds,
		fn:                       fn,
	}
}

type Pipeline struct {
	Fn         []PipelineFunc
	mutexState sync.Mutex
	stopping   bool
}

func (p *Pipeline) Start() {
	log.Info("activating pipelines...")
	p.stopping = false

	for i, f := range p.Fn {
		if f.pauseAfterInMilliseconds == 0 {
			log.Info("pause after function call of '%s' is not set, using default value of 10s", f.name)
			f.pauseAfterInMilliseconds = 10000
		}
		log.Info("%d. %s", i+1, f.name)
		go p.safeWorker(f)
	}
}

func (p *Pipeline) Stop() {
	log.Info("stopping pipelines...")
	p.mutexState.Lock()
	defer p.mutexState.Unlock()
	p.stopping = true
}

func (p *Pipeline) safeWorker(f PipelineFunc) {
	restartCount := 0
	delay := 2
	for {
		func() {
			defer func() {
				if r := recover(); r != nil {
					restartCount++
					debug.PrintStack()
					log.Errorf("⚠️ panic in '%s' occurred: %v – Restart in %ds", f.name, r, delay)
					time.Sleep(time.Duration(delay) * time.Second)
					delay *= 2
				}
			}()
			if restartCount > 10 {
				log.Info("stopping '%s' after %d restarts", f.name, restartCount)
				return
			}
			if p.stopping {
				log.Info("stopping '%s'", f.name)
				return
			}

			p.infiniteLoop(f)
		}()
		if restartCount > 10 || p.stopping {
			return
		}
	}
}

func (p *Pipeline) infiniteLoop(pFunc PipelineFunc) {
	for {
		if p.stopping {
			return
		}
		pFunc.fn(&p.mutexState)
		time.Sleep(time.Duration(pFunc.pauseAfterInMilliseconds) * time.Millisecond)
	}
}
