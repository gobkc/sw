package sw

import (
	"fmt"
)

func logPrint(format string, values ...interface{}) {
	fmt.Printf("%c[0;40;37m%s %v%c[0m\n", 0x1B, "[SW-LOG] "+format,fmt.Sprint(values...), 0x1B)
}

func logInfo(format string, values ...interface{}) {
	fmt.Printf("%c[0;31;34m%s %v%c[0m\n", 0x1B, "[SW-LOG] "+format,fmt.Sprint(values...), 0x1B)
}

func logDefault(format string, values ...interface{}) {
	fmt.Printf("%c[0;36;34m%s %v%c[0m\n", 0x1B, "[SW-LOG] "+format,fmt.Sprint(values...), 0x1B)
}