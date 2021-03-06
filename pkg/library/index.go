package library

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ambientsound/visp/list"
	"github.com/ambientsound/visp/log"
	`github.com/ambientsound/visp/options`
	"github.com/ambientsound/visp/xdg"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/analysis/token/edgengram"
	"github.com/blevesearch/bleve/v2/analysis/token/lowercase"
	"github.com/blevesearch/bleve/v2/analysis/tokenizer/unicode"
	"github.com/blevesearch/bleve/v2/mapping"
	"github.com/blevesearch/bleve/v2/search"
)

type Index interface {
	Add(list list.List) error
	Query(q string) (list.List, error)
	QueryID(id string) (list.Row, error)
	Close() error
}

type index struct {
	bleve bleve.Index
}

const indexAnalyzerName = "index_analyzer"
const queryAnalyzerName = "query_analyzer"
const edgeNgramTokenFilterName = "edge_ngram_filter"

func indexMapping() mapping.IndexMapping {

	indexmap := bleve.NewIndexMapping()
	// indexmap.DefaultAnalyzer = indexAnalyzerName

	docmap := bleve.NewDocumentMapping()
	docmap.DefaultAnalyzer = indexAnalyzerName

	indexmap.DefaultMapping = docmap

	err := indexmap.AddCustomTokenFilter(edgeNgramTokenFilterName,
		map[string]interface{}{
			"type": edgengram.Name,
			"min":  2.0,
			"max":  15.0,
		})

	if err != nil {
		panic(err)
	}

	err = indexmap.AddCustomAnalyzer(indexAnalyzerName,
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
				edgeNgramTokenFilterName,
			},
		})

	if err != nil {
		panic(err)
	}

	err = indexmap.AddCustomAnalyzer(queryAnalyzerName,
		map[string]interface{}{
			"type":      custom.Name,
			"tokenizer": unicode.Name,
			"token_filters": []string{
				en.PossessiveName,
				lowercase.Name,
			},
		})

	if err != nil {
		panic(err)
	}

	return indexmap
}

// New opens and returns a Bleve index.
// A persistent filesystem-backed index is tried first.
// If this index cannot be opened, we use an in-memory index instead.
func New() (Index, error) {
	path := filepath.Join(xdg.CacheDirectory(), "visp", "library.idx")

	idx, err := bleve.Open(path)
	if err != nil {
		log.Debugf("Failed to open index at %s: %s", path, err)
		idx, err = bleve.New(path, indexMapping())
		if err != nil {
			log.Errorf("Failed to create new index at %s: %s", path, err)
			return NewInMemory()
		}
	}

	return &index{
		bleve: idx,
	}, nil
}

func NewInMemory() (Index, error) {
	idx, err := bleve.NewMemOnly(indexMapping())
	if err != nil {
		return nil, fmt.Errorf("failed to create in-memory index: %w", err)
	}
	return &index{
		bleve: idx,
	}, nil
}

func (idx *index) Close() error {
	return idx.bleve.Close()
}

func (idx *index) Add(dataset list.List) error {
	b := idx.bleve.NewBatch()
	for _, row := range dataset.All() {
		// Only tracks are indexed at the moment
		if row.Kind() != list.DataTypeTrack {
			continue
		}

		id := row.ID()
		data := row.Fields()

		err := b.Index(id, data)
		if err != nil {
			return fmt.Errorf("index '%s' (%+v): %w", id, data, err)
		}

		serialized, err := json.Marshal(row.Fields())
		if err != nil {
			return fmt.Errorf("serialize '%s' (%+v): %w", id, data, err)
		}

		// Store data inside index for later retrieval
		b.SetInternal([]byte(id), serialized)
	}

	return idx.bleve.Batch(b)
}

func (idx *index) Query(q string) (list.List, error) {
	const limit = 500

	query := bleve.NewMatchQuery(q)
	query.Analyzer = queryAnalyzerName
	req := bleve.NewSearchRequestOptions(query, limit, 0, false)

	res, err := idx.bleve.Search(req)
	if err != nil {
		return nil, fmt.Errorf("index query: %w", err)
	}

	return idx.hitsAsList(res)
}

func (idx *index) QueryID(id string) (list.Row, error) {
	fields := make(map[string]string)
	document, err := idx.bleve.GetInternal([]byte(id))
	if err != nil {
		return nil, fmt.Errorf("document(%s) not found: %w", id, err)
	}

	err = json.Unmarshal(document, &fields)
	if err != nil {
		return nil, fmt.Errorf("unmarshal(%s, %s): %w", id, string(document), err)
	}

	return list.NewRow(
		id,
		list.DataTypeTrack,
		fields,
	), nil
}

func (idx *index) hitsAsList(res *bleve.SearchResult) (list.List, error) {
	result := list.New()
	result.SetVisibleColumns(options.GetList(options.ColumnsTracklists))

	for _, hit := range res.Hits {
		row, err := idx.hitAsRow(hit)
		if err != nil {
			return nil, fmt.Errorf("document '%s': %w", hit.ID, err)
		}
		result.Add(row)
	}

	return result, nil
}

func (idx *index) hitAsRow(hit *search.DocumentMatch) (list.Row, error) {
	row, err := idx.QueryID(hit.ID)
	if err != nil {
		return nil, err
	}

	score := fmt.Sprintf("%3.1f%%", hit.Score*100)
	row.Set("score", score)

	return row, nil
}
