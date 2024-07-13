package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Simulator struct {
	input   string
	plotter Plotter
	events  map[int][]Event
}

// The simulation reads an input file that describes the sequence of
// operations for each process, including local events, send, and receive operations.
func NewSimulator(input string) *Simulator {
	return &Simulator{
		input,
		Plotter{},
		make(map[int][]Event),
	}
}

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func parseProcessSequence(procSequence string) [][]string {
	procRows := strings.Split(procSequence, "\n")
	re := regexp.MustCompile(`p\d+: (.*)`)

	var processesOps [][]string
	for _, row := range procRows {
		matches := re.FindStringSubmatch(row)
		if len(matches) < 2 {
			continue
		}
		ops := strings.Fields(matches[1])
		processesOps = append(processesOps, ops)
	}
	return processesOps
}

func handleOperation(proc *Process, op string, processes []*Process) error {
	switch {
	case op == "l":
		proc.Local()
	case strings.HasPrefix(op, "s"):
		dst, err := strconv.Atoi(op[1:])
		if err != nil {
			return fmt.Errorf("error converting send destination: %v", err)
		}
		proc.Send(processes[dst])
	case strings.HasPrefix(op, "r"):
		src, err := strconv.Atoi(op[1:])
		if err != nil {
			return fmt.Errorf("error converting receive source: %v", err)
		}
		proc.Recv(src)
	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
	return nil
}

func (s *Simulator) Simulate() error {
	processSequence, err := readFile(s.input)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}
	processOps := parseProcessSequence(processSequence)

	numProcesses := len(processOps)
	processes := make([]*Process, numProcesses)
	for i := 0; i < numProcesses; i++ {
		processes[i] = NewProcess(i, numProcesses)
	}

	var wg sync.WaitGroup
	wg.Add(numProcesses)

	for i, ops := range processOps {
		go func(procID int, operations []string) {
			defer wg.Done()

			for _, op := range operations {
				handleOperation(processes[procID], op, processes)
			}
		}(i, ops)
	}
	wg.Wait()

	for _, process := range processes {
		s.events[process.ID] = process.events
	}
	return nil
}

// The resulting events are plotted on a space-time diagram to visualize
// the causality between events in different processes.
func (s *Simulator) Draw(outputFile string) error {
	return s.plotter.DrawSpaceTimeDiagram(s.events, outputFile)
}
