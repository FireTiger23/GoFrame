// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package tdengine

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

// Tables retrieves and returns the STABLE names of current TDengine database.
//
// TDengine 的建模核心是 STABLE，gf gen 这里把 STABLE 当作“表”返回，
// 这样可以统一接入 gf gen dao / gf gen pbentity。
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var (
		result gdb.Result
		link   gdb.Link
	)

	if link, err = d.SlaveLink(schema...); err != nil {
		return nil, err
	}

	database := d.currentDatabaseName()
	if database == "" {
		// 兜底：没有数据库名时，只能依赖 DSN 当前库。
		result, err = d.DoSelect(ctx, link, `SHOW STABLES`)
	} else {
		// TDengine 支持 show db.stables；避免 WebSocket 连接没有自动切库导致 [0x388] Database not exist。
		result, err = d.DoSelect(ctx, link, fmt.Sprintf(`SHOW %s.STABLES`, quoteTDengineName(database)))
	}
	if err != nil {
		return nil, err
	}

	tableSet := make(map[string]struct{})
	for _, record := range result {
		name := stableNameFromRecord(record)
		if name == "" {
			continue
		}
		tableSet[name] = struct{}{}
	}

	for table := range tableSet {
		tables = append(tables, table)
	}
	sort.Strings(tables)
	return tables, nil
}

func stableNameFromRecord(record gdb.Record) string {
	for _, key := range []string{"stable_name", "name", "STABLE_NAME", "NAME"} {
		if value, ok := record[key]; ok && value != nil {
			name := strings.TrimSpace(value.String())
			if name != "" && !strings.EqualFold(name, "stable_name") && !strings.EqualFold(name, "name") {
				return name
			}
		}
	}

	// 兜底：不同 TDengine 版本字段名可能不同。取第一个非空字段，但不要遍历多个字段作为多个表名。
	for _, value := range record {
		if value == nil {
			continue
		}
		name := strings.TrimSpace(value.String())
		if name != "" && !strings.EqualFold(name, "stable_name") && !strings.EqualFold(name, "name") {
			return name
		}
	}

	return ""
}
