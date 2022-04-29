module github.com/dez11de/exchange

go 1.17

require (
	github.com/bart613/decimal v1.2.1
	github.com/dez11de/cryptodb v0.0.0-00010101000000-000000000000
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	nhooyr.io/websocket v1.8.7
)

require (
	github.com/BurntSushi/toml v1.0.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	gorm.io/driver/mysql v1.2.3 // indirect
	gorm.io/gorm v1.22.5 // indirect
)

replace github.com/dez11de/cryptodb => ../cryptodb
