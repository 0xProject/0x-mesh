// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

// Package registry provides access to Windows registry.
//
package registry

import (
	"syscall"
	"unsafe"

	"github.com/btcsuite/winsvc/winapi"
)

type Key struct {
	Handle syscall.Handle
}

func OpenKey(parent syscall.Handle, path string) (*Key, error) {
	var h syscall.Handle
	e := syscall.RegOpenKeyEx(
		parent, syscall.StringToUTF16Ptr(path),
		0, syscall.KEY_ALL_ACCESS, &h)
	if e != nil {
		return nil, e
	}
	return &Key{Handle: h}, nil
}

func (k *Key) Close() error {
	return syscall.RegCloseKey(k.Handle)
}

func (k *Key) CreateSubKey(name string) (nk *Key, openedExisting bool, err error) {
	var h syscall.Handle
	var d uint32
	e := winapi.RegCreateKeyEx(
		k.Handle, syscall.StringToUTF16Ptr(name),
		0, nil, winapi.REG_OPTION_NON_VOLATILE,
		syscall.KEY_ALL_ACCESS, nil, &h, &d)
	if e != nil {
		return nil, false, e
	}
	return &Key{Handle: h}, d == winapi.REG_OPENED_EXISTING_KEY, nil
}

func (k *Key) DeleteSubKey(name string) error {
	return winapi.RegDeleteKey(k.Handle, syscall.StringToUTF16Ptr(name))
}

func (k *Key) SetUInt32(name string, value uint32) error {
	return winapi.RegSetValueEx(
		k.Handle, syscall.StringToUTF16Ptr(name),
		0, syscall.REG_DWORD,
		(*byte)(unsafe.Pointer(&value)), uint32(unsafe.Sizeof(value)))
}

func (k *Key) SetString(name string, value string) error {
	buf := syscall.StringToUTF16(value)
	return winapi.RegSetValueEx(
		k.Handle, syscall.StringToUTF16Ptr(name),
		0, syscall.REG_SZ,
		(*byte)(unsafe.Pointer(&buf[0])), uint32(len(buf)*2))
}
