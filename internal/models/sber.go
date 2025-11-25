package models

type SberRequest struct {
	MessageID int64  `json:"messageId"`
	SessionID string `json:"sessionId"`
	Uuid      struct {
		UserID      interface{} `json:"userId"`
		UserChannel string      `json:"userChannel"`
		Sub         string      `json:"sub"`
	} `json:"uuid"`
	Payload struct {
		Message struct {
			OriginalText string `json:"original_text"`
		} `json:"message"`
	} `json:"payload"`
}
