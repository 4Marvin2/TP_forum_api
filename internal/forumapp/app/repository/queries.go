package repository

const (
	FindUserByNicknameQuery        = "SELECT id, nickname, about, email, fullname FROM users WHERE nickname = $1;"
	FindUserByEmailOrNicknameQuery = "SELECT nickname, about, email, fullname FROM users WHERE email = $1 OR nickname = $2;"
	CreateUserQuery                = `INSERT INTO users (nickname, fullname, about, email)
				  			   		  VALUES ($1, $2, $3, $4) RETURNING nickname, fullname, about, email;`
	UpdateUserQuery  = "UPDATE users SET fullname = $2, about = $3, email = $4 WHERE nickname = $1 RETURNING nickname, fullname, about, email;"
	CreateForumQuery = `INSERT INTO forums (title, username, slug)
				  		VALUES ($1, $2, $3) RETURNING title, username, slug, posts, threads;`
	FindForumBySlugQuery = "SELECT id, title, username, slug, posts, threads FROM forums WHERE slug = $1;"
	CreateThreadQuery    = `INSERT INTO threads (title, author, forum, message, slug, created)
								 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, title, author, forum, message, votes, slug, created;`
	UpdateForumsThreadCountQuery = "UPDATE forums SET threads = threads + 1 WHERE slug = $1 RETURNING id;"
	UpdateForumsPostsCountQuery  = "UPDATE forums SET posts = posts + $1 WHERE slug = $2 RETURNING id;"
	FindThreadBySlugQuery        = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE slug = $1;"
	FindThreadBySlugOrIdQuery    = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1 OR (slug = $2 AND slug <> '');"
	FindThreadByIdQuery          = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE id = $1;"
	FindThreadsByForumQuery      = "SELECT id, title, author, forum, message, votes, slug, created FROM threads WHERE forum = $1"
	CreateThreadStartQuery       = "INSERT INTO posts (id, parent, path, author, message, forum, thread, created) VALUES "
	FindParentIdForPostQuery     = "SELECT thread FROM posts WHERE id = $1;"
	FindVoteQuery                = "SELECT id FROM votes WHERE user_id = $1 AND thread_id = $2;"
	UpdateVoteQuery              = "UPDATE votes SET voice = $3 WHERE user_id = $1 AND thread_id = $2 RETURNING id;"
	AddVoteQuery                 = "INSERT INTO votes (user_id, thread_id, voice) VALUES ($1, $2, $3) RETURNING id;"
	GetPostsStartQuery           = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE thread = $1"
	UpdateThreadQuery            = "UPDATE threads SET title = $1, message = $2 WHERE id = $3 RETURNING id, title, author, forum, message, votes, slug, created;"
	GetForumUsersStartQuery      = "SELECT id, nickname, about, email, fullname FROM users WHERE id IN (SELECT user_id FROM forum_users WHERE forum_id = $1)"
	GetPostInfoQuery             = "SELECT id, parent, author, message, isEdited, forum, thread, created FROM posts WHERE id = $1;"
	UpdatePostQuery              = "UPDATE posts SET parent = $2, author = $3, message = $4, isEdited = $5, forum = $6, thread = $7, created = $8 WHERE id = $1 RETURNING id, parent, author, message, isEdited, forum, thread, created;"
	GetServiceStatusQuery        = `SELECT
									(SELECT COUNT(*) FROM forums) AS forum, 
									(SELECT COUNT(*) FROM posts) AS post, 
									(SELECT COUNT(*) FROM threads) AS thread, 
									(SELECT COUNT(*) FROM users) AS user;`
	ClearServiceQuery = "TRUNCATE forums, posts, threads, users, votes CASCADE;"
)
