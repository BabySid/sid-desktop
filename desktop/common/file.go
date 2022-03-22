package common

import (
	"fyne.io/fyne/v2"
)

func ReadURI(closer fyne.URIReadCloser) ([]byte, error) {
	data := make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		switch nr, err := closer.Read(buf[:]); true {
		case nr < 0:
			return nil, err
		case nr == 0:
			return data, nil
		case nr > 0:
			data = append(data, buf[:nr]...)
		}
	}
}
