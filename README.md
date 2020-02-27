# goseed🌱
Go (golang) based SAAS approach to quick start your project backend


## Why goseed?
- [TDD](https://en.wikipedia.org/wiki/Test-driven_development) (Test Driven Design) approach
- SAAS approach API (with tenant)
- API Versioning (using [echo](https://github.com/labstack/echo))
- Clean Infrastructure (inspired by [go-clean-arch](https://github.com/bxcodec/go-clean-arch))
- JWT Authentication method ([jwt-go](https://github.com/dgrijalva/jwt-go))
- Casbin Authorization ([casbin](https://github.com/casbin/casbin))
- Popular ORM ([xorm](https://gitea.com/xorm/xorm))
- Application Configuration ([viper](https://github.com/spf13/viper))
- Logging library ([zap](https://github.com/uber-go/zap))
- and much more...

## Run the application
```bash
# cloning
git clone https://github.com/sieryuu/goseed.git

# setup your database connection in ./appsettings.json
# if you are not using postgreSql, please change the dialect in ./main.go

# run the application
go run .

# unit testing
go test ./...
```

cURL 
```bash
# I am using Windows command prompt for testing

# ping
curl localhost:8888/ping

# login
curl --header "Content-Type: application/json" --request POST --data "{ \"Username\": \"admin\", \"Password\": \"123qwe\"}" localhost:8888/me/login
#output: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTgzMDM5MzQ4fQ.k5RH0bLJaxpfyRdbDXR1-qRWFt-gnhU3YlZV9R_Bcms"}

# using token to request something
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTgzMDM5MzQ4fQ.k5RH0bLJaxpfyRdbDXR1-qRWFt-gnhU3YlZV9R_Bcms" localhost:8888/v1/domain1/articles
```

## Finally
Any advice, contribution, suggestion are welcome 😊