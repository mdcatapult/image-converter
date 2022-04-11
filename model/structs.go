package model

type ConvertRequest struct {
	InputFile     string `json:"input-file" binding:"required"`
	InputMaskFile string `json:"input-mask-file" binding:"required"`
	OutputFile    string `json:"output-file" binding:"required"`
}

type ConvertRequestForFijiMacro struct {
	InputFile     string
	InputFilename string
	InputMaskFile string
	InputMaskFilename string
	OutputFile    string
}
