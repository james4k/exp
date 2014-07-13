package graphics

import "sort"

type renderKey uint32

type renderList struct {
	keys  []renderKey
	state []renderState
}

func (r *renderList) insert(key renderKey, state renderState) {
	r.keys = append(r.keys, key)
	r.state = append(r.state, state)
}

func (r *renderList) sort() {
	sort.Sort(r)
}

func (r *renderList) execute() {
	var cur renderState
	for i := range r.state {
		st := &r.state[i]
		changed := cur.flags ^ st.flags
		if changed != 0 {
			st.commit(changed)
		}
		st.draw()
	}
	r.clean()
}

func (r *renderList) clean() {
	if len(r.keys) < cap(r.keys)/3 {
		r.keys = nil
		r.state = nil
	} else {
		r.keys = r.keys[:0]
		r.state = r.state[:0]
	}
}

func (r *renderList) Len() int {
	return len(r.keys)
}

func (r *renderList) Less(i, j int) bool {
	return r.keys[i] < r.keys[j]
}

func (r *renderList) Swap(i, j int) {
	r.keys[i], r.keys[j] = r.keys[j], r.keys[i]
	r.state[i], r.state[j] = r.state[j], r.state[i]
}

type stateFlags uint32

const (
	stateRGBWrite stateFlags = 1 << iota
	stateAlphaWrite
	stateDepthWrite

	stateDepthTestLess stateFlags = 0x10 + iota
	stateDepthTestLEqual
	stateDepthTestEqual
	stateDepthTestGEqual
	stateDepthTestGreater
	stateDepthTestNotEqual
	stateDepthTestNever
	stateDepthTestAlways
	stateDepthTestMask = 0xf0

	stateBlendZero
	stateMSAA
)

type renderState struct {
	flags stateFlags
}

func (r *renderState) key() renderKey {
	return 0
}

func (r *renderState) cmd() {
	r.commit()
	r.draw()
}

// commit updates GL state for this renderable
func (r *renderState) commit(changed stateFlags) {

}

// draw
func (r *renderState) draw() {
}
