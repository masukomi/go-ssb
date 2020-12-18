// SPDX-License-Identifier: MIT

package blobs

import (
	"context"

	"github.com/cryptix/go/logging"
	"github.com/pkg/errors"

	"go.cryptoscope.co/muxrpc/v2"

	"go.cryptoscope.co/ssb"
)

type addHandler struct {
	bs  ssb.BlobStore
	log logging.Interface
}

func (addHandler) HandleConnect(context.Context, muxrpc.Endpoint) {}

func (h addHandler) HandleCall(ctx context.Context, req *muxrpc.Request, edp muxrpc.Endpoint) {
	// TODO: push manifest check into muxrpc

	src, err := req.GetResponseSource()
	if err != nil {
		err = errors.Wrap(err, "add: couldn't get source")
		checkAndLog(h.log, err)
		req.CloseWithError(err)
		return
	}

	r := muxrpc.NewSourceReader(src)
	ref, err := h.bs.Put(r)
	if err != nil {
		err = errors.Wrap(err, "error putting blob")
		checkAndLog(h.log, err)
		req.CloseWithError(err)
		return
	}

	req.Return(ctx, ref)
}
