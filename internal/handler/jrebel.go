package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"JrebelBrainsLicenseServerforGo/internal/model"
)

// JRebelLeasesHandler 处理JRebel租约请求
func (h *Handler) JRebelLeasesHandler(w http.ResponseWriter, r *http.Request) {
	var params url.Values

	// 根据请求方法和Content-Type判断如何解析参数
	if r.Method == http.MethodPost {
		contentType := r.Header.Get("Content-Type")

		// 检查是表单数据还是JSON数据
		if contentType == "application/x-www-form-urlencoded" ||
			contentType == "application/x-www-form-urlencoded; charset=UTF-8" {
			// 解析表单数据
			if err := r.ParseForm(); err != nil {
				h.logger.Printf("Failed to parse form: %v", err)
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}
			params = r.PostForm
		} else {
			// 尝试解析 JSON body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				h.logger.Printf("Failed to read request body: %v", err)
				http.Error(w, "Failed to read request body", http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			var req model.JRebelLeaseRequest
			if err := json.Unmarshal(body, &req); err != nil {
				h.logger.Printf("Failed to parse JSON: %v", err)
				http.Error(w, "Invalid JSON format", http.StatusBadRequest)
				return
			}

			// 将JSON数据转换为url.Values格式以统一处理
			params = url.Values{}
			params.Set("randomness", req.Randomness)
			params.Set("username", req.Username)
			params.Set("guid", req.GUID)
			params.Set("offline", fmt.Sprintf("%t", req.Offline))
			if req.ClientTime > 0 {
				params.Set("clientTime", strconv.FormatInt(req.ClientTime, 10))
			}
			if req.OfflineDays > 0 {
				params.Set("offlineDays", strconv.Itoa(req.OfflineDays))
			}
		}
	} else {
		// GET 请求，从 URL 参数获取
		params = r.URL.Query()
	}

	// 从统一的params中提取参数
	randomness := params.Get("randomness")
	username := params.Get("username")
	guid := params.Get("guid")
	offline := params.Get("offline") == "true"

	// 注意：实际的JRebel客户端可能不会发送offline参数
	// 而是根据其他条件判断，这里我们默认为false
	if params.Get("offline") == "" {
		offline = false
	}

	var clientTime int64
	if clientTimeStr := params.Get("clientTime"); clientTimeStr != "" {
		if ct, err := strconv.ParseInt(clientTimeStr, 10, 64); err == nil {
			clientTime = ct
		}
	}

	var offlineDays int
	if offlineDaysStr := params.Get("offlineDays"); offlineDaysStr != "" {
		if od, err := strconv.Atoi(offlineDaysStr); err == nil {
			offlineDays = od
		}
	}

	// 验证必需参数
	if randomness == "" || guid == "" {
		h.logger.Printf("Missing required parameters: randomness=%s, guid=%s", randomness, guid)
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// 如果没有username，使用其他字段
	if username == "" {
		username = params.Get("definedUserId")
		if username == "" {
			username = "Administrator"
		}
	}

	var validFrom, validUntil *int64

	if offline && clientTime > 0 {
		validFromVal := clientTime
		// 使用配置的离线天数，如果请求中有则使用请求的
		days := h.config.License.OfflineDays
		if offlineDays > 0 {
			days = offlineDays
		}
		validUntilVal := clientTime + int64(days)*24*60*60*1000
		validFrom = &validFromVal
		validUntil = &validUntilVal
	}

	// 生成签名
	validFromStr := "null"
	validUntilStr := "null"
	if validFrom != nil {
		validFromStr = strconv.FormatInt(*validFrom, 10)
		validUntilStr = strconv.FormatInt(*validUntil, 10)
	}

	signature, err := h.signer.SignLeaseData(
		randomness,
		h.config.License.ServerRandomness,
		guid,
		offline,
		validFromStr,
		validUntilStr,
	)
	if err != nil {
		h.logger.Printf("Failed to generate signature: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := model.JRebelLeaseResponse{
		ServerVersion:         h.config.Server.ServerVersion,
		ServerProtocolVersion: h.config.Server.ProtocolVersion,
		ServerGUID:            h.config.Server.ServerGUID,
		GroupType:             "managed",
		ID:                    1,
		LicenseType:           1,
		EvaluationLicense:     false,
		Signature:             signature,
		ServerRandomness:      h.config.License.ServerRandomness,
		SeatPoolType:          "standalone",
		StatusCode:            "SUCCESS",
		Offline:               offline,
		ValidFrom:             validFrom,
		ValidUntil:            validUntil,
		Company:               username,
		OrderID:               "",
		ZeroIds:               []string{},
		LicenseValidFrom:      1490544001000,
		LicenseValidUntil:     1691839999000,
	}

	h.respondJSON(w, http.StatusOK, response)
}

// JRebelLeases1Handler 处理JRebel租约释放
func (h *Handler) JRebelLeases1Handler(w http.ResponseWriter, r *http.Request) {
	h.logger.Printf("JRebelLeases1Handler called: Method=%s", r.Method)

	username := r.URL.Query().Get("username")

	response := model.JRebelLeaseResponse{
		ServerVersion:         h.config.Server.ServerVersion,
		ServerProtocolVersion: h.config.Server.ProtocolVersion,
		ServerGUID:            h.config.Server.ServerGUID,
		GroupType:             "managed",
		StatusCode:            "SUCCESS",
	}

	if username != "" {
		response.Company = username
	}

	h.respondJSON(w, http.StatusOK, response)
}
