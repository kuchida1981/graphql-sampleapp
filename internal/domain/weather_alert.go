package domain

type WeatherAlert struct {
	ID              string                 `firestore:"id"`
	Title           string                 `firestore:"title"`
	Description     string                 `firestore:"description"`
	RawData         map[string]interface{} `firestore:"rawData"`
	AffectedAreas   []string               `firestore:"affectedAreas"`
	Recommendations []string               `firestore:"recommendations"`
}
