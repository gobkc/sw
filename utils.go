package sw

import "os"

func resolveAddress(addr []string) (newAddr string) {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		logPrint("Environment variable PORT is undefined. Using port :8080 by default")
		newAddr = ":8080"
	case 1:
		newAddr = addr[0]
	default:
		panic("too many parameters")
	}
	return
}
