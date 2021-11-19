package filter

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

type Filter struct {
	DB *gorm.DB
}

func (f *Filter) SetStringFilter(stringField, value string) {
	if value != "" {
		f.DB = f.DB.Where(stringField+" like ?", value)
	}
}

func (f *Filter) SetBooleanFilter(boolField, value string) {
	whereSQL := boolField + " = ?"
	switch value {
	case "true":
		f.DB = f.DB.Where(whereSQL, true)
	case "false":
		f.DB = f.DB.Where(whereSQL, false)
	}
}

func (f *Filter) SetContainsFilter(stringField, value string) {
	if value != "" {
		f.DB = f.DB.Where("',' || " + stringField + " || ',' like '%,' || '" + value + "' || ',%'")
	}
}

func (f *Filter) SetOrder(field, value string) {
	if field != "" && (value == "asc" || value == "desc") {
		f.DB.Order(field + " " + value)
	}
}

func WhereLikes(column, tags string, model *gorm.DB) *gorm.DB {
	if tags != "" {
		tagsLikes := ""
		var tags_ []interface{}
		for _, tag := range strings.Split(tags, `,`) {
			tags_ = append(tags_, `%,`+tag+`,%`)
			tagsLikes = tagsLikes + ` or "` + column + `" like ?`
		}
		model = model.Where(tagsLikes[4:], tags_...)
	}
	return model
}

func WhereIn(column string, array []int, model *gorm.DB) *gorm.DB {
	if array != nil {
		arrayStr := []string{}
		for _, value := range array {
			arrayStr = append(arrayStr, strconv.Itoa(value))
		}
		model = model.Where(`"` + column + `" in (` + strings.Join(arrayStr, `,`) + `)`)
	}
	return model
}
