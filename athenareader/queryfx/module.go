package queryfx

import (
	"database/sql"
	"github.com/uber/athenadriver/athenareader/configfx"
	drv "github.com/uber/athenadriver/go"
	"go.uber.org/fx"
	"io/ioutil"
	"os"
)

var Module = fx.Provide(new)

// Params defines the dependencies or inputs
type Params struct {
	fx.In

	Query     string
	DrvConfig *drv.Config
	OConfig   configfx.ReaderOutputConfig
}

// Result defines output
type Result struct {
	fx.Out

	QAD QueryAndDB
}

type QueryAndDB struct {
	DB    *sql.DB
	Query string
}

func new(p Params) (Result, error) {
	// 2. Open Connection.
	dsn := p.DrvConfig.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	var sqlString = p.Query
	if _, err := os.Stat(p.Query); err == nil {
		b, err := ioutil.ReadFile(p.Query)
		if err == nil {
			sqlString = string(b) // convert content to a 'string'
		}
	}
	qad := QueryAndDB{
		DB:    db,
		Query: sqlString,
	}
	return Result{
		QAD: qad,
	}, nil
}
