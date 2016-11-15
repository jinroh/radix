package radix

import "bytes"

type tnode struct {
	v    interface{} // value (if leaf node)
	mask byte
	l    *tnode // left
	r    *tnode // right
	k    string // key piece
}

// Tree holds a root node of a radix-2 tree.
type Tree struct {
	root *tnode
}

// New returns a new tree handler.
func New() *Tree {
	return &Tree{}
}

// Get does a lookup of the given key and returns the node value that
// match the given key.
func (t *Tree) Get(key string) (interface{}, bool) {
	if match, node, _, _ := t.lookup(key); node != nil && match {
		return node.v, true
	}
	return nil, false
}

// Closeup returns a new view of the tree from the prefix passed as a
// key. If the prefix does not match any prefix of an indexed value of
// the tree, it returns nil.
func (t *Tree) Closeup(key string) *Tree {
	if _, node, _, _ := t.lookup(key); node != nil {
		return &Tree{root: node}
	}
	return nil
}

func (t *Tree) lookup(key string) (match bool, node *tnode, parent *tnode, gparent *tnode) {
	node = t.root

	for {
		if node == nil {
			return
		}

		if node.k != "" {
			if !keyMatch(node.k, key) {
				return
			}

			keylen := len(key) - len(node.k)
			if keylen <= 0 {
				match = keylen == 0
				return
			}

			key = key[len(node.k):]
		}

		parent, gparent = node, parent

		if node.mask == 0 || key[0]&node.mask > 0 {
			node = node.r
		} else {
			node = node.l
		}
	}
}

// Insert adds a new value in the tree at the given key. Returns a
// boolean true iff the key already existed and has been replaced.
func (t *Tree) Insert(key string, v interface{}) (interface{}, bool) {
	var node *tnode
	pnode := &t.root
	repld := true

	for {
		if *pnode == nil {
			repld = false
			*pnode = &tnode{k: key}
		}

		node = *pnode
		if node.k != "" && (len(key) < len(node.k) || !keyMatch(node.k, key)) {
			repld = false

			splitpos, diff := xorStrings(node.k, key)
			mask := calcMask(diff)

			splitkey := key[:splitpos]
			lowerkey := node.k[splitpos:]

			var splitnode *tnode
			var lowernode *tnode

			splitnode = &tnode{
				k:    splitkey,
				mask: mask,
			}

			lowernode = &tnode{
				v:    node.v,
				mask: node.mask,
				l:    node.l,
				r:    node.r,
				k:    lowerkey,
			}

			if mask == 0 || key[splitpos]&mask == 0 {
				splitnode.r = lowernode
			} else {
				splitnode.l = lowernode
			}

			*pnode, node = splitnode, splitnode
		}

		curlen := len(node.k)
		keylen := len(key) - curlen
		if keylen == 0 {
			node.v = v
			return v, repld
		}

		if curlen > 0 {
			key = key[curlen:]
		}

		if node.mask == 0 || key[0]&node.mask > 0 {
			pnode = &node.r
		} else {
			pnode = &node.l
		}
	}
}

// Remove removes the node at the given key.
func (t *Tree) Remove(key string) (interface{}, bool) {
	match, node, parent, gparent := t.lookup(key)
	if !match {
		return nil, false
	}

	// root node
	if parent == nil {
		v := node.v
		if node.l == nil && node.r == nil {
			t.root = nil
		} else {
			node.v = nil
		}
		return v, true
	}

	return t.rm(node, parent, gparent)
}

// RemoveBranch remove all elements of the tree starting with the
// given key.
func (t *Tree) RemoveBranch(key string) bool {
	_, node, parent, gparent := t.lookup(key)
	if node == nil {
		return false
	}

	// root node
	if parent == nil {
		t.root = nil
		return true
	}

	_, ok := t.rm(node, parent, gparent)
	return ok
}

func (t *Tree) rm(node, parent, gparent *tnode) (interface{}, bool) {
	v := node.v

	if parent.l == node {
		parent.l = nil
	} else {
		parent.r = nil
	}

	// direct parent has value or its parent is root so we can not merge
	// anything
	if parent.v != nil || gparent == nil {
		return v, true
	}

	if parent.l != nil {
		node = parent.l
	} else {
		node = parent.r
	}

	if node != nil {
		node.k = parent.k + node.k
	}

	if gparent.l == parent {
		gparent.l = node
	} else {
		gparent.r = node
	}

	return v, true
}

// Foreach is used to iterates of the values of the tree. For each
// value, the given callback is called with the value and key as
// parameters.
func (t *Tree) Foreach(cb func(interface{}, string) error) error {
	var keybuf bytes.Buffer
	return foreach(t.root, keybuf, cb)
}

func (t *Tree) forall(cb func(*tnode, string) error) error {
	var keybuf bytes.Buffer
	return forall(t.root, keybuf, cb)
}

func foreach(node *tnode, keybuf bytes.Buffer, cb func(interface{}, string) error) error {
	if node == nil {
		return nil
	}
	keybuf.WriteString(node.k)
	if node.v != nil {
		if err := cb(node.v, keybuf.String()); err != nil {
			return err
		}
	}
	foreach(node.l, keybuf, cb)
	foreach(node.r, keybuf, cb)
	return nil
}

func forall(node *tnode, keybuf bytes.Buffer, cb func(*tnode, string) error) error {
	if node == nil {
		return nil
	}
	keybuf.WriteString(node.k)
	if err := cb(node, keybuf.String()); err != nil {
		return err
	}
	forall(node.l, keybuf, cb)
	forall(node.r, keybuf, cb)
	return nil
}

func keyMatch(keya, keyb string) (match bool) {
	n := len(keya)
	if len(keyb) < n {
		n = len(keyb)
	}
	return keya[0:n] == keyb[0:n]
}

func calcMask(d byte) (mask byte) {
	if d == 0 {
		return 0
	}
	for mask = 0x80; mask > 0 && mask&d == 0; {
		mask = mask >> 1
	}
	return
}
