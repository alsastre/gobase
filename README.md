# Readme
This project contains the basic structure for any golang project with example of the following libraries:
- Database connections: https://upper.io/v4/
- Logging: https://github.com/uber-go/zap
- HTTP Request: https://github.com/go-chi/chi
- Load configuration: https://github.com/spf13/viper

GOVERSION: 1.15

Start a postgresql DB for testing with the Dummy data:
`docker run -p 5432:5432 --name gobasePostgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres`
