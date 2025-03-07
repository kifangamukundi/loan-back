module github.com/kifangamukundi/gm/loan

go 1.23.4

require (
	github.com/cloudinary/cloudinary-go/v2 v2.9.1
	github.com/gin-contrib/cors v1.7.3
	github.com/gin-gonic/gin v1.10.0
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	github.com/jwambugu/mpesa-golang-sdk v1.0.8
	github.com/kifangamukundi/gm/libs/auths v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/binders v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/parameters v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/queryparams v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/rates v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/repositories v0.0.0-00010101000000-000000000000
	github.com/kifangamukundi/gm/libs/transformations v0.0.0-00010101000000-000000000000
	github.com/robfig/cron/v3 v3.0.1
	github.com/twilio/twilio-go v1.23.12
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/postgres v1.5.11
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

require (
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.24.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/schema v1.4.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/ulule/limiter/v3 v3.11.2 // indirect
	golang.org/x/arch v0.12.0 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sync v0.11.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/kifangamukundi/gm/libs/auths => ../../libs/auths
	github.com/kifangamukundi/gm/libs/binders => ../../libs/binders
	github.com/kifangamukundi/gm/libs/exporters => ../../libs/exporters
	github.com/kifangamukundi/gm/libs/gates => ../../libs/gates
	github.com/kifangamukundi/gm/libs/parameters => ../../libs/parameters
	github.com/kifangamukundi/gm/libs/queryparams => ../../libs/queryparams
	github.com/kifangamukundi/gm/libs/rates => ../../libs/rates
	github.com/kifangamukundi/gm/libs/repositories => ../../libs/repositories
	github.com/kifangamukundi/gm/libs/schedules => ../../libs/schedules
	github.com/kifangamukundi/gm/libs/transformations => ../../libs/transformations
)
