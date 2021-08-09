package service

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/venkata6/helpschool/api/dto"
	"github.com/venkata6/helpschool/api/request"
	"github.com/venkata6/helpschool/api/response"
	"github.com/venkata6/helpschool/api/util"
	"net/http"
	"strconv"
	"time"
)

type UserDonationsService interface {
	CreateUserDonations(w http.ResponseWriter, r *http.Request)
	GetUserDonations(w http.ResponseWriter, r *http.Request)
	UpdateUserDonations(w http.ResponseWriter, r *http.Request)
	DeleteUserDonations(w http.ResponseWriter, r *http.Request)
}

type UserDonationsServiceInternal struct {
	db *pgxpool.Pool
}

func NewUserDonationsService(db *pgxpool.Pool) UserDonationsService {
	return &UserDonationsServiceInternal{db: db}
}

// CreateCountries persists the posted Article and returns it
// back to the client as an acknowledgement.
func (a *UserDonationsServiceInternal) CreateUserDonations(w http.ResponseWriter, r *http.Request) {
	data := &request.UserDonationsRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, util.ErrInvalidRequest(err))
		return
	}

	if len(data.UserEmail) == 0 {
		render.Render(w, r, util.ErrInvalidRequest(errors.New("empty UserEmail")))
		return
	}

	if len(data.SchoolId) == 0 {
		render.Render(w, r, util.ErrInvalidRequest(errors.New("empty SchoolId")))
		return
	}
	if len(data.SupplyId) == 0 {
		render.Render(w, r, util.ErrInvalidRequest(errors.New("empty SupplyId")))
		return
	}
	schoolId, _ := uuid.Parse(data.SchoolId)
	supplyId, _ := uuid.Parse(data.SupplyId)

	if _, err := a.db.Exec(context.Background(),
		`INSERT INTO helpschool.users_donations( user_email,user_id,user_name,school_id,supply_id,quantity,status,tracking_url,extra_info)
				VALUES ( $1, $2, $3, $4, $5, $6,$7,$8,$9)`,data.UserEmail,data.UserId,data.UserName, schoolId, supplyId,
				data.Quantity, data.Status, data.TrackingUrl, data.ExtraInfo); err == nil {
		w.WriteHeader(http.StatusCreated)
		render.DefaultResponder(w, r, render.M{"status": "created"})
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		render.DefaultResponder(w, r, render.M{"status": "not created"})
	}
}
func (a *UserDonationsServiceInternal) GetUserDonations(w http.ResponseWriter, r *http.Request) {

	schoolId := chi.URLParam(r, "schoolId")
	var count int
	rowCount, err := a.db.Query(context.Background(), "select count(*) as count from" +
			"  helpschool.school_supplies where school_id = $1 ",schoolId)
	for rowCount.Next() {
		_ = rowCount.Scan(&count)
		//checkErr(err)
	}
	rows, err := a.db.Query(context.Background(), "select su.title,su.description,su.url,ss.school_id,ss.supply_id," +
		"ss.quantity,ss.fulfilled_count,ss.extra_info,ss.created_date from helpschool.school_supplies as ss" +
		" inner join helpschool.supplies as su on ss.supply_id = su.supply_id where  ss.school_id = $1",schoolId)
	defer rows.Close()

	schoolSupplies := make([]response.SchoolSuppliesResponse, count)
	i := 0
	for rows.Next() {
		// Read
		var title string
		var description string
		var url string
		var supplyId string
		var schoolId string
		var quantity int
		var fulfilledCount int
		var extraInfo string
		var postedDate time.Time

		err = rows.Scan( &title,&description,&url,&schoolId,&supplyId, &quantity, &fulfilledCount, &extraInfo,&postedDate)
		schoolSupplies[i].SchoolSupplies = &dto.SchoolSupplies{} // allocate space
		schoolSupplies[i].Title = title
		schoolSupplies[i].Description = description
		schoolSupplies[i].Url = url
		schoolSupplies[i].SupplyId = supplyId
		schoolSupplies[i].SchoolId = schoolId
		schoolSupplies[i].Quantity = strconv.Itoa(quantity)
		schoolSupplies[i].FulfilledCount = strconv.Itoa(fulfilledCount)
		schoolSupplies[i].ExtraInfo = extraInfo
		schoolSupplies[i].PostedDate = postedDate

		if err != nil {
			return
		}
		i++
	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return
	}
	if err := render.RenderList(w, r, NewSchoolSuppliesListResponse(schoolSupplies)); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}

func (a *UserDonationsServiceInternal) UpdateUserDonations(w http.ResponseWriter, r *http.Request) {

	// let us return 3 newest entries as featured for now

	var count =3
	rows, err := a.db.Query(context.Background(), "select su.title,su.description,su.url,ss.school_id," +
		"ss.supply_id,ss.quantity,ss.fulfilled_count,ss.extra_info,ss.created_date from helpschool.school_supplies as " +
		"ss inner join helpschool.supplies as su on ss.supply_id = su.supply_id order by ss.created_date desc limit 3")
	defer rows.Close()

	schoolSupplies := make([]response.SchoolSuppliesResponse, count)
	i := 0
	for rows.Next() {
		// Read
		var title string
		var description string
		var url string
		var supplyId string
		var schoolId string
		var quantity int
		var fulfilledCount int
		var extraInfo string
		var postedDate time.Time

		err = rows.Scan( &title,&description,&url,&schoolId,&supplyId, &quantity, &fulfilledCount, &extraInfo,&postedDate)
		schoolSupplies[i].SchoolSupplies = &dto.SchoolSupplies{} // allocate space
		schoolSupplies[i].Title = title
		schoolSupplies[i].Description = description
		schoolSupplies[i].Url = url
		schoolSupplies[i].SupplyId = supplyId
		schoolSupplies[i].SchoolId = schoolId
		schoolSupplies[i].Quantity = strconv.Itoa(quantity)
		schoolSupplies[i].FulfilledCount = strconv.Itoa(fulfilledCount)
		schoolSupplies[i].ExtraInfo = extraInfo
		schoolSupplies[i].PostedDate = postedDate

		if err != nil {
			return
		}
		i++
	}
	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		return
	}
	if err := render.RenderList(w, r, NewSchoolSuppliesListResponse(schoolSupplies)); err != nil {
		render.Render(w, r, util.ErrRender(err))
		return
	}
}

func (a *UserDonationsServiceInternal) DeleteUserDonations(w http.ResponseWriter, r *http.Request) {
	//render.RenderList(w, r, NewCountriesListResponse(articles))
}

func NewUserDonationsListResponse(schoolSupplies []response.SchoolSuppliesResponse) []render.Renderer {
	list := []render.Renderer{}
	for _, schoolSupply := range schoolSupplies {
		list = append(list, schoolSupply)
	}
	return list
}
