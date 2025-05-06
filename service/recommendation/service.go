package recommendation

import (
	"fmt"
	"html"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ziddigsm/thoughtHub_Backend/types"
	"github.com/ziddigsm/thoughtHub_Backend/utils"
)

var stopWordsList = getStopwords()

func (h *Handler) ValidateRequestBody(w http.ResponseWriter, r *http.Request) {
	var req types.RecommendationRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("invalid request body: %v", err))
		return
	}

	if req.Text == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("no text provided"))
		return
	}

	if req.BlogID <= 0 {
		utils.ErrorResponse(w, http.StatusBadRequest, fmt.Errorf("user ID is required"))
		return
	}

	response, err := h.GetSimilarBlogs(req)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("failed to get recommendations: %v", err))
		return
	}

	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *Handler) GetSimilarBlogs(req types.RecommendationRequest) (types.RecommendationResponse, error) {
	var blogs []types.BlogWithName

	err := h.getBlogsFromDB(req, &blogs)

	if err != nil {
		return types.RecommendationResponse{}, err
	}

	if len(blogs) == 0 {
		return types.RecommendationResponse{
			Recommendations: []types.BlogRecommendation{},
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

	threshold := 0.3
	selectedBlogs := []types.BlogWithName{}
	for _, idx := range blogIndices {
		if similarities[idx] < threshold {
			break
		}
		selectedBlogs = append(selectedBlogs, blogs[idx])
		if len(selectedBlogs) == numRecommendations {
			break
		}
	}

	var detailedBlogs []types.DetailedBlog
	h.blogHandler.GetLikesAndComments(selectedBlogs, &detailedBlogs, nil)

	recommendations := []types.BlogRecommendation{}
	for i, detailedBlog := range detailedBlogs {
		recommendations = append(recommendations, types.BlogRecommendation{
			BlogData:        detailedBlog.BlogData,
			Likes:           detailedBlog.Likes,
			Comments:        detailedBlog.Comments,
			SimilarityScore: similarities[blogIndices[i]],
		})
	}

	return types.RecommendationResponse{
		Recommendations: recommendations,
	}, nil
}

func (h *Handler) getBlogsFromDB(req types.RecommendationRequest, blogs *[]types.BlogWithName) error {
	threeMonthsAgo := time.Now().AddDate(0, -3, 0)

	err := h.db.Table("blogs").Select("blogs.*, Users.name").Joins("LEFT JOIN users on blogs.user_id = users.id").Where("blogs.created_on >= ? AND blogs.id != ? AND blogs.is_active = ?",
		threeMonthsAgo, req.BlogID, true).
		Limit(90).
		Find(&blogs).Error

	if err != nil {
		return err
	}
	return nil
}

func preprocessText(text string) string {

	text = html.UnescapeString(text)
	htmlTagRe := regexp.MustCompile(`<[^>]*>`)
	text = htmlTagRe.ReplaceAllString(text, "")

	text = strings.ToLower(text)

	re := regexp.MustCompile(`[^\w\s]`)
	text = re.ReplaceAllString(text, " ")

	tokenizer := regexp.MustCompile(`\b\w+\b`)
	words := tokenizer.FindAllString(text, -1)

	filtered := make([]string, 0, len(words))

	for _, word := range words {
		if _, isStopword := stopWordsList[word]; !isStopword && len(word) > 2 {
			filtered = append(filtered, word)
		}
	}

	return strings.Join(filtered, " ")
}
func getStopwords() map[string]struct{} {
	stopwordsList := []string{
		"a", "about", "above", "after", "again", "against", "all", "am", "an", "and", "any", "are", "as", "at", "almost", "also", "afterwards", "across", "actually", "among", "amongst", "along", "already", "although", "always", "another", "anyway", "anywhere", "around", "aside", "away",
		"be", "because", "been", "before", "being", "below", "between", "both", "but", "by", "bottom", "besides", "besides", "beyond",
		"can", "cannot", "could", "couldn't", "call", "calls", "come", "comes", "coming", "could", "co", "cry",
		"did", "do", "does", "doing", "don", "down", "during", "define", "despite", "didn't", "doesn't", "don't", "detail", "details", "due",
		"each", "either", "else", "ever", "every", "eg", "especially", "except", "example", "ex", "exceedingly", "exclusively", "exemplary", "extremely", "elsewhere", "empty", "end", "enough", "especially", "even", "ever", "every", "everybody", "everyone", "everything", "everywhere",
		"few", "for", "from", "further", "four", "forom", "forward", "found", "free", "frequently", "former", "formerly", "five", "first", "firstly", "fill", "following", "furthermore", "find", "found", "founder", "founding", "founded",
		"had", "has", "have", "having", "he", "her", "here", "hers", "herself", "him", "himself", "his", "how",
		"i", "if", "in", "into", "is", "it", "its", "itself",
		"just",
		"me", "more", "most", "my", "myself",
		"no", "nor", "not", "now",
		"of", "off", "on", "once", "only", "or", "other", "our", "ours", "ourselves", "out", "over", "own",
		"per", "perhaps", "please", "put", "previous", "probably",
		"rather", "re", "really", "regarding", "right",
		"same", "she", "should", "so", "some", "such", "someone", "something", "somewhere", "since", "six", "sincere", "similar", "similarly", "several", "sincerely", "similarly", "serious", "seem", "seeming", "seemingly", "seemed", "seems", "suppose", "supposed", "supposing", "supposes", "supposedly",
		"than", "that", "the", "their", "theirs", "them", "themselves", "then", "there", "these", "they", "this", "those", "through", "to", "too", "towards", "thus", "than", "twelve", "twenty", "therefore", "though", "thereafter", "therein", "thereupon",
		"under", "until", "up", "upon", "us",
		"very", "via",
		"was", "we", "were", "what", "when", "where", "which", "while", "who", "whom", "why", "will", "with", "would", "whose", "whoever", "whichever", "whenever", "whosoever", "whether",
		"you", "your", "yours", "yourself", "yourselves", "your's", "yet",
		"kindly", "please", "thank", "thanks", "appreciate", "appreciation",
		"dear", "hello", "hi", "greetings", "regards", "sincerely", "best",
		"also", "but", "however", "therefore", "meanwhile", "furthermore", "moreover",
		"although", "despite", "in spite of", "even though", "while",
		"among", "besides", "except", "including", "like", "near", "outside", "past", "throughout",
		"amongst", "back", "beyond",
		"become", "becomes", "becoming", "became", "become", "becoming",
		"de", "detail", "etc", "etcetera", "i.e.", "e.g.",
		"front", "full", "further", "furthermore", "hence", "hereafter", "herein", "hereupon", "howbeit",
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
			df := docFreq[token]
			if df == 1 || float64(df)/float64(len(texts)) > 0.85 {
				continue
			}
			idf := math.Log(1 + float64(len(texts))/(1+float64(df)))
			tfidf[token] = float64(tf) * idf
		}

		tfidfVectors[i] = normalizeVector(tfidf)
	}

	return tfidfVectors
}

func normalizeVector(vec map[string]float64) map[string]float64 {
	sumSquares := 0.0
	for _, val := range vec {
		sumSquares += val * val
	}
	norm := math.Sqrt(sumSquares)
	if norm == 0 {
		return vec
	}
	for k, v := range vec {
		vec[k] = v / norm
	}
	return vec
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
