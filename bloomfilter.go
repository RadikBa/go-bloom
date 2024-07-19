package bloomfilter

import (
	"math"

	"github.com/bits-and-blooms/bitset"
	"github.com/spaolacci/murmur3"
)

type Filter struct {
	bits *bitset.BitSet
	k    uint32
}

func estimateBloomFilterParams(n int, prob float64) (uint, uint32) {
	m := math.Ceil(float64(n) * math.Log(prob) * (-1) / (math.Pow(math.Log(2), 2)))
	k := math.Ceil(math.Log(2) * m / float64(n))
	return uint(m), uint32(k)
}

func NewFilter(size uint, k uint32) *Filter {
	return &Filter{
		bits: bitset.New(size),
		k:    k,
	}
}

func (bf *Filter) hashFunction(key []byte, seed uint32) uint32 {
	return murmur3.Sum32WithSeed(key, seed)
}

func (bf *Filter) Add(key []byte) {
	for seed := uint32(0); seed < bf.k; seed++ {
		index := uint(bf.hashFunction(key, seed)) % bf.bits.Len()
		bf.bits.Set(index)
	}
}

func (bf *Filter) Test(key []byte) bool {
	for seed := uint32(0); seed < bf.k; seed++ {
		index := uint(bf.hashFunction(key, seed)) % bf.bits.Len()
		if !bf.bits.Test(index) {
			return false
		}
	}
	return true
}
