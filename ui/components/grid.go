package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Joey-2134/SSHade/canvas"
	"github.com/Joey-2134/SSHade/constants"
)

func Grid(
	width, height int,
	renderer *lipgloss.Renderer,
	canvasRef *canvas.Canvas,
	cursorX, cursorY int,
	defaultCellColour string,
) string {
	var b strings.Builder

	scaleByWidth := width / constants.GridSize
	scaleByHeight := 2 * (height / constants.GridSize)
	cellWidth := max(min(scaleByWidth, scaleByHeight), 1)
	linesPerRow := max(cellWidth/2, 1)

	for y := range constants.GridSize {
		for range linesPerRow {
			for x := range constants.GridSize {
				colour := defaultCellColour
				if canvasRef != nil {
					if p, ok := canvasRef.PixelAt(x, y); ok {
						colour = p.ColourHex
					}
				}
				if cursorX == x && cursorY == y {
					style := renderer.NewStyle().Background(lipgloss.Color("241")).SetString(strings.Repeat(" ", cellWidth))
					b.WriteString(style.String())
				} else {
					style := renderer.NewStyle().Background(lipgloss.Color(colour)).SetString(strings.Repeat(" ", cellWidth))
					b.WriteString(style.String())
				}
			}
			b.WriteString("\n")
		}
	}
	return b.String()
}
