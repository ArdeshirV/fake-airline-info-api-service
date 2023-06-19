package service

type color string

const (
	Normal      color = "\033[0m"
	Red         color = "\033[0;31m"
	BoldRed     color = "\033[1;31m"
	Blue        color = "\033[0;34m"
	BoldBlue    color = "\033[1;34m"
	Green       color = "\033[0;32m"
	BoldGreen   color = "\033[1;32m"
	Yellow      color = "\033[0;33m"
	BoldYellow  color = "\033[1;33m"
	Magneta     color = "\033[0;35m"
	BoldMagneta color = "\033[1;35m"
)
