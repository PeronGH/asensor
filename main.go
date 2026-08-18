package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/PeronGH/asensor/sensor"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: asensor <command>")
		fmt.Fprintln(os.Stderr, "  list [-verbose]    list sensors")
		fmt.Fprintln(os.Stderr, "  show <index>       show sensor metadata")
		fmt.Fprintln(os.Stderr, "  read [-timeout d] <index>    read one sensor event")
		fmt.Fprintln(os.Stderr, "  watch [-duration d] <index>  stream sensor events")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "list":
		if err := list(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "asensor list: %v\n", err)
			os.Exit(1)
		}
	case "show":
		if err := show(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "asensor show: %v\n", err)
			os.Exit(1)
		}
	case "read":
		if err := read(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "asensor read: %v\n", err)
			os.Exit(1)
		}
	case "watch":
		if err := watch(os.Args[2:]); err != nil {
			fmt.Fprintf(os.Stderr, "asensor watch: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

func list(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	verbose := fs.Bool("verbose", false, "show full details for each sensor")
	fs.Parse(args)

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

func show(args []string) error {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: asensor show <index>")
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly one sensor index, got %d", fs.NArg())
	}
	index, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid index %q: %w", fs.Arg(0), err)
	}

	manager, err := sensor.GetInstance()
	if err != nil {
		return err
	}
	sensors, err := manager.Sensors()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(sensors) {
		return fmt.Errorf("index %d out of range (0-%d)", index, len(sensors)-1)
	}

	s := sensors[index]

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Name:\t%s\n", s.Name())
	fmt.Fprintf(w, "Vendor:\t%s\n", s.Vendor())
	fmt.Fprintf(w, "Type:\t%s (%d)\n", s.Type().String(), int32(s.Type()))
	fmt.Fprintf(w, "StringType:\t%s\n", s.StringType())
	fmt.Fprintf(w, "MinDelay:\t%d µs\n", s.MinDelay())
	fmt.Fprintf(w, "Resolution:\t%g\n", s.Resolution())
	fmt.Fprintf(w, "Reporting:\t%s\n", s.ReportingMode().String())
	fmt.Fprintf(w, "WakeUp:\t%t\n", s.IsWakeUp())
	fmt.Fprintf(w, "FIFO max:\t%d\n", s.FifoMaxEventCount())
	fmt.Fprintf(w, "FIFO reserved:\t%d\n", s.FifoReservedEventCount())
	fmt.Fprintf(w, "Handle:\t%d\n", s.Handle())
	fmt.Fprintf(w, "DirectRate:\t%d\n", s.HighestDirectReportRateLevel())
	return w.Flush()
}

func read(args []string) error {
	fs := flag.NewFlagSet("read", flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	timeout := fs.Duration("timeout", 5*time.Second, "how long to wait for an event")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: asensor read [-timeout duration] <index>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly one sensor index, got %d", fs.NArg())
	}
	index, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid index %q: %w", fs.Arg(0), err)
	}

	manager, err := sensor.GetInstance()
	if err != nil {
		return err
	}
	sensors, err := manager.Sensors()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(sensors) {
		return fmt.Errorf("index %d out of range (0-%d)", index, len(sensors)-1)
	}
	s := sensors[index]

	ev, err := manager.Read(s, *timeout)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Sensor:\t%s\n", s.Name())
	fmt.Fprintf(w, "Timestamp:\t%d ms\n", ev.Timestamp/1_000_000)
	if ev.Type == sensor.StepCounter {
		fmt.Fprintf(w, "Steps:\t%d\n", ev.StepCount())
	} else {
		fmt.Fprintf(w, "Values:\t%s\n", eventValues(ev))
	}
	return w.Flush()
}

// eventValues formats the meaningful float values of an event.
func eventValues(ev *sensor.Event) string {
	if ev.Type == sensor.StepCounter {
		return fmt.Sprintf("steps=%d", ev.StepCount())
	}
	n := ev.Type.ValueCount()
	vals := make([]string, n)
	for i := 0; i < n; i++ {
		vals[i] = fmt.Sprintf("%g", ev.Data[i])
	}
	return strings.Join(vals, " ")
}

func watch(args []string) error {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	fs.SetOutput(os.Stderr)
	duration := fs.Duration("duration", 10*time.Second, "how long to stream events")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: asensor watch [-duration duration] <index>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	if fs.NArg() != 1 {
		return fmt.Errorf("expected exactly one sensor index, got %d", fs.NArg())
	}
	index, err := strconv.Atoi(fs.Arg(0))
	if err != nil {
		return fmt.Errorf("invalid index %q: %w", fs.Arg(0), err)
	}

	manager, err := sensor.GetInstance()
	if err != nil {
		return err
	}
	sensors, err := manager.Sensors()
	if err != nil {
		return err
	}
	if index < 0 || index >= len(sensors) {
		return fmt.Errorf("index %d out of range (0-%d)", index, len(sensors)-1)
	}
	s := sensors[index]

	fmt.Printf("# %s\n", s.Name())
	fmt.Printf("# ts_ms values\n")
	return manager.Watch(s, *duration, func(ev *sensor.Event) {
		fmt.Printf("%d  %s\n", ev.Timestamp/1_000_000, eventValues(ev))
	})
}
