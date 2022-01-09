package usecase

import (
	"database/sql"
	"forumApp/internal/forumapp/models"
	"forumApp/internal/pkg/arrutils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ForumUsecase struct {
	ForumRepo      models.ForumRepository
	contextTimeout time.Duration
}

func NewUserUsecase(fr models.ForumRepository, timeout time.Duration) models.ForumUsecase {
	return &ForumUsecase{
		ForumRepo:      fr,
		contextTimeout: timeout,
	}
}

func (fu *ForumUsecase) CreateForum(forumData models.Forum) (models.Forum, int, error) {
	findedUser, err := fu.ForumRepo.FindUserByNickname(forumData.User)
	if err != nil {
		return models.Forum{}, http.StatusNotFound, err
	}

	forumData.User = findedUser.Nickname

	createdForum, err := fu.ForumRepo.CreateForum(forumData)
	if err != nil {
		existedForum, err := fu.ForumRepo.FindForumBySlug(forumData.Slug)
		if err != nil {
			return models.Forum{}, http.StatusConflict, err
		}

		return existedForum, http.StatusConflict, nil
	}

	return createdForum, http.StatusCreated, nil
}

func (fu *ForumUsecase) GetForum(slug string) (models.Forum, int, error) {
	findedForum, err := fu.ForumRepo.FindForumBySlug(slug)
	if err != nil {
		return models.Forum{}, http.StatusNotFound, err
	}

	return findedForum, http.StatusOK, nil
}

func (fu *ForumUsecase) CreateThread(slug string, threadData models.Thread) (models.Thread, int, error) {
	_, err := fu.ForumRepo.FindUserByNickname(threadData.Author)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	findedForum, err := fu.ForumRepo.FindForumBySlug(slug)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	threadData.Forum = findedForum.Slug

	if threadData.Slug != "" {
		findedThread, err := fu.ForumRepo.FindThreadBySlug(threadData.Slug)
		if err == nil {
			return findedThread, http.StatusConflict, err
		}
	}

	createdThread, err := fu.ForumRepo.CreateThread(threadData)
	if err != nil {
		return models.Thread{}, http.StatusConflict, err
	}

	return createdThread, http.StatusCreated, nil
}

func (fu *ForumUsecase) GetThreads(slug string, params map[string][]string) ([]models.Thread, int, error) {
	_, err := fu.ForumRepo.FindForumBySlug(slug)
	if err != nil {
		return []models.Thread{}, http.StatusNotFound, err
	}

	limit := "100"
	if len(params["limit"]) > 0 {
		limit = params["limit"][0]
	}
	since := ""
	if len(params["since"]) > 0 {
		since = params["since"][0]
	}
	desc := ""
	comparisonSign := ">="
	if len(params["desc"]) > 0 && params["desc"][0] == "true" {
		desc = "desc"
		comparisonSign = "<="
	}

	findedThreads, err := fu.ForumRepo.FindThreadsBySlugWithParams(slug, limit, since, desc, comparisonSign)
	if err != nil {
		return []models.Thread{}, http.StatusNotFound, err
	}

	return findedThreads, http.StatusOK, nil
}

func (fu *ForumUsecase) CreateUser(userData models.User) ([]models.User, int, error) {
	findedUsers, err := fu.ForumRepo.FindUsersByEmailOrNickname(userData.Email, userData.Nickname)
	if err != nil || len(findedUsers) != 0 {
		return findedUsers, http.StatusConflict, err
	}

	createdUser, err := fu.ForumRepo.CreateUser(userData)
	if err != nil {
		return []models.User{}, http.StatusInternalServerError, err
	}
	var createdUsers []models.User
	createdUsers = append(createdUsers, createdUser)

	return createdUsers, http.StatusCreated, nil
}

func (fu *ForumUsecase) GetUser(nickname string) (models.User, int, error) {
	findedUser, err := fu.ForumRepo.FindUserByNickname(nickname)
	if err != nil {
		return models.User{}, http.StatusNotFound, err
	}

	return findedUser, http.StatusOK, nil
}

func (fu *ForumUsecase) UpdateUser(userData models.User) (models.User, int, error) {
	findedUser, err := fu.ForumRepo.FindUserByNickname(userData.Nickname)
	if err != nil {
		return models.User{}, http.StatusNotFound, err
	}

	if len(userData.Fullname) == 0 {
		userData.Fullname = findedUser.Fullname
	}
	if len(userData.About) == 0 {
		userData.About = findedUser.About
	}
	if len(userData.Email) == 0 {
		userData.Email = findedUser.Email
	}

	updatedUser, err := fu.ForumRepo.UpdateUser(userData)
	if err != nil {
		return models.User{}, http.StatusConflict, err
	}

	return updatedUser, http.StatusOK, nil
}

func (fu *ForumUsecase) CreatesPosts(threadSlugOrId string, postsData []models.Post) ([]models.Post, int, error) {
	threadId, _ := strconv.Atoi(threadSlugOrId)

	findedThread, err := fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return []models.Post{}, http.StatusNotFound, err
	}

	createdPosts, err := fu.ForumRepo.CreatePosts(postsData, findedThread)
	if err != nil {
		if err.Error() == "404" {
			return []models.Post{}, http.StatusNotFound, err
		}

		return []models.Post{}, http.StatusConflict, err
	}

	return createdPosts, http.StatusCreated, nil
}

func (fu *ForumUsecase) VoteThread(threadSlugOrId string, voteData models.Vote) (models.Thread, int, error) {
	threadId, _ := strconv.Atoi(threadSlugOrId)

	findedThread, err := fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	findedUser, err := fu.ForumRepo.FindUserByNickname(voteData.Nickname)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	err = fu.ForumRepo.FindVote(findedUser.Id, findedThread.Id)
	if err != nil && err != sql.ErrNoRows {
		return models.Thread{}, http.StatusNotFound, err
	}
	if err == nil {
		err = fu.ForumRepo.UpdateVoteThread(findedUser.Id, findedThread.Id, voteData.Voice)
		if err != nil {
			return models.Thread{}, http.StatusNotFound, err
		}

		findedThread, err = fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
		if err != nil {
			return models.Thread{}, http.StatusNotFound, err
		}

		return findedThread, http.StatusOK, nil
	}
	err = fu.ForumRepo.AddVoteThread(findedUser.Id, findedThread.Id, voteData.Voice)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	findedThread, err = fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	return findedThread, http.StatusOK, nil
}

func (fu *ForumUsecase) FindThreadBySlugOrId(threadSlugOrId string) (models.Thread, int, error) {
	threadId, _ := strconv.Atoi(threadSlugOrId)

	findedThread, err := fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	return findedThread, http.StatusOK, nil
}

func (fu *ForumUsecase) GetPosts(threadSlugOrId string, params map[string][]string) ([]models.Post, int, error) {
	threadId, _ := strconv.Atoi(threadSlugOrId)

	findedThread, err := fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return []models.Post{}, http.StatusNotFound, err
	}

	limit := "100"
	if len(params["limit"]) > 0 {
		limit = params["limit"][0]
	}
	since := ""
	if len(params["since"]) > 0 {
		since = params["since"][0]
	}
	sort := "flat"
	if len(params["sort"]) > 0 {
		sort = params["sort"][0]
	}
	desc := ""
	comparisonSign := ">"
	if len(params["desc"]) > 0 && params["desc"][0] == "true" {
		desc = "desc"
		comparisonSign = "<"
	}

	findedPosts, err := fu.ForumRepo.GetPosts(findedThread.Id, limit, since, sort, desc, comparisonSign)
	if err != nil {
		return []models.Post{}, http.StatusNotFound, err
	}

	return findedPosts, http.StatusOK, nil
}

func (fu *ForumUsecase) UpdateThread(threadSlugOrId string, newThread models.Thread) (models.Thread, int, error) {
	threadId, _ := strconv.Atoi(threadSlugOrId)

	findedThread, err := fu.ForumRepo.FindThreadBySlugOrId(int64(threadId), threadSlugOrId)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	if len(newThread.Title) == 0 && len(newThread.Message) == 0 {
		return findedThread, http.StatusOK, nil
	}

	if len(newThread.Title) == 0 {
		newThread.Title = findedThread.Title
	}
	if len(newThread.Message) == 0 {
		newThread.Message = findedThread.Message
	}

	updatedThread, err := fu.ForumRepo.UpdateThread(findedThread.Id, newThread)
	if err != nil {
		return models.Thread{}, http.StatusNotFound, err
	}

	return updatedThread, http.StatusOK, nil
}

func (fu *ForumUsecase) GetForumUsers(forumSlug string, params map[string][]string) ([]models.User, int, error) {
	findedForum, err := fu.ForumRepo.FindForumBySlug(forumSlug)
	if err != nil {
		return []models.User{}, http.StatusNotFound, err
	}

	limit := "100"
	if len(params["limit"]) > 0 {
		limit = params["limit"][0]
	}
	since := ""
	if len(params["since"]) > 0 {
		since = params["since"][0]
	}
	desc := ""
	comparisonSign := ">"
	if len(params["desc"]) > 0 && params["desc"][0] == "true" {
		desc = "desc"
		comparisonSign = "<"
	}

	findedUsers, err := fu.ForumRepo.GetForumUsers(findedForum.Id, limit, since, desc, comparisonSign)
	if err != nil {
		return []models.User{}, http.StatusNotFound, err
	}

	return findedUsers, http.StatusOK, nil
}

func (fu *ForumUsecase) GetPostInfo(id string, params map[string][]string) (models.PostFull, int, error) {
	postId, _ := strconv.Atoi(id)
	withUser, withForum, withThread := false, false, false
	related := params["related"]
	if len(related) > 0 {
		parsedRelated := strings.Split(related[0], ",")

		if arrutils.StringSliceHas(parsedRelated, "user") {
			withUser = true
		}
		if arrutils.StringSliceHas(parsedRelated, "forum") {
			withForum = true
		}
		if arrutils.StringSliceHas(parsedRelated, "thread") {
			withThread = true
		}
	}

	findedPostInfo, err := fu.ForumRepo.GetPostInfo(int64(postId), withUser, withForum, withThread)
	if err != nil {
		return models.PostFull{}, http.StatusNotFound, err
	}

	return findedPostInfo, http.StatusOK, nil
}

func (fu *ForumUsecase) UpdatePost(id string, newPost models.Post) (models.Post, int, error) {
	postId, _ := strconv.Atoi(id)

	findedPost, err := fu.ForumRepo.FindPost(int64(postId))
	if err != nil {
		return models.Post{}, http.StatusNotFound, err
	}

	if len(newPost.Message) != 0 {
		if newPost.Message != findedPost.Message {
			findedPost.IsEdited = true
		}
		findedPost.Message = newPost.Message
	}

	updatedPost, err := fu.ForumRepo.UpdatePost(findedPost)
	if err != nil {
		return models.Post{}, http.StatusNotFound, err
	}

	return updatedPost, http.StatusOK, nil
}

func (fu *ForumUsecase) ServiceStatus() (models.Status, int, error) {
	curServiceStatis, err := fu.ForumRepo.ServiceStatus()
	if err != nil {
		return models.Status{}, http.StatusInternalServerError, err
	}

	return curServiceStatis, http.StatusOK, nil
}

func (fu *ForumUsecase) ServiceClear() (int, error) {
	err := fu.ForumRepo.ServiceClear()
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
