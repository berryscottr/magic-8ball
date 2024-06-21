package selenium

// import (
// 	"fmt"
// 	"strings"
// 	"time"
// 	"database/sql"
// 	"log"
// 	"net/http"

// 	"github.com/go-sql-driver/mysql"
// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/sclevine/agouti"
// )



// func main() {
// 	writeMessageWithBreak("[F]ull Capture")
// 	writeMessageWithBreak("[M]essaging Test")
// 	writeMessageWithBreak("[L]ogging Test")
// 	writeMessageWithBreak("[D]ebug Mode")
// 	writeMessageWithBreak("[Q]uit")

// 	var confirmValue string
// 	fmt.Scanln(&confirmValue)
// 	confirmValue = strings.ToLower(confirmValue)

// 	switch confirmValue {
// 	case "f":
// 		fullCapture()
// 	case "m":
// 		testMessaging()
// 	case "l":
// 		testImportLogging()
// 	case "d":
// 		// DebugDatabaseData()
// 	case "q":
// 	default:
// 		fmt.Println("Invalid option.")
// 	}

// 	fmt.Println("Program ended.")
// }

// func writeMessageWithBreak(message string) {
// 	fmt.Println(message)
// }

// func fullCapture() error {
// 	// Initialize Chrome WebDriver options
// 	chromeOptions := agouti.ChromeOptions("args", []string{
// 		"--headless",
// 		"--log-level=3",
// 		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.6167.85 Safari/537.36",
// 	})

// 	// Create a WebDriver instance
// 	webDriver = agouti.ChromeDriver(chromeOptions)
// 	if err := webDriver.Start(); err != nil {
// 		return fmt.Errorf("failed to start WebDriver: %v", err)
// 	}
// 	defer webDriver.Stop()

// 	// Perform full capture operations
// 	if err := performFullCapture(); err != nil {
// 		return fmt.Errorf("full capture operation failed: %v", err)
// 	}

// 	return nil
// }

// func performFullCapture() error {
// 	utilitiesObj.WriteTitleMessage("FULL CAPTURE - START")

// 	// Navigate to Division Page and perform initial setup
// 	if err := navigateAndSetup(); err != nil {
// 		return fmt.Errorf("navigation and setup failed: %v", err)
// 	}

// 	// Process division nodes
// 	for _, divisionObj := range linksList {
// 		if err := processDivision(divisionObj); err != nil {
// 			fmt.Printf("Error processing division %s: %v\n", divisionObj.divisionIndicator, err)
// 			continue
// 		}
// 	}

// 	// Finalize session and clean up
// 	utilitiesObj.isImportFinished = true
// 	utilitiesObj.isHTMLDataComplete = true
// 	utilitiesObj.SaveImportedSession()

// 	fmt.Println("FULL CAPTURE - END")
// 	return nil
// }

// func navigateAndSetup() error {
// 	// Navigation and setup code from C# method

// 	// Example navigation using WebDriver
// 	page, err := webDriver.NewPage()
// 	if err != nil {
// 		return fmt.Errorf("failed to open new page: %v", err)
// 	}

// 	// Example navigation steps
// 	if err := navigateTo(page, "https://example.com"); err != nil {
// 		return fmt.Errorf("failed to navigate: %v", err)
// 	}

// 	// Additional setup steps

// 	return nil
// }

// func navigateTo(page *agouti.Page, url string) error {
// 	if err := page.Navigate(url); err != nil {
// 		return fmt.Errorf("navigation to %s failed: %v", url, err)
// 	}
// 	time.Sleep(2 * time.Second) // Adjust sleep time as needed
// 	return nil
// }

// func processDivision(divisionObj DivisionDetail) error {
// 	// Process each division as per C# method

// 	// Example: Navigate to Standings page
// 	page, err := webDriver.NewPage()
// 	if err != nil {
// 		return fmt.Errorf("failed to open new page: %v", err)
// 	}

// 	if err := navigateTo(page, "https://example.com/standings"); err != nil {
// 		return fmt.Errorf("failed to navigate to standings: %v", err)
// 	}

// 	// Example: Parse HTML, update data, save data

// 	return nil
// }

// func testMessaging() {
// 	// utilitiesObj.WriteTitleMessage("START")

// 	fmt.Print("TEST WRITE -- 1st line, 1st text")
// 	fmt.Print("----2nd bit of text")
// 	fmt.Println()
// 	fmt.Println("TEST WRITE -- 2nd line")

// 	var input string
// 	fmt.Scanln(&input) // Read user input

// 	// In a real scenario, you might do something with the input
// 	fmt.Println("You entered:", input)

// 	// Wait for user to press Enter before exiting
// 	fmt.Println("Press Enter to exit...")
// 	fmt.Scanln() // Wait for Enter key press
// }

// func testWebPage() {
// 	url := "http://jmersweb.com/_st/test.aspx"

// 	// Make HTTP GET request
// 	response, err := http.Get(url)
// 	if err != nil {
// 		log.Fatalf("Error making GET request: %v", err)
// 	}
// 	defer response.Body.Close()

// 	if response.StatusCode != http.StatusOK {
// 		log.Fatalf("Unexpected status code: %d", response.StatusCode)
// 	}

// 	// Load HTML document
// 	doc, err := goquery.NewDocumentFromReader(response.Body)
// 	if err != nil {
// 		log.Fatalf("Error loading HTML: %v", err)
// 	}

// 	// Extract inner HTML of body element
// 	htmlContent, err := doc.Find("body").Html()
// 	if err != nil {
// 		log.Fatalf("Error extracting HTML: %v", err)
// 	}

// 	// Print or process the extracted HTML content
// 	fmt.Println("HTML Content:")
// 	fmt.Println(htmlContent)
// }

// func testImportLogging() error {
// 	idCounter := 0

// 	// Establish database connection
// 	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/dbname")
// 	if err != nil {
// 		return fmt.Errorf("error connecting to database: %v", err)
// 	}
// 	defer db.Close()

// 	// Clear linksList (assuming it's a global variable or part of a struct)
// 	linksList := make([]DivisionDetail, 0)

// 	// Save import log messages
// 	saveImportLogMessage("Truncated: Import Main Tables")
// 	saveImportLogMessage("Truncated: Import Process Tables")

// 	// Execute stored procedure and process results
// 	rows, err := db.Query("CALL st_Get_Imported_HTML()")
// 	if err != nil {
// 		return fmt.Errorf("error executing stored procedure: %v", err)
// 	}
// 	defer rows.Close()

// 	// Process each row from the result set
// 	for rows.Next() {
// 		var divisionObj DivisionDetail
// 		err := rows.Scan(
// 			&divisionObj.divisionIndicator,
// 			&divisionObj.divisionCode,
// 			&divisionObj.divisionName,
// 		)
// 		if err != nil {
// 			return fmt.Errorf("error scanning row: %v", err)
// 		}

// 		idCounter++
// 		divisionObj.idValue = idCounter
// 		divisionObj.dayOfWeekName = "NONE"
// 		divisionObj.divisionIncomingNameText = ""
// 		linksList = append(linksList, divisionObj)

// 		saveImportLogMessage(fmt.Sprintf("Division Name: %s, Saved", divisionObj.divisionName))
// 	}

// 	// Save session updates
// 	saveImportLogMessage("Import Session Data Updated")

// 	// Process each division object in linksList
// 	for _, divisionObj := range linksList {
// 		saveImportLogMessage(fmt.Sprintf("Division %d: Team HTML Saved", divisionObj.idValue))
// 	}

// 	// Final session update
// 	saveImportLogMessage("Import Session Data Updated")

// 	// Additional actions (e.g., adding blank lines, writing final messages)
// 	fmt.Println()
// 	fmt.Println("End of Test")
// 	fmt.Println()

// 	return nil
// }

// func saveImportLogMessage(message string) {
// 	fmt.Println(message) // Print or log the message
// 	// In a real application, you might save this message to a file or database table
// }

// func enterCredentials() error {
// 	// Set up Agouti WebDriver configuration
// 	options := agouti.ChromeOptions("args", []string{
// 		"--headless",
// 		"--log-level=3",
// 		"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.6167.85 Safari/537.36",
// 	})
// 	driver := agouti.ChromeDriver(options)

// 	// Start the WebDriver session
// 	err := driver.Start()
// 	if err != nil {
// 		return fmt.Errorf("failed to start WebDriver: %v", err)
// 	}
// 	defer driver.Stop()

// 	// Open a new page
// 	page, err := driver.NewPage()
// 	if err != nil {
// 		return fmt.Errorf("failed to open new page: %v", err)
// 	}

// 	// Navigate to the login page
// 	err = page.Navigate("https://example.com/login")
// 	if err != nil {
// 		return fmt.Errorf("failed to navigate to login page: %v", err)
// 	}

// 	// Get credentials from database or configuration
// 	userName := "your_username"
// 	passWord := "your_password"

// 	// Enter credentials and login
// 	err = page.FindByID("email").Fill(userName)
// 	if err != nil {
// 		return fmt.Errorf("failed to fill email field: %v", err)
// 	}

// 	err = page.FindByID("password").Fill(passWord)
// 	if err != nil {
// 		return fmt.Errorf("failed to fill password field: %v", err)
// 	}

// 	err = page.FindByButton("Login").Click()
// 	if err != nil {
// 		return fmt.Errorf("failed to click login button: %v", err)
// 	}

// 	// Handle additional screen with 'Continue' button (if applicable)
// 	// Example: page.FindByButton("Continue").Click()

// 	// Log actions and sleep
// 	fmt.Println("Login Credentials Entered")
// 	durationValue := getRandomSleepValue(500) // Replace with your sleep function
// 	fmt.Printf("Sleeping: %d ms\n", durationValue)
// 	time.Sleep(time.Millisecond * time.Duration(durationValue))

// 	return nil
// }

// // Example function to mimic utilitiesObj.GetRandomSleepValue()
// func getRandomSleepValue(max int) int {
// 	// Replace with your logic to generate random sleep duration
// 	return 500 // Example value
// }

// func navPage(pageTypeIn PageTypes) error {
// 	var breakDurationType time.Duration
// 	navigationMessagePattern := "Navigating to: %v Page, Sleeping: %v ms%v %v"
// 	navigationPageURL := ""
// 	searchText := ""
// 	loopCount := 0
// 	retryAttempts := 5
// 	durationValue := 0

// 	switch pageTypeIn {
// 	case Login:
// 		breakDurationType = time.Millisecond * 500
// 		navigationPageURL = "https://example.com/login"
// 	case Atlanta:
// 		breakDurationType = time.Millisecond * 500
// 		navigationPageURL = "https://example.com/atlanta"
// 	case DivisionMinimal:
// 		breakDurationType = time.Millisecond * 200
// 		navigationPageURL = "https://example.com/division"
// 		searchText = "Division Minimal"
// 	case DivisionShort:
// 		breakDurationType = time.Millisecond * 500
// 		navigationPageURL = "https://example.com/division"
// 		searchText = "Division Short"
// 	case Standings:
// 		breakDurationType = time.Millisecond * 1000
// 		navigationPageURL = standingsURL(currentDivisionId)
// 		searchText = "Standings"
// 	case Roster:
// 		breakDurationType = time.Millisecond * 1000
// 		navigationPageURL = rosterURL(currentDivisionId)
// 		searchText = "Roster"
// 	case Schedule:
// 		breakDurationType = time.Millisecond * 500
// 		navigationPageURL = scheduleURL(currentDivisionId)
// 		searchText = "Schedule"
// 	case MVP:
// 		breakDurationType = time.Millisecond * 500
// 		navigationPageURL = mvpURL(currentDivisionId)
// 		searchText = "MVP"
// 		retryAttempts = 1
// 	}

// 	// Perform the HTTP GET request
// 	client := &http.Client{}
// 	req, err := http.NewRequest("GET", navigationPageURL, nil)
// 	if err != nil {
// 		return fmt.Errorf("error creating request: %v", err)
// 	}

// 	// Set headers, cookies, or other request parameters as needed

// 	// Execute the request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("error making request: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	// Simulate delay
// 	time.Sleep(breakDurationType)

// 	// Check response content if needed
// 	if searchText != "" {
// 		body := make([]byte, 0)
// 		_, err := resp.Body.Read(body)
// 		if err != nil {
// 			return fmt.Errorf("error reading response body: %v", err)
// 		}

// 		for string(body) != searchText && loopCount < retryAttempts {
// 			loopCount++
// 			time.Sleep(time.Millisecond * 1500)
// 			resp, err = client.Do(req)
// 			if err != nil {
// 				return fmt.Errorf("error making retry request: %v", err)
// 			}
// 			defer resp.Body.Close()

// 			// Check response content again
// 			body = make([]byte, 0)
// 			_, err := resp.Body.Read(body)
// 			if err != nil {
// 				return fmt.Errorf("error reading retry response body: %v", err)
// 			}
// 		}
// 	}

// 	fmt.Printf(navigationMessagePattern, pageTypeIn, durationValue, "", "")
// 	fmt.Println()

// 	return nil
// }

// // Helper functions for constructing URLs (Replace with actual URL construction logic)
// func standingsURL(divisionID string) string {
// 	return fmt.Sprintf("https://example.com/standings?division=%s", divisionID)
// }

// func rosterURL(divisionID string) string {
// 	return fmt.Sprintf("https://example.com/roster?division=%s", divisionID)
// }

// func scheduleURL(divisionID string) string {
// 	return fmt.Sprintf("https://example.com/schedule?division=%s", divisionID)
// }

// func mvpURL(divisionID string) string {
// 	return fmt.Sprintf("https://example.com/mvp?division=%s", divisionID)
// }

// func parseSessionText() error {
// 	currentSessionText := "Current Session"
// 	currentSessionErrorText := fmt.Sprintf("No %s Found or Selected - Exiting", currentSessionText)
// 	var currentSessionItemList []string

// 	// Assuming driver is a global variable holding the database connection

// 	// Example SQL query execution
// 	rows, err := dbDriver.Query("SELECT sessionCode, sessionDesc FROM Sessions")
// 	if err != nil {
// 		return fmt.Errorf("error executing query: %v", err)
// 	}
// 	defer rows.Close()

// 	// Build currentSessionItemList from database results
// 	for rows.Next() {
// 		var sessionCode, sessionDesc string
// 		err := rows.Scan(&sessionCode, &sessionDesc)
// 		if err != nil {
// 			return fmt.Errorf("error scanning row: %v", err)
// 		}
// 		if strings.Contains(sessionDesc, currentSessionText) {
// 			currentSessionItemList = append(currentSessionItemList, sessionCode+"|"+sessionDesc)
// 		}
// 	}

// 	switch len(currentSessionItemList) {
// 	case 0:
// 		return fmt.Errorf(currentSessionErrorText) // No Current Session Found Means We Exit

// 	case 1:
// 		// Set utilitiesObj.CurrentSessionCodeText and utilitiesObj.CurrentSessionDescText
// 		currentSessionInfo := strings.Split(currentSessionItemList[0], "|")
// 		currentSessionCodeText = currentSessionInfo[0]
// 		currentSessionDescText = currentSessionInfo[1]
// 		break

// 	default:
// 		// Handle multiple session items
// 		fmt.Println("Choose", currentSessionText, "to Continue:")
// 		for i, sessionItem := range currentSessionItemList {
// 			fmt.Printf("%d) %s\n", i, strings.Replace(sessionItem, "|", " - ", 1))
// 		}

// 		var confirmValueInt int
// 		fmt.Print("Enter your choice: ")
// 		_, err := fmt.Scanln(&confirmValueInt)
// 		if err != nil {
// 			return fmt.Errorf("error reading input: %v", err)
// 		}

// 		if confirmValueInt >= 0 && confirmValueInt < len(currentSessionItemList) {
// 			// Set utilitiesObj.CurrentSessionCodeText and utilitiesObj.CurrentSessionDescText
// 			currentSessionInfo := strings.Split(currentSessionItemList[confirmValueInt], "|")
// 			currentSessionCodeText = currentSessionInfo[0]
// 			currentSessionDescText = currentSessionInfo[1]
// 		} else {
// 			return fmt.Errorf(currentSessionErrorText) // Invalid choice means we exit
// 		}
// 	}

// 	return nil
// }

// func fixDivisionCode(incomingNameText string) string {
// 	returnVal := ""
// 	openParen := strings.Index(incomingNameText, "(")
// 	closeParen := strings.Index(incomingNameText, ")")

// 	if openParen >= 0 && closeParen > openParen {
// 		openParen++
// 		valueLength := closeParen - openParen
// 		if valueLength > 0 {
// 			returnVal = incomingNameText[openParen : openParen+valueLength]
// 		}
// 	}

// 	return returnVal
// }

// func fixDivisionName(divisionCode string, incomingNameText string) string {
// 	return strings.Replace(incomingNameText, "("+divisionCode+")", "", 1)
// }
