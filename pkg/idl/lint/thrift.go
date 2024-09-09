package lint

import "unicode"

type Scanner struct {
	src    []byte
	pos    int
	token  int
	offset int // 偏移量
	ch     rune
}

type Token string

const TokenStruct Token = "struct"
const TokenString Token = "string"
const TokenList Token = "list"
const TokenLeft Token = "<"
const TokenRight Token = ">"
const TokenEq Token = "="
const TokenService Token = "service"
const TokenCommentUnix = "#"

func (s *Scanner) next() rune {
	s.offset = s.offset + 1
	return rune(s.src[s.offset])
}

/*
*
# 框架配置

	struct BizFlexConfig {
	    1: string psm
	    2: string desc
	    3: list<string> owners
	    4: list<PluginInfo> plugins
	    5: list<SceneInfo> scenes
	}
*/
func (s *Scanner) Scan(input byte) (pos int, tok int, lit string) {
	s.trimSpace()
	switch s.ch {
	case '#':
		//return s.offset, TokenCommentUnix, "#"
	default:

	}
	return
}

func (s *Scanner) scanComment() string {
	offs := s.offset
	next := offs
	for {
		if s.src[next] == '\n' {
			break
		}
		next++
	}
	return string(s.src[offs:next])
}

func (s *Scanner) trimSpace() {
	offs := s.offset
	next := offs
	for {
		if unicode.IsSpace(rune(s.src[next])) {
			next++
			continue
		}
		break
	}
	if next > offs {
		s.offset = next - 1
	}
	return
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	next := offs
	for {
		if s.isIdentifier(s.src[next]) {
			next++
			continue
		}
		break
	}
	return string(s.src[offs:next])
}

func (s *Scanner) isIdentifier(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || b == '_' || '0' <= b && b <= '9'
}
