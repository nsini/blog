package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/nsini/blog/app/about"
	"github.com/nsini/blog/app/api"
	"github.com/nsini/blog/app/board"
	"github.com/nsini/blog/app/home"
	"github.com/nsini/blog/app/post"
	"github.com/nsini/blog/config"
	"github.com/nsini/blog/repository"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	defaultPort = "8080"
)

func main() {

	var (
		addr = envString("PORT", defaultPort)

		httpAddr   = flag.String("http.addr", ":"+addr, "HTTP listen address")
		configAddr = flag.String("config.addr", "", "config file")
		//ctx      = context.Background()
	)

	flag.Parse()

	cf := config.NewConfig(*configAddr)

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	logger = log.With(logger, "caller", log.DefaultCaller)

	db, err := repository.NewDb(logger, cf)
	if err != nil {
		_ = logger.Log("db", "connect", "err", err)
		panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	var (
		postRepository  = repository.NewPostRepository(db)
		imageRepository = repository.NewImageRepository(db)
	)

	fieldKeys := []string{"method"}

	var ps post.Service
	var aboutMe about.Service
	var homeSvc home.Service
	var apiSvc api.Service
	var boardSvc board.Service
	// post
	ps = post.NewService(logger, cf, postRepository, repository.NewUserRepository(db), imageRepository)
	ps = post.NewLoggingService(logger, ps)

	// api
	apiSvc = api.NewService(logger, cf, postRepository, repository.NewUserRepository(db), imageRepository)
	apiSvc = api.NewLoggingService(logger, apiSvc)
	apiSvc = api.NewInstrumentingService(
		prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "post_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys),
		prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "post_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys),
		apiSvc,
	)

	// home
	homeSvc = home.NewService(logger)
	homeSvc = home.NewLoggingService(logger, homeSvc)

	// about
	aboutMe = about.NewService(logger)
	aboutMe = about.NewLoggingService(logger, aboutMe)

	// board
	boardSvc = board.NewService(logger)

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()

	mux.Handle("/post", post.MakeHandler(ps, httpLogger))
	mux.Handle("/post/", post.MakeHandler(ps, httpLogger))
	mux.Handle("/about", about.MakeHandler(aboutMe, httpLogger))
	mux.Handle("/api/", api.MakeHandler(apiSvc, httpLogger))
	mux.Handle("/board", board.MakeHandler(boardSvc, httpLogger))
	mux.Handle("/", home.MakeHandler(homeSvc, httpLogger))

	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.Handle("/", accessControl(mux, logger))

	errs := make(chan error, 2)
	go func() {
		_ = logger.Log("transport", "http", "address", *httpAddr, "msg", "listening")
		errs <- http.ListenAndServe(*httpAddr, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	_ = logger.Log("terminated", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func accessControl(h http.Handler, logger log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		_ = logger.Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)

		h.ServeHTTP(w, r)
	})
}
