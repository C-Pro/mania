package dialogflow

// GenerateResponse .
func GenerateResponse(expectUserResponse bool, msg string) Response {
	return Response{
		Payload: ResponsePayload{
			Google: ResponseGoogle{
				ExpectUserResponse: expectUserResponse,
				RichResponse: RichResponse{
					Items: []Item{
						{
							SimpleResponse: SimpleResponse{
								TextToSpeech: msg,
							},
						},
					},
				},
			},
		},
	}
}
