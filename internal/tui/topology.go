package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"

	"dead-drop-dispatch/internal/game"
)

const (
	defaultTopologyWidth  = 56
	defaultTopologyHeight = 18
)

type RouteTopologyView struct {
	State              game.GameState
	Job                *game.Job
	SelectedRouteIndex int
	Width              int
	Height             int
	Styles             Styles
}

func RenderRouteTopology(view RouteTopologyView) string {
	styles := view.Styles
	if styles.Base.GetForeground() == nil {
		styles = NewStyles()
	}

	width := view.Width
	if width <= 0 {
		width = defaultTopologyWidth
	}
	height := view.Height
	if height <= 0 {
		height = defaultTopologyHeight
	}

	districts := districtByID(view.State)
	selectedRoute, hasRoute := selectedTopologyRoute(view.Job, view.SelectedRouteIndex)
	routeDistricts := routeDistrictSet(selectedRoute, hasRoute)

	lines := []string{styles.Accent.Render("District topology")}
	lines = append(lines, renderTopologyGrid(view.Job, districts, routeDistricts, styles)...)

	if view.Job != nil {
		lines = append(lines, styles.PanelText.Render(" "))
		lines = append(lines, renderTopologyJobSummary(*view.Job, selectedRoute, hasRoute, districts, styles)...)
	}

	if len(view.State.Districts) == 0 {
		lines = append(lines, styles.Muted.Render("No district telemetry."))
	}

	return strings.Join(fitTopologyLines(lines, width, height, styles), "\n")
}

func renderTopologyGrid(job *game.Job, districts map[game.DistrictID]game.District, routeDistricts map[game.DistrictID]bool, styles Styles) []string {
	rows := [][]game.DistrictID{
		{"northline", "crown_verge"},
		{"saint_orison_market", "port_kestrel"},
		{"floodglass", "ashgate_yard"},
	}

	lines := []string{}
	for rowIndex, row := range rows {
		left := renderTopologyNode(row[0], job, districts, routeDistricts, styles)
		right := renderTopologyNode(row[1], job, districts, routeDistricts, styles)
		spacer := styles.PanelText.Render(strings.Repeat(" ", max(2, 28-lipgloss.Width(left))))
		lines = append(lines, left+spacer+right)
		if rowIndex < len(rows)-1 {
			lines = append(lines, styles.Muted.Render("     | market/service        | artery/drone"))
		}
	}
	return lines
}

func renderTopologyNode(id game.DistrictID, job *game.Job, districts map[game.DistrictID]game.District, routeDistricts map[game.DistrictID]bool, styles Styles) string {
	district, ok := districts[id]
	if !ok {
		return styles.PanelText.Render("[ ] " + string(id))
	}

	marker := " "
	style := styles.PanelText
	switch {
	case job != nil && district.ID == job.Origin:
		marker = "O"
		style = styles.Accent
	case job != nil && district.ID == job.Destination:
		marker = "D"
		style = styles.Warning
	case routeDistricts[district.ID]:
		marker = "R"
		style = styles.InlineCode
	}

	pressure := ""
	switch {
	case district.Surveillance >= 4:
		pressure = " surveil"
	case district.Danger >= 4:
		pressure = " danger"
	case district.SignalQuality <= 2:
		pressure = " signal"
	}

	return style.Render(fmt.Sprintf("[%s] %s%s", marker, clipText(district.Name, 19), pressure))
}

func renderTopologyJobSummary(job game.Job, route game.Route, hasRoute bool, districts map[game.DistrictID]game.District, styles Styles) []string {
	origin := topologyDistrictName(job.Origin, districts)
	destination := topologyDistrictName(job.Destination, districts)
	lines := []string{
		styles.PanelText.Render(clipText(job.Title, 52)),
		styles.PanelText.Render(clipText(fmt.Sprintf("%s -> %s", origin, destination), 52)),
	}

	if !hasRoute {
		lines = append(lines, styles.Muted.Render("No route selected."))
		return lines
	}

	lines = append(lines,
		styles.PanelText.Render(clipText(fmt.Sprintf("Route: %s  %dT", route.Name, route.TimeCost), 52)),
		styles.PanelText.Render(clipText("Path: "+formatTopologyPath(route.Districts, districts), 52)),
	)
	if len(route.Traits) > 0 {
		lines = append(lines, styles.Muted.Render(clipText("Traits: "+strings.Join(route.Traits, ", "), 52)))
	}
	return lines
}

func selectedTopologyRoute(job *game.Job, selected int) (game.Route, bool) {
	if job == nil || len(job.Routes) == 0 {
		return game.Route{}, false
	}
	return job.Routes[clampIndex(selected, len(job.Routes))], true
}

func routeDistrictSet(route game.Route, ok bool) map[game.DistrictID]bool {
	set := map[game.DistrictID]bool{}
	if !ok {
		return set
	}
	for _, districtID := range route.Districts {
		set[districtID] = true
	}
	return set
}

func districtByID(state game.GameState) map[game.DistrictID]game.District {
	districts := make(map[game.DistrictID]game.District, len(state.Districts))
	for _, district := range state.Districts {
		districts[district.ID] = district
	}
	return districts
}

func topologyDistrictName(id game.DistrictID, districts map[game.DistrictID]game.District) string {
	if district, ok := districts[id]; ok {
		return district.Name
	}
	return string(id)
}

func formatTopologyPath(path []game.DistrictID, districts map[game.DistrictID]game.District) string {
	if len(path) == 0 {
		return "unknown"
	}
	parts := make([]string, 0, len(path))
	for _, districtID := range path {
		parts = append(parts, topologyDistrictName(districtID, districts))
	}
	return strings.Join(parts, " -> ")
}

func fitTopologyLines(lines []string, width int, height int, styles Styles) []string {
	if height <= 0 {
		return nil
	}
	fitted := make([]string, 0, min(len(lines), height))
	for _, line := range lines {
		if lipgloss.Width(line) > width {
			line = lipgloss.NewStyle().MaxWidth(width).Render(line)
		}
		fitted = append(fitted, line)
		if len(fitted) == height {
			return fitted
		}
	}
	return fitted
}
