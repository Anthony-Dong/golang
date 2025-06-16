package static

import (
	"net"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "static",
		Short: "Serve static files",
		RunE: func(cmd *cobra.Command, args []string) error {
			engine := gin.New()
			engine.Use(gin.Recovery())
			engine.Use(gin.Logger())
			registerStatic(engine)
			ips, err := utils.GetAllIP(true)
			if err != nil {
				return err
			}
			logs.Info("listen addr: http://%s", net.JoinHostPort("localhost", "8080"))
			for _, ip := range ips {
				logs.Info("listen addr: http://%s", net.JoinHostPort(ip.String(), "8080"))
			}
			return engine.Run()
		},
	}
	return cmd, nil
}
