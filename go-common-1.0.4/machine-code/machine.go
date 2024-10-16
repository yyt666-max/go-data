package machine_code

import "sync"

var (
	code string
	wg   = sync.WaitGroup{}
)

const (
	salt = "eolink-apinto-business"
)

func GetMachineCode() string {
	wg.Wait()
	return code
}
