package base64

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

func Decode(data []byte) (mail, secret string, err error) {
	i := bytes.IndexByte(data, separator)
	if i < 0 {
		return "", "", fmt.Errorf("wrong separator index %d", i)
	}
	//
	buf := data[:i]
	n, err := base64.StdEncoding.Decode(buf, buf)
	if err != nil {
		return "", "", fmt.Errorf("decode first part: %w", err)
	}

	mail = string(buf[:n])
	//
	buf = data[i+1:]
	n, err = base64.StdEncoding.Decode(buf, buf)
	if err != nil {
		return "", "", fmt.Errorf("decode second part: %w", err)
	}
	secret = string(buf[:n])
	//
	return mail, secret, nil
}
