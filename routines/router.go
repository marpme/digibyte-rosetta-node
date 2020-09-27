package routines

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/marpme/digibyte-rosetta-node/configuration"
)

func StartRouterCapabilities(wg *sync.WaitGroup, cfg *configuration.Config, router http.Handler) {
	log.Println("Listening on ", "0.0.0.0:"+cfg.Server.Port)
	err := http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router)
	if err != nil {
		log.Printf("Digibyte Rosetta Gateway server exited suddenly: %v\n", err)
		wg.Done()
		os.Exit(1)
	}
}
