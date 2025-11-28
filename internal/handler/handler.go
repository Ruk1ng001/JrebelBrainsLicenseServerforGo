package handler

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"JrebelBrainsLicenseServerforGo/internal/config"
	"JrebelBrainsLicenseServerforGo/internal/crypto"

	"github.com/google/uuid"
)

//go:embed templates/*
var templatesFS embed.FS

type Handler struct {
	config   *config.Config
	signer   *crypto.Signer
	logger   *log.Logger
	template *template.Template
}

// NewHandler 创建处理器
func NewHandler(cfg *config.Config, signer *crypto.Signer, logger *log.Logger) *Handler {
	// 解析模板
	tmpl, err := template.ParseFS(templatesFS, "templates/*.html")
	if err != nil {
		logger.Printf("Warning: Failed to parse templates: %v", err)
	}

	return &Handler{
		config:   cfg,
		signer:   signer,
		logger:   logger,
		template: tmpl,
	}
}

// IndexData 首页数据
type IndexData struct {
	ServerURL       string
	ServerVersion   string
	ProtocolVersion string
	DockerImage     string
}

// IndexHandler 处理首页请求
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// 构建服务器 URL
	serverURL := h.getServerURL(r)

	// 准备模板数据
	data := IndexData{
		ServerURL:       serverURL,
		ServerVersion:   h.config.Server.ServerVersion,
		ProtocolVersion: h.config.Server.ProtocolVersion,
		DockerImage:     h.config.Web.DockerImage,
	}

	// 渲染模板
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.template.ExecuteTemplate(w, "index.html", data); err != nil {
		h.logger.Printf("Failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// getServerURL 获取服务器 URL
func (h *Handler) getServerURL(r *http.Request) string {
	// 如果配置了 BASE_URL，直接使用
	if h.config.Web.BaseURL != "" {
		return h.config.Web.BaseURL
	}

	// 否则自动检测
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	// 检查 X-Forwarded-Proto 头（用于反向代理）
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	// 检查 X-Forwarded-Host 头（用于反向代理）
	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	}

	return fmt.Sprintf("%s://%s", scheme, host)
}

// GenerateUUIDHandler 生成 UUID API
func (h *Handler) GenerateUUIDHandler(w http.ResponseWriter, r *http.Request) {
	newUUID := uuid.New().String()
	serverURL := h.getServerURL(r)

	response := map[string]string{
		"uuid":          newUUID,
		"activationUrl": fmt.Sprintf("%s/%s", serverURL, newUUID),
		"serverUrl":     serverURL,
	}

	h.respondJSON(w, http.StatusOK, response)
}

// respondJSON 返回JSON响应
func (h *Handler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Printf("Failed to encode JSON response: %v", err)
	}
}

// respondXML 返回XML响应
func (h *Handler) respondXML(w http.ResponseWriter, status int, data string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(data))
}
