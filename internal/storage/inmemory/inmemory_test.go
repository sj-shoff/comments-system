package inmemory_test

import (
	"comments-system/internal/models"
	"comments-system/internal/storage/inmemory"
	"comments-system/pkg/errors"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInMemoryStorage(t *testing.T) {
	ctx := context.Background()
	storage := inmemory.NewInMemory()

	t.Run("Create and Get Post", func(t *testing.T) {
		post := models.Post{
			Title:   "Test Post",
			Content: "Content",
			Author:  "Author",
		}

		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)
		require.NotEmpty(t, createdPost.ID)

		gotPost, err := storage.GetPost(ctx, createdPost.ID)
		require.NoError(t, err)
		require.Equal(t, createdPost, gotPost)
	})

	t.Run("Get Post Not Found", func(t *testing.T) {
		_, err := storage.GetPost(ctx, "nonexistent")
		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("Update Post", func(t *testing.T) {
		post := models.Post{
			Title:   "Update Test",
			Content: "Before update",
			Author:  "Author",
		}

		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		createdPost.Title = "Updated Title"
		err = storage.UpdatePost(ctx, createdPost)
		require.NoError(t, err)

		updatedPost, err := storage.GetPost(ctx, createdPost.ID)
		require.NoError(t, err)
		require.Equal(t, "Updated Title", updatedPost.Title)
	})

	t.Run("Update Post Not Found", func(t *testing.T) {
		err := storage.UpdatePost(ctx, models.Post{ID: "nonexistent"})
		require.ErrorIs(t, err, errors.ErrNotFound)
	})

	t.Run("Get Posts with pagination", func(t *testing.T) {
		storage = inmemory.NewInMemory()

		for i := 0; i < 3; i++ {
			post := models.Post{
				Title:   fmt.Sprintf("Post %d", i),
				Content: "Content",
				Author:  "Author",
			}
			_, err := storage.CreatePost(ctx, post)
			require.NoError(t, err)
		}

		posts, err := storage.GetPosts(ctx, 2, 0)
		require.NoError(t, err)
		require.Len(t, posts, 2)

		posts, err = storage.GetPosts(ctx, 2, 2)
		require.NoError(t, err)
		require.Len(t, posts, 1)
	})

	t.Run("Create and Get Comment", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for comment",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		comment := models.Comment{
			PostID:  createdPost.ID,
			Author:  "Commenter",
			Content: "Comment",
		}

		createdComment, err := storage.CreateComment(ctx, comment)
		require.NoError(t, err)
		require.NotEmpty(t, createdComment.ID)

		gotComment, err := storage.GetComment(ctx, createdComment.ID)
		require.NoError(t, err)
		require.Equal(t, createdComment, gotComment)
	})

	t.Run("Create Comment with Parent", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for comment tree",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		parentComment := models.Comment{
			PostID:  createdPost.ID,
			Author:  "Parent",
			Content: "Parent comment",
		}
		createdParent, err := storage.CreateComment(ctx, parentComment)
		require.NoError(t, err)

		childComment := models.Comment{
			PostID:   createdPost.ID,
			ParentID: &createdParent.ID,
			Author:   "Child",
			Content:  "Child comment",
		}
		createdChild, err := storage.CreateComment(ctx, childComment)
		require.NoError(t, err)

		replies, err := storage.GetCommentReplies(ctx, createdParent.ID)
		require.NoError(t, err)
		require.Len(t, replies, 1)
		require.Equal(t, createdChild.ID, replies[0].ID)
	})

	t.Run("Create Comment with Nonexistent Parent", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for invalid comment",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		parentID := "nonexistent"
		childComment := models.Comment{
			PostID:   createdPost.ID,
			ParentID: &parentID,
			Author:   "Child",
			Content:  "Child comment",
		}
		_, err = storage.CreateComment(ctx, childComment)
		require.ErrorIs(t, err, errors.ErrParentNotFound)
	})

	t.Run("Get Comments By Post", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for comments",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			comment := models.Comment{
				PostID:  createdPost.ID,
				Author:  fmt.Sprintf("Author %d", i),
				Content: fmt.Sprintf("Comment %d", i),
			}
			_, err := storage.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		comments, err := storage.GetCommentsByPost(ctx, createdPost.ID, 2, 0)
		require.NoError(t, err)
		require.Len(t, comments, 2)

		comments, err = storage.GetCommentsByPost(ctx, createdPost.ID, 2, 2)
		require.NoError(t, err)
		require.Len(t, comments, 1)
	})

	t.Run("Count Comments By Post", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for count",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		for i := 0; i < 3; i++ {
			comment := models.Comment{
				PostID:  createdPost.ID,
				Author:  fmt.Sprintf("Author %d", i),
				Content: fmt.Sprintf("Comment %d", i),
			}
			_, err := storage.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		count, err := storage.CountCommentsByPost(ctx, createdPost.ID)
		require.NoError(t, err)
		require.Equal(t, 3, count)
	})

	t.Run("Get Comment Replies", func(t *testing.T) {
		post := models.Post{
			Title:   "Post for replies",
			Content: "Content",
			Author:  "Author",
		}
		createdPost, err := storage.CreatePost(ctx, post)
		require.NoError(t, err)

		parent := models.Comment{
			PostID:  createdPost.ID,
			Author:  "Parent",
			Content: "Parent comment",
		}
		createdParent, err := storage.CreateComment(ctx, parent)
		require.NoError(t, err)

		for i := 0; i < 2; i++ {
			reply := models.Comment{
				PostID:   createdPost.ID,
				ParentID: &createdParent.ID,
				Author:   fmt.Sprintf("Reply %d", i),
				Content:  fmt.Sprintf("Reply content %d", i),
			}
			_, err := storage.CreateComment(ctx, reply)
			require.NoError(t, err)
		}

		replies, err := storage.GetCommentReplies(ctx, createdParent.ID)
		require.NoError(t, err)
		require.Len(t, replies, 2)
	})
}
