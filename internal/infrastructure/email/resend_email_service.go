package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"luthierSaas/internal/infrastructure/queue"
)

type EmailJob struct {
    To      string
    Subject string
    Body    string
}

type EmailService struct {
    queue *queue.Queue
    apiKey string
}

func NewEmailService(q *queue.Queue) *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		log.Fatal("RESEND_API_KEY is not set in environment variables")
	}
	return &EmailService{
		queue:  q,
		apiKey: apiKey,
	}
}

// Enqueue para agregar un email a la cola
func (es *EmailService) SendEmailAsync(ctx context.Context, job EmailJob) error {
    return es.queue.Enqueue(ctx, job)
}

// Worker para procesar la cola y enviar emails
func (es *EmailService) StartWorker(ctx context.Context) {
    for {
        data, err := es.queue.Dequeue(ctx)
        if err != nil {
            log.Printf("Error getting job from queue: %v", err)
            continue
        }
        if data == nil {
            continue
        }

        var job EmailJob
        if err := json.Unmarshal(data, &job); err != nil {
            log.Printf("Error parsing job: %v", err)
            continue
        }

        if err := es.sendEmail(job); err != nil {
            log.Printf("Error sending email to %s: %v", job.To, err)
            // Aquí podés decidir reintentar o guardar el error
        } else {
            log.Printf("Email sent successfully to %s", job.To)
        }
    }
}

func (es *EmailService) sendEmail(job EmailJob) error {
	payload := map[string]interface{}{
		"from":    "Luthier SaaS <noreply@tomassorgetti.com.ar>",
		"to":      []string{job.To},
		"subject": job.Subject,
		"html":    job.Body,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+es.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Printf("Email sent to %s", job.To)
		return nil
	}

	var bodyBytes bytes.Buffer
	bodyBytes.ReadFrom(resp.Body)
	return fmt.Errorf("failed to send email: %s", bodyBytes.String())
}