module github.com/gogf/gf/cmd/gf/v2

go 1.20

require (
	github.com/gogf/gf/contrib/drivers/mssql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/mysql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/pgsql/v2 v2.1.0
	github.com/gogf/gf/contrib/drivers/sqlite/v2 v2.1.0
	github.com/gogf/gf/v2 v2.1.0
	github.com/olekukonko/tablewriter v0.0.5
	golang.org/x/tools v0.1.11
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/clbanning/mxj/v2 v2.5.5 // indirect
	github.com/denisenkom/go-mssqldb v0.11.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/glebarez/go-sqlite v1.17.3 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-sql/civil v0.0.0-20190719163853-cb61b32ac6fe // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grokify/html-strip-tags-go v0.0.1 // indirect
	github.com/lib/pq v1.10.4 // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
	golang.org/x/text v0.3.8-0.20211105212822-18b340fc7af2 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.22.1 // indirect
)

replace (
	github.com/gogf/gf/contrib/drivers/mssql/v2 => ../../contrib/drivers/mssql/
	github.com/gogf/gf/contrib/drivers/mysql/v2 => ../../contrib/drivers/mysql/
	github.com/gogf/gf/contrib/drivers/pgsql/v2 => ../../contrib/drivers/pgsql/
	github.com/gogf/gf/contrib/drivers/sqlite/v2 => ../../contrib/drivers/sqlite/
	github.com/gogf/gf/v2 => ../../
)
