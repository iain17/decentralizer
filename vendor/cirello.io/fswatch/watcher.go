// Copyright 2017 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func monitorWorkDir(ctx context.Context, wd string, patterns, ignores []string) (<-chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	memo := make(map[string]struct{})
	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			for _, skipDir := range ignores {
				if skipDir == "" {
					continue
				}
				if strings.HasPrefix(path, filepath.Join(wd, skipDir)) {
					return filepath.SkipDir
				}
			}
			return nil
		}
		for _, p := range patterns {
			if match(p, path) {
				dir := filepath.Dir(path)
				if _, ok := memo[dir]; !ok {
					memo[dir] = struct{}{}
					_ = watcher.Add(dir)
				}
				break
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	log.Println("monitoring", len(memo), "directories")

	changeds := consumeFsnotifyEvents(ctx, patterns, watcher)
	triggereds := triggerRestarts(ctx, changeds)
	return triggereds, nil
}

func consumeFsnotifyEvents(ctx context.Context, patterns []string, watcher *fsnotify.Watcher) chan struct{} {
	changeds := make(chan struct{})

	go func() {
		defer watcher.Close()
		for {
			select {
			case <-ctx.Done():
				return
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}
				for _, p := range patterns {
					if match(p, event.Name) {
						select {
						case changeds <- struct{}{}:
						default:
						}
						break
					}
				}
			case err := <-watcher.Errors:
				log.Println("fswatch error:", err)
			}
		}
	}()
	return changeds
}

func triggerRestarts(ctx context.Context, changeds chan struct{}) chan struct{} {
	triggereds := make(chan struct{})
	go func() {
		lastRun := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-changeds:
				triggereds <- struct{}{}
			}
			const coolDownPeriod = 7500 * time.Millisecond
			if sinceLastRun := time.Since(lastRun); sinceLastRun < coolDownPeriod {
				log.Println("too active, pausing restarts")
				time.Sleep(coolDownPeriod - sinceLastRun)
			}
			lastRun = time.Now()
		}
	}()
	return triggereds
}

func match(p, path string) bool {
	base, dir := filepath.Base(path), filepath.Dir(path)
	pbase, pdir := filepath.Base(p), filepath.Dir(p)

	if matched, err := filepath.Match(pbase, base); err != nil || !matched {
		return false
	}

	if pdir == "." {
		return true
	}
	subpatterns := strings.Split(pdir, "**")

	tmp := dir
	for _, subp := range subpatterns {
		if subp == "" {
			continue
		}
		subp = filepath.Clean(subp)
		t := strings.Replace(tmp, subp, "", 1)
		if t == tmp {
			return false
		}
		tmp = t
	}

	return true
}
