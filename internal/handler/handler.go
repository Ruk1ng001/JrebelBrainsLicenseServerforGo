package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"JrebelBrainsLicenseServerforGo/internal/config"
	"JrebelBrainsLicenseServerforGo/internal/crypto"

	"github.com/google/uuid"
)

type Handler struct {
	config *config.Config
	signer *crypto.Signer
	logger *log.Logger
}

// NewHandler 创建处理器
func NewHandler(cfg *config.Config, signer *crypto.Signer, logger *log.Logger) *Handler {
	return &Handler{
		config: cfg,
		signer: signer,
		logger: logger,
	}
}

// IndexHandler 处理首页请求
func (h *Handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	licenseURL := fmt.Sprintf("%s://%s", scheme, r.Host)
	exampleGUID := uuid.New().String()

	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>JRebel & JetBrains License Server</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; line-height: 1.6; }
        h1 { color: #333; }
        h3 { color: #666; border-bottom: 2px solid #ddd; padding-bottom: 10px; }
        .url { color: #d14; font-weight: bold; }
        hr { margin: 30px 0; }
        .section { margin: 20px 0; }
    </style>
</head>
<body>
    <h3>使用说明（Instructions for use）</h3>
    <hr/>
    
    <div class="section">
        <h1>Hello, This is a JRebel & JetBrains License Server!</h1>
        <p>License Server started at <span class="url">%s</span></p>
        <p>JetBrains Activation address: <span class="url">%s/</span></p>
        <p>JRebel 7.1 and earlier version Activation address: <span class="url">%s/{tokenname}</span>, with any email.</p>
        <p>JRebel 2018.1 and later version Activation address: <span class="url">%s/{guid}</span></p>
        <p>Example: <span class="url">%s/%s</span>, with any email.</p>
    </div>
    
    <hr/>
    
    <div class="section">
        <h1>你好，此地址是 JRebel & JetBrains License Server!</h1>
        <p>许可服务器启动于 <span class="url">%s</span></p>
        <p>JetBrains激活地址: <span class="url">%s/</span></p>
        <p>JRebel 7.1 及旧版本激活地址: <span class="url">%s/{tokenname}</span>, 以及任意邮箱地址。</p>
        <p>JRebel 2018.1+ 版本激活地址: <span class="url">%s/{guid}</span></p>
        <p>例如: <span class="url">%s/%s</span>, 以及任意邮箱地址。</p>
    </div>
</body>
</html>
	`, licenseURL, licenseURL, licenseURL, licenseURL, licenseURL, exampleGUID,
		licenseURL, licenseURL, licenseURL, licenseURL, licenseURL, exampleGUID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
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
