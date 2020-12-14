package logrotator

import (
	"compress/gzip"
	"io"
	"os"
)

func Compress(file string, removeOld bool) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	destFile := file + ".gz"
	df, err := os.Create(destFile)
	if err != nil {
		return err
	}

	gw := gzip.NewWriter(df)

	_, err = io.Copy(gw, f)
	if err != nil {
		return err
	}

	if err = gw.Close(); err != nil {
		return err
	}

	if err = df.Close(); err != nil {
		return err
	}

	if err = f.Close(); err != nil {
		return err
	}

	if removeOld {
		if err = os.Remove(file); err != nil {
			return err
		}
	}

	return nil
}


