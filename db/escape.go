package db

import (
	"bufio"
	"bytes"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

// escape replaces ':' with '\c' and '\' with '\\'.
func escape(value []byte) []byte {
	escaped := []byte{}
	for _, b := range value {
		switch b {
		case ':':
			escaped = append(escaped, ([]byte{'\\', 'c'})...)
		case '\\':
			escaped = append(escaped, ([]byte{'\\', b})...)
		default:
			escaped = append(escaped, b)
		}
	}
	return escaped
}

// unescape is the inverse of escape.
func unescape(value []byte) ([]byte, error) {
	reader := bufio.NewReader(bytes.NewBuffer(value))
	unescaped := []byte{}
	for {
		b, err := reader.ReadByte()
		if err != nil {
			// Assume io.EOF error indicating we reached the end of the value.
			break
		}
		if b == '\\' {
			next, err := reader.ReadByte()
			if err != nil {
				// This is only possible if the value was not escaped properly. Should
				// never happen.
				log.WithFields(log.Fields{
					"error": err.Error(),
					"value": hex.Dump(value),
				}).Error("unexpected error in unescape")
				return nil, err
			}
			if next == 'c' {
				unescaped = append(unescaped, ':')
			} else {
				unescaped = append(unescaped, next)
			}
		} else {
			unescaped = append(unescaped, b)
		}
	}
	return unescaped, nil
}
