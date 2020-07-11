package dbutil

import (
	"bytes"
	"fmt"
)

// SelectEngine 查询代理
type SelectEngine struct {
	db           *DB
	fields       []string // 返回结果字段
	ignoreFields []string // 忽略字段，fields为空，则移除相对于实体的字段；如果fields不为空，则移除fields中对于的字段
	distinct     bool     // 去重
	Condition
}

// NewSelectEngine 创建查询代理
func NewSelectEngine(model Model) *SelectEngine {
	return NewSelectEngineDB(model, new(DB))
}

// NewSelectEngineDB 事务
func NewSelectEngineDB(model Model, db *DB) *SelectEngine {

	selectEngine := new(SelectEngine)
	selectEngine.db = db
	selectEngine.distinct = false
	selectEngine.Condition = InitCondition(model)

	return selectEngine
}

// SelectAll 根据条件查询所有记录
func (selectEngine *SelectEngine) SelectAll() ([]map[string]interface{}, error) {
	var (
		rows []map[string]interface{}
		err  error
	)

	selectFieldsSql := selectEngine.getSelectFieldsSql()
	conditionSql, conditionValues := selectEngine.getConditionSql()
	orderSql := selectEngine.getOrderSql()

	selectSql := fmt.Sprintf("select %s from %s %s %s", selectFieldsSql, selectEngine.tableName, conditionSql, orderSql)

	rows, err = selectEngine.db.SelectBySql(selectSql, conditionValues...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// SelectPage 根据条件分页查询
func (selectEngine *SelectEngine) SelectPage() ([]map[string]interface{}, int, error) {
	var (
		count int
		rows  []map[string]interface{}
		err   error
	)

	selectFieldsSql := selectEngine.getSelectFieldsSql()
	conditionSql, conditionValues := selectEngine.getConditionSql()

	orderSql := selectEngine.getOrderSql()
	limitSql := selectEngine.getLimitSql()

	countSql := fmt.Sprintf("select count(1) as count from %s %s", selectEngine.tableName, conditionSql)
	selectSql := fmt.Sprintf("select %s from %s %s %s %s", selectFieldsSql, selectEngine.tableName, conditionSql, orderSql, limitSql)

	_, err = selectEngine.db.DoTransaction(func() error {
		count, err = selectEngine.db.SelectCountBySql(countSql, conditionValues...)
		if err != nil {
			return err
		}

		rows, err = selectEngine.db.SelectBySql(selectSql, conditionValues...)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return rows, count, nil
}

// SelectOne 查询一行记录
func (selectEngine *SelectEngine) SelectOne() (map[string]interface{}, error) {
	var (
		row map[string]interface{}
		err error
	)

	selectFieldsSql := selectEngine.getSelectFieldsSql()
	conditionSql, conditionValues := selectEngine.getConditionSql()

	selectSql := fmt.Sprintf("select %s from %s %s", selectFieldsSql, selectEngine.tableName, conditionSql)

	rows, err := selectEngine.db.SelectBySql(selectSql, conditionValues...)
	if err != nil {
		return nil, err
	}

	if len(rows) > 0 {
		row = rows[0]
	}

	return row, err
}

// 获取select字段的sql语句
func (selectEngine *SelectEngine) getSelectFieldsSql() string {
	fieldLength := len(selectEngine.fields)
	dbFields := make([]string, fieldLength, 1)

	if fieldLength > 0 {
		for _, field := range selectEngine.fields {
			dbFields = append(dbFields, selectEngine.modelFieldMap[field])
		}
	} else {
		dbFields = selectEngine.dbFields
	}

	for _, ignoreField := range selectEngine.ignoreFields {
		ignoreDbField := selectEngine.modelFieldMap[ignoreField]
		for index, dbField := range dbFields {
			if dbField == ignoreDbField {
				dbFields = append(dbFields[:index], dbFields[index+1:]...)
			}
		}
	}

	var selectFieldsBuffer bytes.Buffer
	if selectEngine.distinct {
		selectFieldsBuffer.WriteString("distinct ")
	}

	for _, dbField := range dbFields {
		selectFieldsBuffer.WriteString(dbField)
		selectFieldsBuffer.WriteString(" as ")
		selectFieldsBuffer.WriteString(selectEngine.dbFieldMap[dbField])
		selectFieldsBuffer.WriteString(", ")
	}

	selectFieldsSql := selectFieldsBuffer.String()
	selectFieldsSql = string([]rune(selectFieldsSql)[0 : len(selectFieldsSql)-2])

	return selectFieldsSql
}

// OrderByDesc 倒序排序
func (selectEngine *SelectEngine) OrderByDesc(jsonFields ...string) *SelectEngine {
	for _, jsonField := range jsonFields {
		dbField := selectEngine.modelFieldMap[jsonField]
		if _, isExist := selectEngine.orderBys[dbField]; !isExist {
			selectEngine.orderKeys = append(selectEngine.orderKeys, dbField)
			selectEngine.orderBys[dbField] = "DESC"
		}
	}
	return selectEngine
}

// OrderByAsc 正序排序
func (selectEngine *SelectEngine) OrderByAsc(jsonFields ...string) *SelectEngine {
	for _, jsonField := range jsonFields {
		dbField := selectEngine.modelFieldMap[jsonField]
		if _, isExist := selectEngine.orderBys[dbField]; !isExist {
			selectEngine.orderKeys = append(selectEngine.orderKeys, dbField)
			selectEngine.orderBys[dbField] = "ASC"
		}
	}
	return selectEngine
}

// Limit 分页设置
func (selectEngine *SelectEngine) Limit(currentPage int, pageSize int) *SelectEngine {
	selectEngine.currentPage = currentPage
	selectEngine.pageSize = pageSize
	return selectEngine
}

// getOrderSql 获取order by的sql语句
func (selectEngine *SelectEngine) getOrderSql() string {
	var orderSqlBuffer bytes.Buffer

	if len(selectEngine.orderBys) > 0 {
		orderSqlBuffer.WriteString(" order by ")
		for key, value := range selectEngine.orderBys {
			orderSqlBuffer.WriteString(key)
			orderSqlBuffer.WriteString(" ")
			orderSqlBuffer.WriteString(value)
			orderSqlBuffer.WriteString(", ")
		}

		orderSql := orderSqlBuffer.String()
		orderSql = string([]rune(orderSql)[0 : len(orderSql)-2])
		return orderSql
	}

	return ""
}

// getLimitSql 获取分页的sql语句
func (selectEngine *SelectEngine) getLimitSql() string {
	return fmt.Sprintf(" limit %d, %d", (selectEngine.currentPage-1)*selectEngine.pageSize, selectEngine.pageSize)
}
