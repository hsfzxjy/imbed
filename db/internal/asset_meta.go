package internal

import (
	"encoding/binary"
	"fmt"
	"sync"
	"unsafe"

	"go.etcd.io/bbolt"
)

func (m *AssetMeta) EarliestOID() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.earliestOID
}

func (m *AssetMeta) LatestOID() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.latestOID
}

type AssetMeta struct {
	mu          sync.RWMutex
	dirty       bool
	earliestOID uint64
	latestOID   uint64
}

var assetMetaKey = []byte("1")

func DecodeAssetMeta(m *AssetMeta, buc *bbolt.Bucket) error {
	data := buc.Get(assetMetaKey)
	if len(data) == 0 {
		m.earliestOID = 0
		m.latestOID = 0
		return nil
	}
	if uintptr(len(data)) != 2*unsafe.Sizeof(uint64(0)) {
		return fmt.Errorf("invalid asset meta data")
	}
	m.earliestOID = binary.BigEndian.Uint64(data[:8])
	m.latestOID = binary.BigEndian.Uint64(data[8:])
	return nil
}

var errAssetOIDOverflow = fmt.Errorf("asset OID overflow")

func (m *AssetMeta) isFull() bool {
	return m.latestOID+1 == m.earliestOID
}

func (m *AssetMeta) getNextOID() (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isFull() {
		return 0, errAssetOIDOverflow
	}
	m.dirty = true
	m.latestOID++
	if m.latestOID == 0 {
		m.latestOID++
	}
	if m.earliestOID == 0 {
		m.earliestOID = m.latestOID
	}
	return m.latestOID, nil
}

func AssetMetaNextOID(m *AssetMeta) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.isFull() {
		return 0, errAssetOIDOverflow
	}
	m.dirty = true
	m.latestOID++
	if m.latestOID == 0 {
		m.latestOID++
	}
	if m.earliestOID == 0 {
		m.earliestOID = m.latestOID
	}
	return m.latestOID, nil
}

func WriteAssetMeta(m *AssetMeta, buc *bbolt.Bucket) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if !m.dirty {
		return nil
	}
	var b [16]byte
	binary.BigEndian.PutUint64(b[:8], m.earliestOID)
	binary.BigEndian.PutUint64(b[8:], m.latestOID)
	return buc.Put(assetMetaKey, b[:])
}

func SetAssetMeta(m *AssetMeta, newEarliest, newLatest uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.earliestOID = newEarliest
	m.latestOID = newLatest
	m.dirty = true
}
