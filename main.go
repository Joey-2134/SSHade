package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	gocrypto "golang.org/x/crypto/ssh"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
	"github.com/Joey-2134/SSHade/db"
	ui "github.com/Joey-2134/SSHade/ui/screens"
)

const (
	defaultHost        = "0.0.0.0"
	defaultPort        = "2222"
	defaultDBPath      = db.DefaultPath
	defaultHostKeyPath = ".ssh/id_ed25519"
)

//go:embed banner.txt
var banner string

var (
	activeSessions = make(map[string]struct{})
	activeMu       sync.Mutex
)

func main() {
	host := getenvOrDefault("SSHADE_HOST", defaultHost)
	port := getenvOrDefault("SSHADE_PORT", defaultPort)
	dbPath := getenvOrDefault("SSHADE_DB_PATH", defaultDBPath)
	hostKeyPath := getenvOrDefault("SSHADE_HOST_KEY_PATH", defaultHostKeyPath)

	database, err := db.Open(dbPath)
	if err != nil {
		log.Fatal("Failed to open database", "error", err)
	}
	defer database.Close()

	c := canvas.New(constants.GridSize, constants.GridSize)
	if err := c.LoadFromDB(context.Background(), database); err != nil {
		log.Fatal("Failed to load canvas from DB", "error", err)
	}
	broadcaster := canvas.NewBroadcaster()
	c.SetBroadcaster(broadcaster)

	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithBannerHandler(func(ctx ssh.Context) string {
			return fmt.Sprintf("%s", banner)
		}),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithMiddleware(
			logging.Middleware(),
			bubbletea.Middleware(func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
				return ui.TeaHandler(sess, c, database, broadcaster)
			}),
			activeterm.Middleware(),
			oneSessionPerFingerprintMiddleware,
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

func getenvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func oneSessionPerFingerprintMiddleware(next ssh.Handler) ssh.Handler {
	return func(sess ssh.Session) {
		fingerprint := ""
		if pk := sess.PublicKey(); pk != nil {
			fingerprint = gocrypto.FingerprintSHA256(pk)
		}

		activeMu.Lock()
		if _, ok := activeSessions[fingerprint]; ok {
			activeMu.Unlock()
			wish.Fatalln(sess, "You already have an active session. Only one session per SSH key is allowed.")
			return
		}
		activeSessions[fingerprint] = struct{}{}
		activeMu.Unlock()

		defer func() {
			activeMu.Lock()
			delete(activeSessions, fingerprint)
			activeMu.Unlock()
		}()

		next(sess)
	}
}
