package xiaomi

// EventType defines the conversion event types
type EventType string

const (
	EventRegister EventType = "APP_REGISTER" // 注册
	EventPay      EventType = "APP_PAY"      // 付费
	// Add more as needed
)

// ReportRequest represents the payload for event reporting
// Note: Xiaomi's exact API might use query params or JSON.
// Based on common practices and the image mentioning "v4/api/exact", it likely uses query params or a specific JSON structure.
// We will assume a flexible map or struct that can be converted to the required format.
type ReportRequest struct {
	OAID        string    `json:"oaid,omitempty"`
	IMEI        string    `json:"imei,omitempty"`
	CallbackURL string    `json:"callback_url,omitempty"` // For old API or specific flows
	EventType   EventType `json:"event_type"`
	Amount      int64     `json:"amount,omitempty"` // For payment events (in cents usually)
	Timestamp   int64     `json:"timestamp"`
	Signature   string    `json:"sign"`
}
