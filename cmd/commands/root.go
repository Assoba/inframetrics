/*
 * Copyright 2017  Assoba S.A.S.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 *
 */

package commands

import (
	"github.com/spf13/cobra"
	"github.com/assoba.fr/inframetrics"
	"net/http"
	"strconv"
	"fmt"
)

var port uint16

var RootCmd = &cobra.Command{
	Use:   "assoba-metrics",
	Short: "Metrics collector for assoba nodes",
	Long:  `Metrics collector for assoba nodes, uses cadvisor, consul, weave and nomad local metrics, and exposes a useful subset for prometheus 2.0 collection`,
	Run: func(cmd *cobra.Command, args []string) {
		go inframetrics.RunStats()
		http.HandleFunc("/metrics", inframetrics.Handler)
		fmt.Printf("Listening on 0.0.0.0:%d \n", port)
		http.ListenAndServe(":"+strconv.Itoa(int(port)), nil)
	},
}

func init() {
	cobra.OnInitialize()
	RootCmd.PersistentFlags().Uint16VarP(&port, "port", "p", 4890, "Http port to listen on")
}
