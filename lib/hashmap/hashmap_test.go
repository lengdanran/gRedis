// Package hashmap lengdanran 2024/4/1 16:31
package hashmap

import (
	"github.com/google/uuid"
	"log"
	"testing"
)

func constructKeysAndValues(size int) []Entry {
	var entries []Entry
	for i := 0; i < size; i++ {
		id := uuid.New().String()
		entries = append(entries, Entry{Key: id, Value: id})
	}
	return entries
}

func TestHashMap_Put(t *testing.T) {
	hashMap := NewHashMap()
	ens := constructKeysAndValues(2000)
	for i := 0; i < len(ens); i++ {
		hashMap.Put(ens[i])
	}
	for i := 0; i < len(ens); i++ {
		val := hashMap.Get(ens[i].Key)
		valStr := val.(string)
		// log.Printf("Key=%s, Val=%s, ValInMap=%s, Eqs=%v\n", ens[i].Key, ens[i].Value, valStr, valStr == ens[i].Value.(string))
		if valStr != ens[i].Value.(string) {
			log.Printf("Not PASS ===> Key=%s, Val=%s, ValInMap=%s, Eqs=%v\n", ens[i].Key, ens[i].Value, valStr, valStr == ens[i].Value.(string))
			t.Failed()
		}
	}
}

func TestHashMap_Del(t *testing.T) {
	hashMap := NewHashMap()
	ens := constructKeysAndValues(500)
	for i := 0; i < len(ens); i++ {
		hashMap.Put(ens[i])
	}
	for i := 0; i < len(ens); i++ {
		if i%5 == 0 {
			hashMap.Del(ens[i].Key)
		}
	}
	for i := 0; i < len(ens); i++ {
		if i%5 == 0 {
			if hashMap.Contains(ens[i].Key) {
				log.Printf("Key=%s delete failed!!!", ens[i].Key)
				t.Fail()
			}
		} else {
			if !hashMap.Contains(ens[i].Key) {
				log.Printf("Key=%s shouldn't be deleted !!!", ens[i].Key)
				t.Fail()
			}
		}
	}
}

func TestHashMap_Keys(t *testing.T) {
	hashMap := NewHashMap()
	ens := constructKeysAndValues(50)
	ensMaps := make(map[string]interface{})
	for _, en := range ens {
		ensMaps[en.Key] = en.Value
		hashMap.Put(en)
	}
	keys := hashMap.Keys()
	for _, k := range keys {
		_, ok := ensMaps[k]
		if !ok {
			log.Printf("Key %s should in hashmap!", k)
			t.Fail()
		}
	}
}

func TestHashMap_Entries(t *testing.T) {
	hashMap := NewHashMap()
	ens := constructKeysAndValues(50)
	ensMaps := make(map[string]interface{})
	for _, en := range ens {
		ensMaps[en.Key] = en.Value
		hashMap.Put(en)
	}
	entries := hashMap.Entries()
	for _, entry := range entries {
		mapEntryVal, ok := ensMaps[entry.Key]
		if !ok {
			log.Printf("Entry %s should in hashmap!", entry.Key)
			t.Fail()
		}
		entryString := entry.Value.(string)
		if mapEntryVal != entryString {
			log.Printf("EntryValue should be the same! HashMapValue=%s, RightValue=%s", entryString, mapEntryVal)
			t.Fail()
		}
	}
}
