package app

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-auth/cookie"
	"go-auth/database"
	"go-auth/mail"
	"go-auth/oauth"
	"go-auth/server"
	"go-auth/server/repository"
	"go-auth/utils/templates"
	"go.uber.org/zap"
)

func Start() {
	HTTP := flag.Bool("http", false, "run server in HTTP mode")

	//

	logger, conf, fileClose := parse()
	defer func() {
		if err := fileClose(); err != nil {
			logger.Error("close log file", zap.Error(err))
		}
	}()

	// //

	logger.Info("Loading OAuth")

	t := time.Now()

	auth, err := oauth.New(conf.Oauth)
	if err != nil {
		logger.Error("oauth", zap.Error(err))

		return
	}

	logger.Info("Done", zap.Duration("time", time.Since(t)))

	// //

	ts := templates.New()

	logger.Info("Loading templates")
	t = time.Now()

	err = templates.Load(ts, conf.Host.Templates)
	if err != nil {
		logger.Error("Error", zap.Error(err))

		return
	}

	logger.Info("Done", zap.Duration("time", time.Since(t)))

	// //

	logger.Info("Dialing smtp service")
	t = time.Now()

	smtp, err := mail.NewMail(conf.SMTP)
	if err != nil {
		logger.Error("Error", zap.Error(err))

		return
	}

	logger.Info("Done", zap.Duration("time", time.Since(t)))

	// //

	logger.Info("Initializing repository")
	t = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	db, err := database.New(ctx, conf)
	if err != nil {
		logger.Error("Error", zap.Error(err))

		return
	}

	logger.Info("Done", zap.Duration("time", time.Since(t)))

	defer db.Close(logger)

	// //
	rep := repository.Repository{
		Memcached: repository.Memcached{
			State: db.Memcached.State,
			Mail:  db.Memcached.Mail,
		},
		//
		Redis: repository.Redis{
			Session: db.Redis.Session,
		},
		//
		Postgres: repository.Postgres{
			User: db.Postgres.User,
		},
		//
		Oauth: repository.Oauth{
			Google: auth.Google,
			Github: auth.Github,
		},
		//
		Cookie: cookie.New(conf.Cookie),
		//
		Mail: smtp,
	}

	// //

	srv := server.New(conf.Host, repository.NewModel(logger, rep, ts, conf.Host))

	ch := make(chan os.Signal, 1)

	go func() {
		logger.Info("Server started at " + srv.Addr)

		if *HTTP {
			err = srv.ListenAndServe()
		} else {
			err = srv.ListenAndServeTLS(conf.Host.Cert, conf.Host.Key)
		}

		if err != nil {
			if err != http.ErrServerClosed {
				logger.Error("Server", zap.Error(err))
			}
			ch <- nil
		}
	}()

	// //

	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	<-ch
	_, _ = os.Stdout.WriteString("\n")

	//

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	logger.Info("Shutting down server")

	t = time.Now()

	//

	err = srv.Shutdown(ctx)
	if err != nil {
		logger.Error("Error", zap.Error(err))
	} else {
		logger.Info("Done", zap.Duration("time", time.Since(t)))
	}
}
