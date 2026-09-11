package console

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/larsonzh/prfrail/internal/adapters"
)

type secretStatus struct {
	Reference string `json:"reference"`
	Exists    bool   `json:"exists"`
}

func (cli CLI) executeSecret(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "usage: prfrail secret <set|status|delete> --ref <reference> [--json]")
		return exitUsage
	}
	switch args[0] {
	case "set":
		return cli.executeSecretSet(ctx, args[1:], stdout, stderr)
	case "status":
		return cli.executeSecretStatus(ctx, args[1:], stdout, stderr)
	case "delete":
		return cli.executeSecretDelete(ctx, args[1:], stdout, stderr)
	default:
		fmt.Fprintln(stderr, "usage: prfrail secret <set|status|delete> --ref <reference> [--json]")
		return exitUsage
	}
}

func (cli CLI) executeSecretSet(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	const usage = "usage: prfrail secret set --ref <reference> [--json]"
	set, parseErr := newFlagSet("secret set")
	reference := set.String("ref", "", "opaque secret reference")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		return writeUsageError(stdout, stderr, "secret set", *jsonOutput, parseErr.String(), usage)
	}
	if len(set.Args()) != 0 || *reference == "" {
		return writeUsageError(stdout, stderr, "secret set", *jsonOutput, "--ref is required and positional values are forbidden", usage)
	}
	if cli.secretManager == nil || cli.readSecret == nil {
		return writeCommandError(stdout, stderr, "secret set", *jsonOutput, adapters.ErrAISecretStore)
	}
	secret, err := cli.readSecret()
	if err != nil {
		return writeCommandError(stdout, stderr, "secret set", *jsonOutput, err)
	}
	defer clear(secret)
	if len(secret) == 0 {
		return writeCommandError(stdout, stderr, "secret set", *jsonOutput, errors.New("empty secret is not allowed"))
	}
	if err := cli.secretManager.Set(ctx, *reference, secret); err != nil {
		return writeCommandError(stdout, stderr, "secret set", *jsonOutput, err)
	}
	data := secretStatus{Reference: *reference, Exists: true}
	if *jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "secret set", OK: true, ExitCode: exitSuccess, Data: data})
	}
	fmt.Fprintf(stdout, "secret set: %s\n", *reference)
	return exitSuccess
}

func (cli CLI) executeSecretStatus(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	const usage = "usage: prfrail secret status --ref <reference> [--json]"
	reference, jsonOutput, ok := parseSecretReference(args, "secret status", usage, stdout, stderr)
	if !ok {
		return exitUsage
	}
	if cli.secretManager == nil {
		return writeCommandError(stdout, stderr, "secret status", jsonOutput, adapters.ErrAISecretStore)
	}
	exists, err := cli.secretManager.Exists(ctx, reference)
	if err != nil {
		return writeCommandError(stdout, stderr, "secret status", jsonOutput, err)
	}
	data := secretStatus{Reference: reference, Exists: exists}
	if jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "secret status", OK: true, ExitCode: exitSuccess, Data: data})
	}
	fmt.Fprintf(stdout, "secret status: %s exists=%t\n", reference, exists)
	return exitSuccess
}

func (cli CLI) executeSecretDelete(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	const usage = "usage: prfrail secret delete --ref <reference> [--json]"
	reference, jsonOutput, ok := parseSecretReference(args, "secret delete", usage, stdout, stderr)
	if !ok {
		return exitUsage
	}
	if cli.secretManager == nil {
		return writeCommandError(stdout, stderr, "secret delete", jsonOutput, adapters.ErrAISecretStore)
	}
	if err := cli.secretManager.Delete(ctx, reference); err != nil && !errors.Is(err, adapters.ErrAISecretNotFound) {
		return writeCommandError(stdout, stderr, "secret delete", jsonOutput, err)
	}
	data := secretStatus{Reference: reference, Exists: false}
	if jsonOutput {
		return writeJSONResponse(stdout, commandResponse{Command: "secret delete", OK: true, ExitCode: exitSuccess, Data: data})
	}
	fmt.Fprintf(stdout, "secret delete: %s\n", reference)
	return exitSuccess
}

func parseSecretReference(args []string, command, usage string, stdout, stderr io.Writer) (string, bool, bool) {
	set, parseErr := newFlagSet(command)
	reference := set.String("ref", "", "opaque secret reference")
	jsonOutput := set.Bool("json", false, "output as JSON")
	if err := set.Parse(args); err != nil {
		writeUsageError(stdout, stderr, command, *jsonOutput, parseErr.String(), usage)
		return "", *jsonOutput, false
	}
	if len(set.Args()) != 0 || *reference == "" {
		writeUsageError(stdout, stderr, command, *jsonOutput, "--ref is required and positional values are forbidden", usage)
		return "", *jsonOutput, false
	}
	return *reference, *jsonOutput, true
}
