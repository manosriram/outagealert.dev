package monitor

import (
	"encoding/base64"
	"fmt"
	"time"
)

func FormatTimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration.Seconds() < 60:
		return "just now"
	case duration.Minutes() < 60:
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes()))
	case duration.Hours() < 24:
		return fmt.Sprintf("%d hours ago", int(duration.Hours()))
	case duration.Hours() < 48:
		return "yesterday"
	default:
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24))
	}
}

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

func makeSlackState(projectId, monitorId string) string {
	return fmt.Sprintf("%s;%s", projectId, monitorId)
}
