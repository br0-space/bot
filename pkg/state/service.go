package state

import (
	"sync"
	"time"

	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/br0-space/bot/interfaces"
)

var getLastPostLock = &sync.Mutex{}

type Service struct {
	log              logger.Interface
	userStatsRepo    interfaces.UserStatsRepoInterface
	messageStatsRepo interfaces.MessageStatsRepoInterface
	lastPost         map[int64]time.Time
}

func NewService(
	userStatsRepo interfaces.UserStatsRepoInterface,
	messageStatsRepo interfaces.MessageStatsRepoInterface,
) *Service {
	state := &Service{
		log:              logger.New(),
		userStatsRepo:    userStatsRepo,
		messageStatsRepo: messageStatsRepo,
		lastPost:         make(map[int64]time.Time),
	}
	state.init()

	return state
}

func (s *Service) ProcessMessage(messageIn telegramclient.WebhookMessageStruct) {
	s.updateUserStats(messageIn)
	s.updateMessageStats(messageIn)
}

func (s *Service) GetLastPost(userID int64) *time.Time {
	getLastPostLock.Lock()
	defer getLastPostLock.Unlock()

	if lastPost, ok := s.lastPost[userID]; ok {
		return &lastPost
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (s *Service) init() {
	users, err := s.userStatsRepo.GetKnownUsers()
	if err != nil {
		s.log.Error("Error while getting known users from DB:", err)

		return
	}

	for _, user := range users {
		s.lastPost[user.ID] = user.LastPost
	}
}

func (s *Service) updateUserStats(messageIn telegramclient.WebhookMessageStruct) {
	s.lastPost[messageIn.From.ID] = time.Now()

	if err := s.userStatsRepo.UpdateStats(
		messageIn.From.ID,
		messageIn.From.UsernameOrName(),
	); err != nil {
		s.log.Error("Error while updating user stats in DB:", err)
	}
}

func (s *Service) updateMessageStats(messageIn telegramclient.WebhookMessageStruct) {
	if err := s.messageStatsRepo.InsertMessageStats(
		messageIn.From.ID,
		messageIn.WordCount(),
	); err != nil {
		s.log.Error("Error while inserting message stats in DB:", err)
	}
}
