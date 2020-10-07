package dbutil

import (
	"bytes"
	"fmt"
	"log"
)

// InsertEngine struct
type InsertEngine struct {
	db           *DB
	ignoreFields []string // 忽略字段
	ModelReflect
}

// BatchInsertEngine struct
type BatchInsertEngine struct {
	db            *DB
	modelReflects []ModelReflect
}

// NewInsertEngine 创建InsertEngine
func NewInsertEngine(model Model) *InsertEngine {
	return NewInsertEngineDB(model, new(DB))
}

// NewInsertEngineDB 创建InsertEngine(事务)
func NewInsertEngineDB(model Model, db *DB) *InsertEngine {

	insertEngine := new(InsertEngine)
	insertEngine.db = db
	insertEngine.ModelReflect = InitModelReflect(model)

	return insertEngine
}

// NewBatchInsertEngine 创建批量新增代理
func NewBatchInsertEngine(models ...Model) *BatchInsertEngine {
	return NewBatchInsertEngineDB(models, new(DB))
}

// NewBatchInsertEngineDB 创建批量新增代理(事务)
func NewBatchInsertEngineDB(models []Model, db *DB) *BatchInsertEngine {
	batchInsertEngine := new(BatchInsertEngine)
	batchInsertEngine.db = db

	for _, model := range models {
		batchInsertEngine.modelReflects = append(batchInsertEngine.modelReflects, InitModelReflect(model))
	}

	return batchInsertEngine
}

// Insert 新增
func (insertEngine *InsertEngine) Insert() error {
	var (
		err error
	)

	insertFieldsSql, valueFieldsSql, values := insertEngine.getInsertFieldsSql()
	insertSql := fmt.Sprintf("insert into %s (%s) values (%s)", insertEngine.tableName, insertFieldsSql, valueFieldsSql)

	_, err = execPrePareSql(insertEngine.db, insertSql, values...)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// InsertMany 批量新增
func (engine *BatchInsertEngine) InsertMany() error {
	var (
		insertDbFields  []string
		valueFieldsSqls string
		err             error
	)

	var insertFieldsBuffer bytes.Buffer

	dbFields := engine.modelReflects[0].dbFields
	modelReflect := engine.modelReflects[0]

	for _, dbField := range dbFields {

		if modelReflect.isNotZeroValue(modelReflect.dbFieldMap[dbField]) {

			insertDbFields = append(insertDbFields, dbField)
			insertFieldsBuffer.WriteString(dbField)
			insertFieldsBuffer.WriteString(", ")
		}
	}

	insertFieldsSql := insertFieldsBuffer.String()
	insertFieldsSql = string([]rune(insertFieldsSql)[0 : len(insertFieldsSql)-2])

	valuesArray := make([]interface{}, 0, 1)
	for _, modelReflect := range engine.modelReflects {
		var valueFieldsBuffer bytes.Buffer
		for _, dbField := range insertDbFields {

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

	insertSql := fmt.Sprintf("insert into %s (%s) values %s", modelReflect.tableName, insertFieldsSql, valueFieldsSqls)

	_, err = execPrePareSql(engine.db, insertSql, valuesArray...)

	return err
}

// 获取insert字段的sql语句
func (insertEngine *InsertEngine) getInsertFieldsSql() (string, string, []interface{}) {
	dbFields := insertEngine.dbFields

	for _, ignoreField := range insertEngine.ignoreFields {
		ignoreDbField := insertEngine.modelFieldMap[ignoreField]
		for index, dbField := range dbFields {
			if dbField == ignoreDbField {
				dbFields = append(dbFields[:index], dbFields[index+1:]...)
			}
		}
	}

	var insertFieldsBuffer bytes.Buffer
	var valueFieldsBuffer bytes.Buffer
	values := make([]interface{}, 0, 1)
	for _, dbField := range dbFields {

		jsonField := insertEngine.dbFieldMap[dbField]
		if insertEngine.isNotZeroValue(jsonField) {

			insertFieldsBuffer.WriteString(dbField)
			insertFieldsBuffer.WriteString(", ")

			valueFieldsBuffer.WriteString("?")
			valueFieldsBuffer.WriteString(",")

			values = append(values, insertEngine.getFieldValue(jsonField))
		}
	}

	insertFieldsSql := insertFieldsBuffer.String()
	insertFieldsSql = string([]rune(insertFieldsSql)[0 : len(insertFieldsSql)-2])

	valueFieldsSql := valueFieldsBuffer.String()
	valueFieldsSql = string([]rune(valueFieldsSql)[0 : len(valueFieldsSql)-1])

	return insertFieldsSql, valueFieldsSql, values
}
