module github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/sorting-service

go 1.16

replace github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen => ../gen

require (
	github.com/angelRaynov/golang-at-ocado/week-03/fulfillment-skeleton/sort/gen v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.37.1
	google.golang.org/protobuf v1.26.0 // indirect
)
