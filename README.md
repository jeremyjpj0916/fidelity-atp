# Fidelity Active Trader Pro - Automated Trading Bot

A powerful automation tool designed to interact with Fidelity's Active Trader Pro platform for executing trades quickly and efficiently. This project includes two implementations:

1. **Python Implementation** (`trading_bot.py`) - Uses PyAutoGUI for GUI automation
2. **Go Implementation** (`trading_bot.go`) - Higher performance alternative using AppleScript for macOS

## Features

- Trade stocks and options with customizable parameters
- Support for multiple account types (Roth, Traditional, HSA, Brokeragelink, Individual)
- Trading by share quantity
- Market and limit order support
- Extended hours trading with Day+ option
- Order chunking for breaking large orders into multiple smaller batches
- Interactive position recording for GUI element locations
- Logging of all trading activities

## Prerequisites

### For Python Implementation
- Python 3.6+
- PyAutoGUI library: `pip install pyautogui`

### For Go Implementation
- Go 1.16+
- Build the executable: `go build trading_bot.go`

## Setup Guide

Before using either trading bot, you must record the positions of UI elements in the Fidelity Active Trader Pro interface.

### 1. Launch Fidelity Active Trader Pro

Ensure the application is running and visible on your primary monitor.

### 2. Record UI Element Positions

Both the Python and Go implementations now use the same **click-based** recording process.

**Steps:**

1.  Run the desired implementation in recording mode:
    -   **Python:** `python trading_bot.py record_positions`
    -   **Go:** `go build trading_bot.go && ./trading_bot -mode record_positions`
2.  The script will prompt you for each UI element needed (e.g., "Click the Stocks button").
3.  Move your mouse cursor over the corresponding element in the Fidelity ATP window.
4.  **Left-click** your mouse.
5.  The script will detect the click, record the coordinates, and automatically proceed to the next element.
6.  Repeat this process for all elements.
7.  The recorded positions will be saved to `click_positions.json`.

This method ensures accurate coordinates are captured exactly where you click, providing a consistent experience on both macOS and Windows.

## Usage Examples

### Python Implementation

1. **Basic Market Order (Buy Shares)**
   ```bash
   python trading_bot.py trade --trade_type stocks --account Roth --ticker AAPL --action buy --amount 10 --order_type market
   ```

2. **Limit Order with Extended Hours Trading**
   ```bash
   python trading_bot.py trade --trade_type stocks --account HSA --ticker NVDA --action buy --amount 5 --order_type limit --limit_price 850.50 --extended_hours
   ```

3. **Repeating an Order Multiple Times**
   ```bash
   python trading_bot.py trade --trade_type stocks --account Traditional --ticker AMZN --action sell --amount 100 --order_type limit --limit_price 180.25 --repeat 5 --min_pause 1 --max_pause 3
   ```

### Go Implementation

1. **Basic Market Order (Buy Shares)**
   ```bash
   ./trading_bot -mode trade -trade_type stocks -account Roth -ticker AAPL -action buy -amount 10 -order_type market
   ```

2. **Limit Order with Extended Hours Trading**
   ```bash
   ./trading_bot -mode trade -trade_type stocks -account HSA -ticker NVDA -action buy -amount 5 -order_type limit -limit_price 850.50 -extended_hours
   ```

3. **Repeating an Order Multiple Times**
   ```bash
   ./trading_bot -mode trade -trade_type stocks -account Traditional -ticker AMZN -action sell -amount 100 -order_type limit -limit_price 180.25 -repeat 5 -min_pause 1 -max_pause 3
   ```

## Command-Line Arguments

| Argument | Description | Values |
|----------|-------------|--------|
| `mode` | Operation mode | `record_positions`, `trade` |
| `--trade_type` | Type of security to trade | `stocks`, `options` |
| `--account` | Account to trade in | `Roth`, `Traditional`, `HSA`, `Brokeragelink`, `Individual` |
| `--ticker` | Ticker symbol | Any valid stock/option symbol |
| `--action` | Buy or sell | `buy`, `sell` |
| `--amount` | Number of shares to trade | Any positive number |
| `--order_type` | Type of order | `market`, `limit` |
| `--limit_price` | Price for limit orders | Any positive number |
| `--extended_hours` | Enable extended hours trading | Flag (no value needed) |
| `--repeat` | Number of times to repeat the same order | Integer > 0 (default: 1) |
| `--min_pause` | Min seconds between repeats | Number (default: 1.0) |
| `--max_pause` | Max seconds between repeats | Number (default: 3.0) |

## Important Notes

- The trading bot interacts with Fidelity's UI, so it needs to be run on a machine where Fidelity ATP is installed and logged in
- **Fidelity-specific restrictions apply:**
  - For extended hours trading (pre-market or after-hours), only limit orders are supported by Fidelity
  - The bot will return an error if you try to use market orders with the `--extended_hours` option
- Price and volume information may have 15-20 minute delays unless you have real-time data enabled in Fidelity
- Do not move the Fidelity ATP window after recording positions
- Keep the Fidelity ATP window visible and active during trading
- The bot logs all activities to `trading_bot.log`
- For limit orders, the `--limit_price` parameter is required
- The Go implementation generally offers better performance for frequent trades

## Troubleshooting

1. **Position Recording Issues**
   - If you need to stop the recording process, press Ctrl+C and restart
   - Make sure Fidelity ATP is on your primary monitor
   - Ensure the window is fully visible and not obscured
   - **Permissions (macOS):** The first time you run the recording mode, macOS might ask for "Input Monitoring" (Python/pynput) or "Accessibility" (Go/robotgo) permissions. Grant these permissions in System Settings > Privacy & Security.
   - **Mistakes:** If you click the wrong spot, press `Ctrl+C` to stop the script, delete the `click_positions.json` file (if created), and restart the recording process.
   - **Visibility:** Ensure Fidelity ATP is on your primary monitor and the window is fully visible and not obscured.

2. **Execution Errors**
   - Check the log file (`trading_bot.log`) for detailed error messages
   - Verify that all UI elements are visible and accessible
   - Re-record positions if the Fidelity ATP interface has changed

3. **Command Errors**
   - Ensure all required parameters are provided
   - Verify parameter values are in the correct format

## Performance Considerations

- The Go implementation typically offers better performance for repeated trades
- For larger orders, using the `--repeat` parameter can help distribute your orders over time
- Adjust the `--min_pause` and `--max_pause` values based on market conditions

## Security

- The bot does not store or transmit your Fidelity credentials
- All operations rely on an already authenticated Fidelity ATP session
- Consider using a dedicated device for automated trading to minimize security risks
