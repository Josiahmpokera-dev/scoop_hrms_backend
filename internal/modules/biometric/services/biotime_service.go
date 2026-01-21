package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Josiahmpokera-dev/hrms-backend/internal/config"
	biometricRepos "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/biometric/repositories"
	"gorm.io/gorm"
)

// BioTimeService handles BioTime API integration
type BioTimeService struct {
	baseURL    string
	username   string
	password   string
	enabled    bool
	client     *http.Client
	configRepo *biometricRepos.BioTimeConfigRepository
	tenantID   *uint
	queue      QueuePublisher
}

// QueuePublisher interface for publishing messages to queue
type QueuePublisher interface {
	PublishTransaction(tenantID *uint, transactions []interface{}) error
}

// NewBioTimeService creates a new BioTime service
func NewBioTimeService(tenantID *uint) *BioTimeService {
	cfg := config.AppConfig
	var queue QueuePublisher
	
	// Initialize queue if enabled
	if cfg != nil && cfg.RabbitMQ.Enabled {
		if q, err := NewQueueClient(); err == nil {
			queue = q
		}
	}

	if cfg == nil {
		// Fallback if config not loaded
		return &BioTimeService{
			baseURL:  "http://10.4.9.24:8087",
			username: "Developer",
			password: "Developer@123",
			enabled:  true,
			client:   &http.Client{Timeout: 30 * time.Second},
			configRepo: biometricRepos.NewBioTimeConfigRepository(),
			tenantID: tenantID,
			queue:    queue,
		}
	}

	return &BioTimeService{
		baseURL:  cfg.BioTime.BaseURL,
		username: cfg.BioTime.Username,
		password: cfg.BioTime.Password,
		enabled:  cfg.BioTime.Enabled,
		client:   &http.Client{Timeout: 30 * time.Second},
		configRepo: biometricRepos.NewBioTimeConfigRepository(),
		tenantID: tenantID,
		queue:    queue,
	}
}

// QueueTransactions queues transactions for background processing
func (s *BioTimeService) QueueTransactions(tenantID *uint, transactions []interface{}) error {
	if s.queue == nil {
		return fmt.Errorf("queue is not initialized")
	}
	return s.queue.PublishTransaction(tenantID, transactions)
}

// AuthRequest represents the authentication request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token   string `json:"token"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// TestConnectionResponse represents the test connection response
type TestConnectionResponse struct {
	Success      bool   `json:"success"`
	Message     string `json:"message"`
	BaseURL     string `json:"base_url"`
	Authenticated bool `json:"authenticated"`
	Token       string `json:"token,omitempty"`
	Error       string `json:"error,omitempty"`
}

// Authenticate authenticates with BioTime API and gets JWT token
func (s *BioTimeService) Authenticate() (string, error) {
	if !s.enabled {
		return "", errors.New("BioTime integration is disabled")
	}

	// Try to get config from database
	biometricConfig, err := s.configRepo.FindByTenantID(s.tenantID)
	if err == nil && biometricConfig != nil {
		// Check if stored token is still valid
		if biometricConfig.IsTokenValid() && biometricConfig.Token != nil {
			return *biometricConfig.Token, nil
		}
		// Use config from database if available
		if biometricConfig.BaseURL != "" {
			s.baseURL = biometricConfig.BaseURL
		}
		if biometricConfig.Username != "" {
			s.username = biometricConfig.Username
		}
		if biometricConfig.Password != "" {
			s.password = biometricConfig.Password
		}
		s.enabled = biometricConfig.Enabled
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Config doesn't exist, create it from environment
		biometricConfig, err = s.configRepo.FindOrCreate(s.tenantID, s.baseURL, s.username, s.password)
		if err != nil {
			return "", fmt.Errorf("failed to create BioTime config: %w", err)
		}
		// If config was just created, continue to authenticate
	} else {
		// Other error - log but continue with environment config
		_ = err
	}

	// Prepare authentication request
	authURL := fmt.Sprintf("%s/jwt-api-token-auth/", s.baseURL)
	authReq := AuthRequest{
		Username: s.username,
		Password: s.password,
	}

	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth request: %w", err)
	}

	// Make authentication request
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to BioTime API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("authentication failed with status %d: %s", resp.StatusCode, string(body))
		_ = s.configRepo.UpdateLastError(s.tenantID, errorMsg)
		return "", fmt.Errorf(errorMsg)
	}

	// Parse response
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		// Try to parse as simple token string
		var tokenStr string
		if err2 := json.Unmarshal(body, &tokenStr); err2 == nil {
			// Store token in database
			expiry := time.Now().Add(24 * time.Hour)
			_ = s.configRepo.UpdateToken(s.tenantID, tokenStr, expiry)
			return tokenStr, nil
		}
		errorMsg := fmt.Sprintf("failed to parse auth response: %v", err)
		_ = s.configRepo.UpdateLastError(s.tenantID, errorMsg)
		return "", fmt.Errorf(errorMsg)
	}

	if authResp.Token == "" {
		if authResp.Error != "" {
			// Update last error in database
			_ = s.configRepo.UpdateLastError(s.tenantID, authResp.Error)
			return "", fmt.Errorf("authentication error: %s", authResp.Error)
		}
		_ = s.configRepo.UpdateLastError(s.tenantID, "authentication failed: no token received")
		return "", errors.New("authentication failed: no token received")
	}

	// Store token in database with expiry (assume 24 hours, adjust based on actual API response)
	expiry := time.Now().Add(24 * time.Hour)
	if err := s.configRepo.UpdateToken(s.tenantID, authResp.Token, expiry); err != nil {
		// Log error but return token anyway
		_ = err
	}

	return authResp.Token, nil
}

// TestConnection tests the connection to BioTime API
func (s *BioTimeService) TestConnection() (*TestConnectionResponse, error) {
	if !s.enabled {
		return &TestConnectionResponse{
			Success:      false,
			Message:     "BioTime integration is disabled",
			BaseURL:     s.baseURL,
			Authenticated: false,
		}, nil
	}

	// Try to authenticate
	token, err := s.Authenticate()
	if err != nil {
		return &TestConnectionResponse{
			Success:      false,
			Message:     "Failed to authenticate with BioTime API",
			BaseURL:     s.baseURL,
			Authenticated: false,
			Error:       err.Error(),
		}, nil
	}

	return &TestConnectionResponse{
		Success:      true,
		Message:     "Successfully connected and authenticated with BioTime API",
		BaseURL:     s.baseURL,
		Authenticated: true,
		Token:       token,
	}, nil
}

// MakeRequest makes an authenticated request to BioTime API
func (s *BioTimeService) MakeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	if !s.enabled {
		return nil, errors.New("BioTime integration is disabled")
	}

	// Ensure we have a valid token
	token, err := s.Authenticate()
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Build URL
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make request
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

// GetToken returns the current authentication token
func (s *BioTimeService) GetToken() (string, error) {
	if !s.enabled {
		return "", errors.New("BioTime integration is disabled")
	}
	return s.Authenticate()
}

// RefreshToken forces a new token to be fetched from BioTime API
// This is useful if you want to manually refresh the token even if it's still valid
func (s *BioTimeService) RefreshToken() (string, error) {
	if !s.enabled {
		return "", errors.New("BioTime integration is disabled")
	}

	// Get config from database to use stored credentials
	biometricConfig, err := s.configRepo.FindByTenantID(s.tenantID)
	if err == nil && biometricConfig != nil {
		// Use config from database
		if biometricConfig.BaseURL != "" {
			s.baseURL = biometricConfig.BaseURL
		}
		if biometricConfig.Username != "" {
			s.username = biometricConfig.Username
		}
		if biometricConfig.Password != "" {
			s.password = biometricConfig.Password
		}
		s.enabled = biometricConfig.Enabled
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Config doesn't exist, create it from environment
		biometricConfig, err = s.configRepo.FindOrCreate(s.tenantID, s.baseURL, s.username, s.password)
		if err != nil {
			return "", fmt.Errorf("failed to create BioTime config: %w", err)
		}
	}

	// Force new authentication (don't check existing token)
	authURL := fmt.Sprintf("%s/jwt-api-token-auth/", s.baseURL)
	authReq := AuthRequest{
		Username: s.username,
		Password: s.password,
	}

	jsonData, err := json.Marshal(authReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth request: %w", err)
	}

	// Make authentication request
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to connect to BioTime API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorMsg := fmt.Sprintf("authentication failed with status %d: %s", resp.StatusCode, string(body))
		_ = s.configRepo.UpdateLastError(s.tenantID, errorMsg)
		return "", fmt.Errorf(errorMsg)
	}

	// Parse response
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		// Try to parse as simple token string
		var tokenStr string
		if err2 := json.Unmarshal(body, &tokenStr); err2 == nil {
			// Store token in database
			expiry := time.Now().Add(24 * time.Hour)
			_ = s.configRepo.UpdateToken(s.tenantID, tokenStr, expiry)
			return tokenStr, nil
		}
		errorMsg := fmt.Sprintf("failed to parse auth response: %v", err)
		_ = s.configRepo.UpdateLastError(s.tenantID, errorMsg)
		return "", fmt.Errorf(errorMsg)
	}

	if authResp.Token == "" {
		if authResp.Error != "" {
			// Update last error in database
			_ = s.configRepo.UpdateLastError(s.tenantID, authResp.Error)
			return "", fmt.Errorf("authentication error: %s", authResp.Error)
		}
		_ = s.configRepo.UpdateLastError(s.tenantID, "authentication failed: no token received")
		return "", errors.New("authentication failed: no token received")
	}

	// Store new token in database with expiry (assume 24 hours, adjust based on actual API response)
	expiry := time.Now().Add(24 * time.Hour)
	if err := s.configRepo.UpdateToken(s.tenantID, authResp.Token, expiry); err != nil {
		// Log error but return token anyway
		_ = err
	}

	return authResp.Token, nil
}

// IntOrString supports JSON numbers that may come as number, string, or null.
// BioTime sometimes returns numeric fields as strings, e.g. "1".
type IntOrString struct {
	Value int
	Valid bool
}

func (v *IntOrString) UnmarshalJSON(b []byte) error {
	raw := strings.TrimSpace(string(b))
	if raw == "" || raw == "null" {
		v.Value = 0
		v.Valid = false
		return nil
	}

	// JSON string
	if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			v.Value = 0
			v.Valid = false
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		v.Value = n
		v.Valid = true
		return nil
	}

	// JSON number
	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	v.Value = n
	v.Valid = true
	return nil
}

func (v IntOrString) MarshalJSON() ([]byte, error) {
	if !v.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(v.Value)
}

// TerminalArea represents the area information for a terminal
type TerminalArea struct {
	ID       IntOrString `json:"id"`
	AreaCode string `json:"area_code"`
	AreaName string `json:"area_name"`
}

// Terminal represents a BioTime terminal device
type Terminal struct {
	ID                IntOrString   `json:"id"`
	SerialNo          string        `json:"sn"`
	IPAddress          string        `json:"ip_address"`
	Alias             string        `json:"alias"`
	TerminalName      *string       `json:"terminal_name"`
	FirmwareVersion   *string       `json:"fw_ver"`
	PushVersion       *string       `json:"push_ver"`
	State             IntOrString   `json:"state"`
	TerminalTimezone  IntOrString   `json:"terminal_tz"`
	Area              TerminalArea  `json:"area"`
	LastActivity      string        `json:"last_activity"`
	UserCount         IntOrString   `json:"user_count"`
	FingerprintCount  IntOrString   `json:"fp_count"`
	FaceCount         IntOrString   `json:"face_count"`
	PalmCount         IntOrString   `json:"palm_count"`
	TransactionCount  IntOrString   `json:"transaction_count"`
	PushTime          *string       `json:"push_time"`
	TransferTime      string        `json:"transfer_time"`
	TransferInterval  IntOrString   `json:"transfer_interval"`
	IsAttendance      IntOrString   `json:"is_attendance"`
	AreaName          string        `json:"area_name"`
}

// TerminalsResponse represents the BioTime API response for terminals
type TerminalsResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Message  string     `json:"msg"`
	Code     int        `json:"code"`
	Data     []Terminal `json:"data"`
}

// GetTerminals fetches the list of terminals from BioTime API
func (s *BioTimeService) GetTerminals() (*TerminalsResponse, error) {
	if !s.enabled {
		return nil, errors.New("BioTime integration is disabled")
	}

	// Make authenticated request to BioTime API
	resp, err := s.MakeRequest("GET", "/iclock/api/terminals/", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch terminals: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch terminals: status %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse response - BioTime API returns wrapped response
	var terminalsResp TerminalsResponse
	if err := json.Unmarshal(body, &terminalsResp); err != nil {
		return nil, fmt.Errorf("failed to parse terminals response: %w", err)
	}

	// Check if API returned an error code
	if terminalsResp.Code != 0 {
		return nil, fmt.Errorf("BioTime API error: code %d, message: %s", terminalsResp.Code, terminalsResp.Message)
	}

	return &terminalsResp, nil
}

// Transaction represents a BioTime transaction/attendance record
type Transaction struct {
	ID                 IntOrString `json:"id"`
	EmpCode            string      `json:"emp_code"`
	FirstName          string      `json:"first_name"`
	LastName           string      `json:"last_name"`
	Department         string      `json:"department"`
	Position           string      `json:"position"`
	PunchTime          string      `json:"punch_time"`
	PunchState         string      `json:"punch_state"`
	PunchStateDisplay  string      `json:"punch_state_display"`
	VerifyType         IntOrString `json:"verify_type"`
	VerifyTypeDisplay  string      `json:"verify_type_display"`
	WorkCode           string      `json:"work_code"`
	GPSLocation        string      `json:"gps_location"`
	AreaAlias          *string     `json:"area_alias"`
	TerminalSN         string      `json:"terminal_sn"`
	Temperature        float64     `json:"temperature"`
	TerminalAlias      *string     `json:"terminal_alias"`
	UploadTime         string      `json:"upload_time"`
}

// TransactionsResponse represents the BioTime API response for transactions
type TransactionsResponse struct {
	Count    int          `json:"count"`
	Next     *string      `json:"next"`
	Previous *string      `json:"previous"`
	Message  string       `json:"msg"`
	Code     int          `json:"code"`
	Data     []Transaction `json:"data"`
}

// GetTransactionsParams represents query parameters for fetching transactions
type GetTransactionsParams struct {
	Page          *int    `json:"page,omitempty"`
	PageSize      *int    `json:"page_size,omitempty"`
	EmpCode       *string `json:"emp_code,omitempty"`
	TerminalSN    *string `json:"terminal_sn,omitempty"`
	TerminalAlias *string `json:"terminal_alias,omitempty"`
	StartTime     *string `json:"start_time,omitempty"`
	EndTime       *string `json:"end_time,omitempty"`
}

// GetTransactions fetches the list of transactions from BioTime API
func (s *BioTimeService) GetTransactions(params *GetTransactionsParams) (*TransactionsResponse, error) {
	if !s.enabled {
		return nil, errors.New("BioTime integration is disabled")
	}

	// Build query string
	queryParams := make([]string, 0)
	if params != nil {
		if params.Page != nil {
			queryParams = append(queryParams, fmt.Sprintf("page=%d", *params.Page))
		}
		if params.PageSize != nil {
			queryParams = append(queryParams, fmt.Sprintf("page_size=%d", *params.PageSize))
		}
		if params.EmpCode != nil && *params.EmpCode != "" {
			queryParams = append(queryParams, fmt.Sprintf("emp_code=%s", *params.EmpCode))
		}
		if params.TerminalSN != nil && *params.TerminalSN != "" {
			queryParams = append(queryParams, fmt.Sprintf("terminal_sn=%s", *params.TerminalSN))
		}
		if params.TerminalAlias != nil && *params.TerminalAlias != "" {
			queryParams = append(queryParams, fmt.Sprintf("terminal_alias=%s", *params.TerminalAlias))
		}
		if params.StartTime != nil && *params.StartTime != "" {
			queryParams = append(queryParams, fmt.Sprintf("start_time=%s", *params.StartTime))
		}
		if params.EndTime != nil && *params.EndTime != "" {
			queryParams = append(queryParams, fmt.Sprintf("end_time=%s", *params.EndTime))
		}
	}

	// Build endpoint URL with query parameters
	endpoint := "/iclock/api/transactions/"
	if len(queryParams) > 0 {
		endpoint += "?" + strings.Join(queryParams, "&")
	}

	// Make authenticated request to BioTime API
	resp, err := s.MakeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transactions: status %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse response - BioTime API returns wrapped response
	var transactionsResp TransactionsResponse
	if err := json.Unmarshal(body, &transactionsResp); err != nil {
		return nil, fmt.Errorf("failed to parse transactions response: %w", err)
	}

	// Check if API returned an error code
	if transactionsResp.Code != 0 {
		return nil, fmt.Errorf("BioTime API error: code %d, message: %s", transactionsResp.Code, transactionsResp.Message)
	}

	return &transactionsResp, nil
}

// GetTransactionByID fetches a single transaction by ID from BioTime API
func (s *BioTimeService) GetTransactionByID(id string) (*Transaction, error) {
	if !s.enabled {
		return nil, errors.New("BioTime integration is disabled")
	}

	if id == "" {
		return nil, errors.New("transaction ID is required")
	}

	// Build endpoint URL
	endpoint := fmt.Sprintf("/iclock/api/transactions/%s/", id)

	// Make authenticated request to BioTime API
	resp, err := s.MakeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch transaction: status %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse response - BioTime API returns single transaction object
	var transaction Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		return nil, fmt.Errorf("failed to parse transaction response: %w", err)
	}

	return &transaction, nil
}
