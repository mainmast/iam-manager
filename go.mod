module mainmast/iam-manager

require (
	github.com/golang-migrate/migrate/v4 v4.7.0
	github.com/lib/pq v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/valyala/fasthttp v1.6.0
	gitlab.com/mainmast/microservices/cm-http.git v1.0.4
	gitlab.com/mainmast/microservices/iam/iam-models.git v1.0.1
	gitlab.com/mainmast/microservices/iam/migrations.git v1.0.0
	golang.org/x/crypto v0.0.0-20190426145343-a29dc8fdc734
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

go 1.13
