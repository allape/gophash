package gophash

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"math"
	"os"
	"testing"
)

const (
	LogoPNG       = "./testdata/pHash.png"
	LogoPNGSha256 = "f916d5b00bf7dff2331f96cbcaabe31e4e15d1d66d0b96660f6d6d3de51d0714"
	LogoPNGPHash  = "0000013f89c1e3920400fe233b9e303fe2e325471dd95119f2a552b2dc550ab4b16ad96d60acef67151c39c4ab2068f70ce7323dc0db012492f27278380cdbb2dac7edc4bdaf6cdb"
)

const (
	MaskedLogoPNG       = "./testdata/pHash.masked.png"
	MaskedLogoPNGSha256 = "aa285957c47a556335d76d08db187b2e3365fdab171708fd3cd3ce8bff52619e"
	MaskedLogoPNGPHash  = "0001c03f89c1e392042486c7eff9f9fe7ee32436c0000000000152b2b6c000000000016960b6c00000000000b62431ff00ee701c0f9400b252f272783c0c00b2dac7edc4bdaf6d24"
)

const DistanceBetweenLogoAndMaskedLogo = 0.293403

func validateTestPNG(t *testing.T, imageFile, imageFileSha256 string) {
	testPNGData, err := os.ReadFile(imageFile)
	if err != nil {
		t.Error(err)
	}

	sha256hasher := sha256.New()
	sha256hasher.Write(testPNGData)
	sha256hash := hex.EncodeToString(sha256hasher.Sum(nil))

	if sha256hash != imageFileSha256 {
		t.Errorf("Expected sha256 hash to be %s, got %s", imageFileSha256, sha256hash)
	}
}

func TestLogoPNG(t *testing.T) {
	validateTestPNG(t, LogoPNG, LogoPNGSha256)

	phash := New(LogoPNG, log.New(os.Stderr, "[phash]", log.LstdFlags))

	hexedHash := hex.EncodeToString(phash.Sum(nil))

	if hexedHash != LogoPNGPHash {
		dis, _ := DistanceBetweenHexString(LogoPNGPHash, hexedHash)
		t.Errorf("Expected hash to be %s, got %s, hash distance: %f", LogoPNGPHash, hexedHash, dis)
	}
}

func TestMaskedLogoPNG(t *testing.T) {
	validateTestPNG(t, MaskedLogoPNG, MaskedLogoPNGSha256)

	phash := New(MaskedLogoPNG, log.New(os.Stderr, "[phash.mask]", log.LstdFlags))

	hexedHash := hex.EncodeToString(phash.Sum(nil))

	if hexedHash != MaskedLogoPNGPHash {
		dis, _ := DistanceBetweenHexString(MaskedLogoPNGPHash, hexedHash)
		t.Errorf("Expected hash to be %s, got %s, hash distance: %f", MaskedLogoPNGPHash, hexedHash, dis)
	}
}

func TestDistance(t *testing.T) {
	dis, err := DistanceBetweenHexString(LogoPNGPHash, MaskedLogoPNGPHash)
	if err != nil {
		t.Error(err)
	}

	dis = math.Round(dis*1000000) / 1000000

	if dis != DistanceBetweenLogoAndMaskedLogo {
		t.Errorf("Expected distance to be %f, but got %f", DistanceBetweenLogoAndMaskedLogo, dis)
	}
}

func TestA(t *testing.T) {
	dis, err := DistanceBetweenHexString(
		"0001c03f89c1e392042486c7eff9f9fe7ee32436c0000000000152b2b6c000000000016960b6c00000000000b62431ff00ee701c0f9400b252f272783c0c00b2dac7ecc4bdaf6d24",
		"000002df89c1e212310d6e233b9e302fe2e376471dd95119f2a552b2dc550ab4b16ad96df0acef67151c39c4ab4461f70ce73238b900024b12f2727cb82404b2dac7adc4bdaf64db",
	)
	if err != nil {
		t.Error(err)
	}
	t.Log("distance: ", dis)
}
