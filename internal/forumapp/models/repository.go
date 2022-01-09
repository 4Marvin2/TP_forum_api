package models

type ForumRepository interface {
	FindUserByNickname(nickname string) (User, error)
	FindUsersByEmailOrNickname(email string, nickname string) ([]User, error)
	CreateUser(userData User) (User, error)
	UpdateUser(userData User) (User, error)

	CreateForum(forumData Forum) (Forum, error)
	FindForumBySlug(slug string) (Forum, error)

	CreateThread(threadData Thread) (Thread, error)
	FindThreadBySlug(slug string) (Thread, error)
	FindThreadsBySlugWithParams(slug string, limit string, since string, desc string, comparisonSign string) ([]Thread, error)
	FindThreadBySlugOrId(id int64, slug string) (Thread, error)
	CreatePosts(posts []Post, thread Thread) ([]Post, error)
	FindVote(userId int64, threadId int64) error
	UpdateVoteThread(userId int64, threadId int64, voice int32) error
	AddVoteThread(userId int64, threadId int64, voice int32) error
	GetPosts(threadId int64, limit string, since string, sort string, desc string, comparisonSign string) ([]Post, error)
	UpdateThread(threadId int64, threadData Thread) (Thread, error)
	GetForumUsers(forumId int64, limit string, since string, desc string, comparisonSign string) ([]User, error)
	GetPostInfo(postId int64, withUser bool, withForum bool, withThread bool) (PostFull, error)
	FindPost(postId int64) (Post, error)
	UpdatePost(postData Post) (Post, error)
	ServiceStatus() (Status, error)
	ServiceClear() error
}
