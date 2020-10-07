package dbutil

import (
	"database/sql/driver"
	"reflect"
	"time"
)

// Model 接口
type Model interface {
	// 返回数据库表名
	GetTableName() string
}

// LocalTime 自定义时间类型
type LocalTime time.Time

// FormatLayout 日期时间格式化格式
const FormatLayout = "2006-01-02 15:04:05"

// UnmarshalJSON unmarshal json方法
func (localTime *LocalTime) UnmarshalJSON(data []byte) (err error) {
	parsedTime, err := time.ParseInLocation(`"`+FormatLayout+`"`, string(data), time.Local)
	*localTime = LocalTime(parsedTime)
	return nil
}

// MarshalJSON marshal json方法
func (localTime LocalTime) MarshalJSON() ([]byte, error) {
	timeByte := make([]byte, 0, len(FormatLayout)+2)
	timeByte = append(timeByte, '"')
	timeByte = time.Time(localTime).AppendFormat(timeByte, FormatLayout)
	timeByte = append(timeByte, '"')
	return timeByte, nil
}

// Value 插入mysql使用
func (localTime LocalTime) Value() (driver.Value, error) {
	if time.Time(localTime).IsZero() {
		return nil, nil
	}

	return []byte(time.Time(localTime).Format(FormatLayout)), nil
}

// String 字符串方法
func (localTime LocalTime) String() string {
	if time.Time(localTime).IsZero() {
		return ""
	}
	return time.Time(localTime).Format(FormatLayout)
}

// ModelReflect 实体
type ModelReflect struct {
	tableName         string
	dbFields          []string          // 数据库字段
	dbFieldMap        map[string]string // db字段名Map-json字段名, key:db字段名, value: json字段名 (tag为dbField)
	dbFieldTypeMap    map[string]string //
	modelValue        reflect.Value     //
	modelFieldMap     map[string]string //
	modelFieldNameMap map[string]string //
}

// InitModelReflect 初始化ModelReflect
func InitModelReflect(model Model) ModelReflect {
	var modelReflect ModelReflect

	modelReflect.tableName = model.GetTableName()
	modelReflect.dbFieldMap = make(map[string]string)
	modelReflect.modelValue = reflect.ValueOf(model)
	modelReflect.modelFieldMap = make(map[string]string)
	modelReflect.modelFieldNameMap = make(map[string]string)
	modelReflect.dbFieldTypeMap = make(map[string]string)

	modelReflect.handleModelReflect(reflect.TypeOf(model))

	return modelReflect
}

func (modelReflect *ModelReflect) handleModelReflect(rt reflect.Type) {
	numField := rt.NumField()
	for i := 0; i < numField; i++ {
		if rt.Field(i).Anonymous {
			modelReflect.handleModelReflect(rt.Field(i).Type)
			continue
		}

		fieldTag := rt.Field(i).Tag
		if fieldTag.Get("db") != "" && fieldTag.Get("json") != "" {
			modelReflect.modelFieldMap[fieldTag.Get("json")] = fieldTag.Get("db")
			modelReflect.modelFieldNameMap[fieldTag.Get("json")] = rt.Field(i).Name
			modelReflect.dbFieldTypeMap[fieldTag.Get("json")] = fieldTag.Get("dbType")

			// 只有带有dbField且值且true的为数据库字段
			if fieldTag.Get("dbField") == "true" {
				modelReflect.dbFieldMap[fieldTag.Get("db")] = fieldTag.Get("json")
				modelReflect.dbFields = append(modelReflect.dbFields, fieldTag.Get("db"))
			}
		}
	}
}

// 根据json字段检查字段是否为零值（true：是零值；false：不是零值）
func (modelReflect *ModelReflect) isZeroValue(jsonField string) bool {
	fieldValue := modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField])
	if fieldValue.IsZero() {
		return true
	}
	return false
}

// isNotZeroValue 根据json字段检查字段值是否不为零值（true：不是零值；false：是零值）
func (modelReflect *ModelReflect) isNotZeroValue(jsonField string) bool {
	return !modelReflect.isZeroValue(jsonField)
}

// 根据json字段获取字段值，返回接口类型
func (modelReflect *ModelReflect) getFieldValue(jsonField string) interface{} {
	dbType := modelReflect.dbFieldTypeMap[jsonField]

	value := modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField])
	switch dbType {
	case "datetime":
		if value.IsZero() {
			return nil
		}
		return value.Interface()
	default:
		return value.Interface()
	}
}

// 根据json字段获取字段值, 返回string类型
func (modelReflect *ModelReflect) getFieldStringValue(jsonField string) string {
	if modelReflect.isZeroValue(jsonField) {
		return ""
	}
	return modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField]).String()
}
