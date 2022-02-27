package main

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
)

type bloomFilter interface {
	add(item string)

	// `false` means the item is definitely not in the set
	// `true` means the item might be in the set
	maybeContains(item string) bool

	// Number of bytes used in any underlying storage
	memoryUsage() int
}

type trivialBloomFilter struct {
	data  []uint64
	hash1 hash.Hash64
	hash2 hash.Hash64
	len   uint64
}

func newTrivialBloomFilter() *trivialBloomFilter {
	return &trivialBloomFilter{
		data:  make([]uint64, 10000),
		hash1: fnv.New64(),
		hash2: fnv.New64a(),
		len:   10000,
	}
}

func (b *trivialBloomFilter) add(item string) {
	k1, k2 := b.getHashes(item)

	b.data[k1/64] |= (1 << (k1 % 64))
	b.data[k2/64] |= (1 << (k2 % 64))
}

func (b *trivialBloomFilter) maybeContains(item string) bool {
	k1, k2 := b.getHashes(item)

	v1 := b.data[k1/64] & (1 << (k1 % 64))
	v2 := b.data[k2/64] & (1 << (k2 % 64))

	return v1 != 0 && v2 != 0
}

func (b *trivialBloomFilter) getHashes(item string) (uint64, uint64) {
	modVal := b.len * 64

	b.hash1.Write([]byte(item))
	k1 := b.hash1.Sum64() % modVal
	b.hash1.Reset()

	b.hash2.Write([]byte(item))
	k2 := b.hash2.Sum64() % modVal
	b.hash2.Reset()

	return k1, k2
}

func (b *trivialBloomFilter) memoryUsage() int {
	return binary.Size(b.data)
}
