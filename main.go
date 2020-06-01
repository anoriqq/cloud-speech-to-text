package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	speech "cloud.google.com/go/speech/apiv1p1beta1"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1p1beta1"
)

func main() {
	args := os.Args[1:]

	file, err := os.Create(`./output.txt`)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
		return
	}
	defer file.Close()

	ctx := context.Background()
	client, err := speech.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return
	}

	gcsURI := args[0]
	sendGCS(file, client, gcsURI)
}

func sendGCS(w io.Writer, client *speech.Client, gcsURI string) error {
	ctx := context.Background()

	// Send the contents of the audio file with the encoding and
	// and sample rate information to be transcripted.
	req := &speechpb.LongRunningRecognizeRequest{
		Config: &speechpb.RecognitionConfig{
			Encoding:          speechpb.RecognitionConfig_FLAC,
			SampleRateHertz:   44100,
			AudioChannelCount: 2,
			LanguageCode:      "ja_JP",
		},
		Audio: &speechpb.RecognitionAudio{
			AudioSource: &speechpb.RecognitionAudio_Uri{Uri: gcsURI},
		},
	}

	op, err := client.LongRunningRecognize(ctx, req)
	if err != nil {
		return err
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		return err
	}

	// Print the results.
	for _, result := range resp.Results {
		for _, alt := range result.Alternatives {
			fmt.Fprintf(w, "\"%v\" (confidence=%3f)\n", alt.Transcript, alt.Confidence)
		}
	}

	return nil
}
