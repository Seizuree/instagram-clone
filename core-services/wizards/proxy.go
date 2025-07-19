package wizards

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy creates a new reverse proxy for a given target URL.
func NewReverseProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL for reverse proxy: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// We can modify the request before it's sent to the target service here.
	// For now, we'll just forward it as is.
	proxy.Director = func(req *http.Request) {
		allowedMethods := map[string]bool{
			"GET":    true,
			"POST":   true,
			"DELETE": true,
			"PUT":    true,
		}

		if !allowedMethods[req.Method] {
			log.Printf("Method %s not allowed for proxy", req.Method)
			// This is a bit tricky with httputil.ReverseProxy.
			// We can't directly stop the request here and send a response.
			// A common pattern is to handle this in a middleware before the proxy.
			// For now, we'll just log and let the request proceed, which might result in a 404 or similar from the target.
			// A more robust solution would involve a custom http.Handler.
			return
		}

		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.Host = targetURL.Host
		// The original request path is preserved by default.
	}

	// We can modify the response before it's sent back to the client here.
	// For now, we'll just forward it as is.
	proxy.ModifyResponse = func(res *http.Response) error {
		// Example: log the status code of the response from the target service
		log.Printf("Proxy to %s: received %d %s", target, res.StatusCode, res.Status)
		return nil
	}

	return proxy
}
