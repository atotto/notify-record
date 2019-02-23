package main

import (
	"fmt"
	"log"
	"os"

	"github.com/atotto/notify-record/message"

	"github.com/godbus/dbus"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	if err = conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='method_call',path='/org/freedesktop/Notifications',member='Notify',eavesdrop='true'").Err; err != nil {
		log.Fatal(err)
	}
	mch := make(chan *dbus.Message, 10)
	conn.Eavesdrop(mch)

	for v := range mch {
		log.Print(message.String(v))
	}
}
