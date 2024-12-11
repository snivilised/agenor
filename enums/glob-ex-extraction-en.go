package enums

//go:generate stringer -type=GlobExtraction -linecomment -trimprefix=GlobExtraction -output glob-ex-en-auto.go

type GlobExtraction uint

const (
	GlobExtractionUndefined GlobExtraction = iota // undefined

	// GlobExtractionStarDot denotes a pattern leading with "*.";
	// eg "*.jpg"
	//
	GlobExtractionStarDot // star-dot

	// GlobExtractionStarOnly denotes a pattern leading with "*", without dot;
	// eg "*jpg"
	//
	GlobExtractionStarOnly // star-only

	// GlobExtractionLeadingDotOnly denotes a pattern leading with ".", without star;
	// eg ".jpg"
	//
	GlobExtractionLeadingDotOnly // leading-dot-only

	// GlobExtractionBare denotes a pattern with neither a "*" or ".";
	// eg "jpg"
	//
	GlobExtractionBare // bare
)
