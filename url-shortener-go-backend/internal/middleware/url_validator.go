package middleware

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type URLValidator struct {
	whitelist      map[string]bool
	blacklist      map[string]bool
	blacklistRegex []*regexp.Regexp
	mu             sync.RWMutex
	config         *URLValidatorConfig
}

type URLValidatorConfig struct {
	AllowPrivateIPs   bool
	AllowLocalhost    bool
	RequireHTTPS      bool
	MaxRedirects      int
	ValidateSSL       bool
	AllowedDomains    []string
	BlockedDomains    []string
	BlockedPatterns   []string
	AllowedProtocols  []string
	BlockedExtensions []string
	MaxURLLength      int
}

func NewURLValidator(config *URLValidatorConfig) *URLValidator {
	if config == nil {
		config = DefaultConfig()
	}

	v := &URLValidator{
		whitelist: make(map[string]bool),
		blacklist: make(map[string]bool),
		config:    config,
	}

	for _, domain := range config.AllowedDomains {
		v.whitelist[strings.ToLower(domain)] = true
	}

	for _, domain := range config.BlockedDomains {
		v.blacklist[strings.ToLower(domain)] = true
	}

	for _, pattern := range config.BlockedPatterns {
		if re, err := regexp.Compile(pattern); err == nil {
			v.blacklistRegex = append(v.blacklistRegex, re)
		}
	}

	return v
}

func DefaultConfig() *URLValidatorConfig {
	return &URLValidatorConfig{
		AllowPrivateIPs:  false,
		AllowLocalhost:   false,
		RequireHTTPS:     false,
		MaxRedirects:     5,
		ValidateSSL:      true,
		AllowedProtocols: []string{"http", "https"},
		BlockedDomains: []string{
			"bit.ly",
			"tinyurl.com",
			"goo.gl",
		},
		BlockedPatterns: []string{
			`.*\.tk$`,
			`.*\.ml$`,
			`.*\.ga$`,
			`.*\.cf$`,
			`.*phishing.*`,
			`.*malware.*`,
			`.*virus.*`,
			`.*hack.*`,
			`[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}`,
		},
		BlockedExtensions: []string{
			".exe", ".bat", ".cmd", ".com", ".pif", ".scr",
			".app", ".dmg",
			".deb", ".rpm",
			".jar", ".jnlp",
			".swf",
			".ps1", ".vbs", ".js", ".jse", ".vbe",
			".msi", ".msp",
			".zip", ".rar", ".7z",
		},
		MaxURLLength: 2048,
	}
}

func (v *URLValidator) ValidateURL(rawURL string) error {
	if v.config.MaxURLLength > 0 && len(rawURL) > v.config.MaxURLLength {
		return fmt.Errorf("URL exceeds maximum length of %d characters", v.config.MaxURLLength)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if err := v.validateProtocol(parsedURL); err != nil {
		return err
	}

	if err := v.checkBlacklistPatterns(rawURL); err != nil {
		return err
	}

	if err := v.validateDomain(parsedURL); err != nil {
		return err
	}

	if err := v.checkBlockedExtensions(parsedURL.Path); err != nil {
		return err
	}

	if err := v.validateIPRestrictions(parsedURL.Host); err != nil {
		return err
	}

	if err := v.checkURLShortenerChain(parsedURL.Host); err != nil {
		return err
	}

	if err := v.performSecurityChecks(parsedURL); err != nil {
		return err
	}

	return nil
}

func (v *URLValidator) validateProtocol(u *url.URL) error {
	if v.config.RequireHTTPS && u.Scheme != "https" {
		return fmt.Errorf("HTTPS is required")
	}

	if len(v.config.AllowedProtocols) > 0 {
		allowed := false
		for _, protocol := range v.config.AllowedProtocols {
			if u.Scheme == protocol {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("protocol '%s' is not allowed", u.Scheme)
		}
	}

	return nil
}

func (v *URLValidator) validateDomain(u *url.URL) error {
	domain := strings.ToLower(u.Hostname())

	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.blacklist[domain] {
		return fmt.Errorf("domain '%s' is blacklisted", domain)
	}

	for blacklisted := range v.blacklist {
		if strings.Contains(domain, blacklisted) {
			return fmt.Errorf("domain contains blacklisted domain '%s'", blacklisted)
		}
	}

	if len(v.whitelist) > 0 {
		if !v.isWhitelisted(domain) {
			return fmt.Errorf("domain '%s' is not whitelisted", domain)
		}
	}

	return nil
}

func (v *URLValidator) isWhitelisted(domain string) bool {
	if v.whitelist[domain] {
		return true
	}

	for whitelisted := range v.whitelist {
		if strings.HasSuffix(domain, "."+whitelisted) {
			return true
		}
	}

	return false
}

func (v *URLValidator) checkBlacklistPatterns(rawURL string) error {
	lowercaseURL := strings.ToLower(rawURL)

	for _, re := range v.blacklistRegex {
		if re.MatchString(lowercaseURL) {
			return fmt.Errorf("URL matches blacklisted pattern")
		}
	}

	return nil
}

func (v *URLValidator) checkBlockedExtensions(path string) error {
	lowercasePath := strings.ToLower(path)

	for _, ext := range v.config.BlockedExtensions {
		if strings.HasSuffix(lowercasePath, ext) {
			return fmt.Errorf("file extension '%s' is not allowed", ext)
		}
	}

	return nil
}

func (v *URLValidator) validateIPRestrictions(host string) error {
	hostname, _, _ := net.SplitHostPort(host)
	if hostname == "" {
		hostname = host
	}

	ip := net.ParseIP(hostname)
	if ip == nil {
		return nil
	}

	if !v.config.AllowLocalhost && (ip.IsLoopback() || hostname == "localhost") {
		return fmt.Errorf("localhost URLs are not allowed")
	}

	if !v.config.AllowPrivateIPs && isPrivateIP(ip) {
		return fmt.Errorf("private IP addresses are not allowed")
	}

	if !v.config.AllowPrivateIPs && !v.config.AllowLocalhost {
		return fmt.Errorf("direct IP addresses are not allowed")
	}

	return nil
}

func (v *URLValidator) checkURLShortenerChain(host string) error {
	knownShorteners := []string{
		"bit.ly", "tinyurl.com", "goo.gl", "ow.ly", "is.gd",
		"buff.ly", "adf.ly", "bit.do", "mcaf.ee", "su.pr",
	}

	for _, shortener := range knownShorteners {
		if strings.Contains(strings.ToLower(host), shortener) {
			return fmt.Errorf("URL shortener chains are not allowed")
		}
	}

	return nil
}

func (v *URLValidator) performSecurityChecks(u *url.URL) error {
	suspiciousParams := []string{
		"redirect", "url", "next", "continue", "return",
		"goto", "target", "dest", "destination",
	}

	query := u.Query()
	for _, param := range suspiciousParams {
		if query.Get(param) != "" {
			return fmt.Errorf("suspicious redirect parameter detected")
		}
	}

	if u.Scheme == "data" || u.Scheme == "javascript" {
		return fmt.Errorf("data and javascript URLs are not allowed")
	}

	if strings.Contains(u.String(), "%00") || strings.Contains(u.String(), "\x00") {
		return fmt.Errorf("null bytes detected in URL")
	}

	if strings.Count(u.String(), "%25") > 2 {
		return fmt.Errorf("multiple URL encoding detected")
	}

	return nil
}

func (v *URLValidator) AddToWhitelist(domain string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.whitelist[strings.ToLower(domain)] = true
}

func (v *URLValidator) RemoveFromWhitelist(domain string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.whitelist, strings.ToLower(domain))
}

func (v *URLValidator) AddToBlacklist(domain string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.blacklist[strings.ToLower(domain)] = true
}

func (v *URLValidator) RemoveFromBlacklist(domain string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	delete(v.blacklist, strings.ToLower(domain))
}

func isPrivateIP(ip net.IP) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",
		"fe80::/10",
	}

	for _, cidr := range privateRanges {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
