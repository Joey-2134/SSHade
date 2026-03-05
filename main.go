package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/db"
	"github.com/Joey-2134/SSHade/ui"
)

const (
	host = "localhost"
	port = "23234"
)

//go:embed banner.txt
var banner string
var users = map[string]string{
	"notjoey": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIst2WyjFbYpwCoyDwKNe46cLYIoh76ZBq1Q5zvuLb74 joey@Joeys-MacBook-Air.local",
	"joeypc":  "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIB1mkr8Tr0+/LZ1yOWOj7Wqu68jjXY+LWOIdSnzlhzk2 joeyg@Joeys-PC",
}

func main() {
	database, err := db.Open(db.DefaultPath)
	if err != nil {
		log.Fatal("Failed to open database", "error", err)
	}
	defer database.Close()

	c := canvas.New(ui.CanvasWidth, ui.CanvasHeight)
	if err := c.LoadFromDB(context.Background(), database); err != nil {
		log.Fatal("Failed to load canvas from DB", "error", err)
	}

	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithBannerHandler(func(ctx ssh.Context) string {
			return fmt.Sprintf(banner, ctx.User())
		}),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(sess ssh.Session) {
					for name, pubkey := range users {
						parsed, _, _, _, _ := ssh.ParseAuthorizedKey(
							[]byte(pubkey),
						)
						if ssh.KeysEqual(sess.PublicKey(), parsed) {
							wish.Println(sess, fmt.Sprintf("Hey %s!\n", name))
							next(sess)
							return
						}
					}
					wish.Println(sess, "Hey, I don't know who you are!")
					next(sess)
				}
			},
			logging.Middleware(),
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return ui.TeaHandler(sess, c, database)
			}),
			activeterm.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Starting SSH server", "host", host, "port", port)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	log.Info("Stopping SSH server")
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
