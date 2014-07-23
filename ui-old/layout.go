// +build ignore

package ui

import (
	"errors"
	"image"
)

type Dropable interface {
	Drop()
}

type LayoutData struct {
}

func (t *LayoutData) Drop() {
}

type LinkComponent struct {
	Parent NodeKey
}

type BoxComponent struct {
	Local image.Rectangle
	World image.Rectangle
}

type layoutBucket struct {
	nodes map[NodeKey]int
	links []LinkComponent
	boxes []BoxComponent
}

func (l *layoutBucket) add(node NodeKey) int {
	_, ok := l.nodes[node]
	if ok {
		panic(errors.New("ui: layoutBucket.add - node exists"))
	}
	idx := len(l.links)
	l.nodes[node] = idx
	l.links = append(l.links, LinkComponent{})
	l.boxes = append(l.boxes, BoxComponent{})
	return idx
}

func (l *layoutBucket) remove(node NodeKey) {
	idx, ok := l.nodes[node]
	if !ok {
		panic(errors.New("ui: layoutBucket.remove - node does not exist"))
	}
	delete(l.nodes, node)
	l.links[idx] = l.links[len(l.links)-1]
	l.links = l.links[:len(l.links)-1]
	l.boxes[idx] = l.boxes[len(l.boxes)-1]
	l.boxes = l.boxes[:len(l.boxes)-1]
}

type layoutData struct {
	nodeDepth  map[NodeKey]int
	depthFirst []layoutBucket
}

type layoutDataRef struct {
	handoffc chan *layoutData
	*layoutData
}

func (l *layoutDataRef) release() {
	p := l.layoutData
	l.layoutData = nil
	l.handoffc <- p
}

type LayoutSystem struct {
	datac chan *layoutData
}

func (l *LayoutSystem) acquire() layoutDataRef {
	return layoutDataRef{
		handoffc:   l.datac,
		layoutData: <-l.datac,
	}
}

func (l *LayoutSystem) init() {
	l.datac = make(chan *layoutData, 1)
	l.datac <- &layoutData{}
}

func (l *LayoutSystem) NodeEnter(node NodeKey, arch Archetype) {
	// setup data for new node..no parents
}

func (l *LayoutSystem) NodeExit(nodes NodeKey) {
}

func (l *LayoutSystem) Run(scene *Scene, frame FrameLoop) {
	l.init()
	for frame.Next() {
		l.update()
	}
}

func (l *LayoutSystem) update() {
	data := l.acquire()
	defer data.release()
	prevb := &data.depthFirst[0]
	for i := range prevb.boxes {
		prevb.boxes[i].World.Min = prevb.boxes[i].Local.Min
	}
	for i := 1; i < len(data.depthFirst); i++ {
		bkt := &data.depthFirst[i]
		for i := range bkt.boxes {
			parent := prevb.nodes[bkt.links[i].Parent]
			bkt.boxes[i].World.Min = bkt.boxes[i].Local.Min.Add(prevb.boxes[parent].World.Min)
		}
		prevb = bkt
	}
}

func (l *LayoutSystem) LinkNode(node, parent NodeKey) {
	/*
		data := l.acquire()
		defer data.release()
		depth, ok := data.nodeDepth[node]
		if !ok {
			panic(errors.New("ui: node not found"))
		}
		parentDepth, ok := data.nodeDepth[parent]
		if !ok {
			panic(errors.New("ui: node not found"))
		}
		bkt := &data.depthFirst[parentDepth+1]
		data.nodeDepth[node] = parentDepth + 1
		var idx int
		if depth == parentDepth+1 {
			idx = bkt.nodes[node]
			bkt.links[idx].Parent = parent
			return
		}
		data.depthFirst[depth].remove(node)
		idx = bkt.add(node)
	*/
}

func (l *LayoutSystem) UnlinkNode(node NodeKey) {
	/*
		data := l.acquire()
		defer data.release()
		idx, ok := data.nodes[node]
		if !ok {
			panic(errors.New("ui: node not found"))
		}
		data.links[idx] = NullNode
	*/
}

func (l *LayoutSystem) NodeParent(node NodeKey) NodeKey {
	data := l.acquire()
	defer data.release()
	depth, ok := data.nodeDepth[node]
	if !ok {
		panic(errors.New("ui: node not found"))
	}
	bkt := &data.depthFirst[depth]
	idx, ok := bkt.nodes[node]
	if !ok {
		panic(errors.New("ui: node not found in layoutBucket"))
	}
	return bkt.links[idx].Parent
}
