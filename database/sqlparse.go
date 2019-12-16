package database

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	sqlPrefix    = "-- +migration "
	sqlDelimiter = ";"
)

type direction int

const (
	dirUp direction = iota
	dirDown
	dirNone
)

func parseMigration(name string, r io.Reader) (Migration, error) {
	scanner := bufio.NewScanner(r)
	parsed := Migration{Name: name}

	buf := &bytes.Buffer{}
	currentDir := dirNone

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, sqlPrefix) {
			sqlDirection := strings.TrimSpace(strings.TrimPrefix(line, sqlPrefix))
			buf.Reset()

			switch sqlDirection {
			case "Up":
				currentDir = dirUp
			case "Down":
				currentDir = dirDown
			}

			continue
		}

		if _, err := buf.WriteString(line + "\n"); err != nil {
			return parsed, err
		}

		if strings.HasSuffix(strings.TrimSpace(line), sqlDelimiter) {
			switch currentDir {
			case dirUp:
				parsed.Up = append(parsed.Up, buf.String())
			case dirDown:
				parsed.Down = append(parsed.Down, buf.String())
			}
			buf.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return parsed, err
	}

	return parsed, nil
}
