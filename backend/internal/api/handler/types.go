package handler

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/tmythicator/ticker-rush/backend/internal/domain"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/exchange/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/ladder/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/leaderboard/v1"
	"github.com/tmythicator/ticker-rush/backend/internal/proto/user/v1"
	redis "github.com/tmythicator/ticker-rush/backend/internal/repository/redis"
)

// ToExternalQuoteFromValkey maps a ValkeyQuote to a Protobuf Quote.
func ToExternalQuoteFromValkey(vq *redis.ValkeyQuote) *exchange.Quote {
	if vq == nil {
		return nil
	}

	return &exchange.Quote{
		Symbol:        vq.Symbol,
		Price:         vq.Price,
		Change:        vq.Change,
		ChangePercent: vq.ChangePercent,
		Timestamp:     timestamppb.New(time.Unix(vq.Timestamp, 0)),
		Source:        vq.Source,
		IsClosed:      vq.IsClosed,
	}
}

// ToExternalUser maps a domain User to a Protobuf User.
func ToExternalUser(u *domain.User) *user.User {
	if u == nil {
		return nil
	}

	pUser := &user.User{
		Username:        u.Username,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Website:         u.Website,
		IsPublic:        u.IsPublic,
		IsAdmin:         u.IsAdmin,
		IsBanned:        u.IsBanned,
		CreatedAt:       timestamppb.New(u.CreatedAt),
		Balance:         u.Balance.InexactFloat64(),
		IsParticipating: u.IsParticipating,
	}

	if u.Portfolio != nil {
		pUser.Portfolio = make(map[string]*user.PortfolioItem)
		for k, v := range u.Portfolio {
			pUser.Portfolio[k] = ToExternalPortfolioItem(v)
		}
	}

	return pUser
}

// ToExternalPortfolioItem maps a domain PortfolioItem to a Protobuf PortfolioItem.
func ToExternalPortfolioItem(item domain.PortfolioItem) *user.PortfolioItem {
	return &user.PortfolioItem{
		StockSymbol:  item.StockSymbol,
		Quantity:     item.Quantity.InexactFloat64(),
		AveragePrice: item.AveragePrice.InexactFloat64(),
	}
}

// ToExternalLadder maps a domain Ladder to a Protobuf Ladder.
func ToExternalLadder(l *domain.Ladder) *ladder.Ladder {
	if l == nil {
		return nil
	}

	allowed := make([]*ladder.TickerInfo, len(l.AllowedTickers))
	for i, t := range l.AllowedTickers {
		allowed[i] = &ladder.TickerInfo{
			Symbol: t.Symbol,
			Source: t.Source,
		}
	}

	return &ladder.Ladder{
		Id:             l.ID,
		Name:           l.Name,
		Type:           l.Type,
		StartTime:      timestamppb.New(l.StartTime),
		EndTime:        timestamppb.New(l.EndTime),
		IsActive:       l.IsActive,
		CreatedAt:      timestamppb.New(l.CreatedAt),
		InitialBalance: l.InitialBalance.InexactFloat64(),
		AllowedTickers: allowed,
	}
}

// ToExternalQuote maps a domain Quote to a Protobuf Quote.
func ToExternalQuote(q *domain.Quote) *exchange.Quote {
	if q == nil {
		return nil
	}

	return &exchange.Quote{
		Symbol:        q.Symbol,
		Price:         q.Price.InexactFloat64(),
		Change:        q.Change.InexactFloat64(),
		ChangePercent: q.ChangePercent.InexactFloat64(),
		Timestamp:     timestamppb.New(q.Timestamp),
		Source:        q.Source,
		IsClosed:      q.IsClosed,
	}
}

// ToExternalPublicProfile maps a domain User to a Protobuf PublicProfile.
func ToExternalPublicProfile(u *domain.User) *user.PublicProfile {
	if u == nil {
		return nil
	}

	pProfile := &user.PublicProfile{
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Website:   u.Website,
		Balance:   u.Balance.InexactFloat64(),
		IsPublic:  u.IsPublic,
	}

	if u.Portfolio != nil {
		pProfile.Portfolio = make(map[string]*user.PortfolioItem)
		for k, v := range u.Portfolio {
			pProfile.Portfolio[k] = ToExternalPortfolioItem(v)
		}
	}

	return pProfile
}

// ToExternalLeaderboardEntry maps a domain LeaderboardEntry to a Protobuf LeaderboardEntry.
func ToExternalLeaderboardEntry(entry domain.LeaderboardEntry) *leaderboard.LeaderboardEntry {
	return &leaderboard.LeaderboardEntry{
		User:  ToExternalPublicProfile(&entry.User),
		Rank:  entry.Rank,
		Score: entry.Score,
	}
}

// ToExternalLeaderboardResponse maps a domain LeaderboardResponse to a Protobuf GetLeaderboardResponse.
func ToExternalLeaderboardResponse(lr *domain.LeaderboardResponse) *leaderboard.GetLeaderboardResponse {
	if lr == nil {
		return nil
	}
	entries := make([]*leaderboard.LeaderboardEntry, len(lr.Entries))
	for i, e := range lr.Entries {
		entries[i] = ToExternalLeaderboardEntry(e)
	}

	return &leaderboard.GetLeaderboardResponse{
		Entries:    entries,
		TotalCount: lr.TotalCount,
		LastUpdate: lr.LastUpdate,
	}
}
