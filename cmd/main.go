package main

import (
	"forumApp/configs"
	"forumApp/internal/forumapp/app/delivery"
	"forumApp/internal/forumapp/app/repository"
	"forumApp/internal/forumapp/app/usecase"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	configs.SetConfig()

	router := mux.NewRouter()

	repo, err := repository.NewPostgresUserRepository(configs.Postgres)
	if err != nil {
		log.Fatal(err)
	}

	timeoutContext := configs.Timeouts.ContextTimeout

	usecase := usecase.NewUserUsecase(repo, timeoutContext)

	delivery.SetUserRouting(router, usecase)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":5000",
		WriteTimeout: http.DefaultClient.Timeout,
		ReadTimeout:  http.DefaultClient.Timeout,
	}

	log.Fatal(srv.ListenAndServe())
}
