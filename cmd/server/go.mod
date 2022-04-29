module github.com/dez11de/trademan/cmd/server

go 1.17

replace github.com/dez11de/exchange => ../../exchange

replace github.com/dez11de/cryptodb => ../../cryptodb

require (
	github.com/BoRuDar/configuration/v3 v3.1.0
	github.com/bart613/decimal v1.2.1
	github.com/dez11de/cryptodb v0.0.0-00010101000000-000000000000
	github.com/dez11de/exchange v0.0.0-00010101000000-000000000000
	github.com/julienschmidt/httprouter v1.3.0
	gorm.io/gorm v1.22.5
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.4 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.10.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/stretchr/testify v1.7.1-0.20210427113832-6241f9ab9942 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/driver/mysql v1.2.3 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)
