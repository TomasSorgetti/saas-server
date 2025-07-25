package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// Función que usa la API de Resend para mandar email
func (es *EmailService) sendEmail(job EmailJob) error {
    // Ejemplo simplificado con un request HTTP (usá tu cliente HTTP preferido)

    if es.apiKey == "" {
        return errors.New("API Key is not set")
    }

    // Aquí iría la lógica real para llamar a la API de Resend
    // Por ejemplo, un POST a https://api.resend.com/emails con job.To, job.Subject y job.Body

    fmt.Printf("Sending email to %s with subject %s\n", job.To, job.Subject)

    // Simulamos éxito
    return nil
}