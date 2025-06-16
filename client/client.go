package pclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Inst       Instance
	Interval   time.Duration

	cancelHB context.CancelFunc
	mu       sync.Mutex
}

func NewClient(baseURL string, inst Instance, hbInterval time.Duration) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		Inst:       inst,
		Interval:   hbInterval,
	}
}

func (c *Client) Register(ctx context.Context) error {
	body, _ := json.Marshal(c.Inst)
	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/register", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("register failed: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Deregister(ctx context.Context) error {
	payload := map[string]string{
		"serviceName": c.Inst.ServiceName,
		"instanceID":  c.Inst.InstanceID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "DELETE", c.BaseURL+"/api/deregister", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("deregister failed: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) Heartbeat(ctx context.Context) error {
	payload := map[string]string{
		"serviceName": c.Inst.ServiceName,
		"instanceID":  c.Inst.InstanceID,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "PUT", c.BaseURL+"/api/heartbeat", bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("heartbeat failed: %d", resp.StatusCode)
	}
	return nil
}

func (c *Client) StartHeartbeat() (stop func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cancelHB != nil {
		return c.cancelHB
	}

	ctx, cancel := context.WithCancel(context.Background())
	c.cancelHB = cancel

	ticker := time.NewTicker(c.Interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := c.Heartbeat(ctx); err != nil {
					fmt.Printf("heartbeat error: %v\n", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				_ = c.Deregister(context.Background())
				return
			}
		}
	}()

	return cancel
}
