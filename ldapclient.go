package main

import (
	"context"
	"crypto/tls"
	"fmt"

	ldap "github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog/log"
)

type LDAPClient struct {
	cfg  *Config
	conn *ldap.Conn
}

func NewLDAPClient(cfg *Config) (*LDAPClient, error) {
	// Dial each request to keep isolation and avoid pool complexity;
	// 如果需要高性能可实现连接池复用。
	address := cfg.LDAPURL

	// Use goroutine + channel to implement connection timeout
	type dialResult struct {
		conn *ldap.Conn
		err  error
	}
	dialCh := make(chan dialResult, 1)

	go func() {
		var l *ldap.Conn
		var err error
		if cfg.UseStartTLS {
			// plain dial then starttls
			l, err = ldap.DialURL(address)
			if err != nil {
				dialCh <- dialResult{nil, err}
				return
			}
			if err = l.StartTLS(&tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}); err != nil {
				l.Close()
				dialCh <- dialResult{nil, err}
				return
			}
		} else if cfg.UseLDAPS {
			tlsCfg := &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}
			l, err = ldap.DialURL(address, ldap.DialWithTLSConfig(tlsCfg))
			if err != nil {
				dialCh <- dialResult{nil, err}
				return
			}
		} else {
			l, err = ldap.DialURL(address)
			if err != nil {
				dialCh <- dialResult{nil, err}
				return
			}
		}
		dialCh <- dialResult{l, nil}
	}()

	// Wait for connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnTimeout)
	defer cancel()

	var l *ldap.Conn
	select {
	case <-ctx.Done():
		log.Error().Str("address", address).Dur("timeout", cfg.ConnTimeout).Msg("LDAP connection timeout")
		return nil, NewLDAPErrorWithCause(ErrConnectionTimeout, "connection timeout", ctx.Err())
	case result := <-dialCh:
		if result.err != nil {
			if cfg.UseStartTLS {
				log.Error().Err(result.err).Str("address", address).Msg("failed to dial LDAP server for StartTLS")
				return nil, NewLDAPErrorWithCause(ErrConnectionFailed, "failed to dial LDAP server", result.err)
			} else if cfg.UseLDAPS {
				log.Error().Err(result.err).Str("address", address).Msg("failed to dial LDAPS server")
				return nil, NewLDAPErrorWithCause(ErrConnectionFailed, "failed to dial LDAPS server", result.err)
			} else {
				log.Error().Err(result.err).Str("address", address).Msg("failed to dial LDAP server")
				return nil, NewLDAPErrorWithCause(ErrConnectionFailed, "failed to dial LDAP server", result.err)
			}
		}
		l = result.conn
	}

	// Optional: service bind for search operations
	if cfg.BindDN != "" && cfg.BindPassword != "" {
		if err := l.Bind(cfg.BindDN, cfg.BindPassword); err != nil {
			l.Close()
			log.Error().Err(err).Str("bindDN", cfg.BindDN).Msg("failed to bind with service account")
			return nil, NewLDAPErrorWithCause(ErrBindFailed, "failed to bind with service account", err)
		}
		log.Debug().Str("bindDN", cfg.BindDN).Msg("service account bind successful")
	}

	if cfg.UseStartTLS {
		log.Debug().Str("address", address).Msg("LDAP connection established with StartTLS")
	} else if cfg.UseLDAPS {
		log.Debug().Str("address", address).Msg("LDAP connection established with LDAPS")
	} else {
		log.Debug().Str("address", address).Msg("LDAP connection established")
	}

	return &LDAPClient{cfg: cfg, conn: l}, nil
}

func (c *LDAPClient) Close() {
	if c.conn != nil {
		_ = c.conn.Close()
	}
}

// FindUserDN uses configured searchBase & filter to find the user's DN and attributes
func (c *LDAPClient) FindUserDN(ctx context.Context, username string) (string, map[string]string, error) {
	// prepare filter
	filter := fmt.Sprintf(c.cfg.UserSearchFilter, ldap.EscapeFilter(username))
	searchReq := ldap.NewSearchRequest(
		c.cfg.UserSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 1, int(c.cfg.RequestTimeout.Seconds()), false,
		filter,
		c.cfg.ReturnAttributes,
		nil,
	)
	// Use Context deadline via goroutine and channel because go-ldap doesn't accept context directly
	type result struct {
		res *ldap.SearchResult
		err error
	}
	ch := make(chan result, 1)
	go func() {
		res, err := c.conn.Search(searchReq)
		ch <- result{res: res, err: err}
	}()

	select {
	case <-ctx.Done():
		log.Warn().Str("username", username).Msg("user search timeout")
		return "", nil, NewLDAPErrorWithCause(ErrSearchTimeout, "user search timeout", ctx.Err())
	case r := <-ch:
		if r.err != nil {
			log.Error().Err(r.err).Str("username", username).Msg("user search failed")
			return "", nil, NewLDAPErrorWithCause(ErrSearchFailed, "user search failed", r.err)
		}
		if len(r.res.Entries) == 0 {
			log.Debug().Str("username", username).Msg("user not found in LDAP")
			return "", nil, NewLDAPError(ErrUserNotFound, "user not found")
		}
		ent := r.res.Entries[0]
		attrs := map[string]string{}
		for _, a := range c.cfg.ReturnAttributes {
			if len(ent.GetAttributeValues(a)) > 0 {
				attrs[a] = ent.GetAttributeValue(a)
			}
		}
		// If cfg.UserDNAttr is set, prefer it; otherwise use entry.DN
		dn := ent.DN
		if c.cfg.UserDNAttr != "" {
			if v := ent.GetAttributeValue(c.cfg.UserDNAttr); v != "" {
				dn = v
			}
		}
		log.Debug().Str("dn", dn).Msg("user found")
		return dn, attrs, nil
	}
}

// AuthenticateWithDN attempts bind using user DN + password
func (c *LDAPClient) AuthenticateWithDN(ctx context.Context, userDN, password string) error {
	type res struct{ err error }
	ch := make(chan res, 1)
	go func() {
		err := c.conn.Bind(userDN, password)
		ch <- res{err: err}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case r := <-ch:
		return r.err
	}
}
