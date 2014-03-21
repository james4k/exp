package scene

type Tx struct {
	w *World
}

func (t Tx) Create(data interface{}) Ref {
	return t.w.create(data)
}

func (t Tx) Delete(node Ref) {
	t.w.delete(node)
}

// Get retrieves values for node by copying to the pointers passed
// into data. Panics if node is invalid or data for all types was not
// found for the node.
func (t Tx) Get(node Ref, data ...interface{}) {
}

// Put assigns data for the specified node. Only data of types used in
// the creation of the node may be used. Panics if node is invalid or
// data for all types was not found for the node.
func (t Tx) Put(node Ref, data ...interface{}) {
}

// Commit completes the transaction.
func (t *Tx) Commit() {
	w := t.w
	t.w = nil
	w.start()
}
