package handler

import (
	"fmt"
	"net/http"
)

// PingHandler 处理ping请求
func (h *Handler) PingHandler(w http.ResponseWriter, r *http.Request) {
	salt := r.URL.Query().Get("salt")
	if salt == "" {
		http.Error(w, "Missing salt parameter", http.StatusForbidden)
		return
	}

	xmlContent := fmt.Sprintf(
		`<PingResponse><message></message><responseCode>OK</responseCode><salt>%s</salt></PingResponse>`,
		salt,
	)

	signature, err := h.signer.SignXML(xmlContent)
	if err != nil {
		h.logger.Printf("Failed to sign XML: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("<!-- %s -->\n%s", signature, xmlContent)
	h.respondXML(w, http.StatusOK, response)
}

// ObtainTicketHandler 处理获取ticket请求
func (h *Handler) ObtainTicketHandler(w http.ResponseWriter, r *http.Request) {
	salt := r.URL.Query().Get("salt")
	username := r.URL.Query().Get("userName")

	if salt == "" || username == "" {
		http.Error(w, "Missing required parameters", http.StatusForbidden)
		return
	}

	xmlContent := fmt.Sprintf(
		`<ObtainTicketResponse><message></message><prolongationPeriod>%s</prolongationPeriod><responseCode>OK</responseCode><salt>%s</salt><ticketId>1</ticketId><ticketProperties>licensee=%s	licenseType=0	</ticketProperties></ObtainTicketResponse>`,
		h.config.License.ProlongationPeriod,
		salt,
		username,
	)

	signature, err := h.signer.SignXML(xmlContent)
	if err != nil {
		h.logger.Printf("Failed to sign XML: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("<!-- %s -->\n%s", signature, xmlContent)
	h.respondXML(w, http.StatusOK, response)
}

// ReleaseTicketHandler 处理释放ticket请求
func (h *Handler) ReleaseTicketHandler(w http.ResponseWriter, r *http.Request) {
	salt := r.URL.Query().Get("salt")
	if salt == "" {
		http.Error(w, "Missing salt parameter", http.StatusForbidden)
		return
	}

	xmlContent := fmt.Sprintf(
		`<ReleaseTicketResponse><message></message><responseCode>OK</responseCode><salt>%s</salt></ReleaseTicketResponse>`,
		salt,
	)

	signature, err := h.signer.SignXML(xmlContent)
	if err != nil {
		h.logger.Printf("Failed to sign XML: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("<!-- %s -->\n%s", signature, xmlContent)
	h.respondXML(w, http.StatusOK, response)
}
