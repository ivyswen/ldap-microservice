package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	// 可选：client 可传入 searchBase / filter 等覆盖默认配置（谨慎允许）
}

type AuthResponse struct {
	Ok     bool              `json:"ok"`
	User   map[string]string `json:"user,omitempty"`
	Error  string            `json:"error,omitempty"`
	Detail string            `json:"detail,omitempty"`
}

// POST /v1/auth
func AuthHandler(cfg *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, AuthResponse{Ok: false, Error: "invalid_json", Detail: err.Error()})
			return
		}
		if req.Username == "" || req.Password == "" {
			respondJSON(w, http.StatusBadRequest, AuthResponse{Ok: false, Error: "missing_credentials"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), cfg.RequestTimeout)
		defer cancel()

		client, err := NewLDAPClient(cfg)
		if err != nil {
			log.Error().Err(err).Msg("failed to create ldap client")
			respondJSON(w, http.StatusInternalServerError, AuthResponse{Ok: false, Error: "ldap_client_error"})
			return
		}
		defer client.Close()

		// 先尝试通过 service bind + search to get user DN (如果配置了)
		userDN, attrs, err := client.FindUserDN(ctx, req.Username)
		if err != nil {
			// 不泄露太多细节给外部
			log.Debug().Err(err).Str("user", req.Username).Msg("FindUserDN failed")
			respondJSON(w, http.StatusUnauthorized, AuthResponse{Ok: false, Error: "invalid_credentials"})
			return
		}

		// 再用用户 DN bind 校验密码
		if err := client.AuthenticateWithDN(ctx, userDN, req.Password); err != nil {
			log.Debug().Err(err).Str("userDN", userDN).Msg("user bind failed")
			respondJSON(w, http.StatusUnauthorized, AuthResponse{Ok: false, Error: "invalid_credentials"})
			return
		}

		// 成功 — 返回用户基本信息
		resp := AuthResponse{
			Ok:   true,
			User: attrs,
		}
		respondJSON(w, http.StatusOK, resp)
	}
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	// 简单 ready 检查：可以扩展为尝试连接 LDAP 服务
	respondJSON(w, http.StatusOK, map[string]string{"ready": "true"})
}

func respondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
