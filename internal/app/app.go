package app

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/passsquale/chat-server/internal/config"
	"github.com/passsquale/chat-server/internal/interceptor"
	"github.com/passsquale/chat-server/pkg/chat_v1"
	_ "github.com/passsquale/chat-server/statik" // инициализация шаблона swagger

	"github.com/passsquale/platform_common/pkg/closer"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const APISwaggerPath = "/api.swagger.json"

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run http server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run http server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {

	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, fn := range inits {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load environments: %v", err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(a.serviceProvider.AccessChecker(ctx).AccessCheck, interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	chat_v1.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatImplementation(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("Listen and serve at %s\n", a.serviceProvider.GRPCConfig().Address())

	return a.grpcServer.Serve(lis)
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := chat_v1.RegisterChatV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-type", "Content-length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           corsMiddleware.Handler(mux),
		Addr:              a.serviceProvider.HTTPConfig().Address(),
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("http listen and serve at %s\n", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFS)))
	mux.HandleFunc(APISwaggerPath, serveSwaggerFile(APISwaggerPath))

	a.swaggerServer = &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           mux,
		Addr:              a.serviceProvider.SwaggerConfig().Address(),
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFS, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("opening swagger file %s", path)

		file, err := statikFS.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			err = file.Close()
			if err != nil {
				log.Printf("error at reading swagger file: %s", path)
			}
		}()

		log.Printf("reading swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("served swagger file: %s", path)
	}
}

func (a *App) runSwaggerServer() error {
	log.Printf("swagger listen and serve at %s\n", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
