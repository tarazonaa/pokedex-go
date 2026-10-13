package pokedexcache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddAndGet(t *testing.T) {
	const interval = 5 * time.Second
	cache := NewCache(interval)
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "someurl",
			val: []byte("Some dataa"),
		},
	}

	for i, test := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache.Add(test.key, test.val)
			if val, ok := cache.Get(test.key); !ok {
				t.Errorf("Expected to find key %s", test.key)
				return
			} else {
				if string(val) != string(test.val) {
					t.Errorf("Mismatched values: got %s and expected %s", string(val), string(test.val))
					return
				}
			}
		})
	}
}
