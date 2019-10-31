package main

import (
	"errors"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/publicsuffix"
)

type urlChecker struct {
	timeout             time.Duration
	documentRoot        string
	excludedPattern     *regexp.Regexp
	excludePrivateHosts bool
	semaphore           semaphore
}

func newURLChecker(t time.Duration, d string, r *regexp.Regexp, excludePrivateHosts bool, s semaphore) urlChecker {
	return urlChecker{t, d, r, excludePrivateHosts, s}
}

func (c urlChecker) Check(u string, f string) error {
	u, local, err := c.resolveURL(u, f)
	if err != nil {
		return err
	}

	if !local && c.excludePrivateHosts {
		uu, _ := url.Parse(u)
		host := uu.Hostname()
		if ip := net.ParseIP(host); ip != nil {
			if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || isPrivate(ip) {
				return nil
			}
		} else if _, icann := publicsuffix.PublicSuffix(host); !icann {
			return nil // private domain
		}
	}

	if c.excludedPattern != nil && c.excludedPattern.MatchString(u) {
		return nil
	}

	if local {
		_, err := os.Stat(u)
		return err
	}

	c.semaphore.Request()
	defer c.semaphore.Release()

	if c.timeout == 0 {
		_, _, err := fasthttp.Get(nil, u)
		return err
	}

	_, _, err = fasthttp.GetTimeout(nil, u, c.timeout)
	return err
}

func (c urlChecker) CheckMany(us []string, f string, rc chan<- urlResult) {
	wg := sync.WaitGroup{}

	for _, s := range us {
		wg.Add(1)

		go func(s string) {
			rc <- urlResult{s, c.Check(s, f)}
			wg.Done()
		}(s)
	}

	wg.Wait()
	close(rc)
}

func (c urlChecker) resolveURL(u string, f string) (string, bool, error) {
	uu, err := url.Parse(u)

	if err != nil {
		return "", false, err
	}

	if uu.Scheme != "" {
		return u, false, nil
	}

	if !path.IsAbs(uu.Path) {
		return path.Join(filepath.Dir(f), uu.Path), true, nil
	}

	if c.documentRoot == "" {
		return "", false, errors.New("document root directory is not specified")
	}

	return path.Join(c.documentRoot, uu.Path), true, nil
}

// isPrivate reports whether `ip' is a local address, according to
// RFC 1918 (IPv4 addresses) and RFC 4193 (IPv6 addresses).
// xref: https://go-review.googlesource.com/c/go/+/162998/
// xref: https://github.com/golang/go/issues/29146
func isPrivate(ip net.IP) bool {
	if ip4 := ip.To4(); ip4 != nil {
		// Local IPv4 addresses are defined in https://tools.ietf.org/html/rfc1918
		return ip4[0] == 10 ||
			(ip4[0] == 172 && ip4[1]&0xf0 == 16) ||
			(ip4[0] == 192 && ip4[1] == 168)
	}
	// Local IPv6 addresses are defined in https://tools.ietf.org/html/rfc4193
	return len(ip) == net.IPv6len && ip[0]&0xfe == 0xfc
}
