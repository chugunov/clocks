package main

import (
	"fmt"
	"image/color"
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Labels struct {
	XYs    plotter.XYs
	Labels []string
}

func (l Labels) Len() int {
	return len(l.XYs)
}

func (l Labels) XY(i int) (float64, float64) {
	return l.XYs[i].X, l.XYs[i].Y
}

func (l Labels) Label(i int) string {
	return l.Labels[i]
}

type Plotter struct{}

func (pl *Plotter) DrawSpaceTimeDiagram(history map[int][]Event, filename string) error {
	p := plot.New()

	p.Title.Text = "Space-Time Diagram"
	p.Title.Padding = vg.Points(20)

	timestamps := make(map[int64]struct{})
	for _, events := range history {
		for _, event := range events {
			timestamps[event.timestamp] = struct{}{}
		}
	}

	var sortedTimestamps []int64
	for ts := range timestamps {
		sortedTimestamps = append(sortedTimestamps, ts)
	}
	sort.Slice(
		sortedTimestamps,
		func(i, j int) bool { return sortedTimestamps[i] < sortedTimestamps[j] },
	)

	// add horizontal dashed lines for each timestamp between min and max with 0.5 offset
	var minTimestamp, maxTimestamp int64
	if len(sortedTimestamps) > 0 {
		minTimestamp = sortedTimestamps[0]
		maxTimestamp = sortedTimestamps[len(sortedTimestamps)-1]
	}
	for ts := minTimestamp; ts <= maxTimestamp; ts++ {
		hline := plotter.NewFunction(func(x float64) float64 { return float64(ts) + 0.5 })
		hline.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
		hline.Color = color.Black
		p.Add(hline)
	}

	// add vertical dashed lines for each process
	for pid := range history {
		vline, err := plotter.NewLine(plotter.XYs{
			{X: float64(pid), Y: float64(sortedTimestamps[0]) - 0.3},
			{X: float64(pid), Y: float64(sortedTimestamps[len(sortedTimestamps)-1]) + 1},
		})
		if err != nil {
			return fmt.Errorf("could not create vertical line: %v", err)
		}
		vline.Color = color.Black
		p.Add(vline)
	}

	// draw connections between data nodes
	arrowColors := []color.Color{
		color.RGBA{R: 255, G: 0, B: 0, A: 255},
		color.RGBA{G: 255, B: 0, A: 255},
		color.RGBA{B: 255, A: 255},
	}
	colorIndex := 0
	for _, events := range history {
		for _, event := range events {
			if event.t == sent {
				for _, recvEvent := range history[event.dst] {
					if recvEvent.t == recv && recvEvent.src == event.src &&
						recvEvent.timestamp > event.timestamp {
						color := arrowColors[colorIndex%len(arrowColors)]
						addLine(
							p,
							float64(event.src),
							float64(event.timestamp),
							float64(event.dst),
							float64(recvEvent.timestamp),
							color,
						)
						colorIndex++
						break
					}
				}
			}
		}
	}

	// add vertical labels for processes
	for pid := range history {
		vertLabel := Labels{
			XYs: plotter.XYs{{
				X: float64(pid),
				Y: float64(sortedTimestamps[len(sortedTimestamps)-1]) + 1,
			}},
			Labels: []string{fmt.Sprintf("Process %d", pid)},
		}
		vertLabelPlotter, err := plotter.NewLabels(vertLabel)
		if err != nil {
			return fmt.Errorf("could not create vertical labels plotter: %v", err)
		}
		vertLabelPlotter.TextStyle[0].Rotation = math.Pi / 2
		p.Add(vertLabelPlotter)
	}

	// draw data nodes and labels for them
	var eventsData plotter.XYs
	var labels Labels
	for pid, events := range history {
		for _, event := range events {
			xy := plotter.XY{
				X: float64(pid),
				Y: float64(event.timestamp),
			}
			eventsData = append(eventsData, xy)
			labels.XYs = append(labels.XYs, xy)
			labels.Labels = append(labels.Labels, fmt.Sprintf("%d", event.timestamp))
		}
	}

	scatter, err := plotter.NewScatter(eventsData)
	if err != nil {
		return fmt.Errorf("could not create scatter plot: %v", err)
	}

	// customize scatter plot
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	scatter.GlyphStyle.Color = color.Black

	p.Add(scatter)

	// draw labels
	labelPlotter, err := plotter.NewLabels(labels)
	if err != nil {
		return fmt.Errorf("could not create labels plotter: %v", err)
	}
	p.Add(labelPlotter)

	// remove axes and legend
	p.HideAxes()

	if err := p.Save(6*vg.Inch, 10*vg.Inch, filename); err != nil {
		return fmt.Errorf("could not save plot: %v", err)
	}
	return nil
}

func addLine(p *plot.Plot, x1, y1, x2, y2 float64, color color.Color) error {
	line, err := plotter.NewLine(plotter.XYs{{X: x1, Y: y1}, {X: x2, Y: y2}})
	if err != nil {
		return fmt.Errorf("could not create line: %v", err)
	}
	line.Color = color
	line.Width = vg.Points(2)

	p.Add(line)
	return nil
}
