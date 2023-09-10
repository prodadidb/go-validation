package utils_test

import (
	"testing"

	"github.com/prodadidb/go-validation/is/utils"
)

func BenchmarkContains(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		utils.Contains("a0b01c012deffghijklmnopqrstu0123456vwxyz", "0123456789")
	}
}

func BenchmarkMatches(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		utils.Matches("alfkjl12309fdjldfsa209jlksdfjLAKJjs9uJH234", "[\\w\\d]+")
	}
}
