package queryfx

import (
	"database/sql"
	"fmt"
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

	query     string
	DrvConfig *drv.Config
	OConfig   configfx.ReaderOutputConfig
}

// Result defines output
type Result struct {
	fx.Out
}

func new(p Params) {
	// 2. Open Connection.
	dsn := p.DrvConfig.Stringify()
	db, _ := sql.Open(drv.DriverName, dsn)
	// 3. Query and print results
	var sqlString = p.query
	if _, err := os.Stat(p.query); err == nil {
		b, err := ioutil.ReadFile(p.query)
		if err != nil {
			fmt.Print(err)
		}
		sqlString = string(b) // convert content to a 'string'
	}
	rows, err := db.Query(sqlString)
	if err != nil {
		println(err.Error())
		return
	}
	defer rows.Close()
	if p.OConfig.Rowonly {
		drv.PrettyPrintSQLRows(rows, p.OConfig.Style)
		return
	}
	drv.PrettyPrintSQLColsRows(rows, p.OConfig.Style)
}
