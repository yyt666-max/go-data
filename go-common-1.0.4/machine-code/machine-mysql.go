package machine_code

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/cftool"
	"github.com/eolinker/go-common/store"
)

type dBConfig struct {
	UserName string `yaml:"user_name"`
	//Password string `yaml:"password"`
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Db   string `yaml:"db"`
}
type machineInit struct {
	*dBConfig `autowired:""`
	db        store.IDB `autowired:""`
}

func (m *machineInit) OnComplete() {

	sqlDb, _ := m.db.DB(context.Background()).DB()
	var variableName, serverUUID string
	row := sqlDb.QueryRow(`SHOW VARIABLES LIKE "server_uuid"`)
	err := row.Scan(&variableName, &serverUUID)
	if err != nil {
		panic(err)
	}
	machineCodeRaw := fmt.Sprintf("%s:%d.%s.%s.%s", m.Ip, m.Port, m.UserName, m.Db, serverUUID)
	code = Md5(fmt.Sprintf("%s%s", salt, machineCodeRaw))
	wg.Done()

}

func init() {
	wg.Add(1)
	cftool.Register[dBConfig]("mysql")
	autowire.Autowired(new(machineInit))

}
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
