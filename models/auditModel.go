package models

type AuditRecord struct {
	Action     string      `json:"Action"`
	DocumentID string      `json:"DocumentId"`
	Timestamp  interface{} `json:"Timestamp"`
	UserID     string      `json:"UserId"`
}
