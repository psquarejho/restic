package restore

import (
	"time"

	"github.com/restic/restic/internal/ui"
)

type jsonPrinter struct {
	terminal term
}

func NewJSONProgress(terminal term) ProgressPrinter {
	return &jsonPrinter{
		terminal: terminal,
	}
}

func (t *jsonPrinter) print(status interface{}) {
	t.terminal.Print(ui.ToJSONString(status))
}

func (t *jsonPrinter) Update(filesFinished, filesTotal, allBytesWritten, allBytesTotal uint64, duration time.Duration) {
	status := statusUpdate{
		MessageType:    "status",
		SecondsElapsed: uint64(duration / time.Second),
		TotalFiles:     filesTotal,
		FilesDone:      filesFinished,
		TotalBytes:     allBytesTotal,
		BytesDone:      allBytesWritten,
	}

	if allBytesTotal > 0 {
		status.PercentDone = float64(allBytesWritten) / float64(allBytesTotal)
	}

	t.print(status)
}

func (t *jsonPrinter) Finish(filesFinished, filesTotal, allBytesWritten, allBytesTotal uint64, duration time.Duration) {
	status := summaryOutput{
		MessageType:    "summary",
		SecondsElapsed: uint64(duration / time.Second),
		TotalFiles:     filesTotal,
		FilesDone:      filesFinished,
		TotalBytes:     allBytesTotal,
		BytesDone:      allBytesWritten,
	}
	t.print(status)
}

type statusUpdate struct {
	MessageType    string  `json:"message_type"` // "status"
	SecondsElapsed uint64  `json:"seconds_elapsed,omitempty"`
	PercentDone    float64 `json:"percent_done"`
	TotalFiles     uint64  `json:"total_files,omitempty"`
	FilesDone      uint64  `json:"files_done,omitempty"`
	TotalBytes     uint64  `json:"total_bytes,omitempty"`
	BytesDone      uint64  `json:"bytes_done,omitempty"`
}

type summaryOutput struct {
	MessageType    string `json:"message_type"` // "summary"
	SecondsElapsed uint64 `json:"seconds_elapsed,omitempty"`
	TotalFiles     uint64 `json:"total_files,omitempty"`
	FilesDone      uint64 `json:"files_done,omitempty"`
	TotalBytes     uint64 `json:"total_bytes,omitempty"`
	BytesDone      uint64 `json:"bytes_done,omitempty"`
}
