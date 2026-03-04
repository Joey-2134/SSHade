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

	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/logging"
)

const (
	host = "localhost"
	port = "23234"
)

//go:embed banner.txt
var banner string
var users = map[string]string{
	// You can add add your name and public key here :)
	"notjoey": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIIst2WyjFbYpwCoyDwKNe46cLYIoh76ZBq1Q5zvuLb74 joey@Joeys-MacBook-Air.local",
}

func main() {
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
					// if the current session's user public key is one of the
					// known users, we greet them and return.
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
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	// Before starting our server, we create a channel and listen for some
	// common interrupt signals.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// We then start the server in a goroutine, as we'll listen for the done
	// signal later.
	go func() {
		log.Info("Starting SSH server", "host", host, "port", port)
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			// We ignore ErrServerClosed because it is expected.
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	// Here we wait for the done signal: this can be either an interrupt, or
	// the server shutting down for any other reason.
	<-done

	// When it arrives, we create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()

	// When we start the shutdown, the server will no longer accept new
	// connections, but will wait as much as the given context allows for the
	// active connections to finish.
	// After the timeout, it shuts down anyway.
	log.Info("Stopping SSH server")
	if err := srv.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
