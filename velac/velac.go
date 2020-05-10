package velac

import (
	"fmt"
	"os"
)

const usage = `VELA Version 2

Usage: vela [operation] [arguments]

Operations:
	status     Show circuit information    (vela status)
	up         Enable a VELA circuit       (vela up vela0)
	down       Disable a VELA circuit      (vela down vela0)
	add        Create a VELA circuit       (vela add [network] [remote] [vid])
	delete     Delete a VELA circuit       (vela delete vela0)
	service    Start/stop the VELA daemon  (vela service [up/down])`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "status":
		fmt.Println("Status")
	default:
		fmt.Println(usage)
	}
}
