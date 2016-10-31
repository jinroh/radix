package radix

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestSimple(t *testing.T) {
	tree := New()

	tree.Insert([]byte(""), "root")
	tree.Insert([]byte("a"), "a")
	tree.Insert([]byte("b"), "b")
	tree.Insert([]byte("d"), "d")
	tree.Insert([]byte("c"), "c")
	tree.Insert([]byte("e"), "e")
	tree.Insert([]byte("f"), "f")
	tree.Insert([]byte("aa"), "aa")

	tree.Insert([]byte("/"), "root")
	tree.Insert([]byte("/aa"), "aa")
	tree.Insert([]byte("/folder/"), "folder1")
	tree.Insert([]byte("/folder/"), "folder1-")
	tree.Insert([]byte("/folder/file1"), "file1")
	tree.Insert([]byte("/folder/file2"), "file2")
	tree.Insert([]byte("/folder/file3"), "file3")
	tree.Insert([]byte("/folder/file4"), "file4")
	tree.Insert([]byte("/folder2"), "folder2")

	tree.Insert([]byte("/foo/bar"), "coucou2")
	tree.Insert([]byte("/wesh"), "wesh")
	tree.Insert([]byte(""), "empty")

	if v, ok := tree.Get([]byte("")); ok && v.(string) != "empty" {
		t.Fatal("empty")
	}
	if v, ok := tree.Get([]byte("a")); ok && v.(string) != "a" {
		t.Fatal("a")
	}
	if v, ok := tree.Get([]byte("b")); ok && v.(string) != "b" {
		t.Fatal("b")
	}
	if _, ok := tree.Get([]byte("/fold")); ok {
		t.Fatal("nil")
	}
	if v, ok := tree.Get([]byte("/aa")); ok && v.(string) != "aa" {
		t.Fatal("/aa")
	}
	if v, ok := tree.Get([]byte("/folder/")); ok && v.(string) != "folder1-" {
		t.Fatal("folder1-")
	}
	if v, ok := tree.Get([]byte("/folder/file1")); ok && v.(string) != "file1" {
		t.Fatal("file1")
	}
	if v, ok := tree.Get([]byte("/folder/file4")); ok && v.(string) != "file4" {
		t.Fatal("file4")
	}
}

func TestRandomString(t *testing.T) {
	l := 100000
	m := make([]string, l)
	var tree *Tree

	tree = New()
	for i := 0; i < l; i++ {
		v := randBytes(30, "ab", false)
		tree.Insert(v, string(v))
		m[i] = string(v)
	}

	for _, s := range m {
		v, ok := tree.Get([]byte(s))
		if !ok || v == nil || v.(string) != s {
			t.Fatal("simple rand on key " + s)
		}
	}

	tree = New()
	for i := 0; i < l; i++ {
		v := randBytes(15, "", true)
		tree.Insert(v, string(v))
		m[i] = string(v)
	}

	for _, s := range m {
		v, ok := tree.Get([]byte(s))
		if !ok || v == nil || v.(string) != s {
			t.Fatal("complex rand on key " + s)
		}
	}
}

func TestCloseup(t *testing.T) {
	t.Skip()

	tree := New()

	tree.Insert([]byte("/folder"), "folder")
	tree.Insert([]byte("/folder/file1"), "file1")
	tree.Insert([]byte("/folder/file2"), "file2")
	tree.Insert([]byte("/folder/file3"), "file3")
	tree.Insert([]byte("/folder/file4"), "file4")
	tree.Insert([]byte("/folder2"), "folder2")

	closeup := tree.Closeup([]byte("/folder/"))
	closeup.Foreach(func(val interface{}, key []byte) error {
		fmt.Println(string(key), val.(string))
		return nil
	})
}

func TestFiles(t *testing.T) {
	t.Skip()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		t.Fatal(err)
	}

	tree := New()
	err = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		tree.Insert([]byte(path), f)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	tree.Foreach(func(val interface{}, key []byte) error {
		fmt.Println(string(key), val.(os.FileInfo).IsDir())
		return nil
	})
}

func randBytes(n int, letterBytes string, randlen bool) []byte {
	var b []byte
	if randlen {
		b = make([]byte, rand.Intn(n))
	} else {
		b = make([]byte, n)
	}
	for i := range b {
		if len(letterBytes) > 0 {
			b[i] = letterBytes[rand.Intn(len(letterBytes))]
		} else {
			b[i] = byte(rand.Intn(1 << 8))
		}
	}
	return b
}
