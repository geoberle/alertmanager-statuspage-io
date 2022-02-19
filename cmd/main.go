/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	alertprocessor "alertmanager-statuspage-io/alertprocessor"
)

type Options struct {
	listenPort int
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cmd := &cobra.Command{
		Use:              "alertmanager-statuspage-io",
		SilenceUsage:     true,
		TraverseChildren: true,
	}

	var opts Options
	cmd.Flags().IntVar(&opts.listenPort, "port", 1236, "The port to listen on for Alertmanager notifications")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		processor, err := alertprocessor.NewAlertProcessor()
		if err != nil {
			log.Fatalf("Failed to create alert process: %s", err)
		}

		mux := http.NewServeMux()
		mux.Handle("/receiver", processor)
		server := &http.Server{Addr: fmt.Sprintf(":%v", opts.listenPort), Handler: mux}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("listen:%+s\n", err)
			}
		}()
		log.Printf("Listening on port %v\n", opts.listenPort)

		<-cmd.Context().Done()

		log.Printf("server stopped")

		ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer func() {
			cancel()
		}()

		if err := server.Shutdown(ctxShutDown); err != nil {
			log.Fatalf("server Shutdown Failed:%+s", err)
		}
		log.Printf("server exited properly")
	}

	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		log.Printf("Aborted signal received ...")
		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		log.Fatalf("%v\n", err)
		os.Exit(1)
	}
}
