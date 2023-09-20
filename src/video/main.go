package video

import (
	"fmt"
	"os/exec"

	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

func CreateVideo(videopath string, audioFile string, outputVideo string) {
	utils.Log("Criando vídeo...")

	cmd := exec.Command(
		"ffmpeg",
		"-i", videopath,
		"-i", audioFile,
		"-vf", "scale=1080:1920",
		"-c:v", "libx264",
		"-profile:v", "high", // Perfis de codificação de vídeo (ajuste conforme necessário)
		"-b:v", "2M", // Taxa de bits de vídeo (ajuste conforme necessário)
		"-r", "30", // Taxa de quadros (ajuste conforme necessário)
		"-c:a", "aac",
		"-b:a", "192k",
		"-shortest",
		"-map", "0:v:0",
		"-map", "1:a:0",
		outputVideo,
	)

	err := cmd.Run()

	if err != nil {
		utils.Fatal("Erro ao criar vídeo com FFMPEG", err)
	}

	utils.Log(fmt.Sprintf("Vídeo %s criado!", outputVideo))
}
