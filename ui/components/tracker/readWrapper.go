package tracker

import (
	"io"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	mp3 "github.com/dece2183/go-stream-mp3"
	"github.com/dece2183/yamusic-tui/log"
	"github.com/dece2183/yamusic-tui/stream"
)

const (
	_PROGRESS_UPDATE_PERIOD = 33 * time.Millisecond
)

type readWrapper struct {
	mu             sync.Mutex
	program        *tea.Program
	decoder        *mp3.Decoder
	trackBuffer    *stream.BufferedStream
	trackBuffered  bool
	trackDone      bool
	lastUpdateTime time.Time
}

func (w *readWrapper) NewReader(reader *stream.BufferedStream) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.trackBuffered = false
	w.trackDone = false
	w.trackBuffer = reader

	decoder, err := mp3.NewDecoder(reader)
	if err != nil {
		w.trackBuffer = nil
		w.decoder = nil
		log.Print(log.LVL_ERROR, "failed to create mp3 decoder: %s", err)
		return
	}

	w.decoder = decoder
	w.lastUpdateTime = time.Now()
}

func (w *readWrapper) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.decoder != nil {
		w.decoder.Seek(0, io.SeekStart)
	}

	if w.trackBuffer != nil {
		w.trackBuffer.Close()
	}
}

func (w *readWrapper) Read(dest []byte) (n int, err error) {
	w.mu.Lock()
	decoder := w.decoder
	buffer := w.trackBuffer
	w.mu.Unlock()

	if buffer == nil || decoder == nil {
		err = io.EOF
		return
	}

	n, err = decoder.Read(dest)
	if err != nil && err != io.EOF {
		if buffer.Error() != nil {
			err = buffer.Error()
			log.Print(log.LVL_ERROR, "buffering error: %s", err)
			go w.program.Send(STOP)
			return
		}
		// bypass mp3 decoding error after rewinding
		log.Print(log.LVL_WARNIGN, "mp3 decoding error: %s", err)
		err = nil
	}

	w.mu.Lock()
	if w.trackBuffer != buffer {
		w.mu.Unlock()
		return
	}

	if buffer.IsBuffered() && !w.trackBuffered {
		w.trackBuffered = true
		go w.program.Send(BUFFERING_COMPLETE)
	}

	if buffer.IsDone() && !w.trackDone {
		w.trackDone = true
		decoder.Seek(0, io.SeekStart)
		w.mu.Unlock()
		buffer.Close()
		go w.program.Send(NEXT)
		return
	}

	if !w.trackDone && time.Since(w.lastUpdateTime) > _PROGRESS_UPDATE_PERIOD {
		w.lastUpdateTime = time.Now()
		w.mu.Unlock()
		fraction := ProgressControl(buffer.Progress())
		go w.program.Send(fraction)
		return
	}

	w.mu.Unlock()
	return
}

func (w *readWrapper) Seek(offset int64, whence int) (int64, error) {
	w.mu.Lock()
	decoder := w.decoder
	if decoder == nil {
		w.mu.Unlock()
		return 0, io.EOF
	}
	w.lastUpdateTime = time.Now()
	w.mu.Unlock()
	return decoder.Seek(offset, whence)
}

func (w *readWrapper) Length() int64 {
	w.mu.Lock()
	buffer := w.trackBuffer
	w.mu.Unlock()
	if buffer == nil {
		return 0
	}
	return buffer.Length()
}

func (w *readWrapper) Progress() float64 {
	w.mu.Lock()
	buffer := w.trackBuffer
	w.mu.Unlock()
	if buffer == nil {
		return 0
	}
	return buffer.Progress()
}
