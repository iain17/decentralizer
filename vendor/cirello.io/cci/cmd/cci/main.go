// Copyright 2018 github.com/ucirello
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command cci implements a simple and dirty CI service.
package main // import "cirello.io/cci/cmd/cci"

import (
	"log"
	"os"

	"cirello.io/cci/pkg/ui/cli"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	log.SetPrefix("cci: ")
	log.SetFlags(0)
	fn := "cci.db"
	if envFn := os.Getenv("CCI_DB"); envFn != "" {
		fn = envFn
	}
	db := openDB(fn)
	cli.Run(db)
}

func openDB(fn string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", fn)
	if err != nil {
		log.Fatalln("cannot open database:", err)
	}
	return db
}
