package base64

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func FuzzBase64(f *testing.F) {
	f.Add("mail", "password")
	f.Fuzz(func(t *testing.T, mail string, secret string) {
		r := require.New(t)

		data, err := Encode(mail, []byte(secret))
		r.NoError(err, "encode")

		m, s, err := Decode(data)
		r.NoError(err, "decode")

		r.Equal(mail, m, "mail")
		r.Equal(secret, s, "secret")
	})
}
