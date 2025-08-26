package pkg

import "database/sql"

type TxOpts = sql.TxOptions

const (
	LevelDefault        sql.IsolationLevel = sql.LevelDefault
	LevelReadUncommited sql.IsolationLevel = sql.LevelReadUncommitted
	LevelReadCommitted  sql.IsolationLevel = sql.LevelReadCommitted
	LevelWriteCommitted sql.IsolationLevel = sql.LevelWriteCommitted
	LevelRepeatableRead sql.IsolationLevel = sql.LevelRepeatableRead
	LevelSnapshot       sql.IsolationLevel = sql.LevelSnapshot
	LevelSerializable   sql.IsolationLevel = sql.LevelSerializable
	LevelLinearizable   sql.IsolationLevel = sql.LevelLinearizable
)
