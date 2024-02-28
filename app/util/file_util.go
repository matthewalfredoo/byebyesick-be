package util

import (
	"fmt"
	"io/ioutil"
	"os"
)

func WriteTempFile(fileData []byte, fileExtension string) (*os.File, error) {
	tempFile, err := ioutil.TempFile("", fmt.Sprintf("tmp-*.%s", fileExtension))
	if err != nil {
		err = tempFile.Close()
		if err != nil {
			return nil, err
		}

		err = os.Remove(tempFile.Name())
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	_, err = tempFile.Write(fileData)
	if err != nil {
		err = tempFile.Close()
		if err != nil {
			return nil, err
		}

		err = os.Remove(tempFile.Name())
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	return tempFile, nil
}
