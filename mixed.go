package servers

import (
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	appserver "github.com/codemodify/systemkit-appserver"
	reflection "github.com/codemodify/systemkit-helpers-reflection"
	logging "github.com/codemodify/systemkit-logging"
)

// MixedServer -
type MixedServer struct {
	servers []appserver.IServer
}

// NewMixedServer -
func NewMixedServer(servers []appserver.IServer) appserver.IServer {
	return &MixedServer{
		servers: servers,
	}
}

// Run - Implement `IServer`
func (thisRef *MixedServer) Run(ipPort string, enableCORS bool) error {
	listener, err := net.Listen("tcp4", ipPort)
	if err != nil {
		return err
	}

	var router = mux.NewRouter()
	thisRef.PrepareRoutes(router)
	thisRef.RunOnExistingListenerAndRouter(listener, router, enableCORS)

	return nil
}

// PrepareRoutes - Implement `IServer`
func (thisRef *MixedServer) PrepareRoutes(router *mux.Router) {
	for _, server := range thisRef.servers {
		server.PrepareRoutes(router)
	}
}

// RunOnExistingListenerAndRouter - Implement `IServer`
func (thisRef *MixedServer) RunOnExistingListenerAndRouter(listener net.Listener, router *mux.Router, enableCORS bool) {
	if enableCORS {
		corsSetterHandler := cors.Default().Handler(router)
		err := http.Serve(listener, corsSetterHandler)
		if err != nil {
			logging.Fatalf("%s, from %s", err.Error(), reflection.GetThisFuncName())

			os.Exit(-1)
		}
	} else {
		err := http.Serve(listener, router)
		if err != nil {
			logging.Fatalf("%s, from %s", err.Error(), reflection.GetThisFuncName())

			os.Exit(-1)
		}
	}
}
