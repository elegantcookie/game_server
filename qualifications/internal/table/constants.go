package table

import "time"

const (
	notifyMangerURL    = "http://localhost:10007/api/manager/"
	createTicketURL    = "http://localhost:10004/api/tickets"
	addTicketToUserURL = "http://localhost:10002/api/users/tickets"
	typeQualifications = "qualifications"
	ticketPrize        = 108
	playersAmount      = 12
	timeDelta          = 6 * time.Hour
)
