module github.com/okcredit/billscan/service

go 1.13

require (
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/golang/protobuf v1.3.3
	github.com/lib/pq v1.8.0
	github.com/okcredit/billscan/api v0.0.0
	github.com/okcredit/go-common v1.8.0
	github.com/okcredit/mixpanel-go v0.1.1
	github.com/okcredit/nap v0.1.6
	github.com/satori/go.uuid v1.2.0
)

replace github.com/okcredit/billscan/api => ./../api
