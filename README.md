# Blog Store

<a href="https://gitpod.io/#https://github.com/dracory/blogstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/blogstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/blogstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/blogstore)](https://goreportcard.com/report/github.com/dracory/blogstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/blogstore)](https://pkg.go.dev/github.com/dracory/blogstore)

Stores blog posts to a database table.

## MCP (Model Context Protocol)

Blog Store includes an MCP (Model Context Protocol) HTTP handler that allows LLM clients (for example Windsurf) to manage blog posts via JSON-RPC tools.

- The MCP handler lives in the `mcp` package
- It supports MCP JSON-RPC methods (`initialize`, `tools/list`, `tools/call`) and legacy aliases (`list_tools`, `call_tool`)
- It exposes tools such as `post_list`, `post_create`, `post_get`, `post_update`, and `post_delete`

See the detailed documentation and examples in: `mcp/README.md`

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). You can find a copy of the license at [https://www.gnu.org/licenses/agpl-3.0.en.html](https://www.gnu.org/licenses/agpl-3.0.txt)

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/dracory/blogstore
```

## Setup

```go
blogStore = blogstore.NewStore(blog.NewStoreOptions{
	DB:                 databaseInstance,
	PostTableName:     "blog_post",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
