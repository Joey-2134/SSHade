package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Joey-2134/SSHade/canvas"
)

func Grid(
	width, height int,
	canvasWidth, canvasHeight int,
	renderer *lipgloss.Renderer,
	canvasRef *canvas.Canvas,
	cursorX, cursorY int,
	defaultCellColour string,
) string {
	var b strings.Builder

	scaleByWidth := width / canvasWidth
	scaleByHeight := 2 * (height / canvasHeight)
	cellWidth := max(min(scaleByWidth, scaleByHeight), 1)
	linesPerRow := max(cellWidth/2, 1)

	for y := range canvasHeight {
		for range linesPerRow {
			for x := range canvasWidth {
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
