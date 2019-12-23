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
	sc := bufio.NewScanner(r)
	m := Migration{Name: name}

	buf := &bytes.Buffer{}
	dir := dirNone

	for sc.Scan() {
		line := sc.Text()

		if strings.HasPrefix(line, sqlPrefix) {
			d := strings.TrimSpace(strings.TrimPrefix(line, sqlPrefix))
			buf.Reset()

			switch d {
			case "Up":
				dir = dirUp
			case "Down":
				dir = dirDown
			}

			continue
		}

		if _, err := buf.WriteString(line + "\n"); err != nil {
			return m, err
		}

		if strings.HasSuffix(strings.TrimSpace(line), sqlDelimiter) {
			switch dir {
			case dirUp:
				m.Up = append(m.Up, buf.String())
			case dirDown:
				m.Down = append(m.Down, buf.String())
			}
			buf.Reset()
		}
	}

	if err := sc.Err(); err != nil {
		return m, err
	}

	return m, nil
}
