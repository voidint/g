// Copyright (c) 2025 voidint <voidint@126.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"github.com/voidint/g/build"
)

func runMcpServer(_ *cli.Context) (err error) {
	transportServer := transport.NewStdioServerTransport()

	mcpServer, err := server.NewServer(transportServer, server.WithServerInfo(protocol.Implementation{
		Name:    "g",
		Version: build.ShortVersion,
	}))
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	envTool, err := protocol.NewTool("env", "Show environment variables of g", struct{}{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	cleanTool, err := protocol.NewTool("clean", "Delete the cached installation package files", struct{}{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	lsTool, err := protocol.NewTool("ls", "List installed go sdk versions", struct{}{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	lsRemoteTool, err := protocol.NewTool("ls-remote", "List remote go sdk versions available for install", struct{}{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	installTool, err := protocol.NewTool("install", "Download and install a go sdk version", InstallReq{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	uninstallTool, err := protocol.NewTool("uninstall", "Uninstall a go sdk version", UninstallReq{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	useTool, err := protocol.NewTool("use", "Switch to specified go sdk version", UseReq{})
	if err != nil {
		return cli.Exit(errstring(err), 1)
	}

	mcpServer.RegisterTool(envTool, envHandler)
	mcpServer.RegisterTool(cleanTool, cleanHandler)
	mcpServer.RegisterTool(lsTool, lsHandler)
	mcpServer.RegisterTool(lsRemoteTool, lsRemoteHandler)
	mcpServer.RegisterTool(installTool, installHandler)
	mcpServer.RegisterTool(uninstallTool, uninstallHandler)
	mcpServer.RegisterTool(useTool, useHandler)

	if err = mcpServer.Run(); err != nil {
		return cli.Exit(errstring(err), 1)
	}
	return err
}

func envHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "g", "env")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

func cleanHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "g", "clean")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

func lsHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "g", "ls", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

func lsRemoteHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	cmd := exec.CommandContext(ctx, "g", "ls-remote", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var items []struct {
		Version   string `json:"version"`
		InUse     bool   `json:"inUse"`
		Installed bool   `json:"installed"`
	}
	if err = json.Unmarshal([]byte(output), &items); err != nil {
		return nil, errors.WithStack(err)
	}

	if output, err = json.Marshal(items); err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

type InstallReq struct {
	Version string `json:"version" description:"go sdk version keywords" required:"true"`
	Nouse   bool   `json:"nouse" description:"don't use the version after installed" required:"false"`
}

func installHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var installReq InstallReq
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &installReq); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "g", "install", fmt.Sprintf("--nouse=%t", installReq.Nouse), installReq.Version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if !strings.Contains(string(output), "installed") {
		if output, err = exec.CommandContext(ctx, "go", "version").Output(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

type UninstallReq struct {
	Version string `json:"version" description:"go sdk version" required:"true"`
}

func uninstallHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var uninstallReq UninstallReq
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &uninstallReq); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "g", "uninstall", uninstallReq.Version)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}

type UseReq struct {
	Version string `json:"version" description:"go sdk version" required:"true"`
}

func useHandler(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var useReq UseReq
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &useReq); err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "g", "use", useReq.Version)
	output, err := cmd.Output()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: string(output),
			},
		},
	}, nil
}
