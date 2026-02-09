module github.com/sebmaspd/go-ecs-log

go 1.22.5

replace github.com/seblkma/go-ecs-log/util => ./util

require (
	github.com/elastic/go-elasticsearch/v7 v7.17.10
	github.com/nxadm/tail v1.4.11
	github.com/seblkma/go-ecs-log/util v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.4
	go.elastic.co/ecslogrus v1.0.0
	gopkg.in/go-extras/elogrus.v7 v7.3.0
)

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1 // indirect
	github.com/magefile/mage v1.9.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
)
