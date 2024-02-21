package mongoapp

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"time"
)

type App struct {
	log           *zap.Logger
	connectionURI string
	client        *mongo.Client
}

func New(log *zap.Logger, connectionURI string) *App {
	return &App{log: log, connectionURI: connectionURI}
}

func (a *App) Client() *mongo.Client {
	if a.client != nil {
		return a.client
	}

	a.log.Panic("attempt to retrieve uninitialized mongo client")
	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) MustRunConcurrently() <-chan struct{} {
	started := make(chan struct{})
	go func() {
		if err := a.Run(); err != nil {
			panic(err)
		}
		started <- struct{}{}
	}()
	return started
}

func (a *App) Run() error {
	a.log.Info("attempting to connect to mongo")

	// creating connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(a.connectionURI))
	if err != nil {
		return fmt.Errorf("failed to connect to mongoapp: %w", err)
	}

	a.log.Info("checking mongo connection")

	// checking connection
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("failed to ping mongoapp: %w", err)
	}

	a.log.Info("successfully connected to mongo")
	a.client = client
	return nil
}

func (a *App) Stop() error {
	a.log.Info("starting mongo gracefully shutdown")
	if err := a.client.Disconnect(context.Background()); err != nil {
		a.log.Error("failed disconnecting from mongo", zap.Error(err))
		return fmt.Errorf("failed to disconnect from mongoapp: %w", err)
	}

	a.client = nil
	a.log.Info("mongo shut down gracefully")
	return nil
}
