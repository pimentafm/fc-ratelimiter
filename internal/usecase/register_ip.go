package usecase

import (
	"context"
	"log"

	"github.com/pimentafm/fc-ratelimiter/config"
	"github.com/pimentafm/fc-ratelimiter/internal/dto"
	"github.com/pimentafm/fc-ratelimiter/internal/entity"
)

type RegisterIP struct {
	ipRepo entity.IPRepository
	config *config.Config
}

func NewRegisterIPUseCase(
	ipRepo entity.IPRepository,
	config *config.Config,
) *RegisterIP {
	return &RegisterIP{
		ipRepo: ipRepo,
		config: config,
	}
}

// Execute This saves a new request depending on IP.
// If we have a request from an endpoint that does not have API Key, we persist the request using the IP
//  1. We need to confirm if the IP has been blocked for exceeding the maximum number of requests.
//  2. We save the IP in the database using the RateLimiter.Allow() and using the environmental variables
//     or if we already have the same IP, we update the requests array.
//  3. Finally, we execute the validation and update/insert the data in the database.
//  4. If we have reached the maximum request amount, we save the IP into the database.
func (ipr *RegisterIP) Execute(
	ctx context.Context,
	input dto.IPRequestSave,
) (dto.IPRequestResult, error) {
	status, blockedErr := ipr.ipRepo.GetBlockedDuration(ctx, input.IP)
	if blockedErr != nil {
		return dto.IPRequestResult{}, blockedErr
	}

	if status == entity.StatusIPBlocked {
		log.Println("ip is blocked due to exceeding the maximum number of requests")
		return dto.IPRequestResult{}, entity.ErrIPExceededAmountRequest
	}

	getRequest, getReqErr := ipr.ipRepo.GetRequest(ctx, input.IP)
	if getReqErr != nil {
		log.Printf("Error getting IP requests: %s \n", getReqErr.Error())
		return dto.IPRequestResult{}, getReqErr
	}

	getRequest.TimeWindowSec = ipr.config.RateLimiter.ByIP.TimeWindow
	getRequest.MaxRequests = ipr.config.RateLimiter.ByIP.MaxRequests
	if valErr := getRequest.Validate(); valErr != nil {
		log.Printf("Error validation in rate limiter: %s \n", valErr.Error())
		return dto.IPRequestResult{}, valErr
	}

	getRequest.AddRequests(input.TimeAdd)
	isAllowed := getRequest.Allow(input.TimeAdd)
	if upsertErr := ipr.ipRepo.UpsertRequest(ctx, input.IP, getRequest); upsertErr != nil {
		log.Printf("Error updating/inserting rate limit: %s \n", upsertErr.Error())
		return dto.IPRequestResult{}, upsertErr
	}

	if !isAllowed {
		if saveErr := ipr.ipRepo.SaveBlockedDuration(
			ctx,
			input.IP,
			ipr.config.RateLimiter.ByIP.BlockedDuration,
		); saveErr != nil {
			return dto.IPRequestResult{}, saveErr
		}
	}

	return dto.IPRequestResult{
		Allow: isAllowed,
	}, nil
}
