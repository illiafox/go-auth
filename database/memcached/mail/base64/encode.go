package base64

import (
	"bytes"
	"encoding/base64"
	"fmt"
)

const separator = byte(':')

func Encode(mail string, secret []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	encoder := base64.NewEncoder(base64.StdEncoding, buf)

	_, err := encoder.Write([]byte(mail))
	if err != nil {
		return nil, fmt.Errorf("encode mail: %w", err)
	}

	err = encoder.Close()
	if err != nil {
		return nil, fmt.Errorf("close encoder: %w", err)
	}

	err = buf.WriteByte(separator)
	if err != nil {
		return nil, fmt.Errorf("write separator: %w", err)
	}

	encoder = base64.NewEncoder(base64.StdEncoding, buf)

	_, err = encoder.Write(secret)
	if err != nil {
		return nil, fmt.Errorf("encode secret: %w", err)
	}

	err = encoder.Close()
	if err != nil {
		return nil, fmt.Errorf("close encoder: %w", err)
	}

	return buf.Bytes(), nil
}
