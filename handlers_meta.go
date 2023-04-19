package main

import (
	"fmt"
	"github.com/prasmussen/gdrive/cli"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
)

func printVersion(ctx cli.Context) {
	fmt.Printf("%s: %s\n", Name, Version)
	fmt.Printf("Golang: %s\n", runtime.Version())
	fmt.Printf("OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}

func printHelp(ctx cli.Context) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "%s usage:\n\n", Name)

	for _, h := range ctx.Handlers() {
		fmt.Fprintf(w, "%s %s\t%s\n", Name, h.Pattern, h.Description)
	}

	w.Flush()
}

func printCommandHelp(ctx cli.Context) {
	args := ctx.Args()
	printCommandPrefixHelp(ctx, args.String("command"))
}

func printSubCommandHelp(ctx cli.Context) {
	args := ctx.Args()
	printCommandPrefixHelp(ctx, args.String("command"), args.String("subcommand"))
}

func printCommandPrefixHelp(ctx cli.Context, prefix ...string) {
	handler := getHandler(ctx.Handlers(), prefix)

	if handler == nil {
		ExitF("Command not found")
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintf(w, "%s\n", handler.Description)
	fmt.Fprintf(w, "%s %s\n", Name, handler.Pattern)
	for _, group := range handler.FlagGroups {
		fmt.Fprintf(w, "\n%s:\n", group.Name)
		for _, flag := range group.Flags {
			boolFlag, isBool := flag.(cli.BoolFlag)
			if isBool && boolFlag.OmitValue {
				fmt.Fprintf(w, "  %s\t%s\n", strings.Join(flag.GetPatterns(), ", "), flag.GetDescription())
			} else {
				fmt.Fprintf(w, "  %s <%s>\t%s\n", strings.Join(flag.GetPatterns(), ", "), flag.GetName(), flag.GetDescription())
			}
		}
	}

	w.Flush()
}

func getHandler(handlers []*cli.Handler, prefix []string) *cli.Handler {
	for _, h := range handlers {
		pattern := stripOptionals(h.SplitPattern())

		if len(prefix) > len(pattern) {
			continue
		}

		if equal(prefix, pattern[:len(prefix)]) {
			return h
		}
	}

	return nil
}

// Strip optional groups (<...>) from pattern
func stripOptionals(pattern []string) []string {
	newArgs := []string{}

	for _, arg := range pattern {
		if strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]") {
			continue
		}
		newArgs = append(newArgs, arg)
	}
	return newArgs
}
