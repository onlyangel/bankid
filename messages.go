package bankid

// Very simple i18n language mapping

import (
	"fmt"
	"strings"
)

// Messages
const (
	RFA1    = "RFA1"
	RFA2    = "RFA2"
	RFA3    = "RFA3"
	RFA4    = "RFA4"
	RFA5    = "RFA5"
	RFA6    = "RFA6"
	RFA8    = "RFA8"
	RFA9    = "RFA9"
	RFA13   = "RFA13"
	RFA14_A = "RFA14_A"
	RFA14_B = "RFA14_B"
	RFA15_A = "RFA15_A"
	RFA15_B = "RFA15_B"
	RFA16   = "RFA16"
	RFA17_A = "RFA17_A"
	RFA17_B = "RFA17_B"
	RFA18   = "RFA18"
	RFA19   = "RFA19"
	RFA20   = "RFA20"
	RFA21   = "RFA21"
	RFA22   = "RFA22"
)

// Messages in Swedish
var messages_SE = map[string]string{
	RFA1:    "Starta BankID-appen",
	RFA2:    "Du har inte BankID-appen installerad. Kontakta din internetbank.",
	RFA3:    "Åtgärden avbruten. Försök igen.",
	RFA4:    "En identifiering eller underskrift för det här personnumret är redan påbörjad. Försök igen.",
	RFA5:    "Internt tekniskt fel. Försök igen.",
	RFA6:    "Åtgärden avbruten.",
	RFA8:    "BankID-appen svarar inte. Kontrollera att den är startad och att du har internetanslutning. Om du inte har något giltigt BankID kan du hämta ett hos din Bank. Försök sedan igen.",
	RFA9:    "Skriv in din säkerhetskod i BankID- appen och välj Legitimera eller Skriv under.",
	RFA13:   "Försöker starta BankID-appen.",
	RFA14_A: "Söker efter BankID, det kan ta en liten stund...\nOm det har gått några sekunder och inget BankID har hittats har du sannolikt inget BankID som går att använda för den aktuella identifieringen/underskriften i den här datorn. Om du har ett BankID- kort, sätt in det i kortläsaren. Om du inte har något BankID kan du hämta ett hos din internetbank. Om du har ett BankID på en annan enhet kan du starta din BankID-app där.",
	RFA14_B: "Söker efter BankID, det kan ta en liten stund...\nOm det har gått några sekunder och inget BankID har hittats har du sannolikt inget BankID som går att använda för den aktuella identifieringen/underskriften i den här enheten. Om du inte har något BankID kan du hämta ett hos din internetbank. Om du har ett BankID på en annan enhet kan du starta din BankID-app där.",
	RFA15_A: "Söker efter BankID, det kan ta en liten stund...\nOm det har gått några sekunder och inget BankID har hittats har du sannolikt inget BankID som går att använda för den aktuella identifieringen/underskriften i den här datorn. Om du har ett BankID- kort, sätt in det i kortläsaren. Om du inte har något BankID kan du hämta ett hos din internetbank.",
	RFA15_B: "Söker efter BankID, det kan ta en liten stund...\nOm det har gått några sekunder och inget BankID har hittats har du sannolikt inget BankID som går att använda för den aktuella identifieringen/underskriften i den här enheten. Om du inte har något BankID kan du hämta ett hos din internetbank.",
	RFA16:   "Det BankID du försöker använda är för gammalt eller spärrat. Använd ett annat BankID eller hämta ett nytt hos din internetbank.",
	RFA17_A: "BankID-appen verkar inte finnas i din dator eller telefon. Installera den och hämta ett BankID hos din internetbank. Installera appen från din appbutik eller https://install.bankid.com.",
	RFA17_B: "Misslyckades att läsa av QR koden. Starta BankID-appen och läs av QR koden.\nOm du inte har BankID-appen måste du installera den och hämta ett BankID hos din internetbank. Installera appen från din appbutik eller https://install.bankid.com.",
	RFA18:   "Starta BankID-appen",
	RFA19:   "Vill du identifiera dig eller skriva under med BankID på den här datorn eller med ett Mobilt BankID?",
	RFA20:   "Vill du identifiera dig eller skriva under med ett BankID på den här enheten eller med ett BankID på en annan enhet?",
	RFA21:   "Identifiering eller underskrift pågår.",
	RFA22:   "Okänt fel. Försök igen.",
}

// Messages in English
var messages_EN = map[string]string{
	RFA1:    "Start your BankID app.",
	RFA2:    "The BankID app is not installed. Please contact your internet bank.",
	RFA3:    "Action cancelled. Please try again.",
	RFA4:    "An identification or signing for this personal number is already started. Please try again.",
	RFA5:    "Internal error. Please try again.",
	RFA6:    "Action cancelled.",
	RFA8:    "The BankID app is not responding. Please check that the program is started and that you have internet access. If you don’t have a valid BankID you can get one from your bank. Try again.",
	RFA9:    "Enter your security code in the BankID app and select Identify or Sign.",
	RFA13:   "Trying to start your BankID app.",
	RFA14_A: "Searching for BankID:s, it may take a little while...\nIf a few seconds have passed and still no BankID has been found, you probably don’t have a BankID which can be used for this identification/signing on this computer. If you have a BankID card, please insert it into your card reader. If you don’t have a BankID you can order one from your internet bank. If you have a BankID on another device you can start the BankID app on that device.",
	RFA14_B: "Searching for BankID:s, it may take a little while...\nIf a few seconds have passed and still no BankID has been found, you probably don’t have a BankID which can be used for this identification/signing on this device. If you don’t have a BankID you can order one from your internet bank. If you have a BankID on another device you can start the BankID app on that device.",
	RFA15_A: "Searching for BankID:s, it may take a little while...\nIf a few seconds have passed and still no BankID has been found, you probably don’t have a BankID which can be used for this identification/signing on this computer. If you have a BankID card, please insert it into your card reader. If you don’t have a BankID you can order one from your internet bank.",
	RFA15_B: "Searching for BankID:s, it may take a little while...\nIf a few seconds have passed and still no BankID has been found, you probably don’t have a BankID which can be used for this identification/signing on this device. If you don’t have a BankID you can order one from your internet bank.",
	RFA16:   "The BankID you are trying to use is revoked or too old. Please use another BankID or order a new one from your internet bank.",
	RFA17_A: "The BankID app couldn’t be found on your computer or mobile device. Please install it and order a BankID from your internet bank. Install the app from your app store or https://install.bankid.com.",
	RFA17_B: "Failed to scan the QR code. Start the BankID app and scan the QR code. If you don't have the BankID app, you need to install it and order a BankID from your internet bank. Install the app from your app store or https://install.bankid.com.",
	RFA18:   "Start the BankID app",
	RFA19:   "Would you like to identify yourself or sign with a BankID on this computer or with a Mobile BankID?",
	RFA20:   "Would you like to identify yourself or sign with a BankID on this device or with a BankID on another device?",
	RFA21:   "Identification or signing in progress.",
	RFA22:   "Unknown error. Please try again.",
}

// Messages - keep track of the user facing messages for the language we choose
type Messages struct {
	msgs map[string]string
}

// NewMessages - instance with messages in the provided language
func NewMessages(lang string) (*Messages, error) {
	var messages map[string]string = nil

	switch strings.ToLower(lang) {
	case "se":
		messages = messages_SE
	case "en":
		messages = messages_EN
	default:
		return nil, fmt.Errorf("%s it not a supported language", lang)
	}

	return &Messages{
		msgs: messages,
	}, nil
}

// Msg - pick out the messages string for the provided key
// Note: No error handling here, missing keys will return an empty string
func (m *Messages) Msg(key string) string {
	return m.msgs[key]
}
