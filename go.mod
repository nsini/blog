module github.com/nsini/blog

require (
	github.com/flosch/pongo2 v0.0.0-20190505152737-8914e1cf9164
	github.com/go-kit/kit v0.8.0
	github.com/gorilla/mux v1.7.2
	github.com/jinzhu/gorm v1.9.8
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.3
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/guregu/null.v3 v3.4.0
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1

replace golang.org/x/tools v0.0.0-20181221001348-537d06c36207 => github.com/golang/tools v0.0.0-20181221001348-537d06c36207

replace golang.org/x/crypto v0.0.0-20190325154230-a5d413f7728c => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c

replace golang.org/x/sync v0.0.0-20181108010431-42b317875d0f => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f

replace golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a => github.com/golang/sync v0.0.0-20181108010431-42b317875d0f

replace golang.org/x/net v0.0.0-20181114220301-adae6a3d119a => github.com/golang/net v0.0.0-20181114220301-adae6a3d119a
