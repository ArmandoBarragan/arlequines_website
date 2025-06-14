package structs

type StripeWebhook struct {
	AmountOfTickets int `json:"amount_of_tickets"`
	PresentationID  int `json:"presentation_id"`
}
