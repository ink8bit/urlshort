package shorten_test

import (
	"testing"

	"urlshort/internal/shorten"
)

func TestGenRandomStr(t *testing.T) {
	for range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		t.Run(t.Name(), func(t *testing.T) {
			str1 := shorten.GenRandomStr()
			str2 := shorten.GenRandomStr()
			if str1 == str2 {
				t.Errorf("Generated strings should differ; got %q value for both strings", str1)
			}
		})
	}
}

func BenchmarkGenRandomStr(b *testing.B) {
	for n := 0; n < b.N; n++ {
		shorten.GenRandomStr()
	}
}
