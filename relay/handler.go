package relay

import (
	"io/ioutil"
	"mime"
	"net/http"
	"time"

	"github.com/opentarock/http-relay/vars"
	"github.com/opentarock/service-api/go/client"
	"github.com/opentarock/service-api/go/proto_headers"
	"github.com/opentarock/service-api/go/reqcontext"

	"code.google.com/p/go.net/context"
)

const defaultRequestTimeout = 10 * time.Second

func newContext(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx = reqcontext.WithCorrId(ctx, *proto_headers.NewRequestCorrelationHeader())
	ctx = reqcontext.WithAuth(ctx, "1", "token")
	return context.WithTimeout(ctx, defaultRequestTimeout)
}

func isJsonRequest(v string) bool {
	mt, _, err := mime.ParseMediaType(v)
	return err == nil && mt == "application/json"
}

func NewRelayHandler(mhClient client.MsgHandlerClient) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := newContext(context.Background())
		defer cancel()

		logger := reqcontext.ContextLogger(ctx, "name", vars.ModuleName, "ip", r.RemoteAddr)

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			logger.Error("Only POST method allowed", "method", r.Method)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if !isJsonRequest(contentType) {
			w.WriteHeader(http.StatusBadRequest)
			logger.Error("Content type must be application/json", "mime", contentType)
			return
		}

		logger.Info("Sending message to router")

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("Error reading request body", "error", err)
			return
		}

		result, err := mhClient.RouteMessage(ctx, string(data))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Error("Error routing the message", "error", err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(result.GetData()))
	})
}
