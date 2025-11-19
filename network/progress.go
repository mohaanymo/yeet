package network

import (
    "fmt"
    "strings"
    "time"
)

const (
    colorReset = "\033[0m"
    colorCyan  = "\033[36m"
    colorGray  = "\033[90m"
)

type ProgressBar struct {
    Total     int64
    Current   int64
    Width     int
    StartTime time.Time
}

func NewProgressBar(total int64) *ProgressBar {
    return &ProgressBar{
        Total:     total,
        Width:     40,
        StartTime: time.Now(),
    }
}

func (pb *ProgressBar) Update(current int64) {
    pb.Current = current
    pb.render()
}

func (pb *ProgressBar) render() {
    percent := float64(pb.Current) / float64(pb.Total)
    filled := int(percent * float64(pb.Width))
    
    // Create progress bar
    bar := colorCyan + strings.Repeat("━", filled) + colorReset +
           colorGray + strings.Repeat("━", pb.Width-filled) + colorReset
    
    // Format file sizes
    currentMB := float64(pb.Current) / (1024 * 1024)
    totalMB := float64(pb.Total) / (1024 * 1024)
    
    // Calculate speed and ETA
    elapsed := time.Since(pb.StartTime).Seconds()
    var speed float64
    var eta string
    
    if elapsed > 0 && pb.Current > 0 {
        speed = float64(pb.Current) / elapsed / (1024 * 1024) // MB/s
        
        if pb.Current < pb.Total {
            remaining := float64(pb.Total-pb.Current) / (speed * 1024 * 1024)
            eta = formatDuration(time.Duration(remaining) * time.Second)
        } else {
            eta = "00:00"
        }
    } else {
        eta = "--:--"
    }
    
    // Print progress bar
    fmt.Printf("\r[ %s ] %.1f%%  %.2f/%.2f MB  %.2f MB/s  ETA: %s",
        bar, percent*100, currentMB, totalMB, speed, eta)
}

func (pb *ProgressBar) Finish() {
    pb.Update(pb.Total)
    fmt.Println() // New line after completion
}

func formatDuration(d time.Duration) string {
    d = d.Round(time.Second)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    d -= m * time.Minute
    s := d / time.Second
    
    if h > 0 {
        return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
    }
    return fmt.Sprintf("%02d:%02d", m, s)
}
