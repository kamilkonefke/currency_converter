package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/joho/godotenv"
)

var (
    currency_from string
    currency_to string
    rate_from float64
    rate_to float64
    amount string
)

type ExchangeRates struct {
	Disclaimer string             `json:"disclaimer"`
	License    string             `json:"license"`
	Timestamp  int64              `json:"timestamp"`
	Base       string             `json:"base"`
	Rates      map[string]float64 `json:"rates"`
}

var currencies = huh.NewOptions(
    "EUR",
    "USD",
    "GBP",
    "JPY",
    "PLN",
    "CZK",
    "CHF",
    "RUB",
)

var api string;

func main() {
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewNote().Description("Currency converter app!"),
            huh.NewSelect[string]().Title("From").
            Options(currencies...).
            Height(8).
            Value(&currency_from),

            huh.NewSelect[string]().Title("To").
            Options(currencies...).
            Height(8).
            Value(&currency_to),

            huh.NewInput().Value(&amount).Title("Amount"),
        ),
    )

    err := form.Run()
    if err != nil {
        fmt.Println("Uh oh:", err)
        os.Exit(1);
    }

    godotenv.Load();
    app_id := os.Getenv("APP_ID")

    // Send request
    resp, err := http.Get(fmt.Sprintf("https://openexchangerates.org/api/latest.json?app_id=%s", app_id))
    if err != nil {
        log.Fatal(err)
    } 
    defer resp.Body.Close()

    // Status 200
    if resp.StatusCode == http.StatusOK {
        body_bytes, err := io.ReadAll(resp.Body)
        if err != nil {
            log.Fatal(err)
        }

        var data ExchangeRates
        if err := json.Unmarshal([]byte(body_bytes), &data); err != nil {
            log.Fatal(err)
        }

        for currency, rate := range data.Rates {
            if currency == currency_from {
                rate_from = rate
            }
            if currency == currency_to {
                rate_to = rate
            }
        }

        amount_float, err := strconv.ParseFloat(amount, 64);
        if err != nil {
            fmt.Println("Amount is not a number!\n", err)
            return;
        }

        amount_result := rate_to / rate_from * amount_float
        fmt.Printf("%s %s -> %.2f %s \n", amount, currency_from, amount_result, currency_to)
    }
}
