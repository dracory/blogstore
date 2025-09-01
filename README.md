# Blog Store

<a href="https://gitpod.io/#https://github.com/dracory/blogstore" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

[![Tests Status](https://github.com/dracory/blogstore/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/dracory/blogstore/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/blogstore)](https://goreportcard.com/report/github.com/dracory/blogstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/blogstore)](https://pkg.go.dev/github.com/dracory/blogstore)

Stores blog posts to a database table.

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
