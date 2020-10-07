package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	// mysql 驱动
	_ "github.com/go-sql-driver/mysql"
)

var dbPool *sql.DB

// InitPool 初始化数据库链接池
// username 用户名 [mysql]
// password 密码 [mysql]
// ip IP地址 [mysql]
// port 端口 [mysql]
// dbname 数据库名 [mysql]
// charset 编码 [mysql]
func InitPool(username string, password string, ip string, port int, dbname string, charset string) {
	var (
		dataSourceName string
		err            error
	)

	dataSourceName = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=%s", username, password, ip, port, dbname, charset)

	dbPool, err = sql.Open("mysql", dataSourceName)

	if err != nil {
		log.Println("数据库连接失败！", err.Error())
		os.Exit(1)
	}

	dbPool.SetMaxOpenConns(50)
	dbPool.SetMaxIdleConns(50)
	dbPool.SetConnMaxLifetime(3600)

	err = dbPool.Ping()
	if err != nil {
		log.Println("数据库连接失败！", err.Error())
		os.Exit(1)
	}

	log.Println("数据库连接成功！")
}

// DB 数据库链接
type DB struct {
	// 开启事务 bool 零值为false
	isTransaction bool
	// 事务
	tx *sql.Tx
}

// SelectCountBySql 查询记录数
func (db *DB) SelectCountBySql(countSql string, params ...interface{}) (int, error) {
	var (
		count int
		err   error
	)

	log.Println("执行SQL：" + countSql)
	if db.isTransaction {
		err = db.tx.QueryRow(countSql, params...).Scan(&count)
	} else {
		err = dbPool.QueryRow(countSql, params...).Scan(&count)
	}

	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectBySql 查询记录
func (db *DB) SelectBySql(selectSql string, params ...interface{}) ([]map[string]interface{}, error) {
	var (
		rows *sql.Rows
		err  error
	)

	log.Println("执行SQL：", selectSql)
	log.Println("参数：", params)
	if db.isTransaction {
		rows, err = db.tx.Query(selectSql, params...)
	} else {
		rows, err = dbPool.Query(selectSql, params...)
	}

	if err != nil {
		return nil, err
	}

	return handleRows(rows)
}

// SelectOneBySql 查询一行记录
func (db *DB) SelectOneBySql(selectSql string, params ...interface{}) (map[string]interface{}, error) {
	var (
		rows []map[string]interface{}
		err  error
	)

	rows, err = db.SelectBySql(selectSql, params...)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		return rows[0], nil
	}

	return nil, nil
}

// DoTransaction 事务操作
func (db *DB) DoTransaction(dbOperates func() error) (bool, error) {
	tx, err := dbPool.Begin()

	if err != nil {
		log.Println("开启事务失败！")
		return false, err
	}

	db.tx = tx
	db.isTransaction = true

	err = dbOperates()
	if err != nil {
		tx.Rollback()
		db.isTransaction = false

		return false, err
	}

	tx.Commit()

	return true, nil
}

func handleRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	records := make([]map[string]interface{}, 0)

	defer rows.Close()
	for rows.Next() {
		record := make(map[string]interface{})

		err := rows.Scan(scanArgs...)

		if err != nil {
			return nil, err
		}

		for i, col := range values {
			if col != nil {

				switch col.(type) {
				case []byte:
					record[columns[i]] = string(col.([]byte))
				case time.Time:
					if col.(time.Time).IsZero() {
						record[columns[i]] = nil
					} else {
						record[columns[i]] = col.(time.Time).Format(FormatLayout)
					}
				default:
					record[columns[i]] = col.(int64)
				}
			}
		}

		records = append(records, record)
	}

	return records, nil
}
