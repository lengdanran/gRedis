// Package hashmap lengdanran 2024/4/1 14:32
package hashmap

import (
	"log"
	"math"
)

type EntryVal interface{}

type Entry struct {
	HashCode int
	Key      string
	Value    EntryVal
	Next     *Entry
}

type HashMap struct {
	LoadFactor float64 // 扩容负载因子
	Size       int
	Cap        int
	Slots      []*Entry
	threshold  int
}

func NewHashMap() *HashMap {
	return &HashMap{
		LoadFactor: 0.75,
		Size:       0,
		Cap:        1 << 4,
	}
}

/**
 * 返回一个 >= cap的最小的 2 的次方数
 *
 * @param cap 给定的容量大小
 * @return >= cap的最小的 2 的次方数
 */
func computeCapacity(cap int) (size int) {
	if cap <= 16 {
		return 16
	}
	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return math.MaxInt32
	}
	return n + 1
}

const prime32 = uint32(16777619)

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (h *HashMap) Contains(key string) bool {
	return h.getEntry(key) != nil
}

func (h *HashMap) ResizeSlots() {
	oldSlots := h.Slots
	oldCap := h.Cap
	newSlots := make([]*Entry, oldCap<<1)
	for i := 0; i < oldCap; i++ {
		highHead := Entry{}
		var highPtr = &highHead
		lowHead := Entry{}
		var lowPtr = &lowHead
		curEntry := oldSlots[i]

		for curEntry != nil {
			highFlag := false

			if curEntry.HashCode&oldCap == 0 {
				highFlag = false
			} else {
				highFlag = true
			}
			if highFlag { // 高位链
				highPtr.Next = curEntry
				highPtr = highPtr.Next
			} else {
				lowPtr.Next = curEntry
				lowPtr = lowPtr.Next
			}
			curEntry = curEntry.Next
		}

		newSlots[i] = lowHead.Next
		newSlots[i+oldCap] = highHead.Next
		lowPtr.Next = nil
		highPtr.Next = nil
	}
	h.Slots = newSlots
	h.Cap = oldCap << 1
	h.threshold = int(float64(len(h.Slots)) * h.LoadFactor)
}

func (h *HashMap) Put(entry Entry) {
	hash := h.hash(entry.Key)
	entry.HashCode = hash
	if h.Slots == nil { // 说明当前的散列表还没初始化，分配对应的存储空间
		h.Slots = make([]*Entry, h.Cap)
		h.threshold = int(float64(len(h.Slots)) * h.LoadFactor)
	}
	index := h.getIndex(hash)
	ptr := h.Slots[index]
	if ptr == nil {
		// 如果当前散列表的桶位中没有元素，则直接将构造好的 Node 对象加入该位置，即可
		h.Slots[index] = &entry
		h.Size++
	} else {
		// 如果第一个元素不为空-->说明发生了 hash 冲突
		for ptr != nil {
			if ptr.Key == entry.Key {
				ptr.Value = entry.Value
				return
			}
			if ptr.Next == nil {
				ptr.Next = &entry
			}
			ptr = ptr.Next
		}
		h.Size++
	}
	if h.Size > h.threshold {
		h.ResizeSlots()
	}
}

func (h *HashMap) getEntry(key string) *Entry {
	if h.Size == 0 {
		return nil
	}
	hash := h.hash(key)
	index := h.getIndex(hash)

	ptr := h.Slots[index]
	if ptr == nil {
		return ptr
	}
	for ptr != nil {
		if ptr.Key == key {
			return ptr
		}
		ptr = ptr.Next
	}
	return ptr
}

func (h *HashMap) Get(key string) EntryVal {
	entry := h.getEntry(key)
	if entry == nil {
		return nil
	}
	return entry.Value
}

func (h *HashMap) Del(key string) EntryVal {
	hash := h.hash(key)
	index := h.getIndex(hash)
	ptr := h.Slots[index]
	if ptr == nil {
		log.Printf("%s is not in current hashmap. do nothing.", key)
		return nil
	}
	prePtr := &Entry{}
	var headPtr = prePtr
	prePtr.Next = ptr
	for ptr != nil {
		if ptr.Key == key {
			// found the entry need to be deleted.
			prePtr.Next = ptr.Next
			h.Size--
			h.Slots[index] = headPtr.Next
			return ptr.Value
		} else {
			prePtr = ptr
			ptr = ptr.Next
		}
	}
	log.Printf("%s is not in current hashmap. do nothing.", key)
	return nil
}

func (h *HashMap) Len() int {
	return h.Size
}

func (h *HashMap) Keys() []string {
	var keys []string
	for i := 0; i < h.Cap; i++ {
		ptr := h.Slots[i]
		if ptr != nil {
			for ptr != nil {
				keys = append(keys, ptr.Key)
				ptr = ptr.Next
			}
		}
	}
	return keys
}

func (h *HashMap) Entries() []*Entry {
	var entries []*Entry
	for i := 0; i < h.Cap; i++ {
		ptr := h.Slots[i]
		if ptr != nil {
			for ptr != nil {
				entries = append(entries, ptr)
				ptr = ptr.Next
			}
		}
	}
	return entries
}

func (h *HashMap) hash(key string) int {
	if len(key) == 0 {
		return 0
	}
	uHashCode := fnv32(key)
	return int(uHashCode) ^ int(uHashCode>>16)
}

func (h *HashMap) getIndex(hash int) int {
	return (len(h.Slots) - 1) & hash
}
