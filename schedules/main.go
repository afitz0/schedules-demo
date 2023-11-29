package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	"go.uber.org/zap/zapcore"

	"github.com/pborman/uuid"

	"schedules"
	"schedules/zapadapter"
)

func main() {
	CreateDailySchedule()
}

func CreateDailySchedule() {
	customer := schedules.Customer{}
	action := &client.ScheduleWorkflowAction{
		ID:                 "recommendations-email-" + uuid.New(),
		Workflow:           schedules.RecommendationsWorkflow,
		TaskQueue:          "schedules-demo",
		Args:               []interface{}{customer},
		WorkflowRunTimeout: 2 * time.Minute,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	logger := zapadapter.NewZapAdapter(zapadapter.NewZapLogger(zapcore.DebugLevel))
	c, err := client.Dial(client.Options{
		Logger: logger,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// Swap out with the desired spec
	//spec := MakeSpecDailyCron()
	//spec := MakeSpecDailyCalendar()
	spec := MakeSpecDailyInterval()

	scheduleID := "daily-recommendations"
	sClient := c.ScheduleClient()
	_, err = sClient.Create(ctx, client.ScheduleOptions{
		ID:             scheduleID,
		Spec:           spec,
		Action:         action,
		Overlap:        enums.SCHEDULE_OVERLAP_POLICY_BUFFER_ONE,
		PauseOnFailure: true,
	})

	if err == temporal.ErrScheduleAlreadyRunning {
		logger.Info("Schedule already registered", "ScheduleID", scheduleID)
	} else if err != nil {
		logger.Error("Failed to create schedule", "Error", err)
	} else {
		logger.Info("Schedule created", "ScheduleID", scheduleID)
	}
	return
}

func MakeSpecDailyCron() client.ScheduleSpec {
	return client.ScheduleSpec{
		CronExpressions: []string{
			"0 0 * * 1-5",
		},
	}
}

func MakeSpecDailyCalendar() client.ScheduleSpec {
	return client.ScheduleSpec{
		Calendars: []client.ScheduleCalendarSpec{
			{
				Hour:      []client.ScheduleRange{{Start: 0, End: 0, Step: 1}}, // Default, but good to see anyways.
				Minute:    []client.ScheduleRange{{Start: 0, End: 0, Step: 1}}, // Default, but good to see anyways.
				DayOfWeek: []client.ScheduleRange{{Start: 1, End: 5}},          // Monday ~ Friday
			},
		},
		TimeZoneName: "US/Pacific",
	}
}

func MakeSpecDailyInterval() client.ScheduleSpec {
	return client.ScheduleSpec{
		Intervals: []client.ScheduleIntervalSpec{
			{
				Every: 24 * time.Hour,
			},
		},
		TimeZoneName: "US/Pacific",
		Skip: []client.ScheduleCalendarSpec{
			{
				DayOfWeek: []client.ScheduleRange{{Start: 0, End: 0, Step: 1}}, // Sunday
			},
			{
				DayOfWeek: []client.ScheduleRange{{Start: 6, End: 6, Step: 1}}, // Saturday
			},
		},
	}
}
