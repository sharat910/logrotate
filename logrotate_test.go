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
	temp := "./temp"
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
	assert.Equal(t, result1, data)

	ffp, err = lr.getFormattedFilepath(now.Add(24 * time.Hour))
	assert.Nil(t, err)

	data, err = ioutil.ReadFile(ffp)
	assert.Nil(t, err)
	result2 := append(header, td2...)
	assert.Equal(t, result2, data)

	err = os.RemoveAll(temp)
	assert.Nil(t, err)
}

func TestLogRotator_getFormattedFilepath(t *testing.T) {
	testTime := time.Date(2020, 10, 9, 19, 30, 0, 0, time.Now().Location())
	temp := os.TempDir()
	fp := filepath.Join(temp, "test.log")
	tests := []struct {
		prependTimeFmt string
		want    string
		wantErr bool
	}{
		{
			"2006-01-02",
			filepath.Join(temp, "2020-10-09_test.log"),
			false,
		},
		{
			"2006-01-02-15",
			filepath.Join(temp, "2020-10-09-19_test.log"),
			false,
		},
		{
			time.RFC3339,
			filepath.Join(temp, testTime.Format(time.RFC3339) + "_test.log"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.prependTimeFmt, func(t *testing.T) {
			lr, err := New(fp, PrependTimeFormat(tt.prependTimeFmt, "_"))
			if err != nil {
				t.Fatalf("unable to create new logrotator: error = %v", err)
			}
			got, err := lr.getFormattedFilepath(testTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFormattedFilepath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getFormattedFilepath() got = %v, want %v", got, tt.want)
			}
		})
	}
}