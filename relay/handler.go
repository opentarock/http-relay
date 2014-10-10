package relay

import (
	"mime"
	"net/http"
	"time"

	"github.com/opentarock/http-relay/vars"
	"github.com/opentarock/service-api/go/proto_headers"
	"github.com/opentarock/service-api/go/reqcontext"

	"code.google.com/p/go.net/context"
)

const defaultRequestTimeout = 2 * time.Second

func newContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx = reqcontext.WithCorrId(ctx, *proto_headers.NewRequestCorrelationHeader())
	return context.WithTimeout(ctx, defaultRequestTimeout)
}

func isJsonRequest(v string) bool {
	mt, _, err := mime.ParseMediaType(v)
	return err == nil && mt == "application/json"
}

func RelayHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := newContext(context.Background())
	defer cancel()

	logger := reqcontext.ContextLogger(ctx, "name", vars.ModuleName, "ip", r.RemoteAddr)

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		logger.Error("Only POST method allowed", "method", r.Method)
		return
	}
	contentType := r.Header.Get("Context-Type")
	if isJsonRequest(contentType) {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error("Context type must be application/json", "mime", contentType)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	logger.Info("Sending message to router")
	w.Write([]byte("{}"))
}
