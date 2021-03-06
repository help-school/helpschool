// Boot the server:
// ----------------
// $ go run main.go

package main

import (
	"context"
	"embed"
	_ "embed"
	"encoding/gob"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/docgen"
	"github.com/go-chi/render"
	"github.com/gorilla/sessions"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/log/log15adapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/venkata6/helpschool/api/auth"
	"github.com/venkata6/helpschool/api/service"
	// "time"
)

//go:embed web
var webFSRoot embed.FS
var webFS = fsMustSub(webFSRoot, "web")

var Store *sessions.FilesystemStore

// User is a user meta retrieved from JWT (Auth0 access token)
type User struct {
	Email string
	Auth0ID string
}

func main() {
	var generateDocs, isProd bool
	flag.BoolVar(&generateDocs, "routes", false, "Generate router documentation")
	flag.BoolVar(&isProd, "prod", false, "Run in production mode")
	flag.Parse()

	ctx := context.Background()
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.NoCache)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	authMiddleware := auth.NewMiddleware(
		"https://helpschool/api",
		"https://helpschool.us.auth0.com/")

	// add CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	Store = sessions.NewFilesystemStore("", []byte("something-very-secret"))
	gob.Register(map[string]interface{}{})

	// connect to the database and setup the connection pool for services to use
	db, err := setUpDatabaseConnection(ctx, isProd)
	if err != nil {
		panic(err)
	}

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

	// requires a valid JWT token
	r.With(authMiddleware.Handler).Get("/api/my-donations", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Donation struct {
			Item      string `json:"item"`
			Timestamp int64  `json:"ts"`
		}

		user := getUserOrFail(w, r, http.StatusBadRequest)
		if user == nil {
			return
		}

		fmt.Printf("user: %#v\n", user)
		render.JSON(w, r, []Donation{
			{"something useful", time.Now().Unix()},
			{"something even more useful", time.Now().Unix()},
		})
	}))

	// Mount the admin sub-router, which btw is the same as:
	// r.Route("/admin", func(r chi.Router) { admin routes here })
	r.Mount("/admin", adminRouter())

	// Public HTML site
	r.Mount("/", http.FileServer(http.FS(webFS)))

	// Passing -routes to the program will generate docs for the above
	// router definition. See the `routes.json` file in this folder for
	// the output.
	if generateDocs {
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

func setUpDatabaseConnection(ctx context.Context, isProd bool) (*pgxpool.Pool, error) {
	// database connection pool setup
	dbPwd := "Nellai987!!!"
	instanceConnectionName := "/cloudsql/helpschool:us-central1:helpschool-db"
	dsn := &url.URL{
		User:     url.UserPassword("postgres", dbPwd),
		Scheme:   "postgres",
		Host:     instanceConnectionName,
		Path:     "helpschool-db",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	if !isProd {
		dbconn := os.Getenv("DB_CONN")
		if dbconn == "" {
			fmt.Println("no $DB_CONN provided, proceeding with default")
			dbconn = "postgresql://postgres:Pass1234@localhost/helpschool"
		}

		var err error
		dsn, err = url.Parse(dbconn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse env $DB_CONN, example postgres://user:password@localhost:5432/helpschool?sslmode=disable")
		}
		fmt.Fprintf(os.Stderr, `Setting up non-prod DB connection using env $DB_CONN = "%s"`, dbconn)
	}

	connectionString := "postgres://postgres:Nellai987!!!@/helpschool?host=/cloudsql/helpschool:us-central1:helpschool-db"
	if !isProd {
		connectionString = dsn.String()
	}

	var poolConfig, err = pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Printf("Unable to parse DATABASE_URL %v \n", err)
		os.Exit(1)
	}
	db, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("db connect: %s", err)
	}

	fmt.Println("Database connection successful!!!")
	return db, nil
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

func getUserOrFail(w http.ResponseWriter, r *http.Request, failStatus int) *User {
	if value := r.Context().Value("user"); value != nil {
		if token, ok := value.(*jwt.Token); ok {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				var user User
				if email, ok := claims["sub"]; ok {
					user.Auth0ID = email.(string)
				}
				if email, ok := claims["https://example.com/email"]; ok {
					user.Email = email.(string)
				}
				return &user
			}
		}
	}
	http.Error(w, "no JWT token", failStatus)
	return nil
}
