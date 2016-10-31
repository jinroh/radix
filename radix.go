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

// NewTree returns a new tree handler.
func New() *Tree {
	return &Tree{}
}

// Get does a lookup of the given key and returns the node value that
// match the given key.
func (t *Tree) Get(key string) (interface{}, bool) {
	if match, node := t.lookup(key); node != nil && match {
		return node.v, true
	}
	return nil, false
}

// Closeup returns a new view of the tree from the prefix passed as a
// key. If the prefix does not match any prefix of an indexed value of
// the tree, it returns nil.
func (t *Tree) Closeup(key string) *Tree {
	if _, node := t.lookup(key); node != nil {
		return &Tree{root: node}
	}
	return nil
}

func (t *Tree) lookup(key string) (match bool, node *tnode) {
	node = t.root

	for {
		if node == nil {
			return false, nil
		}

		if match, _, _ := keyMatch(key, node.k); !match {
			return false, nil
		}

		keylen := len(key) - len(node.k)
		if keylen <= 0 {
			return keylen == 0, node
		}

		key = key[len(node.k):]
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
		match, splitpos, diff := keyMatch(key, node.k)
		if !match || len(key) < len(node.k) {
			repld = false
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

		keylen := len(key) - len(node.k)
		if keylen == 0 {
			node.v = v
			return v, repld
		}

		key = key[len(node.k):]
		if node.mask == 0 || key[0]&node.mask > 0 {
			pnode = &node.r
		} else {
			pnode = &node.l
		}
	}
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

func keyMatch(keya, keyb string) (bool, int, byte) {
	keyalen := len(keya)
	keyblen := len(keyb)
	minlen := min(keyalen, keyblen)

	for i := 0; i < minlen; i++ {
		if d := keya[i] ^ keyb[i]; d != 0 {
			return false, i, d
		}
	}

	return true, minlen, 0
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
