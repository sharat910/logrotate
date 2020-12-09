package logrotator

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogRotator_Rotation(t *testing.T) {
	temp := os.TempDir()
	fp := filepath.Join(temp, "test.log")

	header := []byte("header\n")
	lr, err := New(fp, Header(header))
	assert.Nil(t, err)

	now := time.Now()
	lr.nowFunc = func() time.Time {
		return now
	}
	td1 := []byte("testdata1\n")
	_, err = lr.Write(td1)
	assert.Nil(t, err)
	err = lr.Flush()
	assert.Nil(t, err)

	lr.nowFunc = func() time.Time {
		return now.Add(24 * time.Hour)
	}
	td2 := []byte("testdata2\n")
	_, err = lr.Write(td2)
	assert.Nil(t, err)
	err = lr.Flush()
	assert.Nil(t, err)

	ffp, err := lr.getFormattedFilepath(now)
	assert.Nil(t, err)

	data, err := ioutil.ReadFile(ffp)
	assert.Nil(t, err)
	result1 := append(header, td1...)
	assert.Equal(t, data, result1)

	err = os.Remove(ffp)
	assert.Nil(t, err)

	ffp, err = lr.getFormattedFilepath(now.Add(24 * time.Hour))
	assert.Nil(t, err)

	data, err = ioutil.ReadFile(ffp)
	assert.Nil(t, err)
	result2 := append(header, td2...)
	assert.Equal(t, data, result2)

	err = os.Remove(ffp)
	assert.Nil(t, err)
}
