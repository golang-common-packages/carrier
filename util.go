package carrier

import (
	"bytes"
	"encoding/json"
	"hash/fnv"
	"io"

	"github.com/golang-common-packages/hash"
)

func generateKey(data string) string {
	hash := fnv.New64a()
	hash.Write([]byte(data))

	return string(hash.Sum64())
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

func streamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.String()
}

func hashObject(object interface{}) string {
	hasher := &hash.Client{}
	configAsJSON, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}
	return hasher.SHA1(string(configAsJSON))
}

func bytesToString(data []byte) string {
	return string(data[:])
}
