package dbutil

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// Condition 条件结构体
type Condition struct {
	inSplit        string            // in数据分隔符
	orderBys       map[string]string // 排序字段
	orderKeys      []string          // 排序字段顺序
	currentPage    int
	pageSize       int
	conditionSql   bytes.Buffer  // 条件sql
	conditionValue []interface{} // 条件sql参数值
	ModelReflect
}

// InitCondition 初始化条件结构体
func InitCondition(model Model) Condition {
	var condition Condition

	condition.inSplit = ","
	condition.orderBys = make(map[string]string)
	condition.currentPage = 1
	condition.pageSize = 10
	condition.ModelReflect = InitModelReflect(model)

	return condition
}

// WhereEqs 等于条件
func (condition *Condition) WhereEqs(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {
		if condition.isNotZeroValue(jsonField) {

			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" = ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereNes 不等于条件
func (condition *Condition) WhereNes(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		if condition.isNotZeroValue(jsonField) {
			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" <> ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereLts 小于条件
func (condition *Condition) WhereLts(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		if condition.isNotZeroValue(jsonField) {
			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" < ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereLes 小于等于条件
func (condition *Condition) WhereLes(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {
		if condition.isNotZeroValue(jsonField) {

			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" <= ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereGts 大于条件
func (condition *Condition) WhereGts(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {
		if condition.isNotZeroValue(jsonField) {

			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" > ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereGes 大于等于条件
func (condition *Condition) WhereGes(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {
		if condition.isNotZeroValue(jsonField) {

			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" >= ? ")

			condition.conditionValue = append(condition.conditionValue, condition.getFieldValue(jsonField))
		}
	}

	return condition
}

// WhereIns 包含条件
func (condition *Condition) WhereIns(jsonFields ...string) *Condition {

	for _, jsonField := range jsonFields {
		if condition.isNotZeroValue(jsonField) {

			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" in (")

			value := condition.getFieldStringValue(jsonField)
			valueArray := strings.Split(value, condition.inSplit)
			for index, key := range valueArray {
				condition.conditionSql.WriteString(" ? ")
				if index < len(valueArray)-1 {
					condition.conditionSql.WriteString(",")
				}
				condition.conditionValue = append(condition.conditionValue, key)
			}
			condition.conditionSql.WriteString(" )")
		}
	}

	return condition
}

// WhereLikes 模糊匹配条件
func (condition *Condition) WhereLikes(jsonField ...string) *Condition {
	for _, jsonField := range jsonField {

		if condition.isNotZeroValue(jsonField) {
			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" like ? ")

			condition.conditionValue = append(condition.conditionValue, fmt.Sprintf("%%%s%%", condition.getFieldValue(jsonField)))
		}
	}

	return condition
}

// WhereLeftLikes 左模糊匹配条件
func (condition *Condition) WhereLeftLikes(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		if condition.isNotZeroValue(jsonField) {
			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" like ? ")

			condition.conditionValue = append(condition.conditionValue, fmt.Sprintf("%s%%", condition.getFieldValue(jsonField)))
		}
	}

	return condition
}

// WhereRightLikes 右模糊匹配条件
func (condition *Condition) WhereRightLikes(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		if condition.isNotZeroValue(jsonField) {
			condition.conditionSql.WriteString(" and ")
			condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
			condition.conditionSql.WriteString(" like ? ")

			condition.conditionValue = append(condition.conditionValue, fmt.Sprintf("%%%s", condition.getFieldValue(jsonField)))
		}
	}

	return condition
}

// WhereEqZero 等于0
func (condition *Condition) WhereEqZero(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		condition.conditionSql.WriteString(" and ")
		condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
		condition.conditionSql.WriteString(" = ? ")

		condition.conditionValue = append(condition.conditionValue, 0)
	}

	return condition
}

// WhereNeqZero 不等于0
func (condition *Condition) WhereNeqZero(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		condition.conditionSql.WriteString(" and ")
		condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
		condition.conditionSql.WriteString(" != ? ")

		condition.conditionValue = append(condition.conditionValue, 0)
	}

	return condition
}

// WhereIsNull 是null值
func (condition *Condition) WhereIsNull(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		condition.conditionSql.WriteString(" and ")
		condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
		condition.conditionSql.WriteString(" is null ")
	}

	return condition
}

// WhereIsNotNull 是null值
func (condition *Condition) WhereIsNotNull(jsonFields ...string) *Condition {
	for _, jsonField := range jsonFields {

		condition.conditionSql.WriteString(" and ")
		condition.conditionSql.WriteString(condition.modelFieldMap[jsonField])
		condition.conditionSql.WriteString(" is not null ")
	}

	return condition
}

// 获取where条件sql语句
func (condition *Condition) getConditionSql() (string, []interface{}) {

	conditionSql := condition.conditionSql.String()
	if len(conditionSql) > 0 {
		conditionSql = strings.Replace(conditionSql, "and", "where", 1)
	}

	return conditionSql, condition.conditionValue
}

// 预处理sql语句
func execPrePareSql(db *DB, prepareSql string, args ...interface{}) (sql.Result, error) {
	var (
		stmt   *sql.Stmt
		result sql.Result
		err    error
	)

	log.Println("执行SQL：", prepareSql)
	log.Println("参数：", args)

	stmt, err = dbPool.Prepare(prepareSql)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	if db.isTransaction {
		result, err = db.tx.Stmt(stmt).Exec(args...)
	} else {
		result, err = stmt.Exec(args...)
	}

	return result, err
}
