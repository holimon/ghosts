package main

import (
	"fmt"

	"github.com/holimon/ghosts"
)

func main() {
	domains := []string{"github.com", "assets-cdn.github.com", "codeload.github.com", "github.global.ssl.fastly.net", "global.ssl.fastly.net", "github.githubassets.com"}
	gh, err := ghosts.New(ghosts.Option{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, domain := range domains {
		if addrs, err := gh.Resolve(domain); err == nil {
			gh.Delete(domain)
			for _, addr := range addrs {
				gh.Add(domain, addr)
			}
		}
	}
	if err := gh.WriteBack(); err != nil {
		fmt.Println("Write back error", err)
	}
	if err := gh.DNSFlush(); err != nil {
		fmt.Println("DNS flush error", err)
	}
}
