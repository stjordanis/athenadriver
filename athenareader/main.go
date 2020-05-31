// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/uber/athenadriver/athenareader/configfx"
	"github.com/uber/athenadriver/athenareader/queryfx"
	secret "github.com/uber/athenadriver/examples/constants"
	"go.uber.org/fx"
	"os"
)

var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
var bucket = flag.String("b", secret.OutputBucket, "Athena resultset output bucket")
var database = flag.String("d", "default", "The database you want to query")
var query = flag.String("q", "select 1", "The SQL query string or a file containing SQL string")
var rowOnly = flag.Bool("r", false, "Display rows only, don't show the first row as columninfo")
var moneyWise = flag.Bool("m", false, "Enable moneywise mode to display the query cost as the first line of the output")
var versionFlag = flag.Bool("v", false, "Print the current version and exit")
var admin = flag.Bool("a", false, "Enable admin mode, so database write(create/drop) is allowed at athenadriver level")
func printVersion() {
	println("Current build version: v1.1.6")
}

// main will query Athena and print all columns and rows information in csv format
func main() {
	flag.Usage = func() {
		preBody := "NAME\n\tathenareader - read athena data from command line\n\n"
		desc := "\nEXAMPLES\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\"\n" +
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-03T00:00:00.516940Z,elb_demo_004\n" +
			"\t2015-01-03T00:00:00.902953Z,elb_demo_004\n\n" +
			"\t$ athenareader -d sampledb -q \"select request_timestamp,elb_name from elb_logs limit 2\" -r\n" +
			"\t2015-01-05T20:00:01.206255Z,elb_demo_002\n" +
			"\t2015-01-05T20:00:01.612598Z,elb_demo_008\n\n" +
			"\t$ athenareader -d sampledb -b s3://my-athena-query-result -q tools/query.sql\n" +
			"\trequest_timestamp,elb_name\n" +
			"\t2015-01-06T00:00:00.516940Z,elb_demo_009\n\n" +
			"\n\tAdd '-m' to enable moneywise mode. The first line will display query cost under moneywise mode.\n\n" +
			"\t$ athenareader -b s3://athena-query-result -q 'select count(*) as cnt from sampledb.elb_logs' -m\n" +
			"\tquery cost: 0.00184898369752772851 USD\n" +
			"\tcnt\n" +
			"\t1356206\n\n" +
			"\n\tAdd '-a' to enable admin mode. Database write is enabled at driver level under admin mode.\n\n" +
			"\t$ athenareader -b s3://athena-query-result -q 'DROP TABLE IF EXISTS depreacted_table' -a\n" +
			"\t\n" +
			"AUTHOR\n\tHenry Fuheng Wu (henry.wu@uber.com)\n\n" +
			"REPORTING BUGS\n\thttps://github.com/uber/athenadriver\n"
		fmt.Fprintf(commandLine.Output(), preBody)
		fmt.Fprintf(commandLine.Output(),
			"SYNOPSIS\n\t%s [-v] [-b output_bucket] [-d database_name] [-q query_string_or_file] [-r] [-a] [-m]\n\nDESCRIPTION\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(commandLine.Output(), desc)
	}

	flag.Parse()
	switch {
	case *versionFlag:
		printVersion()
		return
	}
	// 1. Set AWS Credential in Driver Config.
	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	println(bucket, database, query, rowOnly, moneyWise, versionFlag, admin)
	app := fx.New(opts())
	ctx := context.Background()
	app.Start(ctx)
	defer app.Stop(ctx)
}

func opts() fx.Option {
	return fx.Options(
		fx.Provide(func() string { return *query }),
		configfx.Module,
		queryfx.Module,
	)
}
