# Golang Worker Pool

A simple implementation of a worker pool pattern in Go.

## Requirements

- Go 1.16 or higher

## Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/golang-worker-pool.git
cd golang-worker-pool
```

## How to Run

### Basic Usage

To run the application with default settings (10 workers, 50 tasks):

```bash
go run main.go
```

### With Custom Parameters

You can specify the number of workers and tasks using command-line flags:

```bash
go run main.go -workers=10 -tasks=100
```

Available flags:
- `-workers`: Number of workers in the pool (default: 10)
- `-tasks`: Number of tasks to process (default: 50)

## Building the Application

To build an executable:

```bash
go build -o worker-pool
./worker-pool -workers=20 -tasks=200
```