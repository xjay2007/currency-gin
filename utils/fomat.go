package utils

import "fmt"

func FormatFileSize(s int64) string {
	size := float64(s)
	if size < 0 {
		return "Unknown"
	}
	if size == 0 {
		return "0B"
	}
	if size < 1024 {
		return fmt.Sprintf("%dB", int(size))
	}
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%.2fKiB", size)
	}
	size /= 1024
	if size < 1024 {
		return fmt.Sprintf("%.2fMiB", size)
	}
	size /= 1024
	return fmt.Sprintf("%.2fGiB", size)
}