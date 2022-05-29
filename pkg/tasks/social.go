package tasks

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/devict/job-board/pkg/config"
	"github.com/devict/job-board/pkg/data"
	"github.com/devict/job-board/pkg/services"
	"github.com/jmoiron/sqlx"
)

type TaskRunner struct {
	db     *sqlx.DB
	config config.Config
	ctx    context.Context
}

func NewTaskRunner(ctx context.Context, db *sqlx.DB, cfg config.Config) *TaskRunner {
	return &TaskRunner{
		db:     db,
		config: cfg,
		ctx:    ctx,
	}
}

func (tr *TaskRunner) StartBackgroundTasks(wg sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for {
			// Tasks to run in background go here:
			tr.removeOldJobs()
			tr.postUnpublishedJobs()

			select {
			case <-tr.ctx.Done():
				log.Println("shutting down old jobs background process")
				return
			case <-ticker.C:
				continue
			}
		}
	}()
}

func (tr *TaskRunner) removeOldJobs() {
	log.Println("removing old jobs")
	_, err := tr.db.Exec("DELETE FROM jobs WHERE published_at < NOW() - INTERVAL '30 DAYS'")
	if err != nil {
		log.Println(fmt.Errorf("error clearing old jobs: %w", err))
	}
}

func (tr *TaskRunner) postUnpublishedJobs() {
	log.Println("posting jobs to socials")
	unPublishedJobs, err := data.GetJobsToPostOnSocials(tr.db)
	if err != nil {
		log.Println(fmt.Errorf("failed to fetch jobs to post on social media: %w", err))
		return
	}

	for _, job := range unPublishedJobs {
		log.Printf("Posting job: \"%v\" to social platforms.", job.Position)
		tr.publishJobOnSocials(job)
	}

	if len(unPublishedJobs) < 1 {
		log.Println("no jobs to post to social media")
	}
}

func (tr *TaskRunner) publishJobOnSocials(job data.Job) {
	if tr.config.SlackHook != "" {
		if err := services.PostToSlack(job, tr.config); err != nil {
			log.Println(fmt.Errorf("failed to postToSlack: %w", err))
			// continuing...
		}
	}

	if tr.config.Twitter.AccessToken != "" {
		if err := services.PostToTwitter(job, tr.config); err != nil {
			log.Println(fmt.Errorf("failed to postToTwitter: %w", err))
			// continuing...
		}
	}

	// Marking job as published will prevent duplicate social media posts.
	//
	// TODO: If all social posts fail for a job then retry.
	job.SetPublishedToSocials(true, tr.db)
}
