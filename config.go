package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Configuration structure
type Config struct {
	CanPorts            []string
	Port                string
	AutoSetup           bool          // Auto setup CAN interfaces on startup
	Bitrate             int           // Default bitrate for CAN interfaces
	SamplePoint         string        // Default sample point
	RestartMs           int           // Default restart timeout
	SetupRetry          int           // Number of setup retry attempts
	SetupDelay          time.Duration // Delay between setup retries
	EnableFinder        bool          // Enable service finder
	SetupFinderInterval time.Duration // Interval for service finder
	EnableHealthCheck   bool          // Enable health check endpoint
}

// ConfigProvider interface for dependency injection
type ConfigProvider interface {
	GetCanPorts() []string
	GetServerPort() string
	ValidateInterface(ifName string) bool
	GetAutoSetup() bool
	GetDefaultBitrate() int
	GetDefaultSamplePoint() string
	GetDefaultRestartMs() int
	GetSetupRetry() int
	GetSetupDelay() time.Duration
}

// DefaultConfigProvider implements ConfigProvider
type DefaultConfigProvider struct {
	config *Config
}

// NewDefaultConfigProvider creates a new default config provider
func NewDefaultConfigProvider(config *Config) *DefaultConfigProvider {
	return &DefaultConfigProvider{config: config}
}

// GetCanPorts returns configured CAN ports
func (p *DefaultConfigProvider) GetCanPorts() []string {
	return p.config.CanPorts
}

// GetServerPort returns server port
func (p *DefaultConfigProvider) GetServerPort() string {
	return p.config.Port
}

// ValidateInterface checks if interface is in configured ports
func (p *DefaultConfigProvider) ValidateInterface(ifName string) bool {
	for _, port := range p.config.CanPorts {
		if port == ifName {
			return true
		}
	}
	return false
}

// GetAutoSetup returns auto setup configuration
func (p *DefaultConfigProvider) GetAutoSetup() bool {
	return p.config.AutoSetup
}

// GetDefaultBitrate returns default bitrate
func (p *DefaultConfigProvider) GetDefaultBitrate() int {
	return p.config.Bitrate
}

// GetDefaultSamplePoint returns default sample point
func (p *DefaultConfigProvider) GetDefaultSamplePoint() string {
	return p.config.SamplePoint
}

// GetDefaultRestartMs returns default restart timeout
func (p *DefaultConfigProvider) GetDefaultRestartMs() int {
	return p.config.RestartMs
}

// GetSetupRetry returns setup retry count
func (p *DefaultConfigProvider) GetSetupRetry() int {
	return p.config.SetupRetry
}

// GetSetupDelay returns setup retry delay
func (p *DefaultConfigProvider) GetSetupDelay() time.Duration {
	return p.config.SetupDelay
}

func (p *DefaultConfigProvider) GetEnableFinder() bool {
	return p.config.EnableFinder
}

func (p *DefaultConfigProvider) GetSetupFinderInterval() time.Duration {
	return p.config.SetupFinderInterval
}

func (p *DefaultConfigProvider) GetEnableHealthCheck() bool {
	return p.config.EnableHealthCheck
}

// ConfigParser handles parsing configuration from various sources
type ConfigParser struct{}

// NewConfigParser creates a new config parser
func NewConfigParser() *ConfigParser {
	return &ConfigParser{}
}

// ParseConfig parses configuration from command line and environment variables
func (cp *ConfigParser) ParseConfig() (*Config, error) {
	config := &Config{}

	// Command line flags
	var canPortsFlag string
	var serverPort string
	var autoSetup bool
	var bitrate int
	var samplePoint string
	var restartMs int
	var setupRetry int
	var setupDelaySeconds int
	var setupFinderEnabled bool
	var setupFinderInterval int
	var setupHealthCheck bool

	flag.StringVar(&canPortsFlag, "can-ports", "", "Comma-separated list of CAN interfaces (e.g., can0,can1)")
	flag.StringVar(&serverPort, "port", "5260", "HTTP server port")
	flag.BoolVar(&autoSetup, "auto-setup", true, "Automatically setup CAN interfaces on startup")
	flag.IntVar(&bitrate, "bitrate", 1000000, "Default CAN bitrate (bps)")
	flag.StringVar(&samplePoint, "sample-point", "0.75", "Default CAN sample point")
	flag.IntVar(&restartMs, "restart-ms", 100, "Default CAN restart timeout (ms)")
	flag.IntVar(&setupRetry, "setup-retry", 3, "Number of setup retry attempts")
	flag.IntVar(&setupDelaySeconds, "setup-delay", 2, "Delay between setup retries (seconds)")
	flag.BoolVar(&setupFinderEnabled, "enable-finder", true, "Enable service finder")
	flag.IntVar(&setupFinderInterval, "finder-interval", 5, "Interval for service finder in seconds")
	flag.BoolVar(&setupHealthCheck, "enable-healthcheck", true, "Enable health check endpoint")
	flag.Parse()

	// Environment variables (override command line)
	if envPorts := os.Getenv("CAN_PORTS"); envPorts != "" {
		canPortsFlag = envPorts
	}
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		serverPort = envPort
	}
	if envAutoSetup := os.Getenv("CAN_AUTO_SETUP"); envAutoSetup != "" {
		if val, err := strconv.ParseBool(envAutoSetup); err == nil {
			autoSetup = val
		}
	}
	if envBitrate := os.Getenv("CAN_BITRATE"); envBitrate != "" {
		if val, err := strconv.Atoi(envBitrate); err == nil {
			bitrate = val
		}
	}
	if envSamplePoint := os.Getenv("CAN_SAMPLE_POINT"); envSamplePoint != "" {
		samplePoint = envSamplePoint
	}
	if envRestartMs := os.Getenv("CAN_RESTART_MS"); envRestartMs != "" {
		if val, err := strconv.Atoi(envRestartMs); err == nil {
			restartMs = val
		}
	}
	if envSetupRetry := os.Getenv("CAN_SETUP_RETRY"); envSetupRetry != "" {
		if val, err := strconv.Atoi(envSetupRetry); err == nil {
			setupRetry = val
		}
	}
	if envSetupDelay := os.Getenv("CAN_SETUP_DELAY"); envSetupDelay != "" {
		if val, err := strconv.Atoi(envSetupDelay); err == nil {
			setupDelaySeconds = val
		}
	}

	// Parse CAN ports
	if canPortsFlag != "" {
		config.CanPorts = cp.parseCanPorts(canPortsFlag)
	} else {
		// Default to can0 if no ports specified
		config.CanPorts = []string{"can0"}
	}

	// Validate and set configuration
	if serverPort == "" {
		return nil, fmt.Errorf("server port cannot be empty")
	}

	if setupFinderEnabled {
		if setupFinderInterval <= 0 {
			return nil, fmt.Errorf("finder interval must be positive, got %d", config.SetupFinderInterval)
		}
	}

	if setupHealthCheck {
		config.EnableHealthCheck = true
	} else {
		config.EnableHealthCheck = false
	}

	config.Port = serverPort
	config.AutoSetup = autoSetup
	config.Bitrate = bitrate
	config.SamplePoint = samplePoint
	config.RestartMs = restartMs
	config.SetupRetry = setupRetry
	config.SetupDelay = time.Duration(setupDelaySeconds) * time.Second
	config.EnableFinder = setupFinderEnabled
	config.SetupFinderInterval = time.Duration(setupFinderInterval) * time.Second

	return config, nil
}

// parseCanPorts parses comma-separated CAN ports string
func (cp *ConfigParser) parseCanPorts(portsStr string) []string {
	ports := strings.Split(portsStr, ",")
	// Trim whitespace from each port
	for i, port := range ports {
		ports[i] = strings.TrimSpace(port)
	}
	return ports
}

// ValidateConfig validates the configuration
func (cp *ConfigParser) ValidateConfig(config *Config) error {
	if len(config.CanPorts) == 0 {
		return fmt.Errorf("at least one CAN port must be specified")
	}

	for _, port := range config.CanPorts {
		if strings.TrimSpace(port) == "" {
			return fmt.Errorf("CAN port name cannot be empty")
		}
	}

	if config.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	// Validate CAN-specific settings
	if config.Bitrate <= 0 {
		return fmt.Errorf("bitrate must be positive, got %d", config.Bitrate)
	}

	// Common CAN bitrates validation
	validBitrates := []int{
		10000,   // 10 kbps
		20000,   // 20 kbps
		50000,   // 50 kbps
		100000,  // 100 kbps
		125000,  // 125 kbps
		250000,  // 250 kbps
		500000,  // 500 kbps
		1000000, // 1 Mbps
	}

	validBitrate := false
	for _, valid := range validBitrates {
		if config.Bitrate == valid {
			validBitrate = true
			break
		}
	}
	if !validBitrate {
		return fmt.Errorf("bitrate %d is not a standard CAN bitrate. Valid options: %v", config.Bitrate, validBitrates)
	}

	if config.SamplePoint != "" {
		if point, err := strconv.ParseFloat(config.SamplePoint, 64); err != nil {
			return fmt.Errorf("invalid sample point format: %s", config.SamplePoint)
		} else if point <= 0 || point >= 1 {
			return fmt.Errorf("sample point must be between 0 and 1, got %f", point)
		}
	}

	if config.RestartMs < 0 {
		return fmt.Errorf("restart timeout cannot be negative, got %d", config.RestartMs)
	}

	if config.SetupRetry <= 0 {
		return fmt.Errorf("setup retry count must be positive, got %d", config.SetupRetry)
	}

	if config.SetupDelay < 0 {
		return fmt.Errorf("setup delay cannot be negative, got %v", config.SetupDelay)
	}

	return nil
}

// GetConfigSummary returns a summary of the current configuration
func (cp *ConfigParser) GetConfigSummary(config *Config) map[string]interface{} {
	return map[string]interface{}{
		"canPorts":    config.CanPorts,
		"serverPort":  config.Port,
		"autoSetup":   config.AutoSetup,
		"bitrate":     config.Bitrate,
		"samplePoint": config.SamplePoint,
		"restartMs":   config.RestartMs,
		"setupRetry":  config.SetupRetry,
		"setupDelay":  config.SetupDelay.String(),
	}
}

// PrintUsage prints usage information
func PrintUsage() {
	fmt.Println("CAN Communication Service")
	fmt.Println("Usage:")
	fmt.Println("  -can-ports string       Comma-separated list of CAN interfaces (default: can0)")
	fmt.Println("  -port string            HTTP server port (default: 5260)")
	fmt.Println("  -auto-setup             Automatically setup CAN interfaces on startup (default: true)")
	fmt.Println("  -bitrate int            Default CAN bitrate in bps (default: 1000000)")
	fmt.Println("  -sample-point string    Default CAN sample point (default: 0.75)")
	fmt.Println("  -restart-ms int         Default CAN restart timeout in ms (default: 100)")
	fmt.Println("  -setup-retry int        Number of setup retry attempts (default: 3)")
	fmt.Println("  -setup-delay int        Delay between setup retries in seconds (default: 2)")
	fmt.Println("  -enable-finder          Enable service finder (default: true)")
	fmt.Println("  -finder-interval int    Interval for service finder in seconds (default: 5)")
	fmt.Println("  -enable-healthcheck     Enable health check endpoint (default: true)")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("  CAN_PORTS              Comma-separated list of CAN interfaces")
	fmt.Println("  SERVER_PORT            HTTP server port")
	fmt.Println("  CAN_AUTO_SETUP         Automatically setup CAN interfaces (true/false)")
	fmt.Println("  CAN_BITRATE            Default CAN bitrate in bps")
	fmt.Println("  CAN_SAMPLE_POINT       Default CAN sample point")
	fmt.Println("  CAN_RESTART_MS         Default CAN restart timeout in ms")
	fmt.Println("  CAN_SETUP_RETRY        Number of setup retry attempts")
	fmt.Println("  CAN_SETUP_DELAY        Delay between setup retries in seconds")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Basic usage with default settings")
	fmt.Println("  ./can-bridge -can-ports can0,can1")
	fmt.Println("")
	fmt.Println("  # Custom bitrate and sample point")
	fmt.Println("  ./can-bridge -can-ports can0 -bitrate 500000 -sample-point 0.8")
	fmt.Println("")
	fmt.Println("  # Disable auto-setup (manual setup via API)")
	fmt.Println("  ./can-bridge -can-ports can0,can1 -auto-setup=false")
	fmt.Println("")
	fmt.Println("  # Using environment variables")
	fmt.Println("  CAN_PORTS=can0,can1 CAN_BITRATE=500000 ./can-bridge")
	fmt.Println("")
	fmt.Println("  # High availability setup with more retries")
	fmt.Println("  ./can-bridge -can-ports can0,can1 -setup-retry 5 -setup-delay 3")
	fmt.Println("")
	fmt.Println("Valid CAN Bitrates:")
	fmt.Println("  10000, 20000, 50000, 100000, 125000, 250000, 500000, 1000000 (bps)")
	fmt.Println("")
	fmt.Println("API Endpoints:")
	fmt.Println("  GET  /api/setup/config                    - Get setup configuration")
	fmt.Println("  PUT  /api/setup/config                    - Update setup configuration")
	fmt.Println("  GET  /api/setup/available                 - List available CAN interfaces")
	fmt.Println("  POST /api/setup/interfaces/{name}        - Setup specific interface")
	fmt.Println("  DELETE /api/setup/interfaces/{name}      - Teardown specific interface")
	fmt.Println("  POST /api/setup/interfaces/{name}/reset  - Reset specific interface")
	fmt.Println("  GET  /api/setup/interfaces/{name}/state  - Get interface state")
	fmt.Println("  POST /api/setup/interfaces/setup-all     - Setup all interfaces")
	fmt.Println("  POST /api/setup/interfaces/teardown-all  - Teardown all interfaces")
}
