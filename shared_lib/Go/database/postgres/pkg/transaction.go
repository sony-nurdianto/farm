package pkg

import "database/sql"

type SQLTx interface {
	Rollback() error
	Commit() error
	Prepare(query string) (Stmt, error)
	Stmt(st Stmt) Stmt
}

type sqlTx struct {
	tx *sql.Tx
}

func NewSQLTx(tx *sql.Tx) sqlTx {
	return sqlTx{
		tx,
	}
}

func (s sqlTx) Prepare(query string) (Stmt, error) {
	stmt, err := s.tx.Prepare(query)
	if err != nil {
		return nil, err
	}

	return NewStmt(stmt), nil
}

func (s sqlTx) Rollback() error {
	return s.tx.Rollback()
}

func (s sqlTx) Commit() error {
	return s.tx.Commit()
}

func (s sqlTx) Stmt(st Stmt) Stmt {
	statement := st.ToSQLSTMT()
	out := s.tx.Stmt(statement)
	return NewStmt(out)
}
