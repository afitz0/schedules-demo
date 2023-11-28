package schedules

import "context"

type Activities struct{}

func (a *Activities) GatherDataForCustomer(ctx context.Context, customer Customer) (RecommendationsData, error) {
	return RecommendationsData{}, nil
}

func (a *Activities) FilterDataForCustomer(ctx context.Context, unfilteredData RecommendationsData) (RecommendationsData, error) {
	return RecommendationsData{}, nil
}

func (a *Activities) RenderDataForCustomer(ctx context.Context, recommendations RecommendationsData) (string, error) {
	return "", nil
}

func (a *Activities) SendEmail(ctx context.Context, body string) (string, error) {
	return "", nil
}
