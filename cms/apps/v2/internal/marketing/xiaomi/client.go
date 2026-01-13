package xiaomi

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"data_backend/pkg/logger"
)

const (
	BaseURL = "https://trail.e.mi.com/v4/api/exact"
)

type Client struct {
	httpClient *http.Client
	logger     *logger.Logger
}

func NewClient(log *logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		logger:     log,
	}
}

// ReportEvent sends an event to Xiaomi
// If callbackURL is provided, it might be used directly (depending on integration mode).
// Here we implement the standard API approach.
func (c *Client) ReportEvent(req *ReportRequest) error {
	// Construct query parameters
	params := url.Values{}
	if req.OAID != "" {
		params.Add("oaid", req.OAID)
	}
	if req.IMEI != "" {
		params.Add("imei", req.IMEI)
	}
	params.Add("event_type", string(req.EventType))
	params.Add("conv_time", fmt.Sprintf("%d", req.Timestamp))

	// If it's a payment, we might need amount.
	// Note: Xiaomi API specifics for 'amount' might vary (e.g. 'purchase_amount').
	// We'll add it if > 0.
	if req.Amount > 0 {
		params.Add("purchase_amount", fmt.Sprintf("%d", req.Amount))
	}

	// TODO: Add signature generation if required by the specific client account.
	// The image mentions "sign encryption logic offline", implying a simpler or no sign for some cases,
	// or a new signature scheme. For now, we'll send the basic params.

	targetURL := BaseURL + "?" + params.Encode()

	// If a specific callback URL was stored (from the click), we might want to use that instead
	// or append to it.
	// For this implementation, we'll assume we use the standard endpoint with OAID matching.
	// If the user provided a full callback URL in the attribution phase, we could use that:
	if req.CallbackURL != "" {
		// Some flows use the callback URL directly
		targetURL = req.CallbackURL
		// We might need to append event info to it
	}

	c.logger.Infof("Reporting to Xiaomi: %s", targetURL)

	resp, err := c.httpClient.Get(targetURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("xiaomi api returned status: %d", resp.StatusCode)
	}

	return nil
}
