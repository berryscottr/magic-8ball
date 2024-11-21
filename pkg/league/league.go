package league

import (
	// "context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sclevine/agouti"
	// "github.com/chromedp/chromedp"
)

func (league *Data) Login() error {
		// Set up Agouti WebDriver configuration
	options := agouti.ChromeOptions("args", []string{
		"--headless",
		"--log-level=3",
		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.6167.85 Safari/537.36",
	})
	driver := agouti.ChromeDriver(options)

	// Start the WebDriver session
	err := driver.Start()
	if err != nil {
		return fmt.Errorf("failed to start WebDriver: %v", err)
	}
	defer driver.Stop()

	// Open a new page
	page, err := driver.NewPage()
	if err != nil {
		return fmt.Errorf("failed to open new page: %v", err)
	}

	league.User = Usr{
		Email:    os.Getenv("APA_EMAIL"),
		Password: os.Getenv("APA_PASSWORD"),
	}

	// Check for missing environment variables
	if league.User.Email == "" || league.User.Password == "" {
		league.Err = errors.New("missing APA_EMAIL or APA_PASSWORD environment variable")
		log.Error().Err(league.Err).Msg("Login attempt failed due to missing credentials")
		return league.Err
	}

	log.Info().Msg("Starting login attempt")

	// Navigate to the login page
	err = page.Navigate("https://example.com/login")
	if err != nil {
		return fmt.Errorf("failed to navigate to login page: %v", err)
	}

	// Enter credentials and login
	err = page.FindByID("email").Fill(league.User.Email)
	if err != nil {
		return fmt.Errorf("failed to fill email field: %v", err)
	}

	err = page.FindByID("password").Fill(league.User.Password)
	if err != nil {
		return fmt.Errorf("failed to fill password field: %v", err)
	}

	err = page.FindByButton("Login").Click()
	if err != nil {
		return fmt.Errorf("failed to click login button: %v", err)
	}

	// Handle additional screen with 'Continue' button (if applicable)
	// Example: page.FindByButton("Continue").Click()

	// Log actions and sleep
	fmt.Println("Login Credentials Entered")
	durationValue := 2*time.Second // Replace with your sleep function
	time.Sleep(durationValue)
	fmt.Printf("Sleeping: %d ms\n", durationValue)
	time.Sleep(time.Millisecond * time.Duration(durationValue))

	return nil
// 	ctx, cancel := chromedp.NewContext(context.Background())
// 	defer cancel()

// 	// Increase timeout for the login process
// 	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
// 	defer cancel()

// 	league.User = Usr{
// 		Email:    os.Getenv("APA_EMAIL"),
// 		Password: os.Getenv("APA_PASSWORD"),
// 	}

// 	// Check for missing environment variables
// 	if league.User.Email == "" || league.User.Password == "" {
// 		league.Err = errors.New("missing APA_EMAIL or APA_PASSWORD environment variable")
// 		log.Error().Err(league.Err).Msg("Login attempt failed due to missing credentials")
// 		return
// 	}

// 	log.Info().Msg("Starting login attempt")

// 	// Run the login process with more granular steps and additional logging
// 	league.Err = chromedp.Run(ctx,
// 		// Step 1: Navigate to the login page
// 		chromedp.Navigate(`https://accounts.poolplayers.com/login`),
// 		chromedp.Sleep(2*time.Second), // Small delay to allow page to start loading

// 		// Step 2: Wait until the email input field is visible
// 		chromedp.WaitReady(`input[name="email"]`, chromedp.ByQuery),
// 		chromedp.SendKeys(`input[name="email"]`, league.User.Email, chromedp.ByQuery),
// 		logStep("Entered email"),

// 		// Step 3: Wait until the password input field is visible and enter password
// 		chromedp.WaitReady(`input[name="password"]`, chromedp.ByQuery),
// 		chromedp.SendKeys(`input[name="password"]`, league.User.Password, chromedp.ByQuery),
// 		logStep("Entered password"),

// 		// Step 4: Click the login button and wait briefly
// 		chromedp.WaitReady(`button[type="Login"]`, chromedp.ByQuery),
// 		chromedp.Click(`button[type=Login"]`, chromedp.ByQuery),
// 		chromedp.Sleep(3*time.Second), // Allow time for page transition after login
// 		logStep("Clicked login button"),

// 		// Step 5: Wait for a post-login element that signifies successful login
// 		chromedp.WaitReady(`#dashboard`, chromedp.ByID),
// 		logStep("Dashboard element detected after login"),

// 		// Step 6: Capture HTML and Screenshot if login succeeds
// 		chromedp.ActionFunc(func(ctx context.Context) error {
// 			var htmlContent string
// 			if err := chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery).Do(ctx); err != nil {
// 				log.Error().Err(err).Msg("Failed to capture page HTML")
// 				return err
// 			}
// 			log.Info().Msg("Captured page HTML")

// 			// Save HTML to file
// 			if err := os.WriteFile("page.html", []byte(htmlContent), 0644); err != nil {
// 				log.Error().Err(err).Msg("Failed to save HTML to file")
// 				return err
// 			}

// 			// Capture screenshot
// 			var buf []byte
// 			if err := chromedp.FullScreenshot(&buf, 90).Do(ctx); err != nil {
// 				log.Error().Err(err).Msg("Failed to capture screenshot")
// 				return err
// 			}
// 			if err := os.WriteFile("screenshot.png", buf, 0644); err != nil {
// 				log.Error().Err(err).Msg("Failed to save screenshot to file")
// 				return err
// 			}
// 			log.Info().Msg("Captured screenshot")
// 			return nil
// 		}),
// 	)

// 	// Handle any errors from the login process
// 	if league.Err != nil {
// 		log.Error().Err(league.Err).Msg("Login attempt failed")
// 		return
// 	}
// 	log.Info().Msg("Login attempt complete")
// }

// // logStep is a helper function to log steps during the chromedp flow
// func logStep(step string) chromedp.ActionFunc {
// 	return func(ctx context.Context) error {
// 		log.Info().Msg(step)
// 		return nil
// 	}
// }

// // retry is a helper function to retry operations if they fail, with a delay in between attempts
// func retry(attempts int, delay time.Duration, fn func() error) error {
// 	for i := 0; i < attempts; i++ {
// 			if err := fn(); err != nil {
// 					if i == attempts-1 {
// 							return err
// 					}
// 					time.Sleep(delay)
// 					continue
// 			}
// 			return nil
// 	}
// 	return nil
// }

// // captureErrorState captures HTML and screenshot if login fails
// func captureErrorState(ctx context.Context, pageHTML *string) {
// 	// Capture HTML
// 	if err := chromedp.Run(ctx, chromedp.OuterHTML("html", pageHTML)); err != nil {
// 			log.Error().Err(err).Msg("Failed to capture page HTML")
// 	} else {
// 			log.Error().Str("html", *pageHTML).Msg("Captured page HTML for debugging")
// 	}

// 	// Capture Screenshot
// 	saveScreenshot(ctx, "login_error.png")
// }

// // saveScreenshot takes a screenshot and saves it to the provided file name
// func saveScreenshot(ctx context.Context, fileName string) {
// 	var buf []byte
// 	if err := chromedp.Run(ctx, chromedp.FullScreenshot(&buf, 90)); err != nil {
// 			log.Error().Err(err).Msg("Failed to capture screenshot")
// 			return
// 	}
// 	if err := os.WriteFile(fileName, buf, 0644); err != nil {
// 			log.Error().Err(err).Msg("Failed to save screenshot")
// 	} else {
// 			log.Info().Str("file", fileName).Msg("Screenshot saved successfully")
// 	}
// }
}