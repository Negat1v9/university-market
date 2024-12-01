package tgpayment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	baseTgUrl               = "https://api.telegram.org/bot"
	methodCreateInvoiceLink = "/createInvoiceLink"
)

func CreateInvoiceLink(ctx context.Context, paymentID, botToken string, amount int) (string, error) {

	amountValue, err := parseLabelPricesAsJsonValue(amount)
	if err != nil {
		return "", err
	}

	out := url.Values{}
	out.Add("title", "Purchasing a subscription")
	out.Add("description", "You make a one-time payment to renew your subscription, you must purchase the subscription yourself again")
	out.Add("payload", paymentID)
	out.Add("currency", "XTR")
	out.Add("prices", amountValue)

	url := baseTgUrl + botToken + methodCreateInvoiceLink

	cl := http.Client{}
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(out.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := cl.Do(req)
	if err != nil {
		return "", err
	}

	var apiResponse ApiTgResponse

	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return "", err
	}

	if !apiResponse.Ok {
		return "", fmt.Errorf("error create invoice link")
	}

	link, err := apiResponse.Result.MarshalJSON()
	if err != nil {
		return "", err
	}
	// link contain " characters
	return strings.ReplaceAll(string(link), "\"", ""), nil
}

func parseLabelPricesAsJsonValue(amount int) (string, error) {

	val, err := json.Marshal([]LabeledPrice{{Label: "XTR", Amount: amount}})
	if err != nil {
		return "", fmt.Errorf("server error")
	}

	return string(val), nil
}
