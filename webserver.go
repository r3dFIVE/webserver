package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleSignals(signals chan os.Signal) {
	signal := <-signals
	// TODO: Handler graceful shutdown of http/https service!
	log.Printf("SIGNAL RECIEVED: %s", signal)
	os.Exit(1)
}

func serveHTTP(errs chan error, handler func(w http.ResponseWriter, r *http.Request)) {
	log.Printf("Starting http webserver on port 8080\n")
	if err := http.ListenAndServe(":8080", http.HandlerFunc(handler)); err != nil {
		log.Fatal(err)
		errs <- err
	}
}

func serveHTTPS(errs chan error, certFile string, keyFile string) {
	log.Printf("Starting https webserver on port 8443\n")
	if err := http.ListenAndServeTLS(":8443", certFile, keyFile, nil); err != nil {
		log.Fatal(err)
		errs <- err
	}
}

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to force TLS redirect. Reason: Failed to get hostname, %s", err)
		os.Exit(1)
	}
	u := r.URL
	u.Host = net.JoinHostPort(host, "443")
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
}

func Run(certFile *string, keyFile *string) chan error {
	errs := make(chan error)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	if len(*certFile) > 0 && len(*keyFile) > 0 {
		go serveHTTPS(errs, *certFile, *keyFile)
		go serveHTTP(errs, redirectTLS)
	} else {
		go serveHTTP(errs, nil)
	}

	return errs
}

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	go handleSignals(signals)

	certFile := flag.String("cert", "", "SSL Cert/Pem File")
	keyFile := flag.String("key", "", "SSL Key/Pem File")
	flag.Parse()

	errs := Run(certFile, keyFile)

	select {
	case err := <-errs:
		log.Printf("Could not start webserver service due to error: %s", err)
	}
}
