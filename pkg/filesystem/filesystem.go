package filesystem

import (
	"log"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/initialed85/glue/pkg/worker"

	"github.com/initialed85/cameranator/pkg/media/metadata"
)

type Watcher struct {
	path          string
	matcher       *regexp.Regexp
	onFileCreate  func(string, float64)
	onFileWrite   func(string, float64)
	ticker        *time.Ticker
	watcher       *fsnotify.Watcher
	blockedWorker *worker.BlockedWorker
}

func NewWatcher(
	path string,
	matcher *regexp.Regexp,
	onFileCreate func(string, float64),
	onFileWrite func(string, float64),
) *Watcher {
	w := Watcher{
		path:         path,
		matcher:      matcher,
		onFileCreate: onFileCreate,
		onFileWrite:  onFileWrite,
	}

	w.blockedWorker = worker.NewBlockedWorker(
		w.onStart,
		w.work,
		w.onStop,
	)

	return &w
}

func (w *Watcher) onStart() {
	var err error

	for {
		w.watcher, err = fsnotify.NewWatcher()
		if err != nil {
			log.Printf("warning: failed to watch %#+v because %v; will try again...", w.path, err)
			time.Sleep(time.Second)
			continue
		}

		break
	}

	for {
		err = w.watcher.Add(w.path)
		if err != nil {
			log.Printf("warning: failed to watch %#+v because %v; will try again...", w.path, err)
			time.Sleep(time.Second)
			continue
		}

		break
	}

	w.ticker = time.NewTicker(time.Second)
}

func (w *Watcher) handleEvent(event fsnotify.Event) {
	path := event.Name

	if w.matcher != nil {
		if !w.matcher.Match([]byte(path)) {
			return
		}
	}

	create := event.Op&fsnotify.Create == fsnotify.Create
	write := false

	if !create {
		write = event.Op&fsnotify.Write == fsnotify.Write
	}

	if !(create || write) {
		return
	}

	size, err := metadata.GetFileSize(path)
	if err != nil {
		log.Printf("warning: failed to get file size for %#+v because %#+v", path, err)
		return
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		// TODO: fix unbounded goroutine usage
		go w.onFileCreate(path, size)
	}

	if event.Op&fsnotify.Write == fsnotify.Write {
		// TODO: fix unbounded goroutine usage
		go w.onFileWrite(path, size)
	}
}

func (w *Watcher) work() {
	select {
	case event, ok := <-w.watcher.Events:
		if !ok {
			log.Printf("warning: event %#+v not ok; retrying...", event)
			time.Sleep(time.Second)
			return
		}

		if w.matcher != nil {
			if !w.matcher.Match([]byte(event.Name)) {
				return
			}
		}

		w.handleEvent(event)
	case err, ok := <-w.watcher.Errors:
		if !ok {
			log.Printf("warning: error %#+v not ok; retrying...", err)
			time.Sleep(time.Second)
			return
		}

		log.Printf("warning: watcher threw %#+v", err) // TODO: probably do something with this
	case <-w.ticker.C:
		time.Sleep(0) // noop for debounce and so we can exit
	}
}

func (w *Watcher) onStop() {
	_ = w.watcher.Remove(w.path)
	w.watcher = nil

	w.ticker.Stop()
	w.ticker = nil
}

func (w *Watcher) Start() {
	w.blockedWorker.Start()
}

func (w *Watcher) Stop() {
	w.blockedWorker.Stop()
}
