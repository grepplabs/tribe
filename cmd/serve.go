package cmd

import (
	"github.com/grepplabs/tribe/api/v1/server"
	"github.com/grepplabs/tribe/config"
	"github.com/grepplabs/tribe/pkg/log"
	"github.com/spf13/cobra"
	"time"
)

var (
	serverConfig = new(config.ServerConfig)
	logConfig    = new(log.Configuration)
)

// serveCmd represents the server command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Tribe Server",
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func initServerFlags(cmd *cobra.Command, sc *config.ServerConfig) {
	cmd.Flags().StringSliceVar(&sc.EnabledListeners, "server-scheme", []string{"http"}, "the listeners to enable, this can be repeated and defaults to the schemes in the swagger spec")
	cmd.Flags().DurationVar(&sc.CleanupTimeout, "server-cleanup-timeout", 10*time.Second, "grace period for which to wait before killing idle connections")
	cmd.Flags().DurationVar(&sc.GracefulTimeout, "server-graceful-timeout", 15*time.Second, "grace period for which to wait before shutting down the server")
	cmd.Flags().IntVar(&sc.MaxHeaderSize, "server-max-header-size", 1048576, "controls the maximum number of bytes the server will read parsing the request header'sc keys and values, including the request line. It does not limit the size of the request body")

	cmd.Flags().StringVar(&sc.Host, "server-host", "localhost", "the IP to listen on")
	cmd.Flags().IntVar(&sc.Port, "server-port", 0, "the port to listen on for insecure connections, defaults to a random value")
	cmd.Flags().IntVar(&sc.ListenLimit, "server-listen-limit", 0, "limit the number of outstanding requests")

	cmd.Flags().DurationVar(&sc.KeepAlive, "server-keep-alive", 3*time.Minute, "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections")
	cmd.Flags().DurationVar(&sc.ReadTimeout, "server-read-timeout", 30*time.Second, "maximum duration before timing out read of the request")
	cmd.Flags().DurationVar(&sc.WriteTimeout, "server-write-timeout", 60*time.Second, "maximum duration before timing out write of the response")

	cmd.Flags().StringVar(&sc.TLSHost, "server-tls-host", "", "the IP to listen on for tls, when not specified it'sc the same as --host")
	cmd.Flags().IntVar(&sc.TLSPort, "server-tls-port", 0, "the port to listen on for secure connections, defaults to a random value")

	cmd.Flags().StringVar(&sc.TLSCertificate, "server-tls-certificate", "", "the certificate to use for secure connections")
	cmd.Flags().StringVar(&sc.TLSCertificateKey, "server-tls-key", "", "the private key to use for secure connections")
	cmd.Flags().StringVar(&sc.TLSCACertificate, "server-tls-ca", "", "he certificate authority file to be used with mutual tls auth")

	cmd.Flags().DurationVar(&sc.TLSKeepAlive, "server-tls-keep-alive", 0, "sets the TCP keep-alive timeouts on accepted connections. It prunes dead TCP connections")
	cmd.Flags().DurationVar(&sc.TLSReadTimeout, "server-tls-read-timeout", 0, "maximum duration before timing out read of the request")
	cmd.Flags().DurationVar(&sc.TLSWriteTimeout, "server-tls-write-timeout", 0, "maximum duration before timing out write of the response")
}

type Server struct {
	*server.Server
}

func NewServer() *Server {
	return &Server{server.NewServer(nil)}
}

func (s *Server) configureServerFlags() {
	s.EnabledListeners = serverConfig.EnabledListeners
	s.CleanupTimeout = serverConfig.CleanupTimeout
	s.GracefulTimeout = serverConfig.GracefulTimeout
	s.MaxHeaderSize = serverConfig.MaxHeaderSize

	s.Host = serverConfig.Host
	s.Port = serverConfig.Port
	s.ListenLimit = serverConfig.ListenLimit
	s.KeepAlive = serverConfig.KeepAlive
	s.ReadTimeout = serverConfig.ReadTimeout
	s.WriteTimeout = serverConfig.WriteTimeout

	s.TLSHost = serverConfig.TLSHost
	s.TLSPort = serverConfig.TLSPort
	s.TLSCertificate = serverConfig.TLSCertificate
	s.TLSCertificateKey = serverConfig.TLSCertificateKey
	s.TLSCACertificate = serverConfig.TLSCACertificate
	s.TLSListenLimit = serverConfig.TLSListenLimit
	s.TLSKeepAlive = serverConfig.TLSKeepAlive
	s.TLSReadTimeout = serverConfig.TLSReadTimeout
	s.TLSWriteTimeout = serverConfig.TLSWriteTimeout
}
