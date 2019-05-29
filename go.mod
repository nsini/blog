module github.com/nsini/blog

require (
	github.com/flosch/pongo2 v0.0.0-20190505152737-8914e1cf9164
	github.com/go-kit/kit v0.8.0
	github.com/gorilla/mux v1.7.2
	github.com/jinzhu/gorm v1.9.8
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/pkg/errors v0.8.0
	github.com/prometheus/client_golang v0.9.3 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	gopkg.in/guregu/null.v3 v3.4.0
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
