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
