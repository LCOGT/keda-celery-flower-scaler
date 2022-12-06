package main

import (
	"context"
	"fmt"
	"net/url"
)

type QueueLengthResponse struct {
	ActiveQueues []struct {
		Name     string `json:"name"`
		Messages int    `json:"messages"`
	} `json:"active_queues"`
}

type TaskDetails struct {
	State string `json:"state"`
}

func (self *CeleryScaler) queryRunningAndPendingTasks(ctx context.Context, flowerURL *url.URL) (int, error) {
	c := self.client.Clone().SetBaseURL(flowerURL.String())

	var queueLengths QueueLengthResponse

	err := c.Get("/api/queues/length").SetResult(&queueLengths).Do(ctx).Err
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	tasks := make(map[string]TaskDetails)

	err = c.Get("/api/tasks").SetResult(&tasks).Do(ctx).Err
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}

	count := 0

	for _, ac := range queueLengths.ActiveQueues {
		count += ac.Messages
	}

	for _, td := range tasks {
		if td.State == "STARTED" || td.State == "RECEIVED" {
			count += 1
		}
	}

	return count, nil
}
