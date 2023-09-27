package main

import (
	"net/http"

	"github.com/movio/bramble"

	"github.com/movio/bramble/portal"
	// import plugins to run the `init` and register
	_ "github.com/movio/bramble/plugins"
)

// Adds the tenant header to outgoing requests and registers the playground
type authPlugin struct {
	bramble.BasePlugin
}

func (*authPlugin) ID() string {
	return "auth"
}

func (*authPlugin) ApplyMiddlewarePublicMux(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTenant := r.Header.Get("X-Movio-Tenant")
		ctx := r.Context()
		ctx = bramble.AddOutgoingRequestsHeaderToContext(ctx, "X-Movio-Tenant", requestTenant)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

type portalPlugin struct {
	bramble.BasePlugin
}

func (*portalPlugin) ID() string {
	return "portal"
}

func (*portalPlugin) SetupPublicMux(mux *http.ServeMux) {
	mux.HandleFunc("/", portal.Handler("Core Portal", "/query"))
}

var (
	_ bramble.Plugin = new(authPlugin)
	_ bramble.Plugin = new(portalPlugin)
)

func main() {
	bramble.RegisterPlugin(&authPlugin{})
	bramble.RegisterPlugin(&portalPlugin{})
	bramble.Main()
}
