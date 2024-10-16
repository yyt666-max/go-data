package store_postgres

import (
	"context"
	slog "log"
	"os"
	"time"

	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/cftool"
	"github.com/eolinker/go-common/store"
	"gorm.io/driver/postgres" // 更改为postgres驱动程序
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	_ store.IDB = (*storeDB)(nil)
)

type storeDB struct {
	db *gorm.DB
}
type postgresInit struct {
	config *DBConfig `autowired:""`
}

var _ store.IDB = (*storeDB)(nil)

func init() {
	cftool.Register[DBConfig]("postgres") // 修改为postgres
	autowire.Autowired(new(postgresInit))

}

func (m *storeDB) DB(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return m.db.WithContext(context.Background())
	}
	if tx, ok := ctx.Value(store.TxContextKey).(*gorm.DB); ok {
		return tx
	}
	return m.db.WithContext(ctx)
}
func (m *storeDB) IsTxCtx(ctx context.Context) bool {
	if _, ok := ctx.Value(store.TxContextKey).(*gorm.DB); ok {
		return ok
	}
	return false
}

func (m *postgresInit) OnComplete() {
	m.InitDb()
}
func (m *postgresInit) InitDb() {
	dialector := postgres.Open(m.config.getDBNS()) // 修改为postgres的连接
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(slog.New(os.Stderr, "\r\n", slog.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	})
	if err != nil {
		slog.Fatal(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		slog.Fatal(err)
	}
	sqlDb.SetConnMaxLifetime(time.Second * 9)
	sqlDb.SetMaxOpenConns(200)
	sqlDb.SetMaxIdleConns(200)

	autowire.Autowired[store.IDB](&storeDB{db: db})

}
