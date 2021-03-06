package repository

import (
	"errors"
	"fmt"
	"forumApp/configs"
	"forumApp/internal/forumapp/models"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
)

type PostgreForumRepo struct {
	Conn *pgx.ConnPool
}

func NewPostgresUserRepository(config configs.PostgresConfig) (models.ForumRepository, error) {
	ConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		config.User,
		config.DBName,
		config.Password,
		config.Host,
		config.Port)

	pgxConnectionConfig, err := pgx.ParseConnectionString(ConnStr)
	if err != nil {
		log.Fatalf("Invalid config string: %s", err)
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxConnectionConfig,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatalf("Error %s occurred during connection to database", err)
	}

	return &PostgreForumRepo{pool}, nil
}

func (pfr *PostgreForumRepo) FindUserByNickname(nickname string) (models.User, error) {
	var findedUser models.User
	err := pfr.Conn.QueryRow(FindUserByNicknameQuery, nickname).Scan(&findedUser.Id, &findedUser.Nickname, &findedUser.About, &findedUser.Email, &findedUser.Fullname)
	if err != nil {
		return models.User{}, err
	}
	return findedUser, nil
}

func (pfr *PostgreForumRepo) FindUsersByEmailOrNickname(email string, nickname string) ([]models.User, error) {
	var findedUsers []models.User
	rows, err := pfr.Conn.Query(FindUserByEmailOrNicknameQuery, email, nickname)
	if err != nil {
		return []models.User{}, err
	}
	for rows.Next() {
		var curUser models.User
		err := rows.Scan(&curUser.Nickname, &curUser.About, &curUser.Email, &curUser.Fullname)
		if err != nil {
			return []models.User{}, err
		}
		findedUsers = append(findedUsers, curUser)
	}
	return findedUsers, nil
}

func (pfr *PostgreForumRepo) CreateUser(userData models.User) (models.User, error) {
	var createdUser models.User
	err := pfr.Conn.QueryRow(
		CreateUserQuery,
		userData.Nickname,
		userData.Fullname,
		userData.About,
		userData.Email,
	).Scan(
		&createdUser.Nickname,
		&createdUser.Fullname,
		&createdUser.About,
		&createdUser.Email,
	)

	if err != nil {
		return models.User{}, err
	}
	return createdUser, nil
}

func (pfr *PostgreForumRepo) UpdateUser(userData models.User) (models.User, error) {
	var updatedUser models.User
	err := pfr.Conn.QueryRow(
		UpdateUserQuery,
		userData.Nickname,
		userData.Fullname,
		userData.About,
		userData.Email,
	).Scan(
		&updatedUser.Nickname,
		&updatedUser.Fullname,
		&updatedUser.About,
		&updatedUser.Email,
	)
	if err != nil {
		return models.User{}, err
	}
	return updatedUser, nil
}

func (pfr *PostgreForumRepo) CreateForum(forumData models.Forum) (models.Forum, error) {
	var createdForum models.Forum
	err := pfr.Conn.QueryRow(
		CreateForumQuery,
		forumData.Title,
		forumData.User,
		forumData.Slug,
	).Scan(
		&createdForum.Title,
		&createdForum.User,
		&createdForum.Slug,
		&createdForum.Posts,
		&createdForum.Threads,
	)
	if err != nil {
		return models.Forum{}, err
	}
	return createdForum, nil
}

func (pfr *PostgreForumRepo) FindForumBySlug(slug string) (models.Forum, error) {
	var findedForum models.Forum
	err := pfr.Conn.QueryRow(
		FindForumBySlugQuery,
		slug,
	).Scan(
		&findedForum.Id,
		&findedForum.Title,
		&findedForum.User,
		&findedForum.Slug,
		&findedForum.Posts,
		&findedForum.Threads,
	)
	if err != nil {
		return models.Forum{}, err
	}
	return findedForum, nil
}

func (pfr *PostgreForumRepo) FindThreadBySlug(slug string) (models.Thread, error) {
	var findedThread models.Thread
	err := pfr.Conn.QueryRow(
		FindThreadBySlugQuery,
		slug,
	).Scan(
		&findedThread.Id,
		&findedThread.Title,
		&findedThread.Author,
		&findedThread.Forum,
		&findedThread.Message,
		&findedThread.Votes,
		&findedThread.Slug,
		&findedThread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}
	return findedThread, nil
}

func (pfr *PostgreForumRepo) CreateThread(threadData models.Thread) (models.Thread, error) {
	var createdThread models.Thread
	if threadData.Created.String() == "" {
		threadData.Created = time.Now()
	}
	err := pfr.Conn.QueryRow(
		CreateThreadQuery,
		threadData.Title,
		threadData.Author,
		threadData.Forum,
		threadData.Message,
		threadData.Slug,
		threadData.Created,
	).Scan(
		&createdThread.Id,
		&createdThread.Title,
		&createdThread.Author,
		&createdThread.Forum,
		&createdThread.Message,
		&createdThread.Votes,
		&createdThread.Slug,
		&createdThread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}

	var forumId int64
	err = pfr.Conn.QueryRow(UpdateForumsThreadCountQuery, threadData.Forum).Scan(&forumId)
	if err != nil {
		return models.Thread{}, err
	}

	return createdThread, nil
}

func (pfr *PostgreForumRepo) FindThreadsBySlugWithParams(slug string, limit string, since string, desc string, comparisonSign string) ([]models.Thread, error) {
	findedThreads := make([]models.Thread, 0)
	customizeQuery := FindThreadsByForumQuery
	if since != "" {
		customizeQuery += fmt.Sprintf(" AND created %s '%s'", comparisonSign, since)
	}
	customizeQuery += fmt.Sprintf(" ORDER BY created %s LIMIT %s;", desc, limit)

	rows, err := pfr.Conn.Query(customizeQuery, slug)
	if err != nil {
		return []models.Thread{}, err
	}
	for rows.Next() {
		var curThread models.Thread
		err := rows.Scan(
			&curThread.Id,
			&curThread.Title,
			&curThread.Author,
			&curThread.Forum,
			&curThread.Message,
			&curThread.Votes,
			&curThread.Slug,
			&curThread.Created,
		)
		if err != nil {
			return []models.Thread{}, err
		}
		findedThreads = append(findedThreads, curThread)
	}
	return findedThreads, nil
}

func (pfr *PostgreForumRepo) FindThreadBySlugOrId(id int64, slug string) (models.Thread, error) {
	var findedThread models.Thread
	err := pfr.Conn.QueryRow(
		FindThreadBySlugOrIdQuery,
		id,
		slug,
	).Scan(
		&findedThread.Id,
		&findedThread.Title,
		&findedThread.Author,
		&findedThread.Forum,
		&findedThread.Message,
		&findedThread.Votes,
		&findedThread.Slug,
		&findedThread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}
	return findedThread, nil
}

func (pfr *PostgreForumRepo) CreatePosts(posts []models.Post, thread models.Thread) ([]models.Post, error) {
	createdPosts := make([]models.Post, 0)

	if len(posts) == 0 {
		return []models.Post{}, nil
	}

	createdTime := time.Now()

	sqlQuery := CreateThreadStartQuery

	var values []interface{}
	paramNumber := 1
	var sb strings.Builder
	for _, post := range posts {
		_, err := pfr.FindUserByNickname(post.Author)
		if err != nil {
			return []models.Post{}, errors.New("404")
		}

		if post.Parent == 0 {
			sb.WriteString("(nextval('posts_id_seq'::regclass)")
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber))
			sb.WriteString(", ARRAY[currval(pg_get_serial_sequence('posts', 'id'))::bigint], $")
			sb.WriteString(strconv.Itoa(paramNumber + 1))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 2))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 3))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 4))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 5))
			sb.WriteString("),")
			paramNumber += 6
			sqlQuery += sb.String()
			sb.Reset()
			values = append(values, post.Parent, post.Author, post.Message, thread.Forum, thread.Id, createdTime)
		} else {
			var parentId int64
			err = pfr.Conn.QueryRow(FindParentIdForPostQuery, post.Parent).Scan(&parentId)
			if err != nil {
				return []models.Post{}, err
			}
			if parentId != thread.Id {
				return []models.Post{}, errors.New("parent post was created in another thread")
			}
			sb.WriteString("(nextval('posts_id_seq'::regclass)")
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber))
			sb.WriteString(", (SELECT path FROM posts WHERE id = $")
			sb.WriteString(strconv.Itoa(paramNumber + 1))
			sb.WriteString(" AND thread = $")
			sb.WriteString(strconv.Itoa(paramNumber + 2))
			sb.WriteString(") || currval(pg_get_serial_sequence('posts', 'id'))::bigint, $")
			sb.WriteString(strconv.Itoa(paramNumber + 3))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 4))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 5))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 6))
			sb.WriteString(", $")
			sb.WriteString(strconv.Itoa(paramNumber + 7))
			sb.WriteString("),")
			paramNumber += 8
			sqlQuery += sb.String()
			sb.Reset()
			values = append(values, post.Parent, post.Parent, thread.Id, post.Author, post.Message, thread.Forum, thread.Id, createdTime)
		}
	}

	sqlQuery = strings.TrimSuffix(sqlQuery, ",")
	sqlQuery += " RETURNING id, parent, author, message, isEdited, forum, thread, created;"

	rows, err := pfr.Conn.Query(sqlQuery, values...)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var curPost models.Post
		err := rows.Scan(
			&curPost.Id,
			&curPost.Parent,
			&curPost.Author,
			&curPost.Message,
			&curPost.IsEdited,
			&curPost.Forum,
			&curPost.Thread,
			&curPost.Created,
		)

		if err != nil {
			return []models.Post{}, err
		}

		createdPosts = append(createdPosts, curPost)
	}

	var forumId int64
	err = pfr.Conn.QueryRow(UpdateForumsPostsCountQuery, len(createdPosts), thread.Forum).Scan(&forumId)
	if err != nil {
		return []models.Post{}, err
	}

	if err != nil {
		return []models.Post{}, err
	}

	return createdPosts, nil
}

func (pfr *PostgreForumRepo) FindVote(userId int64, threadId int64) error {
	var findedVoteId int64
	err := pfr.Conn.QueryRow(FindVoteQuery, userId, threadId).Scan(&findedVoteId)
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PostgreForumRepo) UpdateVoteThread(userId int64, threadId int64, voice int32) error {
	var updatedVoteId int64
	err := pfr.Conn.QueryRow(UpdateVoteQuery, userId, threadId, voice).Scan(&updatedVoteId)
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PostgreForumRepo) AddVoteThread(userId int64, threadId int64, voice int32) error {
	var updatedVoteId int64
	err := pfr.Conn.QueryRow(AddVoteQuery, userId, threadId, voice).Scan(&updatedVoteId)
	if err != nil {
		return err
	}
	return nil
}

func (pfr *PostgreForumRepo) GetPosts(threadId int64, limit string, since string, sort string, desc string, comparisonSign string) ([]models.Post, error) {
	findedPosts := make([]models.Post, 0)
	sqlQuery := GetPostsStartQuery
	switch sort {
	case "flat":
		if since != "" {
			sqlQuery += fmt.Sprintf(" AND id %s %s", comparisonSign, since)
		}
		sqlQuery += fmt.Sprintf(" ORDER BY created %s, id %s LIMIT %s", desc, desc, limit)
	case "tree":
		if since != "" {
			sqlQuery += fmt.Sprintf(" AND path %s (SELECT path FROM posts WHERE id = %s)", comparisonSign, since)
		}
		sqlQuery += fmt.Sprintf(" ORDER BY path[1] %s, path %s LIMIT %s", desc, desc, limit)
	case "parent_tree":
		sqlQuery += " AND path && (SELECT ARRAY (SELECT id FROM posts WHERE thread = $1 AND parent = 0"
		if since != "" {
			sqlQuery += fmt.Sprintf(" AND path %s (SELECT path[1:1] FROM posts WHERE id = %s)", comparisonSign, since)
		}
		sqlQuery += fmt.Sprintf(" ORDER BY path[1] %s, path LIMIT %s)) ORDER BY path[1] %s, path", desc, limit, desc)
	default:
		sqlQuery = "error"
	}

	if sqlQuery == "error" {
		return []models.Post{}, errors.New("undefined sort type")
	}

	rows, err := pfr.Conn.Query(sqlQuery, threadId)
	if err != nil {
		return []models.Post{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var curPost models.Post
		err := rows.Scan(
			&curPost.Id,
			&curPost.Parent,
			&curPost.Author,
			&curPost.Message,
			&curPost.IsEdited,
			&curPost.Forum,
			&curPost.Thread,
			&curPost.Created,
		)
		if err != nil {
			return []models.Post{}, err
		}
		findedPosts = append(findedPosts, curPost)
	}

	return findedPosts, nil
}

func (pfr *PostgreForumRepo) UpdateThread(threadId int64, threadData models.Thread) (models.Thread, error) {
	var updatedThread models.Thread
	err := pfr.Conn.QueryRow(
		UpdateThreadQuery,
		threadData.Title,
		threadData.Message,
		threadId,
	).Scan(
		&updatedThread.Id,
		&updatedThread.Title,
		&updatedThread.Author,
		&updatedThread.Forum,
		&updatedThread.Message,
		&updatedThread.Votes,
		&updatedThread.Slug,
		&updatedThread.Created,
	)
	if err != nil {
		return models.Thread{}, err
	}
	return updatedThread, nil
}

func (pfr *PostgreForumRepo) GetForumUsers(forumId int64, limit string, since string, desc string, comparisonSign string) ([]models.User, error) {
	findedUsers := make([]models.User, 0)
	sqlQuery := GetForumUsersStartQuery
	if since != "" {
		sqlQuery += fmt.Sprintf(" AND nickname %s '%s'", comparisonSign, since)
	}
	sqlQuery += fmt.Sprintf(" ORDER BY nickname %s LIMIT %s", desc, limit)

	rows, err := pfr.Conn.Query(sqlQuery, forumId)
	if err != nil {
		return []models.User{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var curUser models.User
		err := rows.Scan(&curUser.Id, &curUser.Nickname, &curUser.About, &curUser.Email, &curUser.Fullname)
		if err != nil {
			return []models.User{}, err
		}
		findedUsers = append(findedUsers, curUser)
	}
	return findedUsers, nil
}

func (pfr *PostgreForumRepo) GetPostInfo(postId int64, withUser bool, withForum bool, withThread bool) (models.PostFull, error) {
	var findedPostInfo models.PostFull
	var findedPost models.Post
	err := pfr.Conn.QueryRow(
		GetPostInfoQuery,
		postId,
	).Scan(
		&findedPost.Id,
		&findedPost.Parent,
		&findedPost.Author,
		&findedPost.Message,
		&findedPost.IsEdited,
		&findedPost.Forum,
		&findedPost.Thread,
		&findedPost.Created,
	)
	if err != nil {
		return models.PostFull{}, err
	}
	findedPostInfo.Post = &findedPost

	if withUser {
		var findedUser models.User
		err = pfr.Conn.QueryRow(

			FindUserByNicknameQuery,
			findedPost.Author,
		).Scan(
			&findedUser.Id,
			&findedUser.Nickname,
			&findedUser.About,
			&findedUser.Email,
			&findedUser.Fullname,
		)
		if err != nil {
			return models.PostFull{}, err
		}
		findedPostInfo.Author = &findedUser
	}

	if withForum {
		var findedForum models.Forum
		err = pfr.Conn.QueryRow(

			FindForumBySlugQuery,
			findedPost.Forum,
		).Scan(
			&findedForum.Id,
			&findedForum.Title,
			&findedForum.User,
			&findedForum.Slug,
			&findedForum.Posts,
			&findedForum.Threads,
		)
		if err != nil {
			return models.PostFull{}, err
		}
		findedPostInfo.Forum = &findedForum
	}

	if withThread {
		var findedThread models.Thread
		err = pfr.Conn.QueryRow(

			FindThreadByIdQuery,
			findedPost.Thread,
		).Scan(
			&findedThread.Id,
			&findedThread.Title,
			&findedThread.Author,
			&findedThread.Forum,
			&findedThread.Message,
			&findedThread.Votes,
			&findedThread.Slug,
			&findedThread.Created,
		)
		if err != nil {
			return models.PostFull{}, err
		}
		findedPostInfo.Thread = &findedThread
	}

	return findedPostInfo, nil
}

func (pfr *PostgreForumRepo) FindPost(postId int64) (models.Post, error) {
	var findedPost models.Post
	err := pfr.Conn.QueryRow(
		GetPostInfoQuery,
		postId,
	).Scan(
		&findedPost.Id,
		&findedPost.Parent,
		&findedPost.Author,
		&findedPost.Message,
		&findedPost.IsEdited,
		&findedPost.Forum,
		&findedPost.Thread,
		&findedPost.Created,
	)
	if err != nil {
		return models.Post{}, err
	}
	return findedPost, nil
}

func (pfr *PostgreForumRepo) UpdatePost(postData models.Post) (models.Post, error) {
	var updatedPost models.Post
	err := pfr.Conn.QueryRow(
		UpdatePostQuery,
		postData.Id,
		postData.Parent,
		postData.Author,
		postData.Message,
		postData.IsEdited,
		postData.Forum,
		postData.Thread,
		postData.Created,
	).Scan(
		&updatedPost.Id,
		&updatedPost.Parent,
		&updatedPost.Author,
		&updatedPost.Message,
		&updatedPost.IsEdited,
		&updatedPost.Forum,
		&updatedPost.Thread,
		&updatedPost.Created,
	)
	if err != nil {
		return models.Post{}, err
	}
	return updatedPost, nil
}

func (pfr *PostgreForumRepo) ServiceStatus() (models.Status, error) {
	var curServiceStatus models.Status
	err := pfr.Conn.QueryRow(
		GetServiceStatusQuery,
	).Scan(
		&curServiceStatus.Forum,
		&curServiceStatus.Post,
		&curServiceStatus.Thread,
		&curServiceStatus.User,
	)
	if err != nil {
		return models.Status{}, err
	}
	return curServiceStatus, nil
}

func (pfr *PostgreForumRepo) ServiceClear() error {
	_, err := pfr.Conn.Exec(ClearServiceQuery)
	if err != nil {
		return err
	}

	return nil
}
