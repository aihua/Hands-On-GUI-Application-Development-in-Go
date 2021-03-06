package main

import (
	"github.com/andlabs/ui"

	"github.com/PacktPublishing/Hands-On-GUI-Application-Development-in-Go/client"
)

type mainUI struct {
	content                 *ui.Label
	subject, from, to, date *ui.Label
	list                    *ui.Box

	server *client.EmailServer
}

func (m *mainUI) setEmail(e *client.EmailMessage) {
	m.subject.SetText(e.Subject)
	m.to.SetText(e.ToEmailString())
	m.from.SetText(e.FromEmailString())
	m.date.SetText(e.DateString())
	m.content.SetText(e.Content)
}

func (m *mainUI) appendEmail(email *client.EmailMessage) {
	item := ui.NewButton(email.Subject)
	item.OnClicked(func(*ui.Button) {
		m.setEmail(email)
	})
	m.list.Append(item, false)
}

func (m *mainUI) listEmails(list []*client.EmailMessage) {
	for _, email := range list {
		captured := email
		m.appendEmail(captured)
	}
}

func (m *mainUI) incomingEmail(email *client.EmailMessage) {
	ui.QueueMain(func() {
		m.appendEmail(email)
	})
}

func (m *mainUI) buildToolbar() ui.Control {
	toolbar := ui.NewHorizontalBox()

	compose := ui.NewButton("New")
	compose.OnClicked(func(*ui.Button) {
		compose := &composeUI{server: m.server}
		compose.buildUI().Show()
	})

	toolbar.Append(compose, false)
	toolbar.Append(ui.NewButton("Reply"), false)
	toolbar.Append(ui.NewButton("Reply All"), false)

	toolbar.Append(ui.NewLabel(" "), false)
	toolbar.Append(ui.NewVerticalSeparator(), false)
	toolbar.Append(ui.NewLabel(" "), false)
	toolbar.Append(ui.NewButton("Delete"), false)
	toolbar.Append(ui.NewLabel(" "), false)
	toolbar.Append(ui.NewVerticalSeparator(), false)
	toolbar.Append(ui.NewLabel(" "), false)

	toolbar.Append(ui.NewButton("Cut"), false)
	toolbar.Append(ui.NewButton("Copy"), false)
	toolbar.Append(ui.NewButton("Paste"), false)

	return toolbar
}

func (m *mainUI) buildUI() *ui.Window {
	window := ui.NewWindow("GoMail", 600, 400, false)
	window.SetMargined(true)
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})

	m.list = ui.NewVerticalBox()
	inbox := ui.NewGroup("Inbox                         ")
	inbox.SetChild(m.list)

	m.subject = ui.NewLabel("subject")
	m.content = ui.NewLabel("content")
	labels := ui.NewVerticalBox()
	labels.Append(ui.NewLabel("From"), false)
	labels.Append(ui.NewLabel("To"), false)
	labels.Append(ui.NewLabel("Date"), false)

	values := ui.NewVerticalBox()
	m.from = ui.NewLabel("email")
	values.Append(m.from, false)
	m.to = ui.NewLabel("email")
	values.Append(m.to, false)
	m.date = ui.NewLabel("date")
	values.Append(m.date, false)

	meta := ui.NewHorizontalBox()
	meta.SetPadded(true)
	meta.Append(labels, false)
	meta.Append(values, true)

	detail := ui.NewVerticalBox()
	detail.SetPadded(true)
	detail.Append(m.subject, false)
	detail.Append(meta, false)
	detail.Append(ui.NewHorizontalSeparator(), false)
	detail.Append(m.content, true)

	layout := ui.NewVerticalBox()
	layout.SetPadded(true)
	layout.Append(m.buildToolbar(), false)

	content := ui.NewHorizontalBox()
	content.SetPadded(true)
	content.Append(inbox, false)
	content.Append(ui.NewVerticalSeparator(), false)
	content.Append(detail, true)
	layout.Append(content, true)

	window.SetChild(layout)
	return window
}

func main() {
	server := client.NewTestServer()
	err := ui.Main(func() {
		main := new(mainUI)
		main.server = server
		window := main.buildUI()

		main.listEmails(server.ListMessages())
		main.setEmail(server.CurrentMessage())
		go func() {
			incoming := server.Incoming()
			for email := range incoming {
				main.incomingEmail(email)
			}
		}()

		window.Show()
	})
	if err != nil {
		panic(err)
	}
}
