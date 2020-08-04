package main

import (
	"fmt"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
)

type gui struct {
	app      *tview.Application
	hostList *tview.List
	details  *tview.TextView

	nics map[string]*NIC
}

func newGUI() *gui {
	g := &gui{
		app:      tview.NewApplication(),
		hostList: tview.NewList(),
		details:  tview.NewTextView(),
		nics:     make(map[string]*NIC),
	}

	flex := tview.NewFlex()

	flex.SetDirection(tview.FlexColumn)
	flex.AddItem(g.hostList, 0, 30, true)
	flex.AddItem(g.details, 0, 70, false)

	g.hostList.SetBorder(true)
	g.hostList.SetBorderColor(tcell.ColorGray)
	g.hostList.SetTitle(" Stations ")
	g.hostList.SetTitleColor(tcell.ColorGreenYellow)

	g.details.SetBorder(true)
	g.details.SetBorderColor(tcell.ColorGray)
	g.details.SetTitleColor(tcell.ColorGreenYellow)

	g.app.SetRoot(flex, true)

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

	sec := fmt.Sprintf("%v", nic.IPs)

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
