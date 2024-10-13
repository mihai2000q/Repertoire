package song

type CreateSongRequest struct {
	Title      string `validate:"required,max=100"`
	IsRecorded *bool
}
