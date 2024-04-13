package utils

import (
	"fmt"
	"strings"
)

func FormatSize(size int) string {
	trimSize := float64(size)
	units := []string{"B", "KB", "M", "G"}
	unitsIndex := 0
	for {
		tmpSize := trimSize / 1024
		if tmpSize < 1 || unitsIndex >= len(units)-1 {
			data := fmt.Sprintf("%.2f", trimSize)
			data = strings.TrimSuffix(data, "0")
			data = strings.TrimSuffix(data, "0")
			data = strings.TrimSuffix(data, ".")
			return fmt.Sprintf("%s%s", data, units[unitsIndex])
		}
		trimSize = tmpSize
		unitsIndex = unitsIndex + 1
	}
}
