package pkg

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	dbg "github.com/visionmedia/go-debug"
)

var (
	infileName       = filepath.Join("testdata", "infile.txt")
	fileCompressed   = filepath.Join("testdata", "compressed.txt.sz")
	fileDecompressed = filepath.Join("testdata", "decompressed.txt")
	debug            = dbg.Debug("TEST")
)

type limitWriter struct {
	w io.Writer
	n int
}

func (w *limitWriter) Write(p []byte) (n int, err error) {
	debug("n: %d", w.n)

	if len(p) > w.n {
		p = p[:w.n]
	}
	if len(p) > 0 {
		n, err = w.w.Write(p)
		w.n -= n
	}
	if w.n == 0 {
		err = errors.New("past write limit")
	}
	debug("n: %d", w.n)
	return
}

func TestSnap(t *testing.T) {

	file := make([]byte, 1<<24)

	_, err := rand.Read(file)
	if err != nil {
		panic(err)
	}

	h := md5.New()

	buf := bytes.NewBuffer(file)

	io.Copy(h, buf)
	expected := h.Sum(nil)
	t.Logf("MD5 sum of original: %x", expected)

	out := bytes.NewBuffer(nil)

	s := NewSnapper(buf, out)

	assert.NoError(t, s.Snap(), "Error while compressing")

	buf.Reset()

	u := NewSnapper(bytes.NewReader(out.Bytes()), buf)
	assert.NoError(t, u.Unsnap(), "Error while decompressing")

	io.Copy(h, buf)

	result := h.Sum(nil)

	assert.Equal(t, expected, result)

}

// func TestLimitedSpaceSnap(t *testing.T) {

// 	outBuf := bytes.Buffer{}

// 	lw := &limitWriter{w: &outBuf, n: 1}

// 	in, err := os.Open(infileName)
// 	if err != nil {
// 		t.Fatalf("unable to open testfile '%s': %v", infileName, err)
// 	}
// 	defer in.Close()

// 	s := NewSnapper(in, lw)

// 	if err := s.Snap(); err == nil {
// 		t.Error("No error raised compressing to small buffer")
// 	}
// }

// func TestLimitedSpaceUnsnap(t *testing.T) {

// 	/* ----------------- Prepare data ----------------- */
// 	outbuf := &bytes.Buffer{}
// 	in, err := os.Open(infileName)

// 	if err != nil {
// 		t.Fatalf("unable to open testfile '%s': %v", infileName, err)
// 	}

// 	lw := &limitWriter{w: &bytes.Buffer{}, n: 1}

// 	/* ----------------- Testing ----------------- */
// 	s := NewSnapper(in, outbuf)
// 	s.Snap()

// 	u := NewSnapper(outbuf, lw)
// 	if err := u.Unsnap(); err == nil {
// 		t.Error("No error raised decompressing to small buffer")
// 	}
// }
