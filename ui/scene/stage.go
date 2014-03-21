package scene

import "sync"

// TODO: when a system goroutine terminates, we should runlock and
// wg.Done or we will deadlock the stage.
// TODO: there must be a simpler way to do this lockstepping. perhaps
// two channels can do it..

type stage struct {
	id Stage
	mu sync.RWMutex
	stageClient
}

func (s *stage) init(id Stage) {
	s.id = id
	s.cond = sync.NewCond(s.mu.RLocker())
}

func (s *stage) add() {
	s.mu.RLock()
	s.wg.Add(1)
}

// cycle is called by the world to let all systems in this stage run.
func (s *stage) cycle() {
	s.mu.Lock()
	s.wg, s.nextwg = s.nextwg, s.wg
	s.mu.Unlock()
	s.cond.Broadcast()
	s.wg.Wait()
}

func (s *stage) kill() {
	s.mu.Lock()
	s.wg, s.nextwg = s.nextwg, s.wg
	s.dead = true
	s.mu.Unlock()
	s.cond.Broadcast()
	s.wg.Wait()
}

type stageClient struct {
	cond   *sync.Cond
	wg     sync.WaitGroup
	nextwg sync.WaitGroup
	dead   bool
}

// wait is called by each system, and returns when it is time to
// execute. returns false if world has ended.
func (s *stageClient) wait() bool {
	s.wg.Done()
	s.nextwg.Add(1)
	// runlock
	s.cond.Wait()
	// rlock
	if s.dead {
		s.wg.Done()
		return false
	}
	return true
}

type stagesById []stage

func (s stagesById) Len() int {
	return len(s)
}

func (s stagesById) Less(i, j int) bool {
	return s[i].id < s[j].id
}

func (s stagesById) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
