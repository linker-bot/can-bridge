# CAN-Bridge Project Documentation

[中文](README_zhCN.md)

![# CAN Bridge](.github/assets/banner.jpg)

## 🔧Project Overview

**CAN-Bridge** is a hardware adaptation bridge service developed in Golang, specifically designed to support efficient and stable data communication for devices using the CAN (Controller Area Network) protocol.

The project aims to provide an easy-to-use HTTP API interface, allowing users to send and receive CAN messages, manage interface configurations, and monitor device status via network requests, while also supporting dynamic configuration of multiple CAN interfaces.

## ⚙️Features

### ✨New Features

* **Interface Setup Manager**:

  * Automatically executes `ip link set can0 up type can bitrate 1000000`
  * Configurable bitrate, sample point, and restart timeout
  * Provides retry mechanism and error handling
  * Supports interface state querying and validation

* **Enhanced Configuration System**:

  * Supports configuration via command-line parameters and environment variables
  * Provides configuration validation to ensure parameter correctness, including bitrate
  * Includes detailed usage instructions

* **Comprehensive Interface Management API**:

  * Configuration management API (`GET /api/setup/config`, `PUT /api/setup/config`)
  * Interface operation API (setup, shutdown, reset, status query)
  * Batch operation API (setup or teardown all interfaces at once)

### 📦Main Functions

* **Dynamic Interface Configuration**: Supports dynamic configuration of multiple CAN interfaces (e.g., `can0`, `can1`) via command-line or environment variables.
* **HTTP API Service**: Provides easy-to-use RESTful API interfaces.
* **Message Sending and Validation**: Supports sending standard CAN messages and performs validation of data length and interface availability.
* **Health Check and Automatic Recovery**: Built-in interface health monitoring with automatic recovery of failed interfaces.
* **Real-time Monitoring and Statistics**: Offers monitoring of interface status, message statistics, success rates, error logging, and average latency.

### 🛠️Program Features

* ✅ Automatic Initialization: Automatically configures CAN interfaces on program startup.
* ✅ Retry Mechanism: Automatically retries interface setup upon failure.
* ✅ Status Monitoring: Real-time monitoring of interface status and error statistics.
* ✅ Graceful Shutdown: Automatically shuts down interfaces upon program exit.
* ✅ Dependency Injection: Facilitates testing and extension.
* ✅ Error Handling: Comprehensive error handling and logging.

## 📦Installation and Usage

Supports installation on bare-metal systems, Docker containers, and source code builds.

For detailed installation instructions, see the [Installation Guide](docs/install.md).

### 💡Usage Examples

**Basic Usage**

```bash
./can-bridge -can-ports can0,can1
```

**Set Port**

```bash
./can-bridge -port 5260
```

**Disable Automatic Setup (Managed via API)**

```bash
./can-bridge -auto-setup=false
```

**Custom Bitrate**

```bash
./can-bridge -can-ports can0 -bitrate 500000
```

**Sample Point**

```bash
./can-bridge -sample-point 0.75
```

**Restart Timeout**

```bash
./can-bridge -restart-ms 100
```

**Setup Retry**

```bash
./can-bridge -setup-retry 3
```

**Setup Delay**

```bash
./can-bridge -setup-delay 2
```

**Enable Service Finder**

```bash
./can-bridge -enable-finder=true
```

**Service Finder Interval**

```bash
./can-bridge -finder-interval 5
```

**Enable Health Check**

```bash
./can-bridge -enable-healthcheck=true
```

**Configure Interface via API**

```bash
curl -X POST localhost:5260/api/setup/interfaces/can0 \
  -H "Content-Type: application/json" \
  -d '{"bitrate": 500000, "withRetry": true}'
```

## 🌐API Documentation

### 📍Base Path

`http://localhost:5260/api`

### ⭐ Status & Monitoring

APIs for retrieving system status, interface health, and performance metrics.

* `GET /api/status`: Get the complete system status, including uptime, watchdog status, and all interface details.
* `GET /api/interfaces`: Get a list of configured and active interfaces.
* `GET /api/interfaces/:name/status`: Get the detailed status for a specific interface.
* `GET /api/health`: Get a summary of the system's health.
* `GET /api/metrics`: Get detailed metrics formatted for external monitoring systems (e.g., Prometheus).

### ✉️ Message Sending

* `POST /api/can`: Send a single CAN message. The request body should contain the message details (e.g., ID, Data).

### 🔧 Interface Setup Management

APIs for dynamically configuring, starting, stopping, and managing CAN interfaces.

**Configuration Management**:

* `GET /api/setup/config`: Get the current interface setup configuration (e.g., default bitrate, sample point).
* `PUT /api/setup/config`: Update the global configuration for interface setup.

**Interface Operations**:

* `GET /api/setup/available`: Get a list of all available CAN interfaces on the operating system.
* `POST /api/setup/interfaces/{name}`: Set up and bring up a specific CAN interface based on the configuration.
* `DELETE /api/setup/interfaces/{name}`: Bring down and tear down a specific CAN interface.
* `POST /api/setup/interfaces/{name}/reset`: Reset a specific CAN interface (teardown and then setup).
* `GET /api/setup/interfaces/{name}/state`: Get the current setup state of a specific interface (e.g., if it is up, config details).

**Batch Operations**:

* `POST /api/setup/interfaces/setup-all`: Set up all configured interfaces or a specific list of interfaces from the request.
* `POST /api/setup/interfaces/teardown-all`: Tear down all configured interfaces.

### 📡 Message Listening & Retrieval

APIs for capturing, viewing, and managing messages from the CAN bus in real-time.

**Listener Control**:

* `POST /api/messages/:interface/listen/start`: Start listening for CAN messages on a specific interface.
* `POST /api/messages/:interface/listen/stop`: Stop listening for CAN messages on a specific interface.
* `GET /api/messages/:interface/listen/status`: Get the current listening status for a specific interface.
* `GET /api/messages/listen/status`: Get a summary of the listening status for all interfaces.

**Message Retrieval**:

* `GET /api/messages/:interface`: Get all cached messages for a specific interface. Supports filtering by `id` query parameter.
* `GET /api/messages/:interface/recent`: Get the N most recent messages from an interface (specify with the `count` query parameter).
* `GET /api/messages/`: Get all cached messages from all interfaces, grouped by interface.

**Message Management & Statistics**:

* `GET /api/messages/:interface/statistics`: Get message statistics for a specific interface (total received, errors, etc.).
* `DELETE /api/messages/:interface`: Clear the message buffer for a specific interface.
* `GET /api/messages/statistics`: Get global message statistics for all interfaces.
* `DELETE /api/messages/`: Clear the message buffers for all interfaces.

## 🚀Performance Optimization and Stability

* Implements retry mechanisms for reliable message transmission.
* Utilizes mutex locks to ensure thread safety.
* Real-time monitoring of interface health status with automatic recovery.

## 📝Logging and Debugging

Logs are output to the standard output stream in a friendly format, including clear error messages and runtime status information.

## 📦Deployment Recommendations

Deployment using systemd or Docker containers is recommended to ensure long-term stable operation.

## 🤝Contribution Guide

Issues and Pull Requests are welcomed to improve and optimize the project.

## 📄License

This project is licensed under the Apache-2.0 license.
