// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package tdengine implements gdb.Driver for TDengine.
//
// 设计目标：
//  1. 让 gf gen dao / gf gen pbentity 能通过 TDengine STABLE 生成模型。
//  2. 通过 TDengine 官方 driver-go 的 taosWS database/sql 驱动建立连接。
//  3. Tables 返回 STABLE 列表；TableFields 解析 DESCRIBE 输出并包含 TAG 字段。
//
// 注意：
//
//	TDengine 的核心模型是 STABLE / TAGS / 子表 / TS，本驱动主要用于 gf 代码生成和基础
//	Raw SQL 操作；生产写入仍建议在业务仓储层封装子表名、TAGS、批量写入等规则。
package tdengine

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/taosdata/driver-go/v3/taosWS"

	"github.com/gogf/gf/v2/database/gdb"
)

const (
	driverName       = "tdengine"
	driverNameTaos   = "taos"
	driverNameTaosWS = "taosws"
	quoteChar        = "`"
)

// Driver is the driver for TDengine database.
type Driver struct {
	*gdb.Core
}

func init() {
	driverObj := New()
	for _, name := range []string{driverName, driverNameTaos, driverNameTaosWS} {
		if err := gdb.Register(name, driverObj); err != nil {
			panic(err)
		}
	}
}

// New creates and returns a driver that implements gdb.Driver.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for TDengine.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{Core: core}, nil
}

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

// Open creates and returns an underlying sql.DB object for TDengine.
//
// 推荐 link 写法：
//
//	tdengine:root:password@ws(127.0.0.1:56041)/mfk_iot?timezone=Asia%2FShanghai
//
// 兼容写法：
//
//	tdengine:root:password@tcp(127.0.0.1:56041)/mfk_iot
//
// 本驱动会把 @tcp( 自动转换为 @ws(，因为这里使用 taosWS 驱动。
func (d *Driver) Open(config *gdb.ConfigNode) (*sql.DB, error) {
	dsn := buildTaosWSDsn(config)
	return sql.Open("taosWS", dsn)
}

// CheckLocalTypeForField maps TDengine field types to GoFrame local types.
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	t := strings.ToUpper(strings.TrimSpace(fieldType))
	if i := strings.Index(t, "("); i >= 0 {
		t = strings.TrimSpace(t[:i])
	}

	switch t {
	case "BOOL":
		return gdb.LocalTypeBool, nil
	case "TINYINT", "SMALLINT", "INT":
		return gdb.LocalTypeInt, nil
	case "BIGINT", "TIMESTAMP":
		return gdb.LocalTypeInt64, nil
	case "UTINYINT", "USMALLINT", "UINT":
		return gdb.LocalTypeUint, nil
	case "UBIGINT":
		return gdb.LocalTypeUint64, nil
	case "FLOAT":
		return gdb.LocalTypeFloat32, nil
	case "DOUBLE":
		return gdb.LocalTypeFloat64, nil
	case "BINARY", "VARCHAR", "NCHAR", "JSON", "GEOMETRY", "VARBINARY":
		return gdb.LocalTypeString, nil
	default:
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}

// buildTaosWSDsn converts GoFrame ConfigNode into taosWS driver DSN.
func buildTaosWSDsn(config *gdb.ConfigNode) string {
	if config.Link != "" {
		link := strings.TrimSpace(config.Link)

		// GoFrame link 通常以 "type:" 开头，而 database/sql 的 DSN 不需要该前缀。
		for _, prefix := range []string{"tdengine:", "taosws:", "taos:"} {
			if strings.HasPrefix(strings.ToLower(link), prefix) {
				link = link[len(prefix):]
				break
			}
		}

		// 兼容用户沿用 GoFrame 习惯写 tcp(...) 的情况。
		link = strings.Replace(link, "@tcp(", "@ws(", 1)

		// TDengine 创建 MFK_IOT 后会显示为 mfk_iot。WebSocket DSN 中库名按小写更稳。
		link = normalizeDsnDatabaseName(link)
		return link
	}

	protocol := strings.TrimSpace(config.Protocol)
	if protocol == "" || strings.EqualFold(protocol, "tcp") {
		protocol = "ws"
	}

	host := strings.TrimSpace(config.Host)
	if host == "" {
		host = "127.0.0.1"
	}

	port := strings.TrimSpace(config.Port)
	if port == "" {
		port = "6041"
	}

	name := normalizeTDengineDatabaseName(config.Name)
	extra := strings.TrimSpace(config.Extra)
	if extra != "" && !strings.HasPrefix(extra, "?") {
		extra = "?" + extra
	}

	return fmt.Sprintf("%s:%s@%s(%s:%s)/%s%s", config.User, url.QueryEscape(config.Pass), protocol, host, port, name, extra)
}

// currentDatabaseName returns current TDengine database name for metadata SQL.
func (d *Driver) currentDatabaseName() string {
	if d != nil && d.GetConfig() != nil {
		if name := normalizeTDengineDatabaseName(d.GetConfig().Name); name != "" {
			return name
		}
		if name := databaseNameFromDsn(d.GetConfig().Link); name != "" {
			return name
		}
	}
	return ""
}

func normalizeTDengineDatabaseName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "`")
	return strings.ToLower(name)
}

func normalizeDsnDatabaseName(dsn string) string {
	slash := strings.LastIndex(dsn, "/")
	if slash < 0 || slash == len(dsn)-1 {
		return dsn
	}

	before := dsn[:slash+1]
	after := dsn[slash+1:]

	query := ""
	if question := strings.Index(after, "?"); question >= 0 {
		query = after[question:]
		after = after[:question]
	}

	if after == "" {
		return dsn
	}

	return before + normalizeTDengineDatabaseName(after) + query
}

func databaseNameFromDsn(dsn string) string {
	if dsn == "" {
		return ""
	}

	for _, prefix := range []string{"tdengine:", "taosws:", "taos:"} {
		if strings.HasPrefix(strings.ToLower(dsn), prefix) {
			dsn = dsn[len(prefix):]
			break
		}
	}

	slash := strings.LastIndex(dsn, "/")
	if slash < 0 || slash == len(dsn)-1 {
		return ""
	}

	name := dsn[slash+1:]
	if question := strings.Index(name, "?"); question >= 0 {
		name = name[:question]
	}

	return normalizeTDengineDatabaseName(name)
}

func quoteTDengineName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "`")
	return "`" + strings.ReplaceAll(name, "`", "``") + "`"
}
