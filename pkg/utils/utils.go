package utils

import (
	"bufio"
	"io"
)

type LineHandler func(string) error

func ScanLines(r io.Reader, handlers ...LineHandler) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for _, handler := range handlers {
			err := handler(line)
			if err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
