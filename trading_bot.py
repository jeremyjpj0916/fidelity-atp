import pyautogui
from pynput import mouse
from pynput import keyboard
import time
import json
import argparse
import random
import os
from concurrent.futures import ThreadPoolExecutor
import logging

# Setup logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    filename='trading_bot.log'
)

CONFIG_FILE = "click_positions.json"

def load_config():
    """Load predefined click positions from a JSON file."""
    try:
        with open(CONFIG_FILE, "r") as file:
            return json.load(file)
    except FileNotFoundError:
        print("Config file not found. Please run in record_positions mode to configure UI coordinates.")
        exit(1)

def save_config(config):
    """Save UI click positions to a JSON file."""
    with open(CONFIG_FILE, "w") as file:
        json.dump(config, file, indent=4)
    print("Configuration saved successfully.")

def record_positions():
    """Record UI click positions by waiting for real mouse clicks (cross-platform)."""
    config = {}
    elements = [
        ("screen_focus", "Click the screen to focus it"),
        ("options_button", "Click the Options button"),
        ("stocks_button", "Click the Stocks button"),
        ("account_dropdown", "Click the account dropdown"),
        ("account_roth", "Click the Roth account option"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_traditional", "Click the Traditional account option"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_hsa", "Click the HSA account option"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_brokeragelink", "Click the Brokeragelink account option"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_dropdown", "Click the account dropdown AGAIN to reopen it"),
        ("account_individual", "Click the Individual account option"),
        ("ticker_box", "Click inside the ticker input box"),
        ("buy_button", "Click the Buy button"),
        ("sell_button", "Click the Sell button"),
        ("amount_box", "Click inside the amount input box"),
        ("market_order", "Click the Market order option"),
        ("limit_order", "Click the Limit order option"),
        ("limit_price_box", "Click inside the Limit price input box"),
        ("day_dropdown", "Click the day dropdown"),
        ("day_option", "Click the day option"),
        ("day_dropdown", "Click the day dropdown AGAIN to reopen it"),
        ("day_dropdown", "Click the day dropdown AGAIN to reopen it"),
        ("day_plus_option", "Click the Day+ option"),
        ("place_order_button", "Click the Place Order button"),
    ]

    print("Move the mouse to each element and LEFT-click to record.")
    print("Press ESC to skip optional elements or Ctrl+C to abort the entire process.\n")

    for idx, (key, prompt) in enumerate(elements, 1):
        print(f"[{idx}/{len(elements)}] {prompt} (press ESC to skip) – waiting for input…")

        recorded = False

        def on_press(key):
            nonlocal recorded
            try:
                if key == keyboard.Key.esc:
                    print(f"Skipped {prompt}\n")
                    recorded = True
                    return False  # stop listener
            except:
                pass  # for any other key press exceptions

        def on_click(x, y, button, pressed):
            nonlocal recorded
            if pressed:
                config[key] = [x, y]
                print(f"Recorded ({x}, {y})\n")
                recorded = True
                return False  # stop listener

        # Use both keyboard and mouse listeners
        with mouse.Listener(on_click=on_click) as m_listener, \
             keyboard.Listener(on_press=on_press) as k_listener:
            while not recorded:
                time.sleep(0.05)

    save_config(config)
    print("All positions saved to click_positions.json")

def human_like_delay(min_time=0.1, max_time=0.8):
    """Pause for a random time to simulate human behavior."""
    time.sleep(random.uniform(min_time/2, max_time/2))  # Reduced by half for faster operation

def click(position, n=1):
    """Perform mouse click(s) at the given position."""
    for _ in range(n):
        pyautogui.click(position[0], position[1])
        human_like_delay(0.01, 0.1)  # Reduced by half

def type_text(text, press_enter=True):
    """Type text into the currently active input field."""
    pyautogui.typewrite(str(text), interval=random.uniform(0.02, 0.05))  # Faster typing
    if press_enter:
        human_like_delay(0.05, 0.15)  # Reduced delay before pressing enter
        pyautogui.press('enter')
    human_like_delay(0.05, 0.4)  # Reduced default delay after typing

def select_account(config, account_type):
    """Select the specified trading account."""
    account_type = account_type.lower()
        
    account_key = f"account_{account_type}"
    if account_key not in config:
        logging.error(f"Account type '{account_type}' not found in configuration")
        print(f"Error: Account type '{account_type}' not found in configuration")
        exit(1)
    
    click(config["account_dropdown"])
    human_like_delay(0.05, 0.4)  # Reduced delay
    click(config[account_key])
    human_like_delay(0.15, 0.3)  # Faster account selection delay

def select_ticker(config, ticker):
    """Enter and select a ticker symbol."""
    click(config["ticker_box"])
    type_text(ticker)
    human_like_delay(0.25, 0.5)  # Faster ticker loading delay

def select_trade_action(config, action):
    """Select buy or sell action."""
    if action.lower() == "buy":
        click(config["buy_button"])
    elif action.lower() == "sell":
        click(config["sell_button"])
    else:
        logging.error(f"Invalid action: {action}")
        print(f"Error: Invalid action {action}. Must be 'buy' or 'sell'")
        exit(1)
    
    human_like_delay(0.15, 0.3)

def enter_amount(config, amount):
    """Enter the trading amount (shares)."""
    click(config["amount_box"])
    # Format as integer for shares
    type_text(int(amount), press_enter=False)
    human_like_delay(0.15, 0.3)

def select_order_type(config, order_type, limit_price=None, extended_hours=False):
    """Select the order type and related settings."""
    if order_type.lower() == "market":
        click(config["market_order"])
    elif order_type.lower() == "limit":
        click(config["limit_order"])
        # Add a longer delay here to give the UI time to respond
        human_like_delay(0.2, 0.4)
        if limit_price is not None:
            click(config["limit_price_box"])
            type_text(limit_price, press_enter=True)
        else:
            logging.error("Limit price is required for limit orders")
            print("Error: Limit price is required for limit orders")
            exit(1)
    else:
        logging.error(f"Invalid order type: {order_type}")
        print(f"Error: Invalid order type {order_type}. Must be 'market' or 'limit'")
        exit(1)
    
    # Handle extended hours trading if requested
    if extended_hours:
        click(config["day_dropdown"])
        human_like_delay(0.15, 0.3)
        click(config["day_plus_option"])
        human_like_delay(0.15, 0.3)
    else:
        # For regular trading hours, explicitly select the Day option
        click(config["day_dropdown"])
        human_like_delay(0.15, 0.3)
        click(config["day_option"])
        human_like_delay(0.15, 0.3)
    
    human_like_delay(0.15, 0.3)

def submit_order(config):
    """Submit the order and confirm it."""
    click(config["place_order_button"])
    # Interface automatically returns to the buying screen, no need to click "new order" button

def execute_single_trade(config, args, batch_num=None, batch_amount=None):
    """Execute a single trade with the given parameters."""
    try:
        amount_to_use = batch_amount if batch_amount is not None else args.amount
        
        # Log the trade
        batch_info = f" (Batch {batch_num})" if batch_num is not None else ""
        logging.info(f"Executing {args.action} order{batch_info}: "
                    f"{int(amount_to_use)} shares of {args.ticker} as {args.order_type} order")
        
        # Select trade type
        click(config["stocks_button"] if args.trade_type == "stocks" else config["options_button"], n=2)
        
        # Select account
        select_account(config, args.account)
        
        # Enter ticker
        select_ticker(config, args.ticker)
        
        # Select buy/sell
        select_trade_action(config, args.action)
        
        # Enter amount (shares)
        enter_amount(config, amount_to_use)
        
        # Select order type and set limit price if needed
        select_order_type(config, args.order_type, args.limit_price, args.extended_hours)
        
        # Submit the order
        submit_order(config)
        
        return True
    except Exception as e:
        logging.error(f"Error executing trade: {str(e)}")
        print(f"Error executing trade: {str(e)}")
        return False

def execute_trade(args):
    """Execute the stock trade based on CLI inputs, potentially repeating the order multiple times."""
    config = load_config()
    
    # If repeat is enabled, execute the same trade multiple times
    if args.repeat > 1:
        print(f"Executing the same {args.action} order for {args.amount} shares {args.repeat} times")
         
        for i in range(args.repeat):
            print(f"Executing order {i+1}/{args.repeat} for {args.amount} shares")
             
            success = execute_single_trade(config, args, i+1, args.amount)
             
            if not success:
                print(f"Failed at repeat {i+1}. Stopping.")
                break
             
            # Pause between repeats
            if i < args.repeat - 1:
                pause_time = random.uniform(args.min_pause, args.max_pause)
                print(f"Pausing for {pause_time:.2f} seconds before next repeat...")
                time.sleep(pause_time)
    else:
        # Execute as a single trade
        execute_single_trade(config, args)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Automated Stock Trading Bot")
    parser.add_argument("mode", choices=["trade", "record_positions"], help="Mode to run the script in")
    
    # Trade parameters
    parser.add_argument("--trade_type", choices=["stocks", "options"], help="Trade Stocks or Options")
    parser.add_argument("--account", choices=["Roth", "Traditional", "HSA", "Brokeragelink", "Individual", "Individual_TOD"], 
                        help="Select account type")
    parser.add_argument("--ticker", type=str, help="Stock ticker symbol")
    parser.add_argument("--action", choices=["buy", "sell"], help="Buy or Sell action")
    parser.add_argument("--amount", type=float, help="Number of shares to trade")
    parser.add_argument("--order_type", choices=["market", "limit"], help="Order type")
    parser.add_argument("--limit_price", type=float, help="Limit price if using a limit order")
    parser.add_argument("--extended_hours", action="store_true", help="Enable extended hours trading (Day+)")
    
    # Repeating parameters
    parser.add_argument("--repeat", type=int, default=1, help="Number of times to repeat the order")
    parser.add_argument("--min_pause", type=float, default=1.0, help="Minimum pause time between repeats (seconds)")
    parser.add_argument("--max_pause", type=float, default=3.0, help="Maximum pause time between repeats (seconds)")

    args = parser.parse_args()

    if args.mode == "record_positions":
        record_positions()
    elif args.mode == "trade":
        # Validate required parameters
        required_params = ["trade_type", "account", "ticker", "action", "amount", "order_type"]
        missing_params = [param for param in required_params if getattr(args, param) is None]
        
        if missing_params:
            print(f"Error: Missing required parameters: {', '.join(missing_params)}")
            parser.print_help()
            exit(1)
            
        # Validate limit price for limit orders
        if args.order_type == "limit" and args.limit_price is None:
            print("Error: Limit price is required for limit orders")
            exit(1)
            
        execute_trade(args)
