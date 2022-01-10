package models

type ForumUsecase interface {
	CreateUser(userData User) (Users, int, error)
	GetUser(nickname string) (User, int, error)
	UpdateUser(userData User) (User, int, error)

	CreateForum(forumData Forum) (Forum, int, error)
	GetForum(slug string) (Forum, int, error)

	CreateThread(slug string, threadData Thread) (Thread, int, error)
	GetThreads(slug string, params map[string][]string) (Threads, int, error)

	CreatesPosts(threadSlugOrId string, postsData []Post) (Posts, int, error)
	VoteThread(threadSlugOrId string, voteData Vote) (Thread, int, error)
	FindThreadBySlugOrId(threadSlugOrId string) (Thread, int, error)
	GetPosts(threadSlugOrId string, params map[string][]string) (Posts, int, error)
	UpdateThread(threadSlugOrId string, newThread Thread) (Thread, int, error)
	GetForumUsers(forumSlug string, params map[string][]string) (Users, int, error)
	GetPostInfo(id string, params map[string][]string) (PostFull, int, error)
	UpdatePost(id string, newPost Post) (Post, int, error)
	ServiceStatus() (Status, int, error)
	ServiceClear() (int, error)
}
