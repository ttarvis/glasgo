package main

import(
	"math"
	"regexp"
)

// computes Shannon Entropy, H. Note this is
// computed with a slice of bytes, not a string
// H = -Sum(i in {1,n}) ((counti/N) * log2(counti/N))
// counti is count of character i
// n is number of different characters in a string s
// N is length of string s
func H(s []byte) (float64, float64) {
        var entropy, sum float64;
        // may want to map runes to float at some point
        count := make(map[byte]float64);

        for _, b := range s {
                count[b]++;
        }

        N := float64(len(s));

        for _, ni := range count {
                sum += ni * math.Log2(ni);
        }
        entropy = math.Log2(N) - sum/N;
	Hn := entropy / math.Log2(float64(len(count)));

        return entropy, Hn;
}

// normalised specific entropy
func H2(s []byte) float64 {
	entropy, ni := H(s);

	return entropy / math.Log2(ni);
}

func isHex(s string) bool {
        matched, err := regexp.MatchString("^(0x|0X)?[a-fA-F0-9]+$", s);
        if err != nil {
                warnf("In isHex, %v", err);
                return false;
        }
        if matched {
                return true
        }
        return false;
}

func isBase64(s string) bool {
        isMatch, err := regexp.MatchString("^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$", s);
        if err != nil {
                warnf("In isBase64, %v", err);
                return false;
        }
        if isMatch {
                return true;
        }
        return false;
}

func isAlphanum(s string) (bool, error) {
	b := []byte(s);
	isMatch, err := regexp.Match("[a-zA-Z0-9]$", b);
	if err != nil {
		return false, err;
	}
	if isMatch {
		return true, nil;
	}
	return false, nil;
}	
