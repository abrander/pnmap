package main

import (
	"fmt"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

const ips = "ips"
const hostnames = "hostnames"
const useragents = "useragents"
const vendor = "vendor"
const applications = "applications"
const seen = "seen"
const lastseen = "lastseen"
const firstseen = "firstseen"

type gui struct {
	app       *tview.Application
	hostList  *tview.List
	details   *tview.TextView
	help      *tview.TextView
	secondary string

	nics map[string]*NIC
}

func newGUI() *gui {
	g := &gui{
		app:       tview.NewApplication(),
		hostList:  tview.NewList(),
		details:   tview.NewTextView(),
		secondary: ips,
		nics:      make(map[string]*NIC),
	}

	flex := tview.NewFlex()

	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(g.hostList, 0, 30, true)
	flex.AddItem(g.details, 0, 70, false)

	g.hostList.SetBorder(true)
	g.hostList.SetBorderColor(tcell.ColorGray)
	g.hostList.SetTitle(" Stations ")
	g.hostList.SetTitleColor(tcell.ColorGreenYellow)
	g.hostList.SetMainTextColor(tcell.ColorGreen)
	g.hostList.SetSecondaryTextColor(tcell.ColorWhite)
	g.hostList.SetSelectedTextColor(tcell.ColorBlack)
	g.hostList.SetSelectedBackgroundColor(tcell.ColorGreen)

	g.details.SetBorder(true)
	g.details.SetBorderColor(tcell.ColorGray)
	g.details.SetTitleColor(tcell.ColorGreenYellow)
	g.details.SetTextColor(tcell.ColorWhite)
	g.details.SetDynamicColors(true)

	g.app.SetRoot(flex, true)

	g.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			switch event.Rune() {
			case '1':
				g.secondary = ips
			case '2':
				g.secondary = hostnames
			case '3':
				g.secondary = useragents
			case '4':
				g.secondary = vendor
			case '5':
				g.secondary = applications
			case '6':
				g.secondary = seen
			case '7':
				g.secondary = lastseen
			case '8':
				g.secondary = firstseen
			}
			go func() {
				for _, nic := range g.nics {
					g.updateNIC(nic)
				}
				return
			}()

		}
		return event
	})

	return g
}

func (g *gui) Run() error {
	return g.app.Run()
}

func (g *gui) selectHost() {
	selected := g.hostList.GetCurrentItem()

	mac, _ := g.hostList.GetItemText(selected)

	nic, found := g.nics[mac[0:17]]

	if !found {
		// This should not happen
		return
	}

	g.details.SetTitle(" " + nic.MAC + " ")
	g.details.SetText(nic.String())
}

func (g *gui) updateNIC(nic *NIC) {
	defer g.app.Draw()

	var sec string

	switch g.secondary {
	default:
		sec = fmt.Sprintf("  %v", nic.IPs)
	case ips:
		sec = fmt.Sprintf("  %v", nic.IPs)
	case hostnames:
		sec = fmt.Sprintf("  %v", nic.Hostnames)
	case useragents:
		sec = fmt.Sprintf("  %v", nic.UserAgents)
	case vendor:
		sec = fmt.Sprintf("  %v", nic.Vendor)
	case applications:
		sec = fmt.Sprintf("  %v", nic.Applications)
	case seen:
		sec = fmt.Sprintf("  %v", nic.Seen)
	case lastseen:
		sec = fmt.Sprintf("  %v", nic.LastSeen)
	case firstseen:
		sec = fmt.Sprintf("  %v", nic.FirstSeen)
	}

	g.nics[nic.MAC] = nic

	selected := g.hostList.GetCurrentItem()

	for i := 0; i < g.hostList.GetItemCount(); i++ {
		main, _ := g.hostList.GetItemText(i)
		if main[0:17] == nic.MAC {
			g.hostList.RemoveItem(i)
			g.hostList.InsertItem(i, main, sec, 0, g.selectHost)

			g.hostList.SetCurrentItem(selected)

			// If the current item is selected, update details view.
			if selected == i {
				g.selectHost()
			}

			return
		}
	}

	g.hostList.AddItem(nic.MAC+" "+OUIVendor(nic.MAC), sec, 0, g.selectHost)
}
