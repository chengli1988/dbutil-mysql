package dbutil

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// DeleteEngine struct
type DeleteEngine struct {
	db *DB
	Condition
}

// NewDeleteEngine 创建DeleteEngine
func NewDeleteEngine(model Model) *DeleteEngine {
	return NewDeleteEngineDB(model, new(DB))
}

// NewDeleteEngineDB 创建DeleteEngine(事务)
func NewDeleteEngineDB(model Model, db *DB) *DeleteEngine {

	deleteEngine := new(DeleteEngine)
	deleteEngine.db = db
	deleteEngine.ModelReflect = InitModelReflect(model)

	return deleteEngine
}

// Delete 删除
func (engine *DeleteEngine) Delete() (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)

	conditionSql, conditionValues := engine.getConditionSql()

	deleteSql := fmt.Sprintf("delete from %s %s", engine.tableName, conditionSql)
	if strings.IndexAny(deleteSql, "?") == -1 {
		return 0, errors.New("99001, 操作失败，缺少必要参数！")
	}

	if err != nil {
		return 0, err
	}

	result, err = execPrePareSql(engine.db, deleteSql, conditionValues...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}
