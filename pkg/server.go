package pkg

import (
	"context"
	"log"
	"net/http"
	"time"

	content2 "github.com/console-dns/server/pkg/content"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

type ConsoleServer struct {
	ctx    context.Context
	server *http.Server
}

func (r *ConsoleServer) StartAsync(notify func()) {
	go func() {
		<-r.ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(r.ctx, 10*time.Second)
		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("服务在 10s 内未退出,项目强制结束 ")
			}
		}()
		err := r.server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		// 正常结束
		cancel()
	}()
	go func() {
		err := r.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		log.Printf("HTTP Server 已正常结束.")
		notify()
	}()
}

func NewConsoleRoute(content *content2.Content, ctx context.Context) *ConsoleServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(bindValue(content2.CONTENT, content))
	r.Mount("/", Index())
	server := &http.Server{Addr: content.Config.Server.AddressPort(), Handler: r}
	return &ConsoleServer{
		ctx:    ctx,
		server: server,
	}
}

func bindValue(key string, value any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, key, value)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
