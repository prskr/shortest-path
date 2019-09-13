// Copyright © 2019 Peter Kurfer peter.kurfer@googlemail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/baez90/shortest-path/internal/app/config"
	"github.com/baez90/shortest-path/internal/app/crawling"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var (
	rootCmd = &cobra.Command{
		Use:   "shortest-path",
		Args:  cobra.ExactArgs(2),
		Short: "",
		Long:  ``,
		Run:   runTraverseCommand,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("failed to execute command")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("max-hops", "20", "depth of the search")
	rootCmd.PersistentFlags().String("log-level", "info", "log level to use")
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	log.SetLevel(log.InfoLevel)
}

func initConfig() {
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.WithError(err).Error("failed to bind flags to viper")
	}
}

func runTraverseCommand(cmd *cobra.Command, args []string) {

	crawler := crawling.NewWikiCrawler(args[0], args[1], uint16(viper.GetInt("max-hops")))

	start := time.Now()
	if res, err := crawler.SearchShortestPath(); err != nil {
		log.
			WithError(err).
			Error("Failed to resolve shortest path")
		os.Exit(2)
	} else {
		duration := time.Since(start)
		log.Infof("Resolved path in %d ms", duration.Milliseconds())
		for _, visitedPage := range res.VisitedPages() {
			log.Info(visitedPage)
		}
		log.Infof("Visited %d pages to find path", crawler.FetchedPages())
		log.Infof("Discovered %d unique links during search", crawler.DiscoveredPages())
	}
}
