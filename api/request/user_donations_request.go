package request

import (
	"net/http"
)

type UserDonationsRequest struct {
	UserEmail 			string `json:"user_email"`
	UserId    			string `json:"user_id"`
	UserName 			string `json:"user_name"`
	SchoolId       		string `json:"school_id"`
	SupplyId       		string `json:"supply_id"`
	Quantity       		string `json:"quantity"`
	Status		   		string `json:"status"`
	TrackingUrl		   	string `json:"tracking_url"`
	ExtraInfo      		string `json:"extra_info"`
}

func (a *UserDonationsRequest) Bind(r *http.Request) error {
	return nil
}
