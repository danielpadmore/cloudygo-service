package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type File struct {
	path       string
	userConfig interface{}
	watcher    *fsnotify.Watcher
	updated    func()
}

func New(fp string, c interface{}, updated func()) (*File, error) {
	ap, err := filepath.Abs(fp)
	if err != nil {
		return nil, err
	}

	f := &File{path: ap, userConfig: c, updated: updated}
	go f.watch(ap)

	time.Sleep(10 * time.Millisecond)

	return f, f.loadData()
}

func (f *File) Close() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	f.watcher.Close()
}

func (f *File) loadData() error {
	cf, err := os.OpenFile(f.path, os.O_RDONLY, 0655)
	if err != nil {
		return err
	}
	defer cf.Close()

	jd := json.NewDecoder(cf)
	return jd.Decode(f.userConfig)
}

func (f *File) watch(filepath string) {
	var err error
	f.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Println("Error", err)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-f.watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Rename == fsnotify.Rename ||
					event.Op&fsnotify.Remove == fsnotify.Remove ||
					event.Op&fsnotify.Chmod == fsnotify.Chmod {
					err := f.loadData()
					if err != nil {
						log.Println("error", err)
						return
					}

					if f.updated != nil {
						f.updated()
					}
				}
			case err, ok := <-f.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error", err)
			}
		}
	}()

	err = f.watcher.Add(filepath)
	if err != nil {
		log.Println("error", err)
	}
	<-done
}
