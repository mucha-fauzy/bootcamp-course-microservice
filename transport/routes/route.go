package routes

import (
	"bootcamp_course_microservice/internal/services"
	custom_middleware "bootcamp_course_microservice/transport/middleware"
	"net/http"

	"bootcamp_course_microservice/internal/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Router struct {
	Handler        *handlers.Handler
	Authentication *custom_middleware.Authentication
}

func ProvideRouter(service services.Service, auth *custom_middleware.Authentication) *Router {
	handler := handlers.ProvideHandler(service)
	return &Router{
		Handler:        handler,
		Authentication: auth,
	}
}

func (r *Router) SetupRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	mux.Route("/v1", func(rc chi.Router) {
		// Protected endpoints accessible only by teachers
		rc.Group(func(rc chi.Router) {
			rc.Use(r.Authentication.VerifyTeacherJWT)
			rc.Get("/courses/{userid}", r.Handler.ReadCourseByUserID)
			rc.Post("/courses", r.Handler.CreateCourse)
		})
	})
	return mux
}
