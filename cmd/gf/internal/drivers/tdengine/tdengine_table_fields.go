// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package tdengine

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// TableFields retrieves and returns the field information of specified TDengine STABLE.
//
// DESCRIBE STABLE 的输出包含普通列和 TAGS，本函数会把两者都转换成 gdb.TableField，
// 便于 gf gen 生成包含 TAG 字段的模型。
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result gdb.Result
		link   gdb.Link
	)

	if link, err = d.SlaveLink(schema...); err != nil {
		return nil, err
	}

	database := d.currentDatabaseName()
	if database == "" {
		result, err = d.DoSelect(ctx, link, fmt.Sprintf(`DESCRIBE %s`, d.QuoteWord(table)))
	} else {
		// TDengine 官方语法支持 DESCRIBE [db_name.]stb_name。
		result, err = d.DoSelect(ctx, link, fmt.Sprintf(`DESCRIBE %s.%s`, quoteTDengineName(database), d.QuoteWord(table)))
	}
	if err != nil {
		return nil, err
	}

	fields = make(map[string]*gdb.TableField)
	for i, record := range result {
		name := recordString(record, "field", "Field", "name", "Name")
		if name == "" || strings.EqualFold(name, "field") {
			continue
		}

		fieldType := buildTDengineFieldType(record)
		note := recordString(record, "note", "Note")
		isTag := strings.Contains(strings.ToLower(note), "tag")

		key := ""
		extra := ""
		comment := ""

		if strings.EqualFold(name, "ts") {
			key = "PRI"
			comment = "TDengine timestamp column"
		}
		if isTag {
			extra = "TAG"
			if comment == "" {
				comment = "TDengine tag column"
			} else {
				comment += "; TDengine tag column"
			}
		}

		fields[name] = &gdb.TableField{
			Index:   i,
			Name:    name,
			Type:    normalizeTDengineType(fieldType),
			Null:    false,
			Key:     key,
			Default: nil,
			Extra:   extra,
			Comment: comment,
		}
	}

	return fields, nil
}

func recordString(record gdb.Record, names ...string) string {
	for _, name := range names {
		if value, ok := record[name]; ok && value != nil {
			return strings.TrimSpace(value.String())
		}
	}
	for key, value := range record {
		for _, name := range names {
			if strings.EqualFold(key, name) && value != nil {
				return strings.TrimSpace(value.String())
			}
		}
	}
	return ""
}

func buildTDengineFieldType(record gdb.Record) string {
	fieldType := recordString(record, "type", "Type")
	if fieldType == "" {
		return "NCHAR"
	}

	length := recordString(record, "length", "Length")
	if length == "" || length == "-1" || length == "0" {
		return fieldType
	}

	upper := strings.ToUpper(strings.TrimSpace(fieldType))
	if strings.Contains(upper, "(") {
		return fieldType
	}

	switch upper {
	case "NCHAR", "BINARY", "VARCHAR", "VARBINARY":
		return fmt.Sprintf("%s(%s)", upper, length)
	default:
		return fieldType
	}
}

func normalizeTDengineType(fieldType string) string {
	t := strings.TrimSpace(fieldType)
	if t == "" {
		return "NCHAR"
	}

	upper := strings.ToUpper(t)
	base := upper
	if i := strings.Index(base, "("); i >= 0 {
		base = strings.TrimSpace(base[:i])
	}

	switch base {
	case "TIMESTAMP", "BOOL", "TINYINT", "SMALLINT", "INT", "BIGINT", "UTINYINT", "USMALLINT", "UINT", "UBIGINT", "FLOAT", "DOUBLE":
		return base
	case "NCHAR", "BINARY", "VARCHAR", "VARBINARY":
		return upper
	}

	return t
}
