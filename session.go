package opentok

// Session represents an OpenTok session
type Session struct {
	SessionID      string `json:"session_id"`
	PartnerID      string `json:"partner_id"`
	CreatedAt      string `json:"create_dt"`
	MediaServerURL string `json:"media_server_url"`
}
