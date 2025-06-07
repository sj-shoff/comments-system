package pubsub_test

import (
	"comments-system/internal/models"
	"comments-system/internal/pubsub"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPubSub_SubscribeAndPublish(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID := "post1"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := ps.Subscribe(ctx, postID)
	assert.NoError(t, err)

	comment := &models.Comment{ID: "comment1", PostID: postID, Content: "Test comment"}
	ps.Publish(postID, comment)

	select {
	case receivedComment := <-ch:
		assert.Equal(t, comment, receivedComment)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout waiting for comment")
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID := "post1"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := ps.Subscribe(ctx, postID)
	assert.NoError(t, err)

	cancel() // Unsubscribe

	select {
	case _, ok := <-ch:
		assert.False(t, ok, "Channel should be closed")
	case <-time.After(100 * time.Millisecond):

	}
}

func TestPubSub_MultipleSubscribers(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID := "post1"

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	ch1, err := ps.Subscribe(ctx1, postID)
	assert.NoError(t, err)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	ch2, err := ps.Subscribe(ctx2, postID)
	assert.NoError(t, err)

	comment := &models.Comment{ID: "comment1", PostID: postID, Content: "Test comment"}
	ps.Publish(postID, comment)

	select {
	case receivedComment := <-ch1:
		assert.Equal(t, comment, receivedComment)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout waiting for comment on subscriber 1")
	}

	select {
	case receivedComment := <-ch2:
		assert.Equal(t, comment, receivedComment)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout waiting for comment on subscriber 2")
	}
}

func TestPubSub_PublishToMultiplePosts(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID1 := "post1"
	postID2 := "post2"

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()
	ch1, err := ps.Subscribe(ctx1, postID1)
	assert.NoError(t, err)

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	ch2, err := ps.Subscribe(ctx2, postID2)
	assert.NoError(t, err)

	comment1 := &models.Comment{ID: "comment1", PostID: postID1, Content: "Test comment 1"}
	ps.Publish(postID1, comment1)

	comment2 := &models.Comment{ID: "comment2", PostID: postID2, Content: "Test comment 2"}
	ps.Publish(postID2, comment2)

	select {
	case receivedComment := <-ch1:
		assert.Equal(t, comment1, receivedComment)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout waiting for comment on subscriber 1")
	}

	select {
	case receivedComment := <-ch2:
		assert.Equal(t, comment2, receivedComment)
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout waiting for comment on subscriber 2")
	}
}

func TestPubSub_PublishWithoutSubscribers(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID := "post1"

	comment := &models.Comment{ID: "comment1", PostID: postID, Content: "Test comment"}
	ps.Publish(postID, comment)

}

func TestPubSub_BufferedChannel(t *testing.T) {
	ps := pubsub.NewPubSub()
	postID := "post1"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch, err := ps.Subscribe(ctx, postID)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		comment := &models.Comment{ID: string(rune(i)), PostID: postID, Content: "Test comment"}
		ps.Publish(postID, comment)
	}

	for i := 0; i < 10; i++ {
		select {
		case receivedComment := <-ch:
			_ = receivedComment
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Timeout waiting for comment %d", i)
			return
		}
	}
}
