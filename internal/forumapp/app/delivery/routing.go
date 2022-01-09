package delivery

import (
	"forumApp/internal/forumapp/models"

	"github.com/gorilla/mux"
)

func SetUserRouting(router *mux.Router, us models.ForumUsecase) {
	forumHandler := &ForumHandler{
		ForumUsecase: us,
	}

	router.HandleFunc("/api/forum/create", forumHandler.CreateForumHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/forum/{slug}/details", forumHandler.ForumDetailsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/forum/{slug}/create", forumHandler.CreateForumThreadHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/forum/{slug}/users", forumHandler.GetForumUsersHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/forum/{slug}/threads", forumHandler.GetForumThreadsHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/post/{id}/details", forumHandler.PostDetailsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/post/{id}/details", forumHandler.EditPostHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/service/clear", forumHandler.ServiceClearHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/service/status", forumHandler.ServiceStatusHandler).Methods("GET", "OPTIONS")

	router.HandleFunc("/api/thread/{slug_or_id}/create", forumHandler.CreatePostsHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/details", forumHandler.ThreadDetailsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/details", forumHandler.UpdateThreadHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/posts", forumHandler.GetThreadsPostsHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/thread/{slug_or_id}/vote", forumHandler.VoteThreadHandler).Methods("POST", "OPTIONS")

	router.HandleFunc("/api/user/{nickname}/create", forumHandler.CreateUserHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/user/{nickname}/profile", forumHandler.GetUserProfileHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/user/{nickname}/profile", forumHandler.UpdateUserProfileHandler).Methods("POST", "OPTIONS")
}
