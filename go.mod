module github.com/libesz/poolmanager

go 1.24.0

require (
	github.com/auth0/go-jwt-middleware v1.0.1
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/nlopes/slack v0.6.0
	github.com/prometheus/client_golang v1.20.5
	github.com/urfave/negroni v1.0.0
	gopkg.in/yaml.v2 v2.4.0
	periph.io/x/conn/v3 v3.7.1
	periph.io/x/host/v3 v3.8.2
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/gorilla/websocket v1.2.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/errors v0.8.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/shurcooL/httpfs v0.0.0-20230704072500-f1e31cf0ba5c // indirect
	github.com/shurcooL/vfsgen v0.0.0-20230704071429-0000e147ea92 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace github.com/dgrijalva/jwt-go => github.com/golang-jwt/jwt/v4 v4.5.1
