package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/PeronGH/asensor/sensor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: asensor <command>")
		fmt.Fprintln(os.Stderr, "  list [-verbose]    list sensors")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "list":
		if err := list(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "asensor list: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func list(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	verbose := fs.Bool("verbose", false, "show full details for each sensor")
	if err := fs.Parse(args); err != nil {
		return err
	}

	manager, err := sensor.GetInstance()
	if err != nil {
		return err
	}

	sensors, err := manager.Sensors()
	if err != nil {
		return err
	}

	if !*verbose {
		for i, s := range sensors {
			fmt.Printf("%3d  %s\n", i, s.Name())
		}
		return nil
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
