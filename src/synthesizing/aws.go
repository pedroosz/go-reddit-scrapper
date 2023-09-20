package synthesizing

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/pedroosz/go-reddit-scrapper/src/utils"
)

func AWSSynthesizing(filepath string, voice string) {

	text := utils.ReadEntireFile(filepath)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String("us-east-1"),
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("ACCESS_ID_AWS"),
				os.Getenv("SECRET_ID_AWS"),
				"",
			),
		},
	}))

	svc := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(text),
		VoiceId:      aws.String(voice),
		Engine:       aws.String("neural"),
	}

	output, err := svc.SynthesizeSpeech(input)

	if err != nil {
		utils.Fatal("Erro ao sintetizar o texto", err)
	}

	outFile, err := os.Create("teste.mp3")

	if err != nil {
		utils.Fatal("Erro ao tentar criar teste.mp3", err)
	}

	defer outFile.Close()

	_, errAws := io.Copy(outFile, output.AudioStream)

	if errAws != nil {
		utils.Fatal("Erro ao tentar baixar arquivo", err)
	}

}
