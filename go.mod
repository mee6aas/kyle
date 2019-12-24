module github.com/mee6aas/kyle

go 1.12

require (
	github.com/google/uuid v1.1.1
	github.com/lesomnus/go-netstat v0.0.0-20190921134533-e16f81caffca
	github.com/mee6aas/zeep v0.0.0-20190921104938-fd56df8c2bbd
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/mitchellh/go-ps v0.0.0-20190716172923-621e5597135b
	github.com/otiai10/copy v1.0.2
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	golang.org/x/sys v0.0.0-20190924154521-2837fb4f24fe // indirect
	google.golang.org/grpc v1.23.1
	gotest.tools v2.2.0+incompatible
)

replace github.com/mee6aas/zeep => ../zeep
