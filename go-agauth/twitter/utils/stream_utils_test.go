package utils

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var defaultTestTimeout = time.Second * 1

// assertDone asserts that the empty struct channel is closed before the given
// timeout elapses.
func assertDone(t *testing.T, ch <-chan struct{}, timeout time.Duration) {
	select {
	case <-ch:
		_, more := <-ch
		assert.False(t, more)
	case <-time.After(timeout):
		t.Errorf("expected channel to be closed within timeout %v", timeout)
	}
}

func TestStopped(t *testing.T) {
	done := make(chan struct{})
	assert.False(t, Stopped(done))
	close(done)
	assert.True(t, Stopped(done))
}

func TestSleepOrDone_Sleep(t *testing.T) {
	wait := time.Nanosecond * 20
	done := make(chan struct{})
	completed := make(chan struct{})
	go func() {
		SleepOrDone(wait, done)
		close(completed)
	}()
	// wait for goroutine SleepOrDone to sleep
	assertDone(t, completed, defaultTestTimeout)
}

func TestSleepOrDone_Done(t *testing.T) {
	wait := time.Second * 5
	done := make(chan struct{})
	completed := make(chan struct{})
	go func() {
		SleepOrDone(wait, done)
		close(completed)
	}()
	// close done, interrupting SleepOrDone
	close(done)
	// assert that SleepOrDone exited, closing completed
	assertDone(t, completed, defaultTestTimeout)
}

func TestStreamResponseBodyReader(t *testing.T) {
	cases := []struct {
		in   []byte
		want [][]byte
	}{
		{
			in: []byte("foo\r\nbar\r\n"),
			want: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
			},
		},
		{
			in: []byte("foo\nbar\r\n"),
			want: [][]byte{
				[]byte("foo\nbar"),
			},
		},
		{
			in: []byte("foo\r\n\r\n"),
			want: [][]byte{
				[]byte("foo"),
				[]byte(""),
			},
		},
		{
			in: []byte("foo\r\nbar"),
			want: [][]byte{
				[]byte("foo"),
				[]byte("bar"),
			},
		},
		{
			// Message length is more than bufio.MaxScanTokenSize, which can't be
			// parsed by bufio.Scanner with default buffer size.
			in: []byte(strings.Repeat("X", bufio.MaxScanTokenSize+1) + "\r\n"),
			want: [][]byte{
				[]byte(strings.Repeat("X", bufio.MaxScanTokenSize+1)),
			},
		},
	}

	for _, c := range cases {
		body := bytes.NewReader(c.in)
		reader := NewStreamResponseBodyReader(body)

		for i, want := range c.want {
			data, err := reader.ReadNext()
			if err != nil {
				t.Errorf("reader(%q).readNext() * %d: err == %q, want nil", c.in, i, err)
			}
			if !bytes.Equal(data, want) {
				t.Errorf("reader(%q).readNext() * %d: data == %q, want %q", c.in, i, data, want)
			}
		}

		data, err := reader.ReadNext()
		if err != io.EOF {
			t.Errorf("reader(%q).readNext() * %d: err == %q, want io.EOF", c.in, len(c.want), err)
		}
		if len(data) != 0 {
			t.Errorf("reader(%q).readNext() * %d: data == %q, want \"\"", c.in, len(c.want), data)
		}
	}
}
