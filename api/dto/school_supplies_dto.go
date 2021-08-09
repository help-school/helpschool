package dto

import "time"

type UserDonations struct {
	UserEmail 			string `json:"user_email"`
	UserId    			string `json:"user_id"`
	UserName 			string `json:"user_name"`
	SchoolId       		string `json:"school_id"`
	SupplyId       		string `json:"supply_id"`
	Quantity       		string `json:"quantity"`
	Status		   		string `json:"status"`
	TrackingUrl		   	string `json:"tracking_url"`
	CreatedDate		   	string `json:"created_date"`
	ExtraInfo      		string `json:"extra_info"`
	ModifiedDate	   	time.Time `json:"modified_date"`
}
