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
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
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
	cmd.Flags().IntVar(&opts.listenPort, "port", 1234, "The port to listen on for Alertmanager notifications")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		fmt.Printf("listen on port %v\n", opts.listenPort)
		// todo start listing on alertmanager notifications
		return nil
	}

	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nAborted...")
		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
