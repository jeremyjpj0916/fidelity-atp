package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
	"github.com/go-vgo/robotgo"
	"github.com/robotn/gohook"
)

// Config stores the UI click positions
type Config struct {
	StocksButton           []int `json:"stocks_button"`
	OptionsButton          []int `json:"options_button"`
	AccountDropdown        []int `json:"account_dropdown"`
	AccountRoth            []int `json:"account_roth"`
	AccountTraditional     []int `json:"account_traditional"`
	AccountHSA             []int `json:"account_hsa"`
	AccountBrokeragelink   []int `json:"account_brokeragelink"`
	AccountIndividual      []int `json:"account_individual"`
	TickerBox              []int `json:"ticker_box"`
	BuyButton              []int `json:"buy_button"`
	SellButton             []int `json:"sell_button"`
	AmountBox              []int `json:"amount_box"`
	MarketOrder            []int `json:"market_order"`
	LimitOrder             []int `json:"limit_order"`
	LimitPriceBox          []int `json:"limit_price_box"`
	DayDropdown            []int `json:"day_dropdown"`
	DayOption              []int `json:"day_option"`
	DayPlusOption          []int `json:"day_plus_option"`
	PlaceOrderButton       []int `json:"place_order_button"`
}

// TradeParams holds the parameters for a trade operation
type TradeParams struct {
	Mode          string
	TradeType     string
	Account       string
	Ticker        string
	Action        string
	Amount        float64
	OrderType     string
	LimitPrice    float64
	ExtendedHours bool
	Repeat        int
	MinPause      float64
	MaxPause      float64
}

const (
	configFile = "click_positions.json"
	logFile    = "trading_bot.log"
)

var (
	logger *log.Logger
)

func init() {
	// Setup logging
	logFile, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
	logger.Println("Trading bot started")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
}

// runCommand executes a system command and returns its output
func runCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// mouseClick simulates a mouse click at the given coordinates
func mouseClick(x, y int, times int) error {
	// Move mouse to the position
	robotgo.MoveMouseSmooth(x, y, 0.7, 0.7)
	humanLikeDelay(0.05, 0.1)
	
	// Click the specified number of times
	for i := 0; i < times; i++ {
		robotgo.Click("left")
		humanLikeDelay(0.01, 0.05)
	}
	
	// Allow time for the click to register and focus to change
	humanLikeDelay(0.2, 0.3)
	return nil
}

// typeText types the given text with a human-like interval
func typeText(text string, pressEnter bool) error {
	// First, try to focus the "Fidelity Active Trader Pro" window
	// (This will only work on machines where the app is named exactly this)
	robotgo.ActiveName("Fidelity Active Trader Pro")
	humanLikeDelay(0.05, 0.1)
	
	// Type the text with small delays between characters for reliability
	for _, char := range text {
		robotgo.TypeStr(string(char))
		humanLikeDelay(0.02, 0.05)
	}
 
	humanLikeDelay(0.05, 0.15)
 
	if pressEnter {
		robotgo.KeyTap("enter")
	}

	humanLikeDelay(0.05, 0.15)
 
	return nil
}

// humanLikeDelay pauses for a random time between minTime and maxTime seconds
func humanLikeDelay(minTime, maxTime float64) {
	time.Sleep(time.Duration((rand.Float64()*(maxTime-minTime) + minTime) * float64(time.Second)))
}

// loadConfig loads the click positions from the JSON config file
func loadConfig() (*Config, error) {
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var rawConfig map[string][]int
	if err := json.Unmarshal(configData, &rawConfig); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	config := &Config{
		StocksButton:         rawConfig["stocks_button"],
		OptionsButton:        rawConfig["options_button"],
		AccountDropdown:      rawConfig["account_dropdown"],
		AccountRoth:          rawConfig["account_roth"],
		AccountTraditional:   rawConfig["account_traditional"],
		AccountHSA:           rawConfig["account_hsa"],
		AccountBrokeragelink: rawConfig["account_brokeragelink"],
		AccountIndividual:    rawConfig["account_individual"],
		TickerBox:            rawConfig["ticker_box"],
		BuyButton:            rawConfig["buy_button"],
		SellButton:           rawConfig["sell_button"],
		AmountBox:            rawConfig["amount_box"],
		MarketOrder:          rawConfig["market_order"],
		LimitOrder:           rawConfig["limit_order"],
		LimitPriceBox:        rawConfig["limit_price_box"],
		DayDropdown:          rawConfig["day_dropdown"],
		DayOption:            rawConfig["day_option"],
		DayPlusOption:        rawConfig["day_plus_option"],
		PlaceOrderButton:     rawConfig["place_order_button"],
	}

	return config, nil
}

// saveConfig saves the click positions to the JSON config file
func saveConfig(config map[string][]int) error {
	configData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return fmt.Errorf("error encoding config: %v", err)
	}

	if err := ioutil.WriteFile(configFile, configData, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// recordPositions interactively records UI click positions
func recordPositions() error {
	fmt.Println("Move the mouse to each element and LEFT-click to record. Press CTRL-C to abort.\n")

	config := make(map[string][]int)
	elements := []struct {
		key     string
		message string
	}{
		{"screen_focus", "Click the screen to focus it"},
		{"options_button", "Click the Options button"},
		{"stocks_button", "Click the Stocks button"},
		{"account_dropdown", "Click the account dropdown"},
		{"account_roth", "Click the Roth account option"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_traditional", "Click the Traditional account option"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_hsa", "Click the HSA account option"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_brokeragelink", "Click the Brokeragelink account option"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_dropdown", "Click the account dropdown AGAIN to reopen it"},
		{"account_individual", "Click the Individual account option"},
		{"ticker_box", "Click inside the ticker input box"},
		{"buy_button", "Click the Buy button"},
		{"sell_button", "Click the Sell button"},
		{"amount_box", "Click inside the amount input box"},
		{"market_order", "Click the Market order option"},
		{"limit_order", "Click the Limit order option"},
		{"limit_price_box", "Click inside the Limit price input box"},
		{"day_dropdown", "Click the day dropdown"},
		{"day_option", "Click the day option"},
		{"day_dropdown", "Click the day dropdown AGAIN to reopen it"},
		{"day_dropdown", "Click the day dropdown AGAIN to reopen it"},
		{"day_plus_option", "Click the Day+ option"},
		{"place_order_button", "Click the Place Order button"},
	}

	fmt.Println("Starting event listener...")
	evChan := hook.Start()
	defer hook.StopEvent()
	fmt.Println("Event listener started. Waiting for clicks or ESC key to skip...")

	for i, el := range elements {
		fmt.Printf("[%d/%d] %s (press ESC to skip)\n", i+1, len(elements), el.message)

		clicked := false
		for !clicked {
			select {
			case ev := <-evChan:
				if ev.Kind == hook.MouseDown && ev.Button == hook.MouseMap["left"] {
					x, y := robotgo.GetMousePos()
					config[el.key] = []int{x, y}
					fmt.Printf("Recorded (%d, %d)\n\n", x, y)
					clicked = true // Move to the next element
				} else if ev.Kind == hook.KeyDown && ev.Keycode == 53 { // 53 is ESC key on macOS
					fmt.Printf("Skipped %s\n\n", el.message)
					// Don't record position for this element
					clicked = true // Move to the next element
				}
			default:
				time.Sleep(50 * time.Millisecond) // Avoid busy-waiting
			}
		}
	}

	if err := saveConfig(config); err != nil {
		return err
	}

	fmt.Println("All positions saved to click_positions.json")
	return nil
}

// selectAccount selects the specified trading account
func selectAccount(config *Config, accountType string) error {
	var accountCoords []int

	// Choose the right account coordinates based on account type
	switch strings.ToLower(accountType) {
	case "roth":
		accountCoords = config.AccountRoth
	case "traditional":
		accountCoords = config.AccountTraditional
	case "hsa":
		accountCoords = config.AccountHSA
	case "brokeragelink":
		accountCoords = config.AccountBrokeragelink
	case "individual":
		accountCoords = config.AccountIndividual
	default:
		return fmt.Errorf("unknown account type: %s", accountType)
	}

	if len(accountCoords) != 2 {
		return fmt.Errorf("invalid coordinates for account type: %s", accountType)
	}

	if err := mouseClick(config.AccountDropdown[0], config.AccountDropdown[1], 1); err != nil {
		return err
	}
	humanLikeDelay(0.3, 0.6)

	if err := mouseClick(accountCoords[0], accountCoords[1], 1); err != nil {
		return err
	}
	humanLikeDelay(0.3, 0.6)

	return nil
}

// selectTicker enters and selects a ticker symbol
func selectTicker(config *Config, ticker string) error {
	if err := mouseClick(config.TickerBox[0], config.TickerBox[1], 1); err != nil {
		return err
	}

	if err := typeText(ticker, true); err != nil {
		return err
	}

	humanLikeDelay(0.5, 1.0) // Allow time for ticker to load
	return nil
}

// selectTradeAction selects buy or sell action
func selectTradeAction(config *Config, action string) error {
	var actionCoords []int

	switch strings.ToLower(action) {
	case "buy":
		actionCoords = config.BuyButton
	case "sell":
		actionCoords = config.SellButton
	default:
		return fmt.Errorf("invalid action: %s, must be 'buy' or 'sell'", action)
	}

	if err := mouseClick(actionCoords[0], actionCoords[1], 1); err != nil {
		return err
	}

	humanLikeDelay(0.3, 0.6)
	return nil
}

// enterAmount enters the trading amount (shares)
func enterAmount(config *Config, amount float64) error {
	if err := mouseClick(config.AmountBox[0], config.AmountBox[1], 1); err != nil {
		return err
	}

	if err := typeText(fmt.Sprintf("%.0f", amount), true); err != nil {
		return err
	}

	humanLikeDelay(0.3, 0.6)
	return nil
}

// selectOrderType selects the order type and related settings
func selectOrderType(config *Config, orderType string, limitPrice float64, extendedHours bool) error {
	switch strings.ToLower(orderType) {
	case "market":
		if err := mouseClick(config.MarketOrder[0], config.MarketOrder[1], 1); err != nil {
			return err
		}

		//Click twice just incase
		if err := mouseClick(config.MarketOrder[0], config.MarketOrder[1], 1); err != nil {
			return err
		}

		//Click three times just incase
		if err := mouseClick(config.MarketOrder[0], config.MarketOrder[1], 1); err != nil {
			return err
		}

	case "limit":
		if err := mouseClick(config.LimitOrder[0], config.LimitOrder[1], 1); err != nil {
			return err
		}

		//Click twice just incase
		if err := mouseClick(config.LimitOrder[0], config.LimitOrder[1], 1); err != nil {
			return err
		}

		//Click three times just incase
		if err := mouseClick(config.LimitOrder[0], config.LimitOrder[1], 1); err != nil {
			return err
		}
		
		// Add a longer delay here to give the UI time to respond
		humanLikeDelay(0.6, 0.8) // Increased by 1/4 second

		if err := mouseClick(config.LimitPriceBox[0], config.LimitPriceBox[1], 1); err != nil {
			return err
		}

		// Add a small delay before typing to ensure the input field is ready
		humanLikeDelay(0.25, 0.4)

		if err := typeText(fmt.Sprintf("%.2f", limitPrice), true); err != nil {
			return err
		}
		
		// Add a delay after entering the price to allow the UI to process it
		humanLikeDelay(0.25, 0.4)
	default:
		return fmt.Errorf("invalid order type: %s, must be 'market' or 'limit'", orderType)
	}

	// Handle extended hours trading if requested
	if extendedHours {
		if err := mouseClick(config.DayDropdown[0], config.DayDropdown[1], 1); err != nil {
			return err
		}
		humanLikeDelay(0.3, 0.6)

		if err := mouseClick(config.DayPlusOption[0], config.DayPlusOption[1], 1); err != nil {
			return err
		}
		humanLikeDelay(0.3, 0.6)
	} else {
		if err := mouseClick(config.DayDropdown[0], config.DayDropdown[1], 1); err != nil {
			return err
		}
		humanLikeDelay(0.3, 0.6)
		if err := mouseClick(config.DayOption[0], config.DayOption[1], 1); err != nil {
			return err
		}
		humanLikeDelay(0.3, 0.6)
	}

	humanLikeDelay(0.3, 0.6)
	return nil
}

// submitOrder submits the order and confirms it
func submitOrder(config *Config) error {
	if err := mouseClick(config.PlaceOrderButton[0], config.PlaceOrderButton[1], 1); err != nil {
		return err
	}
	// Interface automatically returns to the buying screen, no need to click "new order" button

	return nil
}

// executeSingleTrade executes a single trade with the given parameters
func executeSingleTrade(config *Config, params *TradeParams, batchNum int, batchAmount float64) error {
	amountToUse := params.Amount
	if batchAmount > 0 {
		amountToUse = batchAmount
	}

	// Log the trade
	batchInfo := ""
	if batchNum > 0 {
		batchInfo = fmt.Sprintf(" (Batch %d)", batchNum)
	}
	logger.Printf("Executing %s order%s: %.0f shares of %s as %s order",
		params.Action, batchInfo, amountToUse, params.Ticker, params.OrderType)

	// Select trade type
	var tradeButton []int
	if params.TradeType == "stocks" {
		tradeButton = config.StocksButton
	} else {
		tradeButton = config.OptionsButton
	}

	if err := mouseClick(tradeButton[0], tradeButton[1], 2); err != nil {
		return fmt.Errorf("error clicking trade type button: %v", err)
	}

	// Select account
	if err := selectAccount(config, params.Account); err != nil {
		return fmt.Errorf("error selecting account: %v", err)
	}

	// Enter ticker
	if err := selectTicker(config, params.Ticker); err != nil {
		return fmt.Errorf("error entering ticker: %v", err)
	}

	// Select buy/sell
	if err := selectTradeAction(config, params.Action); err != nil {
		return fmt.Errorf("error selecting trade action: %v", err)
	}

	// Enter amount (shares)
	if err := enterAmount(config, amountToUse); err != nil {
		return fmt.Errorf("error entering amount: %v", err)
	}

	// Select order type and set limit price if needed
	if err := selectOrderType(config, params.OrderType, params.LimitPrice, params.ExtendedHours); err != nil {
		return fmt.Errorf("error setting order type: %v", err)
	}

	// Submit the order
	if err := submitOrder(config); err != nil {
		return fmt.Errorf("error submitting order: %v", err)
	}

	return nil
}

// executeTrade executes the stock trade based on CLI inputs, potentially repeating the same order multiple times
func executeTrade(params *TradeParams) error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// If repeat is enabled, execute the same trade multiple times
	if params.Repeat > 1 {
		fmt.Printf("Executing the same %s order for %.0f shares %d times\n",
			params.Action, params.Amount, params.Repeat)

		for i := 0; i < params.Repeat; i++ {
			fmt.Printf("Executing order %d/%d for %.0f shares\n",
				i+1, params.Repeat, params.Amount)

			if err := executeSingleTrade(config, params, i+1, params.Amount); err != nil {
				fmt.Printf("Failed at repeat %d: %v. Stopping.\n", i+1, err)
				return err
			}

			// Pause between repeats
			if i < params.Repeat-1 {
				pauseTime := rand.Float64()*(params.MaxPause-params.MinPause) + params.MinPause
				fmt.Printf("Pausing for %.2f seconds before next repeat...\n", pauseTime)
				time.Sleep(time.Duration(pauseTime * float64(time.Second)))
			}
		}
	} else {
		// Execute as a single trade
		if err := executeSingleTrade(config, params, 0, 0); err != nil {
			return err
		}
	}

	return nil
}

// validateParams validates the trade parameters
func validateParams(params *TradeParams) error {
	if params.Mode == "trade" {
		// Check required parameters
		if params.TradeType == "" {
			return fmt.Errorf("trade_type is required")
		}
		if params.Account == "" {
			return fmt.Errorf("account is required")
		}
		if params.Ticker == "" {
			return fmt.Errorf("ticker is required")
		}
		if params.Action == "" {
			return fmt.Errorf("action is required")
		}
		if params.Amount <= 0 {
			return fmt.Errorf("amount must be greater than 0")
		}
		if params.OrderType == "" {
			return fmt.Errorf("order_type is required")
		}
		if params.OrderType == "limit" && params.LimitPrice <= 0 {
			return fmt.Errorf("limit_price is required for limit orders")
		}
	}
	return nil
}

func main() {
	// Define command line flags
	modePtr := flag.String("mode", "", "Mode to run in: 'trade' or 'record_positions'")
	tradeTypePtr := flag.String("trade_type", "stocks", "Trade type: 'stocks' or 'options'")
	accountPtr := flag.String("account", "", "Account type: 'Roth', 'Traditional', 'HSA', 'Brokeragelink', 'Individual'")
	tickerPtr := flag.String("ticker", "", "Stock ticker symbol")
	actionPtr := flag.String("action", "", "Action: 'buy' or 'sell'")
	amountPtr := flag.Float64("amount", 0, "Number of shares to trade")
	orderTypePtr := flag.String("order_type", "", "Order type: 'market' or 'limit'")
	limitPricePtr := flag.Float64("limit_price", 0, "Limit price for limit orders")
	extendedHoursPtr := flag.Bool("extended_hours", false, "Enable extended hours trading (Day+)")
	repeatPtr := flag.Int("repeat", 1, "Number of times to repeat the order")
	minPausePtr := flag.Float64("min_pause", 1.0, "Minimum pause time between repeats (seconds)")
	maxPausePtr := flag.Float64("max_pause", 3.0, "Maximum pause time between repeats (seconds)")

	// Parse the flags
	flag.Parse()

	// Check if mode is provided
	if *modePtr == "" {
		flag.Usage()
		fmt.Println("\nError: mode is required")
		os.Exit(1)
	}

	// Create trading parameters
	params := &TradeParams{
		Mode:          *modePtr,
		TradeType:     *tradeTypePtr,
		Account:       *accountPtr,
		Ticker:        *tickerPtr,
		Action:        *actionPtr,
		Amount:        *amountPtr,
		OrderType:     *orderTypePtr,
		LimitPrice:    *limitPricePtr,
		ExtendedHours: *extendedHoursPtr,
		Repeat:        *repeatPtr,
		MinPause:      *minPausePtr,
		MaxPause:      *maxPausePtr,
	}

	// Execute the appropriate mode
	var err error
	switch params.Mode {
	case "record_positions":
		err = recordPositions()
	case "trade":
		if err = validateParams(params); err == nil {
			err = executeTrade(params)
		}
	default:
		err = fmt.Errorf("invalid mode: %s", params.Mode)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		logger.Fatalf("Error: %v", err)
		os.Exit(1)
	}

	fmt.Println("Operation completed successfully")
	logger.Println("Operation completed successfully")
}
