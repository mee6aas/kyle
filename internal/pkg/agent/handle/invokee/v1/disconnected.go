package v1

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// Disconnected is invoked when invokee client disconnected
func (h Handle) Disconnected(addr *net.TCPAddr) {
	log.WithField("addr", addr.String()).Info("Disconnected")
}
