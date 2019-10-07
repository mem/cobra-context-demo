package cmd

import (
	"context"
	"time"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "demo",
		Short: "demo shows how to pass a context.Context to a cobra command",
		Long:  "A simple demo to show how to pass a context.Context to a cobra command",
		Run:   runWithoutContext,
	}

	contextAdder ctxAdder
)

func init() {
	// This is a common pattern for structuring programs using the
	// flag or cobra packages.
	cmd1 := &cobra.Command{
		Use: "cmd1",
		Run: addContext(context.Background(), runWithContext),
	}

	cmd2 := &cobra.Command{
		Use: "cmd2",
		Run: contextAdder.withContext(runWithContext),
	}

	rootCmd.AddCommand(cmd1)

	rootCmd.AddCommand(cmd2)
}

func runWithoutContext(cmd *cobra.Command, args []string) {
	cmd.Printf("called as: %s\n", cmd.CalledAs())
	cmd.Printf("name: %s\n", cmd.Name())
	cmd.Printf("args: %v\n", args)
}

func runWithContext(ctx context.Context, cmd *cobra.Command, args []string) {
	cmd.Printf("called as: %s\n", cmd.CalledAs())
	cmd.Printf("name: %s\n", cmd.Name())
	cmd.Printf("ctx: %s\n", ctx)
	cmd.Printf("args: %v\n", args)

	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for n := 0; n < 10; n++ {
		cmd.Println("Working...")

		select {
		case <-tick.C:

		case <-ctx.Done():
			cmd.Println("Context done")
			return
		}
	}
}

type commandWithContext func(context.Context, *cobra.Command, []string)

func addContext(ctx context.Context, fn commandWithContext) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(ctx, cmd, args)
	}
}

type ctxAdder struct {
	ctx context.Context
}

func (c *ctxAdder) setContext(ctx context.Context) {
	c.ctx = ctx
}

func (c *ctxAdder) withContext(fn commandWithContext) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fn(c.ctx, cmd, args)
	}
}

func Execute(ctx context.Context) error {
	contextAdder.setContext(ctx)

	return rootCmd.Execute()
}
