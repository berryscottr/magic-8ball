package selenium

// import (
// 	"fmt"
// 	"math/rand"
// 	"strings"
// 	"time"
// )


// // NewUtilities initializes a new Utilities instance
// func NewUtilities(connectionString string) *Utilities {
// 	se := NewSqlDataExecute(connectionString)
// 	return &Utilities{
// 		se: se,
// 	}
// }

// // WriteTitleMessage writes a title message to console
// func (u *Utilities) WriteTitleMessage(messageIn string) {
// 	if len(messageIn) > 0 {
// 		titleBorder := fmt.Sprintf("=====%s=====", strings.Repeat("=", len(messageIn)))
// 		fmt.Println()
// 		fmt.Println(titleBorder)
// 		fmt.Println(" :: " + messageIn)
// 		fmt.Println(titleBorder)
// 		fmt.Println()
// 		fmt.Println()
// 	}
// }

// // AddBlankLine adds a blank line to console output
// func (u *Utilities) AddBlankLine() {
// 	fmt.Println()
// }

// // WriteMessage writes a message to console with optional formatting
// func (u *Utilities) WriteMessage(message1In string, message2In string, formatIn MessageFormatOption, includeLineBreak bool) {
// 	switch formatIn {
// 	case TwoElement:
// 		fmt.Printf("%s :: %s", time.Now().Format("2006-01-02 15:04:05"), message1In)
// 	case ThreeElement:
// 		fmt.Printf("%s :: %s :: %s", time.Now().Format("2006-01-02 15:04:05"), message1In, message2In)
// 	case None:
// 		fmt.Print(message1In)
// 	}

// 	if includeLineBreak {
// 		fmt.Println()
// 	}
// }

// // SaveImportLogMessage saves an import log message to the database
// func (u *Utilities) SaveImportLogMessage(messageTextIn string) {
// 	u.se.Parameters.Clear()
// 	u.se.SetParameter("@logMessage", sql.NVarChar, messageTextIn, 100)
// 	u.se.SetParameter("@isImportItem", sql.Bit, true)
// 	u.se.ExecStoredProc("st_Save_Imported_LogItem")
// }

// // SaveImportedSession saves imported session data to the database
// func (u *Utilities) SaveImportedSession() {
// 	u.se.Parameters.Clear()
// 	u.se.SetParameter("@sessionCode", sql.VarChar, u.CurrentSessionCodeText, 10)
// 	if len(u.CurrentSessionDescText) > 0 {
// 		u.se.SetParameter("@sessionDescription", sql.VarChar, u.CurrentSessionDescText, 50)
// 	}
// 	if u.includeDivisionHTML {
// 		u.se.SetParameter("@divisionHTML", sql.VarChar, u.RawDivisionPageHTML, -1)
// 	}
// 	// Set other parameters as needed

// 	u.se.SetParameter("@importedSessionIdOut", sql.Int, -1, 4, sql.ParameterDirectionOutput)
// 	u.se.ExecStoredProc("st_Save_Imported_Session")

// 	if u.getSessionId {
// 		u.importedSessionId = u.se.GetParameter("@importedSessionIdOut").(int)
// 	}

// 	u.SaveImportLogMessage("Import Session Data Updated")
// }

// // TruncateImportTables truncates import-related tables in the database
// func (u *Utilities) TruncateImportTables() {
// 	u.se.Parameters.Clear()

// 	// Truncate main import tables
// 	u.se.ExecStoredProc("st_Truncate_ImportTables")
// 	u.SaveImportLogMessage("Truncated: Import Main Tables")

// 	// Truncate processed import tables
// 	u.se.ExecStoredProc("st_Truncate_ProcessedImportTables")
// 	u.SaveImportLogMessage("Truncated: Import Process Tables")
// }

// // UpdateDivisionDayNumber updates division day number in the database
// func (u *Utilities) UpdateDivisionDayNumber(divisionId int, dayName string) {
// 	u.se.Parameters.Clear()
// 	u.se.SetParameter("@divisionId", sql.Int, divisionId)
// 	u.se.SetParameter("@dayNumber", sql.TinyInt, u.ParseSQLDayNumber(dayName))
// 	u.se.ExecStoredProc("st_Update_Imported_Division_DayNumber")
// }

// // ParseSQLDayNumber parses day name to SQL day number
// func (u *Utilities) ParseSQLDayNumber(dayNameIn string) int {
// 	switch dayNameIn {
// 	case "Sunday":
// 		return 0
// 	case "Monday":
// 		return 1
// 	case "Tuesday":
// 		return 2
// 	case "Wednesday":
// 		return 3
// 	case "Thursday":
// 		return 4
// 	case "Friday":
// 		return 5
// 	case "Saturday":
// 		return 6
// 	default:
// 		return 255
// 	}
// }

// // GetRandomSleepValue returns a random sleep duration based on SleepDuration
// func (u *Utilities) GetRandomSleepValue(durationIn SleepDuration) int {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	var lowerValue, upperValue int

// 	switch durationIn {
// 	case Minimal:
// 		lowerValue, upperValue = 1361, 4521
// 	case Short:
// 		lowerValue, upperValue = 3821, 6344
// 	case Medium:
// 		lowerValue, upperValue = 7027, 14639
// 	case Large:
// 		lowerValue, upperValue = 19711, 34972
// 	case Eternal:
// 		lowerValue, upperValue = 40132, 89390
// 	}

// 	randomValue := r.Intn(upperValue-lowerValue+1) + lowerValue
// 	randomValue *= r.Intn(2) + 1 // Multiply by 1 or 2
// 	return randomValue
// }

// // Sleep pauses the execution for a random duration
// func (u *Utilities) Sleep() {
// 	u.Sleep(Medium)
// }

// // Sleep pauses the execution for a random duration based on SleepDuration
// func (u *Utilities) Sleep(duration SleepDuration) {
// 	randomValue := u.GetRandomSleepValue(duration)
// 	time.Sleep(time.Duration(randomValue) * time.Millisecond)
// }
