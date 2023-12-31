package nfsd

import (
	"net"

	"github.com/go-git/go-billy/v5"
	"github.com/willscott/go-nfs"
	nfshelper "github.com/willscott/go-nfs/helpers"
)

type ROFS struct {
	billy.Filesystem
}

// Capabilities exports the filesystem as readonly
func (ROFS) Capabilities() billy.Capability {
	return billy.ReadCapability | billy.SeekCapability
}

func Serve(addr string, fs billy.Filesystem) error {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	handler := nfshelper.NewNullAuthHandler(ROFS{fs})
	cacheHelper := nfshelper.NewCachingHandler(handler, 1024)
	return nfs.Serve(listener, cacheHelper)
}
