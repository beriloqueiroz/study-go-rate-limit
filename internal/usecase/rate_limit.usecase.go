package usecase

type RateLimitUseCase struct {
	rateLimitRepository RateLimitRepository
}

func NewRateLimitUseCase(rateLimitRepository RateLimitRepository) *RateLimitUseCase {
	return &RateLimitUseCase{
		rateLimitRepository: rateLimitRepository,
	}
}

type RateLimitUseCaseInputDto struct {
	Ip  string
	Key string
}

type RateLimitUseCaseOutputDto struct {
	Allow bool
}

func (uc *RateLimitUseCase) Execute(input RateLimitUseCaseInputDto) {

}
