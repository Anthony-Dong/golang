package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/anthony-dong/golang/pkg/utils"
)

func InitPprof(addr string) error {
	if err := http.ListenAndServe(addr, http.DefaultServeMux); err != nil {
		return err
	}
	ips, err := utils.GetAllIP(true)
	if err != nil {
		return nil
	}
	for _, elem := range ips {
		fmt.Printf("http://%s%s/debug/pprof \n", elem, addr)
	}
	return nil
}
