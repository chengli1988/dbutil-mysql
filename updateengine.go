package dbutil

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

// UpdateEngine 更新Engine
type UpdateEngine struct {
	db *DB
	Condition
}

// BatchUpdateEngine struct
type BatchUpdateEngine struct {
	db            *DB
	modelReflects []ModelReflect
}

// NewUpdateEngine 创建InsertEngine
func NewUpdateEngine(model Model) *UpdateEngine {
	return NewUpdateEngineDB(model, new(DB))
}

// NewUpdateEngineDB 创建UpdateEngine(事务)
func NewUpdateEngineDB(model Model, db *DB) *UpdateEngine {

	updateEngine := new(UpdateEngine)
	updateEngine.db = db
	updateEngine.ModelReflect = InitModelReflect(model)

	return updateEngine
}

// NewBatchUpdateEngine 创建BatchUpdateEngine
func NewBatchUpdateEngine(models ...Model) *BatchUpdateEngine {
	return NewBatchUpdateEngineDB(models, new(DB))
}

// NewBatchUpdateEngineDB 创建BatchUpdateEngine(事务)
func NewBatchUpdateEngineDB(models []Model, db *DB) *BatchUpdateEngine {
	batchUpdateEngine := new(BatchUpdateEngine)
	batchUpdateEngine.db = db

	for _, model := range models {
		batchUpdateEngine.modelReflects = append(batchUpdateEngine.modelReflects, InitModelReflect(model))
	}

	return batchUpdateEngine
}

// Update 更新
func (engine *UpdateEngine) Update() (int, error) {
	var (
		result       sql.Result
		rowsAffected int64
		err          error
	)

	updateFieldsSql, updateValues := engine.getUpdateFieldsSql()
	conditionSql, conditionValues := engine.getConditionSql()

	if strings.IndexAny(conditionSql, "?") == -1 {
		return 0, errors.New("操作失败，缺少必要参数！")
	}

	updateSql := fmt.Sprintf("update %s set %s %s", engine.tableName, updateFieldsSql, conditionSql)

	values := append(updateValues, conditionValues...)

	result, err = execPrePareSql(engine.db, updateSql, values...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil
}

// 获取update字段的sql
func (engine *UpdateEngine) getUpdateFieldsSql() (string, []interface{}) {
	values := make([]interface{}, 0, 1)

	var updateFieldsBuffer bytes.Buffer
	for _, dbField := range engine.dbFields {

		updateFieldsBuffer.WriteString(dbField)
		updateFieldsBuffer.WriteString(" = ?, ")

		values = append(values, engine.getFieldValue(engine.dbFieldMap[dbField]))
	}

	updateFieldsSql := updateFieldsBuffer.String()
	updateFieldsSql = string([]rune(updateFieldsSql)[0 : len(updateFieldsSql)-2])

	return updateFieldsSql, values
}

// ReplaceIntoMany 批量替换新增
// 根据主键或者唯一索引，存在则删除后新增，不存在则新增
func (engine *BatchUpdateEngine) ReplaceIntoMany() error {
	var (
		replaceDbFields []string
		valueFieldsSqls string
		err             error
	)

	var replaceFieldsBuffer bytes.Buffer
	dbFields := engine.modelReflects[0].dbFields
	modelReflect := engine.modelReflects[0]

	for _, dbField := range dbFields {
		if modelReflect.isNotZeroValue(modelReflect.dbFieldMap[dbField]) {
			replaceDbFields = append(replaceDbFields, dbField)
			replaceFieldsBuffer.WriteString(dbField)
			replaceFieldsBuffer.WriteString(", ")
		}
	}

	replaceFieldsSql := replaceFieldsBuffer.String()
	replaceFieldsSql = string([]rune(replaceFieldsSql)[0 : len(replaceFieldsSql)-2])

	valuesArray := make([]interface{}, 0, 1)
	for _, modelReflect := range engine.modelReflects {
		var valueFieldsBuffer bytes.Buffer
		for _, dbField := range replaceDbFields {

			valueFieldsBuffer.WriteString("?")
			valueFieldsBuffer.WriteString(", ")

			fieldValue := modelReflect.getFieldValue(modelReflect.dbFieldMap[dbField])
			valuesArray = append(valuesArray, fieldValue)
		}

		valueFieldsSql := valueFieldsBuffer.String()
		valueFieldsSql = string([]rune(valueFieldsSql)[0 : len(valueFieldsSql)-2])

		valueFieldsSqls = fmt.Sprintf("%s (%s),", valueFieldsSqls, valueFieldsSql)
	}

	valueFieldsSqls = string([]rune(valueFieldsSqls)[0 : len(valueFieldsSqls)-1])

	replaceSql := fmt.Sprintf("replace into %s (%s) values %s", modelReflect.tableName, replaceFieldsSql, valueFieldsSqls)

	_, err = execPrePareSql(engine.db, replaceSql, valuesArray...)

	return err
}
