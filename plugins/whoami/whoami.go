// SPDX-License-Identifier: MIT

package whoami

import (
	"context"
	"fmt"

	"github.com/cryptix/go/logging"
	"github.com/pkg/errors"
	"go.cryptoscope.co/muxrpc/v2"
	refs "go.mindeco.de/ssb-refs"

	"go.cryptoscope.co/ssb"
)

var (
	_      ssb.Plugin = plugin{} // compile-time type check
	method            = muxrpc.Method{"whoami"}
)

func checkAndLog(log logging.Interface, err error) {
	if err != nil {
		if err := logging.LogPanicWithStack(log, "checkAndLog", err); err != nil {
			log.Log("event", "warning", "msg", "faild to write panic file", "err", err)
		}
	}
}

func New(log logging.Interface, id *refs.FeedRef) ssb.Plugin {
	return plugin{handler{
		log: log,
		id:  id,
	}}
}

type plugin struct {
	h handler
}

func (plugin) Name() string { return "whoami" }

func (plugin) Method() muxrpc.Method { return method }

func (wami plugin) Handler() muxrpc.Handler { return wami.h }

func (plugin) WrapEndpoint(edp muxrpc.Endpoint) interface{} {
	return endpoint{edp}
}

type handler struct {
	log logging.Interface
	id  *refs.FeedRef
}

func (handler) HandleConnect(ctx context.Context, edp muxrpc.Endpoint) {}

func (h handler) HandleCall(ctx context.Context, req *muxrpc.Request, edp muxrpc.Endpoint) {
	// TODO: push manifest check into muxrpc
	if req.Type == "" {
		req.Type = "async"
	}
	if req.Method.String() != "whoami" {
		req.CloseWithError(fmt.Errorf("wrong method"))
		return
	}
	type ret struct {
		ID string `json:"id"`
	}

	err := req.Return(ctx, ret{h.id.Ref()})
	checkAndLog(h.log, err)
}

type endpoint struct {
	edp muxrpc.Endpoint
}

func (edp endpoint) WhoAmI(ctx context.Context) (refs.FeedRef, error) {
	var resp struct {
		ID refs.FeedRef `json:"id"`
	}

	err := edp.edp.Async(ctx, &resp, muxrpc.TypeJSON, method)
	if err != nil {
		return refs.FeedRef{}, errors.Wrap(err, "error making async call")
	}

	return resp.ID, nil
}
