package main

import (
	"fmt"
	"strings"
    "./cli"
)

func printVersion(ctx cli.Context) {
    fmt.Printf("%s v%s\n", Name, Version)
}

func printHelp(ctx cli.Context) {
    fmt.Printf("%s usage:\n\n", Name)

    for _, h := range ctx.Handlers() {
        fmt.Printf("%s %s  (%s)\n", Name, h.Pattern, h.Description)
    }
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

    fmt.Printf("%s %s  (%s)\n", Name, handler.Pattern, handler.Description)
    for name, flags := range handler.Flags {
        fmt.Printf("\n%s:\n", name)
        for _, flag := range flags {
            boolFlag, isBool := flag.(cli.BoolFlag)
            if isBool && boolFlag.OmitValue {
                fmt.Printf("  %s  (%s)\n", strings.Join(flag.GetPatterns(), ", "), flag.GetDescription())
            } else {
                fmt.Printf("  %s <%s>  (%s)\n", strings.Join(flag.GetPatterns(), ", "), flag.GetName(), flag.GetDescription())
            }
        }
    }
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
