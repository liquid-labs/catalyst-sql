package catsqlutils

import (
  "file"

  "github.com/go-pg/pg"
)

// InitLogger is a simple query-hook that logs the DB initialization queries to
// for inclusion with the source code as a reference.
type InitLogger struct { }

func (d InitLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d InitLogger) AfterQuery(q *pg.QueryEvent) {
  // TODO: for test and prod, would be nice to compare what gets generated here
  // with what's in the local git clone for integrity.
  if env.IsDev() {
    if file == nil {
      initialize it
    }
  	file.Println(q.FormattedQuery())
  }
}
