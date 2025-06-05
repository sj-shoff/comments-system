package inmemory

import (
	"comments-system/internal/models"
	"comments-system/pkg/errors"
	"comments-system/pkg/utils"
	"context"
	"sort"
	"sync"
	"time"
)

type Storage struct {
	postsMu      sync.RWMutex
	posts        map[string]models.Post
	commentsMu   sync.RWMutex
	comments     map[string]models.Comment
	postComments map[string][]string
	commentTree  map[string][]string
}

func NewInMemory() *Storage {
	return &Storage{
		posts:        make(map[string]models.Post),
		comments:     make(map[string]models.Comment),
		postComments: make(map[string][]string),
		commentTree:  make(map[string][]string),
	}
}

func (s *Storage) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	if post.ID == "" {
		post.ID = utils.GenerateID()
	}
	post.CreatedAt = time.Now()
	s.posts[post.ID] = post
	return post, nil
}

func (s *Storage) GetPosts(ctx context.Context) ([]models.Post, error) {
	s.postsMu.RLock()
	defer s.postsMu.RUnlock()

	posts := make([]models.Post, 0, len(s.posts))
	for _, p := range s.posts {
		posts = append(posts, p)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].CreatedAt.After(posts[j].CreatedAt)
	})

	return posts, nil
}

func (s *Storage) GetPost(ctx context.Context, id string) (models.Post, error) {
	s.postsMu.RLock()
	defer s.postsMu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return models.Post{}, errors.ErrNotFound
	}
	return post, nil
}

func (s *Storage) UpdatePost(ctx context.Context, post models.Post) error {
	s.postsMu.Lock()
	defer s.postsMu.Unlock()

	if _, ok := s.posts[post.ID]; !ok {
		return errors.ErrNotFound
	}
	s.posts[post.ID] = post
	return nil
}

func (s *Storage) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	s.commentsMu.Lock()
	defer s.commentsMu.Unlock()

	if _, err := s.GetPost(ctx, comment.PostID); err != nil {
		return models.Comment{}, err
	}

	if comment.ParentID != nil {
		if parent, ok := s.comments[*comment.ParentID]; !ok || parent.PostID != comment.PostID {
			return models.Comment{}, errors.ErrParentNotFound
		}
	}

	comment.ID = utils.GenerateID()
	comment.CreatedAt = time.Now()
	s.comments[comment.ID] = comment

	s.postComments[comment.PostID] = append(s.postComments[comment.PostID], comment.ID)

	if comment.ParentID != nil {
		s.commentTree[*comment.ParentID] = append(s.commentTree[*comment.ParentID], comment.ID)
	}

	return comment, nil
}

func (s *Storage) GetComment(ctx context.Context, id string) (models.Comment, error) {
	s.commentsMu.RLock()
	defer s.commentsMu.RUnlock()

	comment, ok := s.comments[id]
	if !ok {
		return models.Comment{}, errors.ErrNotFound
	}
	return comment, nil
}

func (s *Storage) GetCommentsByPost(ctx context.Context, postID string, limit, offset int) ([]models.Comment, error) {
	s.commentsMu.RLock()
	defer s.commentsMu.RUnlock()

	commentIDs, ok := s.postComments[postID]
	if !ok {
		return nil, nil
	}

	var rootComments []models.Comment
	for _, id := range commentIDs {
		if comment, ok := s.comments[id]; ok && comment.ParentID == nil {
			rootComments = append(rootComments, comment)
		}
	}

	sort.Slice(rootComments, func(i, j int) bool {
		return rootComments[i].CreatedAt.After(rootComments[j].CreatedAt)
	})

	start := offset
	if start > len(rootComments) {
		return []models.Comment{}, nil
	}
	end := start + limit
	if end > len(rootComments) {
		end = len(rootComments)
	}

	return rootComments[start:end], nil
}

func (s *Storage) CountCommentsByPost(ctx context.Context, postID string) (int, error) {
	s.commentsMu.RLock()
	defer s.commentsMu.RUnlock()

	count := 0
	for _, id := range s.postComments[postID] {
		if comment, ok := s.comments[id]; ok && comment.ParentID == nil {
			count++
		}
	}
	return count, nil
}

func (s *Storage) GetCommentReplies(ctx context.Context, parentID string) ([]models.Comment, error) {
	s.commentsMu.RLock()
	defer s.commentsMu.RUnlock()

	replyIDs, ok := s.commentTree[parentID]
	if !ok {
		return nil, nil
	}

	replies := make([]models.Comment, 0, len(replyIDs))
	for _, id := range replyIDs {
		if comment, ok := s.comments[id]; ok {
			replies = append(replies, comment)
		}
	}

	sort.Slice(replies, func(i, j int) bool {
		return replies[i].CreatedAt.Before(replies[j].CreatedAt)
	})

	return replies, nil
}

func (s *Storage) Close() error {
	return nil
}
