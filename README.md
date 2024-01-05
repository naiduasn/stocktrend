# stocktrend

This Go application fetches data from an API based on configuration from an `.env` file.

## Prerequisites

- Go installed (version 1.21.X)
- `.env` file in the project root directory with the following keys:
  - `URL`: API endpoint URL
  - `TIMEOUT`: Timeout in seconds (numeric value)

## Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/naiduasn/stocktrend.git
    cd stocktrend
    ```

2. **Install necessary packages:**

    Ensure you have installed the required external packages used in the project.

    ```bash
    go get github.com/joho/godotenv
    ```

3. **Create and set environment variables in `.env` file:**

    Create a file named `.env` in the root directory of the project and add the required environment variables:

    ```plaintext
    URL=<API to fetch live stocks data>
    TIMEOUT=60
    ```

## Usage

### Build and Run

Use the following commands to build and run the application:

```bash
go build -o stocktrend main.go
./stocktrend
