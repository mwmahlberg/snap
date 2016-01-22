package internal

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	dbg "github.com/tj/go-debug"
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
	h := md5.New()

	in, err := os.Open(infileName)
	if err != nil {
		t.Fatalf("unable to open testfile '%s': %v", infileName, err)
	}

	io.Copy(h, in)
	expected := h.Sum(nil)
	t.Logf("MD5 sum of original '%s': %x", infileName, expected)

	out, _ := os.OpenFile(fileCompressed, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	in.Seek(0, 0)

	s := NewSnapper(in, out)

	if err := s.Snap(); err != nil {
		t.Error(err)
	}

	in.Close()
	out.Close()

	in, err = os.Open(fileCompressed)
	if err != nil {
		t.Fatalf("unable to open compressed '%s': %v", fileDecompressed, err)
	}

	out, err = os.OpenFile(fileDecompressed, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0640)
	if err != nil {
		t.Fatalf("unable to open target '%s' for decompression: %v", fileDecompressed, err)
	}
	defer out.Close()
	u := NewSnapper(in, out)

	if err := u.Unsnap(); err != nil {
		t.Errorf("error during decompression: %v", err)
	}
	in.Close()

	out.Seek(0, 0)
	h.Reset()

	io.Copy(h, out)

	result := h.Sum(nil)

	if fmt.Sprintf("%x", result) != fmt.Sprintf("%x", expected) {
		t.Errorf("Checksum of uncompressed file (%x) does not match original(%x)", result, expected)
	}

}

func TestLimitedSpaceSnap(t *testing.T) {

	outBuf := bytes.Buffer{}

	lw := &limitWriter{w: &outBuf, n: 1}

	in, err := os.Open(infileName)
	if err != nil {
		t.Fatalf("unable to open testfile '%s': %v", infileName, err)
	}
	defer in.Close()

	s := NewSnapper(in, lw)

	if err := s.Snap(); err == nil {
		t.Error("No error raised compressing to small buffer")
	}
}

func TestLimitedSpaceUnsnap(t *testing.T) {

	/* ----------------- Prepare data ----------------- */
	outbuf := &bytes.Buffer{}
	in, err := os.Open(infileName)

	if err != nil {
		t.Fatalf("unable to open testfile '%s': %v", infileName, err)
	}

	lw := &limitWriter{w: &bytes.Buffer{}, n: 1}

	/* ----------------- Testing ----------------- */
	s := NewSnapper(in, outbuf)
	s.Snap()

	u := NewSnapper(outbuf, lw)
	if err := u.Unsnap(); err == nil {
		t.Error("No error raised decompressing to small buffer")
	}
}
