package integration

import (
	"encoding/base64"
	"fmt"
	"os"
)

func Base64Encode(src string) []byte {
	data := []byte(src)
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
	base64.StdEncoding.Encode(dst, data)
	return dst
}

func Base64Decode(src string) (string, error) {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, []byte(src))
	if err != nil {
		fmt.Println("decode error:", err)
		return "", err
	}
	return string(dst[:n]), nil
}

func makeSlackRedirectUrl(projectId, monitorId string) string {
	return fmt.Sprintf("%s/monitor/%s/%s", os.Getenv("HOST_WITH_SCHEME"), projectId, monitorId)
}
