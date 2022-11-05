package main

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"sync/atomic"
	"unsafe"
)

type vlanMap struct {
	r map[uint32]bool
}

func newVlanMap() *vlanMap {
	return &vlanMap{
		r: make(map[uint32]bool),
	}
}

func (m *vlanMap) Set(vid uint32, action bool) {
	m.r[vid] = action
}

func (m *vlanMap) Get(vid uint32) bool {
	if val, ok := m.r[vid]; ok {
		return val
	}
	return true
}

var vlanMapPointer unsafe.Pointer

func vlanMapLookup(vid uint32) bool {
	return (*vlanMap)(atomic.LoadPointer(&vlanMapPointer)).Get(vid)
}

func vlanMapReload() error {
	data, err := ioutil.ReadFile(flagVlanPath)
	if err != nil {
		return err
	}

	fileMap := make(map[uint32]bool)
	if err := yaml.Unmarshal(data, fileMap); err != nil {
		return err
	}

	l := newVlanMap()
	for vid, action := range fileMap {
		l.Set(vid, action)
	}

	atomic.StorePointer(&vlanMapPointer, unsafe.Pointer(l))

	return nil
}
