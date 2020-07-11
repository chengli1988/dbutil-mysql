package dbutil

import "reflect"

// Model 接口
type Model interface {
	// 返回数据库表名
	GetTableName() string
}

// ModelReflect 实体
type ModelReflect struct {
	tableName         string
	dbFields          []string          // 数据库字段
	dbFieldMap        map[string]string // db字段名Map-json字段名, key:db字段名, value: json字段名 (tag为dbField)
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

			// 只有带有dbField且值且true的为数据库字段
			if fieldTag.Get("dbField") == "true" {
				modelReflect.dbFieldMap[fieldTag.Get("db")] = fieldTag.Get("json")
				modelReflect.dbFields = append(modelReflect.dbFields, fieldTag.Get("db"))
			}
		}
	}
}

// 根据json字段检查字段值是否不为空字符串("")、零(0)
func (modelReflect *ModelReflect) checkFieldValid(jsonField string) bool {
	fieldValue := modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField])

	if !fieldValue.IsValid() || modelReflect.getFieldValue(jsonField) == "" || modelReflect.getFieldValue(jsonField) == 0 {
		return false
	}

	return true
}

// 根据json字段获取字段值，返回接口类型
func (modelReflect *ModelReflect) getFieldValue(jsonField string) interface{} {
	return modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField]).Interface()
}

// 根据json字段获取字段值, 返回string类型
func (modelReflect *ModelReflect) getFieldStringValue(jsonField string) string {
	return modelReflect.modelValue.FieldByName(modelReflect.modelFieldNameMap[jsonField]).String()
}
