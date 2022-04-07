package model

type ConvertRequest struct {
	InputFile  string `json:"input-file" binding:"required"`
	OutputFile string `json:"output-file" binding:"required"`
}
