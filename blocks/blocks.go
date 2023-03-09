package blocks

import (
	"KYVENetwork/kyve-tm-bsync/types"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"github.com/tendermint/tendermint/libs/json"
	"io"
	"net/http"
)

func RetrieveBundle(storageId string) ([]types.DataItem, error) {
	raw, err := downloadFromUrl(fmt.Sprintf("https://arweave.net/%s", storageId))
	if err != nil {
		return nil, err
	}

	fmt.Println(len(raw))
	fmt.Println(createChecksum(raw))

	deflated, err := decompressGzip(raw)
	if err != nil {
		return nil, err
	}

	var bundle []types.DataItem

	if err := json.Unmarshal(deflated[:], &bundle); err != nil {
		return nil, err
	}

	return bundle, nil
}

func downloadFromUrl(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func createChecksum(input []byte) (hash string) {
	h := sha256.New()
	h.Write(input)
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

func decompressGzip(input []byte) ([]byte, error) {
	var out bytes.Buffer
	r, err := gzip.NewReader(bytes.NewBuffer(input))
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(&out, r); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
