package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client представляет HTTP клиент для работы с API
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый API клиент
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// RR Intervals related types

// RRInterval represents a single R-R interval
type RRInterval struct {
	ID        string    `json:"id"`
	RRValue   int64     `json:"rr_value"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID  *string   `json:"device_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// RRIntervalsResponse represents response with R-R intervals data
type RRIntervalsResponse struct {
	RRIntervals []RRInterval `json:"rr_intervals"`
	TotalCount  int64        `json:"total_count"`
	TimeRange   TimeRange    `json:"time_range"`
}

// RRStatisticsResponse represents response with R-R intervals statistics
type RRStatisticsResponse struct {
	Summary    *RRStatisticalSummary `json:"summary"`
	Histogram  *RRHistogramData      `json:"histogram,omitempty"`
	HRVMetrics *HRVMetrics           `json:"hrv_metrics,omitempty"`
	TimeRange  TimeRange             `json:"time_range"`
}

// RRStatisticalSummary represents basic statistics
type RRStatisticalSummary struct {
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"std_dev"`
	Min    int64   `json:"min"`
	Max    int64   `json:"max"`
	Count  int64   `json:"count"`
}

// RRHistogramData represents histogram data
type RRHistogramData struct {
	Bins       []HistogramBin        `json:"bins"`
	TotalCount int64                 `json:"total_count"`
	BinWidth   int64                 `json:"bin_width"`
	Statistics *RRStatisticalSummary `json:"statistics"`
}

// HistogramBin represents one histogram bin
type HistogramBin struct {
	RangeStart int64   `json:"range_start"`
	RangeEnd   int64   `json:"range_end"`
	Count      int64   `json:"count"`
	Frequency  float64 `json:"frequency"`
}

// HRVMetrics represents HRV metrics
type HRVMetrics struct {
	RMSSD           float64 `json:"rmssd"`
	SDNN            float64 `json:"sdnn"`
	PNN50           float64 `json:"pnn50"`
	TriangularIndex float64 `json:"triangular_index"`
	TINN            float64 `json:"tinn"`
	VLFPower        float64 `json:"vlf_power"`
	LFPower         float64 `json:"lf_power"`
	HFPower         float64 `json:"hf_power"`
	LFHFRatio       float64 `json:"lf_hf_ratio"`
	TotalPower      float64 `json:"total_power"`
}

// TimeRange represents time range
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// ScatterplotResponse represents response with scatterplot data
type ScatterplotResponse struct {
	Points     []ScatterplotPoint     `json:"points"`
	TotalCount int64                  `json:"total_count"`
	Statistics *ScatterplotStatistics `json:"statistics"`
	Ellipse    *PoincarePlotEllipse   `json:"ellipse"`
}

// ScatterplotPoint represents a point in scatterplot
type ScatterplotPoint struct {
	RRn  int64 `json:"rr_n"`
	RRn1 int64 `json:"rr_n1"`
}

// ScatterplotStatistics represents scatterplot statistics
type ScatterplotStatistics struct {
	SD1         float64 `json:"sd1"`
	SD2         float64 `json:"sd2"`
	SD1SD2Ratio float64 `json:"sd1_sd2_ratio"`
	CSI         float64 `json:"csi"`
	CVI         float64 `json:"cvi"`
}

// PoincarePlotEllipse represents Poincare plot ellipse parameters
type PoincarePlotEllipse struct {
	CenterX float64 `json:"center_x"`
	CenterY float64 `json:"center_y"`
	SD1     float64 `json:"sd1"`
	SD2     float64 `json:"sd2"`
	Area    float64 `json:"area"`
}

// GetRRIntervals получает R-R интервалы
func (c *Client) GetRRIntervals(apiKey string, from, to time.Time) (*RRIntervalsResponse, error) {
	url := fmt.Sprintf("%s/v1/rr-intervals?from=%s&to=%s",
		c.baseURL, from.Format(time.RFC3339), to.Format(time.RFC3339))

	var response RRIntervalsResponse
	if err := c.makeRequest("GET", url, apiKey, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetRRStatistics получает статистику R-R интервалов
func (c *Client) GetRRStatistics(apiKey string, from, to time.Time, includeHistogram, includeHRV bool, binsCount int) (*RRStatisticsResponse, error) {
	url := fmt.Sprintf("%s/v1/rr-intervals/analytics/statistics?from=%s&to=%s&include_histogram=%t&include_hrv=%t",
		c.baseURL, from.Format(time.RFC3339), to.Format(time.RFC3339), includeHistogram, includeHRV)
	
	if binsCount > 0 {
		url += fmt.Sprintf("&bins_count=%d", binsCount)
	}

	var response RRStatisticsResponse
	if err := c.makeRequest("GET", url, apiKey, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetRRScatterplot получает данные скаттерплота R-R интервалов
func (c *Client) GetRRScatterplot(apiKey string, from, to time.Time) (*ScatterplotResponse, error) {
	url := fmt.Sprintf("%s/v1/rr-intervals/analytics/scatterplot?from=%s&to=%s",
		c.baseURL, from.Format(time.RFC3339), to.Format(time.RFC3339))

	var response ScatterplotResponse
	if err := c.makeRequest("GET", url, apiKey, nil, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// makeRequest выполняет HTTP запрос
func (c *Client) makeRequest(method, url, apiKey string, body interface{}, response interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if response != nil {
		if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
