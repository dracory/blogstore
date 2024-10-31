# Blog Store

<a href="https://gitpod.io/#https://github.com/gouniverse/blogstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/gouniverse/blogstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/gouniverse/blogstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/blogstore)](https://goreportcard.com/report/github.com/gouniverse/blogstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/blogstore)](https://pkg.go.dev/github.com/gouniverse/blogstore)

Stores blog posts to a database table.

## License

This project is licensed under the GNU General Public License version 3 (GPL-3.0). You can find a copy of the license at https://www.gnu.org/licenses/gpl-3.0.en.html

For commercial use, please use my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

## Installation
```
go get -u github.com/gouniverse/blogstore
```

## Setup

```go
blogStore = blogstore.NewStore(blog.NewStoreOptions{
	DB:                 databaseInstance,
	PostTableName:     "blog_post",
	AutomigrateEnabled: true,
	DebugEnabled:       false,
})
```
