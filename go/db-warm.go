package catsql

type dbSetupFunc func(db *sql.DB)

var setupFuncs = make([]dbSetupFunc, 0, 8)

// RegisterPrepareStmts generates the typical CRUD and other resource specific
// prepared statements. The 'SQL' here refers the the SQL query string.
func RegisterPrepareStmts(stmtsSQLMap map(**pg.Stmt)string) {
  // TODO: from latest API... content?
}

func Warm() {
  for _, setupFunc := range setupFuncs {
    setupFunc(db)
  }
}
