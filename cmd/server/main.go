package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"JrebelBrainsLicenseServerforGo/internal/config"
	"JrebelBrainsLicenseServerforGo/internal/crypto"
	"JrebelBrainsLicenseServerforGo/internal/handler"
)

var (
	Version   = "dev"
	BuildDate = "unknown"
	Platform  = "unknown"
)

func main() {
	// 解析命令行参数
	port := flag.String("p", "", "Server port (default: 8081)")
	keyPath := flag.String("key", "", "Path to private key file")
	debug := flag.Bool("debug", false, "Enable debug mode")
	version := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Printf("JRebel License Server\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
		fmt.Printf("Platform:   %s\n", Platform)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		os.Exit(0)
	}

	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 命令行参数覆盖配置
	if *port != "" {
		cfg.Server.Port = *port
	}
	if *keyPath != "" {
		cfg.License.PrivateKeyPath = *keyPath
	}

	// 加载私钥
	privateKey, err := config.LoadPrivateKey(cfg.License.PrivateKeyPath)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	// 创建签名器
	signer := crypto.NewSigner(privateKey)

	// 创建日志器
	logger := log.New(os.Stdout, "[LICENSE-SERVER] ", log.LstdFlags)

	// 打印启动信息
	logger.Printf("Starting JRebel License Server %s", Version)
	logger.Printf("Build Date: %s", BuildDate)
	logger.Printf("Platform: %s", Platform)

	// 创建处理器
	h := handler.NewHandler(cfg, signer, logger)

	// 设置路由
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.IndexHandler)
	mux.HandleFunc("/jrebel/leases", h.JRebelLeasesHandler)
	mux.HandleFunc("/jrebel/leases/1", h.JRebelLeases1Handler)
	mux.HandleFunc("/agent/leases", h.JRebelLeasesHandler)
	mux.HandleFunc("/agent/leases/1", h.JRebelLeases1Handler)
	mux.HandleFunc("/rpc/ping.action", h.PingHandler)
	mux.HandleFunc("/rpc/obtainTicket.action", h.ObtainTicketHandler)
	mux.HandleFunc("/rpc/releaseTicket.action", h.ReleaseTicketHandler)

	// 创建服务器
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)

	var httpHandler http.Handler = mux
	if *debug {
		httpHandler = debugMiddleware(loggingMiddleware(mux, logger), logger)
	} else {
		httpHandler = loggingMiddleware(mux, logger)
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      httpHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	go func() {
		logger.Printf("License Server started at http://localhost:%s", cfg.Server.Port)
		logger.Printf("Debug mode: %v", *debug)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exited")
}

// loggingMiddleware 日志中间件
func loggingMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 记录请求
		logger.Printf("%s %s from %s", r.Method, r.RequestURI, r.RemoteAddr)

		// 调用下一个处理器
		next.ServeHTTP(w, r)

		// 记录耗时
		logger.Printf("Completed in %v", time.Since(start))
	})
}

// debugMiddleware 调试中间件
func debugMiddleware(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 打印所有请求头
		logger.Println("=== Request Headers ===")
		for name, values := range r.Header {
			for _, value := range values {
				logger.Printf("%s: %s", name, value)
			}
		}

		// 打印查询参数
		if len(r.URL.Query()) > 0 {
			logger.Println("=== Query Parameters ===")
			for key, values := range r.URL.Query() {
				for _, value := range values {
					logger.Printf("%s: %s", key, value)
				}
			}
		}

		// 如果是POST请求，打印body
		if r.Method == http.MethodPost && r.Body != nil {
			body, err := io.ReadAll(r.Body)
			r.Body.Close()

			if err == nil && len(body) > 0 {
				logger.Printf("=== Request Body ===\n%s", string(body))
				// 使用 bytes.NewReader 重新创建 body 供后续使用
				r.Body = io.NopCloser(bytes.NewReader(body))
			} else {
				// 如果读取失败或body为空，创建空的reader
				r.Body = io.NopCloser(bytes.NewReader([]byte{}))
			}
		}

		logger.Println("===================")

		next.ServeHTTP(w, r)
	})
}
