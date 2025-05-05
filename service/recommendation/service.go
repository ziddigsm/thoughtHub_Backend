package recommendation

import (
	"encoding/base64"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"gorm.io/gorm"
)

type RecommendationService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *RecommendationService {
	return &RecommendationService{
		db: db,
	}
}

type RecommendationRequest struct {
	Text   string `json:"text"`
	UserID int    `json:"user_id"`
}

type BlogRecommendation struct {
	Title           string  `json:"title"`
	SimilarityScore float64 `json:"similarity_score"`
}

type RecommendationResponse struct {
	Recommendations []BlogRecommendation `json:"recommendations"`
}

func (s *RecommendationService) GetSimilarBlogs(req RecommendationRequest) (RecommendationResponse, error) {
	var blogs []types.Blogs
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)
	
	err := s.db.Where("created_on >= ? AND user_id != ? AND is_active = ?", 
		threeMonthsAgo, req.UserID, true).
		Limit(90).
		Find(&blogs).Error
	
	if err != nil {
		return RecommendationResponse{}, err
	}

	if len(blogs) == 0 {
		return RecommendationResponse{
			Recommendations: []BlogRecommendation{},
		}, nil
	}

	inputProcessed := preprocessText(req.Text)

	blogContents := make([]string, len(blogs))
	for i, blog := range blogs {
		blogContents[i] = preprocessText(blog.Content)
	}

	allTexts := append(blogContents, inputProcessed)
	tfidfMatrix := calculateTFIDF(allTexts)
	
	inputVector := tfidfMatrix[len(tfidfMatrix)-1]
	
	similarities := make([]float64, len(blogs))
	for i := 0; i < len(blogs); i++ {
		similarities[i] = cosineSimilarity(inputVector, tfidfMatrix[i])
	}

	blogIndices := make([]int, len(blogs))
	for i := range blogIndices {
		blogIndices[i] = i
	}
	
	sort.Slice(blogIndices, func(i, j int) bool {
		return similarities[blogIndices[i]] > similarities[blogIndices[j]]
	})

	numRecommendations := 3
	if len(blogIndices) < numRecommendations {
		numRecommendations = len(blogIndices)
	}

	recommendations := make([]BlogRecommendation, numRecommendations)
	for i := 0; i < numRecommendations; i++ {
		idx := blogIndices[i]
		recommendations[i] = BlogRecommendation{
			Title:           blogs[idx].Title,
			SimilarityScore: similarities[idx],
		}
	}

	return RecommendationResponse{
		Recommendations: recommendations,
	}, nil
}

func preprocessText(text string) string {
	text = strings.ToLower(text)
	
	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, "")
	
	words := strings.Fields(text)
	
	stopwords := getStopwords()
	filtered := make([]string, 0, len(words))
	
	for _, word := range words {
		if _, isStopword := stopwords[word]; !isStopword && len(word) > 2 {
			filtered = append(filtered, word)
		}
	}
	
	return strings.Join(filtered, " ")
}

func getStopwords() map[string]struct{} {
	stopwordsList := []string{
		"a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "as", "at",
		"be", "because", "been", "before", "being", "below", "between", "both", "but", "by",
		"can", "did", "do", "does", "doing", "don", "down", "during",
		"each", "few", "for", "from", "further",
		"had", "has", "have", "having", "he", "her", "here", "hers", "herself", "him", "himself", "his", "how",
		"i", "if", "in", "into", "is", "it", "its", "itself",
		"just",
		"me", "more", "most", "my", "myself",
		"no", "nor", "not", "now",
		"of", "off", "on", "once", "only", "or", "other", "our", "ours", "ourselves", "out", "over", "own",
		"same", "she", "should", "so", "some", "such",
		"than", "that", "the", "their", "theirs", "them", "themselves", "then", "there", "these", "they", "this", "those", "through", "to", "too",
		"under", "until", "up",
		"very",
		"was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would",
		"you", "your", "yours", "yourself", "yourselves",
	}
	
	stopwords := make(map[string]struct{}, len(stopwordsList))
	for _, word := range stopwordsList {
		stopwords[word] = struct{}{}
	}
	
	return stopwords
}

func calculateTFIDF(texts []string) []map[string]float64 {
	docFreq := make(map[string]int)
	allTokens := make(map[string]struct{})
	
	tokenizedTexts := make([][]string, len(texts))
	
	for i, text := range texts {
		tokens := strings.Fields(text)
		tokenizedTexts[i] = tokens
		
		seen := make(map[string]bool)
		for _, token := range tokens {
			allTokens[token] = struct{}{}
			if !seen[token] {
				docFreq[token]++
				seen[token] = true
			}
		}
	}
	
	tfidfVectors := make([]map[string]float64, len(texts))
	
	for i, tokens := range tokenizedTexts {
		termFreq := make(map[string]int)
		for _, token := range tokens {
			termFreq[token]++
		}
		
		tfidf := make(map[string]float64)
		for token, tf := range termFreq {
			idf := math.Log(float64(len(texts)) / float64(docFreq[token]))
			tfidf[token] = float64(tf) * idf
		}
		
		tfidfVectors[i] = tfidf
	}
	
	return tfidfVectors
}

func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0
	
	for token, val1 := range vec1 {
		if val2, exists := vec2[token]; exists {
			dotProduct += val1 * val2
		}
		norm1 += val1 * val1
	}
	
	for _, val2 := range vec2 {
		norm2 += val2 * val2
	}
	
	if norm1 == 0 || norm2 == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func encodeImageToBase64(img []byte) string {
	if img == nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(img)
}