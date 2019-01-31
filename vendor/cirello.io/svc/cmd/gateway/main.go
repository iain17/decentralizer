// Command gateway implements an edge reverse-proxy used for serving Go packages and static content.
package main

import (
	"os"
	"sync"
)

const gatewayTokenCookie = "gateway-jwt"

var (
	publicBindIP   = os.Getenv("GATEWAY_PUBLIC_BIND_IP")
	servicesBindIP = os.Getenv("GATEWAY_SERVICES_BIND_IP")
	googleClientID = os.Getenv("GATEWAY_GOOGLE_CLIENT_ID")
	links          = os.Getenv("GATEWAY_LINKS_GIST")
	baseGithubAcct = os.Getenv("GATEWAY_BASE_GITHUB_ACCOUNT")
	frontPkgDomain = os.Getenv("GATEWAY_FRONT_PACKAGE_DOMAIN")
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		services()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		publicSites()
	}()
	wg.Wait()
}
