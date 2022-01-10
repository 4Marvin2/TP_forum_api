package delivery

import (
	"errors"
	"forumApp/internal/forumapp/models"
	"forumApp/internal/pkg/ioutils"
	"net/http"

	"github.com/gorilla/mux"
)

type ForumHandler struct {
	ForumUsecase models.ForumUsecase
}

func (uh *ForumHandler) CreateForumHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newForum models.Forum
	err := ioutils.ReadJSON(r, &newForum)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	forum, code, err := uh.ForumUsecase.CreateForum(newForum)
	if code == http.StatusNotFound {
		ioutils.SendError(w, code, err.Error())
		return
	}
	ioutils.Send(w, code, forum)
}

func (uh *ForumHandler) ForumDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)["slug"]

	findedForum, code, err := uh.ForumUsecase.GetForum(slug)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedForum)
}

func (uh *ForumHandler) CreateForumThreadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)["slug"]

	var newThread models.Thread
	err := ioutils.ReadJSON(r, &newThread)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdThread, code, err := uh.ForumUsecase.CreateThread(slug, newThread)
	if code == http.StatusNotFound {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, createdThread)
}

func (uh *ForumHandler) GetForumUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)["slug"]

	findedUsers, code, err := uh.ForumUsecase.GetForumUsers(slug, r.URL.Query())
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedUsers)
}

func (uh *ForumHandler) GetForumThreadsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)["slug"]

	findedThreads, code, err := uh.ForumUsecase.GetThreads(slug, r.URL.Query())
	if err != nil || code == http.StatusNotFound {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedThreads)
}

func (uh *ForumHandler) PostDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	findedPostIndo, code, err := uh.ForumUsecase.GetPostInfo(id, r.URL.Query())
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedPostIndo)
}

func (uh *ForumHandler) EditPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var newPost models.Post
	err := ioutils.ReadJSON(r, &newPost)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedPost, code, err := uh.ForumUsecase.UpdatePost(id, newPost)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, updatedPost)
}

func (uh *ForumHandler) ServiceClearHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	code, err := uh.ForumUsecase.ServiceClear()
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.SendWithoutBody(w, code)
}

func (uh *ForumHandler) ServiceStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	curServiceStatis, code, err := uh.ForumUsecase.ServiceStatus()
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, curServiceStatis)
}

func (uh *ForumHandler) CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slugOrId := mux.Vars(r)["slug_or_id"]

	var newPosts models.Posts
	err := ioutils.ReadJSON(r, &newPosts)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdPosts, code, err := uh.ForumUsecase.CreatesPosts(slugOrId, newPosts)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, createdPosts)
}

func (uh *ForumHandler) ThreadDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slugOrId := mux.Vars(r)["slug_or_id"]

	findedThread, code, err := uh.ForumUsecase.FindThreadBySlugOrId(slugOrId)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedThread)
}

func (uh *ForumHandler) UpdateThreadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slugOrId := mux.Vars(r)["slug_or_id"]

	var newThread models.Thread
	err := ioutils.ReadJSON(r, &newThread)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedThread, code, err := uh.ForumUsecase.UpdateThread(slugOrId, newThread)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, updatedThread)
}

func (uh *ForumHandler) GetThreadsPostsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slugOrId := mux.Vars(r)["slug_or_id"]

	findedPosts, code, err := uh.ForumUsecase.GetPosts(slugOrId, r.URL.Query())
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, findedPosts)
}

func (uh *ForumHandler) VoteThreadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	slugOrId := mux.Vars(r)["slug_or_id"]

	var newVote models.Vote
	err := ioutils.ReadJSON(r, &newVote)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	threadInfo, code, err := uh.ForumUsecase.VoteThread(slugOrId, newVote)
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, threadInfo)
}

func (uh *ForumHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	var newUser models.User
	err := ioutils.ReadJSON(r, &newUser)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Nickname = nickname

	createdUsers, code, err := uh.ForumUsecase.CreateUser(newUser)
	if code == http.StatusConflict {
		ioutils.Send(w, code, createdUsers)
		return
	}
	if err != nil {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, createdUsers[0])
}

func (uh *ForumHandler) GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	findedUser, code, err := uh.ForumUsecase.GetUser(nickname)
	if code == http.StatusNotFound || err != nil {
		ioutils.SendError(w, http.StatusNotFound, errors.New("Can't find user with nickname #"+nickname+"\n").Error())
		return
	}
	ioutils.Send(w, code, findedUser)
}

func (uh *ForumHandler) UpdateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	nickname := mux.Vars(r)["nickname"]

	var newUser models.User
	err := ioutils.ReadJSON(r, &newUser)
	if err != nil {
		ioutils.SendError(w, http.StatusBadRequest, err.Error())
		return
	}
	newUser.Nickname = nickname

	updatedUser, code, err := uh.ForumUsecase.UpdateUser(newUser)
	if err != nil || code != http.StatusOK {
		ioutils.SendError(w, code, err.Error())
		return
	}

	ioutils.Send(w, code, updatedUser)
}
