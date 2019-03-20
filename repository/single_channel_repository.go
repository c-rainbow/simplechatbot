package repository

import (
	models "github.com/c-rainbow/simplechatbot/models"
)

type SingleChannelRepositoryT interface {
	GetAllBots() []int64
}

// SingleChannelRepository repo for a single channel
type SingleChannelRepository struct {
	chanInfo *models.Channel
	baseRepo BaseRepositoryT
}

var _ SingleChannelRepositoryT = (*SingleChannelRepository)(nil)

func (repo *SingleChannelRepository) GetAllBots() []int64 {
	// Copy BotIDs to prevent unexpected modification from outside
	copiedIDs := append([]int64{}, repo.chanInfo.BotIDs...)
	return copiedIDs
}
