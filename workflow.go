package schedules

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Customer struct{}
type RecommendationsData struct{}

func RecommendationsWorkflow(ctx workflow.Context, customer Customer) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("Started recommnedations workflow", "Customer", customer)

	var a *Activities
	var customerData RecommendationsData
	err := workflow.ExecuteActivity(ctx, a.GatherDataForCustomer).Get(ctx, &customerData)
	if err != nil {
		return err
	}

	var filteredData RecommendationsData
	err = workflow.ExecuteActivity(ctx, a.FilterDataForCustomer, customerData).Get(ctx, &filteredData)
	if err != nil {
		return err
	}

	var htmlEmail string
	err = workflow.ExecuteActivity(ctx, a.RenderDataForCustomer, filteredData).Get(ctx, &htmlEmail)
	if err != nil {
		return err
	}

	// Limit send email activity to only one attempt
	sendEmailOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, sendEmailOptions)
	err = workflow.ExecuteActivity(ctx, a.SendEmail, htmlEmail).Get(ctx, nil)
	if err != nil {
		return err
	}

	logger.Info("Completed recommendations workflow", "Customer", customer)

	return nil
}
