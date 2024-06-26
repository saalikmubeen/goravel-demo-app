package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/saalikmubeen/goravel"
	"github.com/saalikmubeen/goravel-demo-app/handlers"
	"github.com/saalikmubeen/goravel-demo-app/middleware"
	"github.com/saalikmubeen/goravel-demo-app/models"
)

type application struct {
	App        *goravel.Goravel
	Handlers   *handlers.Handlers
	Models     *models.Models
	Middleware *middleware.Middleware
	wg         *sync.WaitGroup
}

func main() {

	app := initGoravel()

	go app.listenForShutDown()
	err := app.App.ListenAndServe()
	if err != nil {
		app.App.ErrorLog.Println(err)
	}

}

func (app *application) listenForShutDown() {
	// Create a quit channel which carries os.Signal values. Use buffered
	// We need to use a buffered channel here because signal.Notify() does not
	// wait for a receiver to be available when sending a signal to the quit channel.
	//  If we had used a regular (non-buffered) channel here instead, a signal could be
	// ‘missed’ if our quit channel is not ready to receive at the exact moment that the
	// signal is sent. By using a buffered channel, we avoid this problem and ensure
	// that we never miss a signal.
	quit := make(chan os.Signal, 1)

	// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and relay
	// them to the quit channel. Any other signal will not be caught by signal.Notify()
	// and will retain their default behavior.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Read the signal from the quit channel. This code will block until a signal is
	// received.
	s := <-quit

	// Log a message to say we caught the signal. Notice that we also call the
	// String() method on the signal to get the signal name and include it in the log
	// entry properties.
	app.App.InfoLog.Printf("caught signal: %s", s.String())

	// ** put any clean up tasks here

	// **

	// Call Wait() to block until our WaitGroup counter is zero. This essentially blocks
	// until the background goroutines have finished. Then we return nil on the shutdownError
	// channel to indicate that the shutdown as compleeted without any issues.
	// Uses sync.WaitGroup to wait for any background goroutines before terminating the application.
	app.wg.Wait()

}
