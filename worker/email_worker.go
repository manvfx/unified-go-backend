package worker

import (
	"context"
	"time"
	"unified-go-backend/config"
	"unified-go-backend/database"
	"unified-go-backend/utils"

	"github.com/go-gomail/gomail"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/errgroup"
)

func SendVerificationEmail(to string, code string, cfg *config.Config) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.SMTPUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/plain", "Your verification code is: "+code)

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)

	return d.DialAndSend(m)
}

func ProcessEmailVerificationJobs(ctx context.Context, cfg *config.Config) error {
	g, ctx := errgroup.WithContext(ctx)

	for {
		result, err := database.RedisClient.BLPop(ctx, 0*time.Second, "email_verification_queue").Result()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return err
		}

		email := result[1]

		g.Go(func() error {
			verificationCode, err := database.RedisClient.Get(ctx, email).Result()
			if err != nil {
				utils.Logger.Errorf("Failed to get verification code for email %s: %v", email, err)
				return err
			}

			if err := SendVerificationEmail(email, verificationCode, cfg); err != nil {
				utils.Logger.Errorf("Failed to send verification email to %s: %v", email, err)
				return err
			}

			utils.Logger.Infof("Verification email sent to %s", email)
			return nil
		})

		if err := g.Wait(); err != nil {
			utils.Logger.Errorf("Failed to process email verification job: %v", err)
		}
	}
}
