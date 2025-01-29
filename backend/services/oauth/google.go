package oauth

import (
	"FullStackOfYear/backend/internal/errors"
	"FullStackOfYear/backend/types/auth"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Configuration constants
const (
	defaultTimeout = 10 * time.Second
	maxRetries     = 3
	retryWaitTime  = 1 * time.Second
	maxConcurrent  = 100
)

// GoogleProvider implements the OAuth Provider interface for Google
type GoogleProvider struct {
	config    *oauth2.Config
	client    *http.Client
	semaphore chan struct{}
	metrics   *Metrics
	mu        sync.RWMutex
}

// Metrics holds provider metrics
type Metrics struct {
	RequestCount     int64
	ErrorCount       int64
	TokenExchanges   int64
	TokenValidations int64
	TokenRevocations int64
}

// NewGoogleProvider creates a new Google OAuth provider instance
func NewGoogleProvider(clientID, clientSecret, redirectURL string) Provider {
	client := &http.Client{
		Timeout: defaultTimeout,
	}

	provider := &GoogleProvider{
		config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		client:    client,
		semaphore: make(chan struct{}, maxConcurrent),
		metrics:   &Metrics{},
	}

	return provider
}

// doRequest performs an HTTP request with retry and timeout
func (p *GoogleProvider) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	// Acquire semaphore
	select {
	case p.semaphore <- struct{}{}:
		defer func() { <-p.semaphore }()
	case <-ctx.Done():
		return nil, errors.NewAppError(errors.InternalError, "Request cancelled")
	}

	// Update metrics
	p.mu.Lock()
	p.metrics.RequestCount++
	p.mu.Unlock()

	for retry := 0; retry <= maxRetries; retry++ {
		if retry > 0 {
			log.Printf("Retrying request (attempt %d/%d)", retry, maxRetries)
			time.Sleep(retryWaitTime * time.Duration(retry))
		}

		resp, err = p.client.Do(req)
		if err == nil && resp.StatusCode < 500 {
			return resp, nil
		}

		if err != nil {
			log.Printf("Request failed: %v", err)
			p.mu.Lock()
			p.metrics.ErrorCount++
			p.mu.Unlock()
		}
	}

	return nil, errors.NewAppError(errors.InternalError, "Max retries exceeded")
}

// GenerateAuthURL generates the Google OAuth authorization URL
func (p *GoogleProvider) GenerateAuthURL(state string) string {
	return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
}

// ExchangeCode exchanges authorization code for OAuth tokens
func (p *GoogleProvider) ExchangeCode(ctx context.Context, code string) (*auth.OAuthToken, error) {
	p.mu.Lock()
	p.metrics.TokenExchanges++
	p.mu.Unlock()

	token, err := p.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to exchange OAuth code")
	}

	return &auth.OAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry.Unix() - time.Now().Unix(),
		ExpiresAt:    token.Expiry,
	}, nil
}

// GetUserInfo fetches user information from Google API
func (p *GoogleProvider) GetUserInfo(ctx context.Context, token *auth.OAuthToken) (*auth.OAuthUserInfo, error) {
	client := p.config.Client(ctx, &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.ExpiresAt,
	})

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to fetch user info")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewAppError(errors.InternalError, "Failed to fetch user info")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to read response body")
	}

	var userInfo auth.OAuthUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to parse user info")
	}

	if userInfo.ID == "" || userInfo.Email == "" {
		return nil, errors.NewAppError(errors.InternalError, "Invalid user info received")
	}

	return &userInfo, nil
}

// RefreshToken refreshes an expired access token
func (p *GoogleProvider) RefreshToken(ctx context.Context, refreshToken string) (*auth.OAuthToken, error) {
	token, err := p.config.TokenSource(ctx, &oauth2.Token{
		RefreshToken: refreshToken,
	}).Token()

	if err != nil {
		return nil, errors.NewAppError(errors.InternalError, "Failed to refresh token")
	}

	return &auth.OAuthToken{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.Expiry.Unix() - time.Now().Unix(),
		ExpiresAt:    token.Expiry,
	}, nil
}

// ValidateToken validates an access token
func (p *GoogleProvider) ValidateToken(ctx context.Context, accessToken string) error {
	p.mu.Lock()
	p.metrics.TokenValidations++
	p.mu.Unlock()

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"https://www.googleapis.com/oauth2/v3/tokeninfo",
		nil,
	)
	if err != nil {
		return errors.NewAppError(errors.InternalError, "Failed to create validation request")
	}

	q := req.URL.Query()
	q.Add("access_token", accessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := p.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewAppError(errors.Unauthorized, "Invalid token")
	}

	return nil
}

// RevokeToken revokes an access token
func (p *GoogleProvider) RevokeToken(ctx context.Context, accessToken string) error {
	p.mu.Lock()
	p.metrics.TokenRevocations++
	p.mu.Unlock()

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"https://oauth2.googleapis.com/revoke",
		strings.NewReader(fmt.Sprintf("token=%s", accessToken)),
	)
	if err != nil {
		return errors.NewAppError(errors.InternalError, "Failed to create revoke request")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.NewAppError(errors.InternalError, "Failed to revoke token")
	}

	return nil
}
