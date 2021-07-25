// Boot the server:
// ----------------
// $ go run main.go

package main

import (
	"context"
	"embed"
	_ "embed"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/venkata6/helpschool/api/service"

	// "time"
)

var routes = flag.Bool("routes", false, "Generate router documentation")
var db *pgxpool.Pool
var bProd = false

//go:embed web
var webFSRoot embed.FS
var webFS = fsMustSub(webFSRoot, "web")

func main() {

	flag.Parse()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.NoCache)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// add CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	// connect to the database and setup the connection pool for services to use
	setUpDatabaseConnection()

	// RESTy routes for "countries" resource
	countryService := service.NewCountriesService(db)

	startedAt := time.Now()
	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "up and running for %s", time.Since(startedAt))
	})

	r.Route("/api/countries", func(r chi.Router) {
		r.With(paginate).Get("/", countryService.GetCountries)
		r.Post("/", countryService.CreateCountries)   // POST /countries
		r.Delete("/", countryService.DeleteCountries) // DELETE /countries
	})

	statesService := service.NewStatesService(db)
	// // RESTy routes for "states" resource
	r.Route("/api/states", func(r chi.Router) {
		r.With(paginate).Get("/", statesService.GetStates)
		r.Post("/", statesService.CreateStates)   // POST /countries
		r.Delete("/", statesService.DeleteStates) // DELETE /countries
	})
	//
	// // RESTy routes for "districts" resource
	districtsService := service.NewDistrictsService(db)
	r.Route("/api/districts", func(r chi.Router) {
		r.With(paginate).Get("/state/{stateId}", districtsService.GetDistricts)
		r.Post("/", districtsService.CreateDistricts)   // POST /countries
		r.Delete("/", districtsService.DeleteDistricts) // DELETE /countries
	})
	//
	// // RESTy routes for "schools" resource
	schoolsService := service.NewSchoolsService(db)
	r.Route("/api/schools", func(r chi.Router) {
		r.With(paginate).Get("/district/{districtId}", schoolsService.GetSchools)
		r.Post("/", schoolsService.CreateSchools)   // POST /countries
		r.Delete("/", schoolsService.DeleteSchools) // DELETE /countries
	})

	// // RESTy routes for "supplies" resource
	suppliesService := service.NewSuppliesService(db)
	r.Route("/api/supplies", func(r chi.Router) {
		r.With(paginate).Get("/", suppliesService.GetSupplies)
		r.Post("/", suppliesService.CreateSupplies)   // POST /countries
		r.Delete("/", suppliesService.DeleteSupplies) // DELETE /countries
	})

	// // RESTy routes for "supplies" resource
	schoolSuppliesService := service.NewSchoolSuppliesService(db)
	r.Route("/api/schools/{schoolId}/supplies", func(r chi.Router) {
		r.With(paginate).Get("/", schoolSuppliesService.GetSchoolSupplies)
		r.Post("/", schoolSuppliesService.CreateSchoolSupplies)   // POST /countries
		r.Delete("/", schoolSuppliesService.DeleteSchoolSupplies) // DELETE /countries
	})

	// // RESTy routes for "featured supplies" resource
	r.Route("/api/schools/supplies", func(r chi.Router) {
		r.With(paginate).Get("/", schoolSuppliesService.GetFeaturedSchoolSupplies)
	})

	// // RESTy routes for POST "teachers requests" resource
	teachersRequestService := service.NewTeachersRequestService(db)
	r.Route("/api/teachers/requests", func(r chi.Router) {
		r.With(paginate).Get("/{id}", teachersRequestService.GetTeachersRequest)
		r.Post("/", teachersRequestService.CreateTeachersRequest) // POST /teachers/requests
	})

	// Mount the admin sub-router, which btw is the same as:
	// r.Route("/admin", func(r chi.Router) { admin routes here })
	r.Mount("/admin", adminRouter())

	// r.Mount("/web/", http.StripPrefix("/", http.FileServer(http.FS(webFS))))
	r.Mount("/", http.FileServer(http.FS(webFS)))

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if *routes {
		fmt.Printf(docgen.JSONRoutesDoc(r))
		fmt.Println(docgen.MarkdownRoutesDoc(r, docgen.MarkdownOpts{
			ProjectPath: "github.com/go-chi/chi",
			Intro:       "Welcome to the chi/_examples/rest generated docs.",
		}))
		return
	}
	defer db.Close() // remove when sql is ready
	fmt.Println("starting the server on :8080")
	fmt.Printf("Server stopped, error: %s\n", http.ListenAndServe(":8080", r))
}


func setUpDatabaseConnection() {

	var dbPwd = ""
	var dsn *url.URL
	var instanceConnectionName = ""
	// database connection pool setup
	if bProd {
		dbPwd = "Nellai987!!!"
		instanceConnectionName = "/cloudsql/helpschool:us-central1:helpschool-db"
		dsn = &url.URL{
			User:     url.UserPassword("postgres", dbPwd),
			Scheme:   "postgres",
			Host:     instanceConnectionName,
			Path:     "helpschool-db",
			RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
		}

		// "postgres://username:password@/databasename?host=/cloudsql/example:us-central1:example123"
		// "postgres://postgres:Nellai987!!!@/helpschool-db?host=/cloudsql/helpschool:us-central1:helpschool-db"

	} else {
		var err error
		dbconn := os.Getenv("DB_CONN")
		if dbconn == "" {
			dbconn = "postgresql://postgres:Pass1234@localhost/helpschool"
		}
		fmt.Fprintf(os.Stderr, `Setting up non-prod DB connection using env $DB_CONN = "%s"`, dbconn)

		dsn, err = url.Parse(dbconn)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to parse env $DB_CONN, example postgres://user:password@localhost:5432/helpschool?sslmode=disable")
			os.Exit(2)
		}
	}
	var connectionString = ""
	if bProd == true {
		connectionString = "postgres://postgres:Nellai987!!!@/helpschool?host=/cloudsql/helpschool:us-central1:helpschool-db"
	} else {
		connectionString = dsn.String()
	}
	var poolConfig, err = pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Printf("Unable to parse DATABASE_URL %v \n", err)
		os.Exit(1)
	}
	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Printf("Unable to create connection pool  %v \n ", err)
		os.Exit(1)
	} else {
		fmt.Printf("Database connection successful!!!   \n ")
	}

	// database connection pool setup
}

// A completely separate router for administrator routes
func adminRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: index"))
	})
	r.Get("/accounts", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("admin: list accounts.."))
	})
	r.Get("/users/{userId}", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("admin: view user id %v", chi.URLParam(r, "userId"))))
	})
	return r
}

// AdminOnly middleware restricts access to just administrators.
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isAdmin, ok := r.Context().Value("acl.admin").(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}

// This is entirely optional, but I wanted to demonstrate how you could easily
// add your own logic to the render.Respond method.
func init() {
	render.Respond = func(w http.ResponseWriter, r *http.Request, v interface{}) {
		if err, ok := v.(error); ok {

			// We set a default error status response code if one hasn't been set.
			if _, ok := r.Context().Value(render.StatusCtxKey).(int); !ok {
				w.WriteHeader(400)
			}

			// We log the error
			fmt.Printf("Logging err: %s\n", err.Error())

			// We change the response to not reveal the actual error message,
			// instead we can transform the message something more friendly or mapped
			// to some code / language, etc.
			render.DefaultResponder(w, r, render.M{"status": "error"})
			return
		}

		render.DefaultResponder(w, r, v)
	}
}

func fsMustSub(root fs.FS, path string) fs.FS {
	sub, err := fs.Sub(webFSRoot, path)
	if err != nil {
		panic(err)
	}
	return sub
}