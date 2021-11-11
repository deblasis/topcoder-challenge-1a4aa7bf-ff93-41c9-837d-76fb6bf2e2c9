package testing

import (
	"net/url"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/mock"
)

type MockApp struct {
	mock.Mock
}

// Create a new window for the application.
// The first window to open is considered the "master" and when closed
// the application will exit.
func (a *MockApp) NewWindow(title string) fyne.Window {
	args := a.Called(title)
	return args.Get(0).(fyne.Window)
}

// Open a URL in the default browser application.
func (a *MockApp) OpenURL(url *url.URL) error {
	args := a.Called(url)
	return args.Error(0)
}

// Icon returns the application icon, this is used in various ways
// depending on operating system.
// This is also the default icon for new windows.
func (a *MockApp) Icon() fyne.Resource {
	args := a.Called()
	return args.Get(0).(fyne.Resource)
}

// SetIcon sets the icon resource used for this application instance.
func (a *MockApp) SetIcon(r fyne.Resource) {
	a.Called(r)
}

// Run the application - this starts the event loop and waits until Quit()
// is called or the last window closes.
// This should be called near the end of a main() function as it will block.
func (a *MockApp) Run() {
	a.Called()
}

// Calling Quit on the application will cause the application to exit
// cleanly, closing all open windows.
// This function does no thing on a mobile device as the application lifecycle is
// managed by the operating system.
func (a *MockApp) Quit() {
	a.Called()
}

// Driver returns the driver that is rendering this application.
// Typically not needed for day to day work, mostly internal functionality.
func (a *MockApp) Driver() fyne.Driver {
	args := a.Called()
	return args.Get(0).(fyne.Driver)
}

// UniqueID returns the application unique identifier, if set.
// This must be set for use of the Preferences() functions... see NewWithId(string)
func (a *MockApp) UniqueID() string {
	args := a.Called()
	return args.String(0)
}

// SendNotification sends a system notification that will be displayed in the operating system's notification area.
func (a *MockApp) SendNotification(n *fyne.Notification) {
	a.Called(n)
}

// Settings return the globally set settings, determining theme and so on.
func (a *MockApp) Settings() fyne.Settings {
	args := a.Called()
	return args.Get(0).(fyne.Settings)
}

// Preferences returns the application preferences, used for storing configuration and state
func (a *MockApp) Preferences() fyne.Preferences {
	args := a.Called()
	return args.Get(0).(fyne.Preferences)
}

// Storage returns a storage handler specific to this application.
func (a *MockApp) Storage() fyne.Storage {
	args := a.Called()
	return args.Get(0).(fyne.Storage)
}

// Lifecycle returns a type that allows apps to hook in to lifecycle events.
func (a *MockApp) Lifecycle() fyne.Lifecycle {
	args := a.Called()
	return args.Get(0).(fyne.Lifecycle)
}

// An App is the definition of a graphical application.
// Apps can have multiple windows, it will exit when the first window to be
// shown is closed. You can also cause the app to exit by calling Quit().
// To start an application you need to call Run() somewhere in your main() function.
// Alternatively use the window.ShowAndRun() function for your main window.
type App interface {
	// Create a new window for the application.
	// The first window to open is considered the "master" and when closed
	// the application will exit.
	NewWindow(title string) fyne.Window

	// Open a URL in the default browser application.
	OpenURL(url *url.URL) error

	// Icon returns the application icon, this is used in various ways
	// depending on operating system.
	// This is also the default icon for new windows.
	Icon() fyne.Resource

	// SetIcon sets the icon resource used for this application instance.
	SetIcon(fyne.Resource)

	// Run the application - this starts the event loop and waits until Quit()
	// is called or the last window closes.
	// This should be called near the end of a main() function as it will block.
	Run()

	// Calling Quit on the application will cause the application to exit
	// cleanly, closing all open windows.
	// This function does no thing on a mobile device as the application lifecycle is
	// managed by the operating system.
	Quit()

	// Driver returns the driver that is rendering this application.
	// Typically not needed for day to day work, mostly internal functionality.
	Driver() fyne.Driver

	// UniqueID returns the application unique identifier, if set.
	// This must be set for use of the Preferences() functions... see NewWithId(string)
	UniqueID() string

	// SendNotification sends a system notification that will be displayed in the operating system's notification area.
	SendNotification(*fyne.Notification)

	// Settings return the globally set settings, determining theme and so on.
	Settings() fyne.Settings

	// Preferences returns the application preferences, used for storing configuration and state
	Preferences() fyne.Preferences

	// Storage returns a storage handler specific to this application.
	Storage() fyne.Storage

	// Lifecycle returns a type that allows apps to hook in to lifecycle events.
	Lifecycle() fyne.Lifecycle
}
