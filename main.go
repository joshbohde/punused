package main

import (
	"context"
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"strings"

	"github.com/bep/punused/internal/lib"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	wd         = flag.String("wd", "", "working directory for the project")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
	}

	// Set a default working directory
	if *wd == "" {
		*wd, _ = os.Getwd()
	}

	// Default to "every go file in the workspace".
	pattern := "**/*.go"

	args := flag.Args()

	if len(args) > 0 {
		pattern = strings.TrimPrefix(args[0], "./")
	}

	ctx := context.Background()

	err := lib.Run(
		ctx,
		lib.RunConfig{
			WorkspaceDir:    *wd,
			FilenamePattern: pattern,
			Out:             os.Stdout,
		},
	)

	pprof.StopCPUProfile()

	if err != nil {
		log.Fatal(err)
	}
}
