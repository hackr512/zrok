package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"github.com/openziti-test-kitchen/zrok/cmd/zrok/endpoint_ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"html/template"
	"net"
	"net/http"
	"time"
)

func init() {
	testCmd.AddCommand(newTestEndpointCommand().cmd)
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Utilities used for testing zrok",
}

type testEndpointCommand struct {
	address string
	port    uint16
	t       *template.Template
	cmd     *cobra.Command
}

func newTestEndpointCommand() *testEndpointCommand {
	cmd := &cobra.Command{
		Use:   "endpoint",
		Short: "Start a simple HTTP endpoint",
		Args:  cobra.ExactArgs(0),
	}
	command := &testEndpointCommand{cmd: cmd}
	var err error
	if command.t, err = template.ParseFS(endpoint_ui.FS, "index.html"); err != nil {
		panic(err)
	}
	cmd.Flags().StringVarP(&command.address, "address", "a", "0.0.0.0", "The address for the HTTP listener")
	cmd.Flags().Uint16VarP(&command.port, "port", "p", 9090, "The port for the HTTP listener")
	cmd.Run = command.run
	return command
}

func (cmd *testEndpointCommand) run(_ *cobra.Command, _ []string) {
	http.HandleFunc("/", cmd.serveIndex)
	http.HandleFunc("/index.html", cmd.serveIndex)
	if err := http.ListenAndServe(fmt.Sprintf("%v:%d", cmd.address, cmd.port), nil); err != nil {
		panic(err)
	}
}

func (cmd *testEndpointCommand) serveIndex(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("%v {%v} -> /index.html", r.RemoteAddr, r.Host)
	ed := &endpointData{
		Now:    time.Now(),
		Host:   r.Host,
		Remote: r.RemoteAddr,
	}
	ed.getIps()
	if err := cmd.t.Execute(w, ed); err != nil {
		log.Error(err)
	}
}

type endpointData struct {
	Now    time.Time
	Host   string
	Remote string
	Ips    string
}

func (ed *endpointData) getIps() {
	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if len(ed.Ips) != 0 {
					ed.Ips += ", "
				}
				ed.Ips += ipnet.IP.String()
			}
		}
	}
}
