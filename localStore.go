package main

import "sync"

// LocalStore - Local Cache for checking whether a URL is already crawled or not
type LocalStore struct {
	sync.RWMutex
	store map[string]bool
}

func (ls *LocalStore) add(url string) bool {
	ls.Lock()
	defer ls.Unlock()
	ls.store[url] = true
	return true
}

func (ls *LocalStore) contains(url string) bool {
	ls.RLock()
	defer ls.RUnlock()
	val, present := ls.store[url]
	return val && present
}
