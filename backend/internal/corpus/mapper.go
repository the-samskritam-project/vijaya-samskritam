package corpus

import (
	"fmt"
)

// FieldMapper maps corpus-specific fields to standard format
type FieldMapper struct {
	DevanagariField  string
	TranslationField string
}

// MapDocument maps a corpus document to a standard SearchResult
func (fm *FieldMapper) MapDocument(corpusName string, doc map[string]interface{}, score float64) (*SearchResult, error) {
	result := &SearchResult{
		CorpusName: corpusName,
		Score:      score,
	}

	// Extract Devanagari text
	if devanagari, ok := doc[fm.DevanagariField].(string); ok {
		result.DevanagariText = devanagari
	}

	// Extract translation
	if translation, ok := doc[fm.TranslationField].(string); ok {
		result.FullTranslation = translation
	}

	// Extract document ID
	if id, ok := doc["_id"].(string); ok {
		result.DocumentID = id
	} else if id := doc["_id"]; id != nil {
		result.DocumentID = fmt.Sprintf("%v", id)
	}

	// Extract chapter number
	if chapterNum, ok := doc["chapter_number"].(int32); ok {
		ch := int(chapterNum)
		result.ChapterNumber = &ch
	} else if chapterNum, ok := doc["chapter_number"].(int); ok {
		result.ChapterNumber = &chapterNum
	} else if chapterNum, ok := doc["chapter_number"].(float64); ok {
		ch := int(chapterNum)
		result.ChapterNumber = &ch
	}

	// Extract verse number
	if verseNum, ok := doc["verse_number"].(string); ok {
		result.VerseNumber = &verseNum
	}

	// Extract prose number
	if proseNum, ok := doc["prose_number"].(string); ok {
		result.ProseNumber = &proseNum
	}

	return result, nil
}
