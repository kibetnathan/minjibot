module github.com/kibetnathan/minjibot/integration_tests

go 1.26.2

require (
	github.com/jackc/pgx/v5 v5.9.2
	github.com/kibetnathan/minjibot v0.0.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	golang.org/x/text v0.40.0 // indirect
)

replace github.com/kibetnathan/minjibot => ../
