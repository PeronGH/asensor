package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/PeronGH/asensor/sensor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: asensor <command>")
		fmt.Fprintln(os.Stderr, "  list    list all sensors")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "list":
		if err := list(); err != nil {
			fmt.Fprintf(os.Stderr, "asensor list: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func list() error {
	manager, err := sensor.GetInstance()
	if err != nil {
		return err
	}

	sensors, err := manager.Sensors()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tVENDOR\tTYPE\tMIN DELAY (µs)\tRESOLUTION\tWAKE-UP")
	for _, s := range sensors {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%g\t%t\n",
			s.Name(),
			s.Vendor(),
			s.Type().String(),
			s.MinDelay(),
			s.Resolution(),
			s.IsWakeUp(),
		)
	}
	return w.Flush()
}
