package main

import (
	_ "context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	_ "time"

	"github.com/NYTimes/gziphandler"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/crypto/acme/autocert"
	"gopkg.in/yaml.v3"
)

type Route struct {
	Host   string `yaml:"host"`
	Path   string `yaml:"path"`
	Target string `yaml:"target"`
}

type TLSConfig struct {
	UseLetsEncrypt bool     `yaml:"use_lets_encrypt"`
	CacheDir       string   `yaml:"cache_dir"`
	Email          string   `yaml:"email"`
	Hosts          []string `yaml:"hosts"`
	CertFile       string   `yaml:"cert_file"`
	KeyFile        string   `yaml:"key_file"`
}

type Config struct {
	ListenHTTP  string    `yaml:"listen_http"`
	ListenHTTPS string    `yaml:"listen_https"`
	Routes      []Route   `yaml:"routes"`
	TLS         TLSConfig `yaml:"tls"`
}

var config atomic.Value // holds *Config

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func watchConfig(path string, reload func(*Config)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				if cfg, err := loadConfig(path); err == nil {
					reload(cfg)
					log.Println("Configuration reloaded.")
				} else {
					log.Println("Error while reloading configuration:", err)
				}
			}
		case err := <-watcher.Errors:
			log.Println("Error while watch file:", err)
		}
	}
}

func buildHandler(cfg *Config) http.Handler {
	proxies := make([]struct {
		host, path string
		proxy      *httputil.ReverseProxy
	}, len(cfg.Routes))

	for i, route := range cfg.Routes {
		targetURL, err := url.Parse(route.Target)
		if err != nil {
			log.Fatalf("URL error: %v", err)
		}
		proxies[i] = struct {
			host, path string
			proxy      *httputil.ReverseProxy
		}{
			host:  route.Host,
			path:  route.Path,
			proxy: httputil.NewSingleHostReverseProxy(targetURL),
		}
	}

	return gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if h, _, err := net.SplitHostPort(r.Host); err == nil {
			host = h
		}
		path := r.URL.Path

		for _, p := range proxies {
			if (p.host == "" || strings.EqualFold(p.host, host)) && strings.HasPrefix(path, p.path) {
				p.proxy.ServeHTTP(w, r)
				return
			}
		}

		http.NotFound(w, r)

	}))
}

func startServers(cfg *Config) {
	handler := buildHandler(cfg)

	// HTTP Server (ie. for redirect to HTTPS)
	if cfg.ListenHTTP != "" {
		go func() {
			log.Printf("HTTP poslouchá na %s", cfg.ListenHTTP)
			log.Fatal(http.ListenAndServe(cfg.ListenHTTP, handler))
		}()
	}

	if cfg.TLS.UseLetsEncrypt {
		m := &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			Cache:      autocert.DirCache(cfg.TLS.CacheDir),
			Email:      cfg.TLS.Email,
			HostPolicy: autocert.HostWhitelist(cfg.TLS.Hosts...),
		}

		server := &http.Server{
			Addr:      cfg.ListenHTTPS,
			Handler:   handler,
			TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		}

		log.Printf("HTTPS listening on %s (Let's Encrypt)", cfg.ListenHTTPS)
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else if cfg.TLS.CertFile != "" && cfg.TLS.KeyFile != "" {
		server := &http.Server{
			Addr:    cfg.ListenHTTPS,
			Handler: handler,
		}
		log.Printf("HTTPS listening on %s (static certificate)", cfg.ListenHTTPS)
		log.Fatal(server.ListenAndServeTLS(cfg.TLS.CertFile, cfg.TLS.KeyFile))
	}
}

func main() {
	cfgPath := "config.yaml"
	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatalf("Error while loading configuration: %v", err)
	}

	config.Store(cfg)
	go watchConfig(cfgPath, func(newCfg *Config) {
		config.Store(newCfg)
	})

	startServers(cfg)

	select {}
}
