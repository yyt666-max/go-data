package log_init

import "io"

type writes []io.Writer

func ToCopyToIoWriter(ws ...io.Writer) io.Writer {
	return writes(ws)
}
func (ws writes) Write(p []byte) (n int, err error) {
	for _, w := range ws {
		n, err = w.Write(p)
	}
	return
}
