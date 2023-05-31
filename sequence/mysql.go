package sequence

import (
	"database/sql"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 建立MySQL连接 执行REPLACE INTO语句
// REPLACE INTO sequence (name) VALUES ('short_url')
// SELECT LAST_INSERT_ID()

const sqlReplaceStub = `REPLACE INTO sequence (stub) VALUES ('a')`

type MySQL struct {
	conn sqlx.SqlConn
}

func NewMySQL(dsn string) Sequence {
	return &MySQL{
		conn: sqlx.NewMysql(dsn),
	}
}

func (m *MySQL) Next() (seq uint64, err error) {
	// 预编译
	var stmt sqlx.StmtSession
	stmt, err = m.conn.Prepare(sqlReplaceStub)
	if err != nil {
		logx.Errorw("conn.Prepare failed", logx.Field("err", err.Error()))
		return
	}
	defer stmt.Close()

	// 执行
	var rest sql.Result
	rest, err = stmt.Exec()
	if err != nil {
		logx.Errorw("stmt.Exec failed", logx.Field("err", err.Error()))
		return
	}

	// 获取主键id
	var lid int64
	lid, err = rest.LastInsertId()
	if err != nil {
		logx.Errorw("rest.LastInsertId failed", logx.Field("err", err.Error()))
		return
	}
	return uint64(lid), nil
}
