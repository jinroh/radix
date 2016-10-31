package radix

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	goradix "github.com/armon/go-radix"
)

var pathsBench map[string]os.FileInfo

func TestSimple(t *testing.T) {
	tree := New()

	tree.Insert("", "root")
	tree.Insert("a", "a")
	tree.Insert("b", "b")
	tree.Insert("d", "d")
	tree.Insert("c", "c")
	tree.Insert("e", "e")
	tree.Insert("f", "f")
	tree.Insert("aa", "aa")

	tree.Insert("/", "root")
	tree.Insert("/aa", "aa")
	tree.Insert("/folder/", "folder1")
	tree.Insert("/folder/", "folder1-")
	tree.Insert("/folder/file1", "file1")
	tree.Insert("/folder/file2", "file2")
	tree.Insert("/folder/file3", "file3")
	tree.Insert("/folder/file4", "file4")
	tree.Insert("/folder2", "folder2")

	tree.Insert("/foo/bar", "coucou2")
	tree.Insert("/wesh", "wesh")
	tree.Insert("", "empty")

	if v, ok := tree.Get(""); ok && v.(string) != "empty" {
		t.Fatal("empty")
	}
	if v, ok := tree.Get("a"); ok && v.(string) != "a" {
		t.Fatal("a")
	}
	if v, ok := tree.Get("b"); ok && v.(string) != "b" {
		t.Fatal("b")
	}
	if _, ok := tree.Get("/fold"); ok {
		t.Fatal("nil")
	}
	if v, ok := tree.Get("/aa"); ok && v.(string) != "aa" {
		t.Fatal("/aa")
	}
	if v, ok := tree.Get("/folder/"); ok && v.(string) != "folder1-" {
		t.Fatal("folder1-")
	}
	if v, ok := tree.Get("/folder/file1"); ok && v.(string) != "file1" {
		t.Fatal("file1")
	}
	if v, ok := tree.Get("/folder/file4"); ok && v.(string) != "file4" {
		t.Fatal("file4")
	}
	tree.Foreach(func(val interface{}, key string) error {
		fmt.Println(key, val.(string))
		return nil
	})
}

func TestRandomString(t *testing.T) {
	l := 100
	m := make([]string, l)
	var tree *Tree

	tree = New()
	for i := 0; i < l; i++ {
		v := randString(30, "ab", false)
		tree.Insert(v, v)
		m[i] = v
	}

	for _, s := range m {
		v, ok := tree.Get(s)
		if !ok || v == nil || v.(string) != s {
			t.Fatal("simple rand on key " + s)
		}
	}

	tree = New()
	for i := 0; i < l; i++ {
		v := randString(15, "abcdefghifklmnopqrstuvwxyz", true)
		tree.Insert(v, v)
		m[i] = v
	}

	tree.Foreach(func(val interface{}, key string) error {
		if key != val.(string) {
			t.Fatal("bad key", key, val)
		}
		return nil
	})

	for _, s := range m {
		v, ok := tree.Get(s)
		if !ok || v == nil || v.(string) != s {
			t.Fatal("complex rand on key " + s)
		}
	}
}

func TestCloseup(t *testing.T) {
	tree := New()

	tree.Insert("/folder", "folder")
	tree.Insert("/folder/file1", "file1")
	tree.Insert("/folder/file2", "file2")
	tree.Insert("/folder/file3", "file3")
	tree.Insert("/folder/file4", "file4")
	tree.Insert("/folder2", "folder2")

	closeup := tree.Closeup("/folder/")
	closeup.Foreach(func(val interface{}, key string) error {
		fmt.Println(key, val.(string))
		return nil
	})
}

var result interface{}

func BenchmarkRandomInsertSelf(b *testing.B) {
	var replaced bool
	tree := New()
	for n := 0; n < b.N; n++ {
		k := randString(15, "ab", true)
		replaced = tree.Insert(k, k)
	}
	result = replaced
}

func BenchmarkRandomInsertGoRadix(b *testing.B) {
	var replaced bool
	tree := goradix.New()
	for n := 0; n < b.N; n++ {
		k := randString(15, "ab", true)
		_, replaced = tree.Insert(k, k)
	}
	result = replaced
}

func BenchmarkFilesSelf(b *testing.B) {
	var tree *Tree
	for n := 0; n < b.N; n++ {
		tree = New()
		for k, v := range pathsBench {
			tree.Insert(k, v)
		}
		for k, _ := range pathsBench {
			v, ok := tree.Get(k)
			if !ok || v == nil {
				os.Exit(1)
			}
		}
	}
}

func BenchmarkFilesGoRadix(b *testing.B) {
	var tree *goradix.Tree
	for n := 0; n < b.N; n++ {
		tree = goradix.New()
		for k, v := range pathsBench {
			tree.Insert(k, v)
		}
		for k, _ := range pathsBench {
			v, ok := tree.Get(k)
			if !ok || v == nil {
				os.Exit(1)
			}
		}
	}
}

func BenchmarkInsertFilesSelf(b *testing.B) {
	var tree *Tree
	for n := 0; n < b.N; n++ {
		tree = New()
		for k, v := range pathsBench {
			tree.Insert(k, v)
		}
	}
}

func BenchmarkInsertFilesGoRadix(b *testing.B) {
	var tree *goradix.Tree
	for n := 0; n < b.N; n++ {
		tree = goradix.New()
		for k, v := range pathsBench {
			tree.Insert(k, v)
		}
	}
}

func TestMain(m *testing.M) {
	benchDir := flag.String("test.benchdir", "", "Directory to run benchmark with")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	isBench := false
	flag.Visit(func(flag *flag.Flag) {
		if flag.Name == "test.bench" {
			isBench = true
		}
	})

	if isBench {
		wd, err := os.Getwd()
		checkError(err)

		if len(*benchDir) > 0 && (*benchDir)[:1] == "~" {
			usr, _ := user.Current()
			*benchDir = filepath.Join(usr.HomeDir, (*benchDir)[1:])
		}

		if *benchDir == "" {
			*benchDir = wd
		} else if !filepath.IsAbs(*benchDir) {
			*benchDir = filepath.Join(wd, *benchDir)
		}

		checkError(err)

		pathsBench = make(map[string]os.FileInfo)
		err = filepath.Walk(*benchDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			pathsBench[path] = f
			return nil
		})
		checkError(err)
	}

	os.Exit(m.Run())
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func randString(n int, letterBytes string, randlen bool) string {
	var b []byte
	if randlen {
		b = make([]byte, rand.Intn(n))
	} else {
		b = make([]byte, n)
	}
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
