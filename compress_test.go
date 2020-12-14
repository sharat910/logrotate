package logrotator

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestCompress(t *testing.T) {
	sf := "./test/test.json"
	df := "./test/test.json.gz"
	err := Compress(sf, false)
	assert.Nil(t, err)
	stats, err := os.Stat(df)
	assert.Nil(t, err)
	assert.False(t, stats.IsDir())
	assert.Equal(t,filepath.Base(df), stats.Name())
	assert.Nil(t, os.Remove(df))
}

func BenchmarkCompress(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Compress("./test/test.json", false)
	}
}
