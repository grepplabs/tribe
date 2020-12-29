package cmd

import (
	"github.com/grepplabs/tribe/config"
	"github.com/spf13/cobra"
)

func initCorsFlags(cmd *cobra.Command, cc *config.CorsConfig) {
	cmd.Flags().BoolVar(&cc.Enabled, "cors-enabled", false, "enables CORS. It will be enabled and preflight-requests (OPTION) will be answered")
	cmd.Flags().StringSliceVar(&cc.Options.AllowedOrigins, "cors-allowed-origins", []string{}, "A list of origins a cross-domain request can be executed from. If the special * value is present in the list, all origins will be allowed. An origin may contain a wildcard (*) to replace 0 or more characters")
	cmd.Flags().StringSliceVar(&cc.Options.AllowedMethods, "cors-allowed-methods", []string{"GET", "POST", "PUT", "PATCH", "DELETE"}, "A list of methods the client is allowed to use with cross-domain requests")
	cmd.Flags().StringSliceVar(&cc.Options.AllowedHeaders, "cors-allowed-headers", []string{"Authorization", "Content-Type", "Origin", "Accept", "X-Requested-With"}, "A list of non simple headers the client is allowed to use with cross-domain requests")
	cmd.Flags().StringSliceVar(&cc.Options.ExposedHeaders, "cors-exposed-headers", []string{}, "Indicates which headers are safe to expose to the API of a CORS API specification")
	cmd.Flags().BoolVar(&cc.Options.AllowCredentials, "cors-allow-credentials", true, "Indicates whether the request can include user credentials like cookies, HTTP authentication or client side SSL certificates")
	cmd.Flags().BoolVar(&cc.Options.OptionsPassthrough, "cors-options-passthrough", false, "Instructs preflight to let other potential next handlers to process the OPTIONS method. Turn this on if your application handles OPTIONS")
	cmd.Flags().IntVar(&cc.Options.MaxAge, "cors-max-age", 0, "Indicates how long (in seconds) the results of a preflight request can be cached. The default is 0 which stands for no max age")
	cmd.Flags().BoolVar(&cc.Options.Debug, "cors-debug", false, "Adds additional output to debug server side CORS issues")
}
