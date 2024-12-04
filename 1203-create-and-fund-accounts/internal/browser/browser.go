package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

type browser struct {
	ctx context.Context
}

func NewBrowser() *browser {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, _ := chromedp.NewContext(allocCtx)

	return &browser{
		ctx: ctx,
	}
}

func (b *browser) RequestTokens(pemFilePath string) error {
	var err error

	// Navigate to unlock page first
	err = chromedp.Run(b.ctx,
		chromedp.Navigate("https://testnet-wallet.multiversx.com/unlock"),
	)
	if err != nil {
		return fmt.Errorf("failed to access unlock page: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Navigate to PEM unlock page
	err = chromedp.Run(b.ctx,
		chromedp.Navigate("https://testnet-wallet.multiversx.com/unlock/pem"),
		chromedp.WaitReady(`input[data-testid="walletFile"]`),
	)
	if err != nil {
		return fmt.Errorf("failed to access PEM unlock page: %w", err)
	}

	// Upload PEM file
	err = chromedp.Run(b.ctx,
		chromedp.SetUploadFiles(`input[data-testid="walletFile"]`, []string{pemFilePath}, chromedp.ByQuery),
	)
	if err != nil {
		return fmt.Errorf("failed to upload PEM file: %w", err)
	}

	// Click access wallet button
	err = chromedp.Run(b.ctx,
		chromedp.WaitVisible(`button[data-testid="submitButton"]`),
		chromedp.Click(`button[data-testid="submitButton"]`),
	)
	if err != nil {
		return fmt.Errorf("failed to click access wallet button: %w", err)
	}

	// Wait for wallet to load
	time.Sleep(1 * time.Second)

	// Navigate to faucet
	err = chromedp.Run(b.ctx,
		chromedp.Navigate("https://testnet-wallet.multiversx.com/faucet"),
		chromedp.WaitVisible(`iframe[title="reCAPTCHA"]`),
	)
	if err != nil {
		return fmt.Errorf("failed to navigate to faucet: %w", err)
	}

	fmt.Println("Waiting captcha...")

	// Click request tokens button
	err = chromedp.Run(b.ctx,
		chromedp.WaitVisible(`button[data-testid="requestTokensButton"]`),
		chromedp.WaitEnabled(`button[data-testid="requestTokensButton"]`),
		chromedp.Click(`button[data-testid="requestTokensButton"]`),
	)
	if err != nil {
		return fmt.Errorf("failed to request tokens: %w", err)
	}

	// Wait for request to complete
	time.Sleep(1 * time.Second)

	// Navigate to logout page
	err = chromedp.Run(b.ctx,
		chromedp.Navigate("https://testnet-wallet.multiversx.com/logout"),
		chromedp.WaitReady("body"),
	)
	if err != nil {
		return fmt.Errorf("failed to navigate to logout page: %w", err)
	}

	return nil
}

func (fr *browser) Close() {
	chromedp.Cancel(fr.ctx)
}
