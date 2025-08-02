package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleConfig struct {
    ClientID     string
    ClientSecret string
    RedirectURL  string
    Scopes       []string
}

func NewGoogleConfig() *GoogleConfig {
    return &GoogleConfig{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
    }
}

func (c *GoogleConfig) OAuth2Config() *oauth2.Config {
    return &oauth2.Config{
        ClientID:     c.ClientID,
        ClientSecret: c.ClientSecret,
        RedirectURL:  c.RedirectURL,
        Scopes:       c.Scopes,
        Endpoint:     google.Endpoint,
    }
}

type GoogleUserInfo struct {
    ID            string `json:"id"`
    Email         string `json:"email"`
    VerifiedEmail bool   `json:"verified_email"`
    Name          string `json:"name"`
    GivenName     string `json:"given_name"`
    FamilyName    string `json:"family_name"`
    Picture       string `json:"picture"`
}

// GetUserInfo obtiene la información del usuario desde Google usando el token de acceso
func GetUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUserInfo, error) {
    client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
    resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
    if err != nil {
        return nil, fmt.Errorf("error al obtener información del usuario: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error al leer respuesta: %w", err)
    }

    var userInfo GoogleUserInfo
    if err := json.Unmarshal(body, &userInfo); err != nil {
        return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
    }

    return &userInfo, nil
}