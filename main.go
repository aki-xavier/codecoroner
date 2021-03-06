package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"sort"
	"strings"

	"github.com/aki-xavier/codecoroner/unused"
	"golang.org/x/tools/go/buildutil"
)

func main() {
	var ignoreList string
	ucf := unused.NewCodeFinder()
	flag.BoolVar(&(ucf.Verbose), "v", false,
		"prints extra information during execution to stderr")
	flag.BoolVar(&(ucf.IncludeTests), "tests", false, "include tests in the analysis")
	flag.StringVar(&(ignoreList), "ignore", "",
		"don't read files that contain the given comma-separated strings (use to avoid /testdata, etc) ")
	// hack for testing code with build flags
	flag.Var((*buildutil.TagsFlag)(&build.Default.BuildTags), "tags", "a list of build tags")
	flag.Parse()
	// handle ignore list
	ucf.Ignore = strings.Split(ignoreList, ",")
	if len(ucf.Ignore) > 0 && ucf.Ignore[0] == "" {
		ucf.Ignore = nil
	}

	if len(flag.Args()) == 0 {
		fmt.Println("Must specify either 'funcs' or 'idents' command. Run with -help for more info.")
		os.Exit(2)
	}
	command := flag.Arg(0)
	switch command {
	case "funcs", "functions":
		ucf.Idents = false
	case "idents", "identifiers":
		ucf.Idents = true
	default:
		fmt.Println("Must specify either 'funcs' or 'idents' command. Run with -help for more info.")
		os.Exit(2)
	}

	unusedObjects, err := ucf.Run(flag.Args()[1:])
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	ucf.Logf("") // ensure a newline before printing results if -v is on

	sort.Sort(unused.ByPosition(unusedObjects))
	for _, o := range unusedObjects {
		fmt.Printf("%s\n", o)
	}
}
