package parse

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/knakk/rdf"
)

func parseAllTTL(s string) (r []rdf.Triple, err error) {
	dec := NewTTLDecoder(bytes.NewBufferString(s))
	for tr, err := dec.DecodeTriple(); err != io.EOF; tr, err = dec.DecodeTriple() {
		if err != nil {
			return r, err
		}
		r = append(r, tr)
	}
	return r, err
}

func TestTTL(t *testing.T) {
	for _, test := range ttlTestSuite[:38] {
		triples, err := parseAllTTL(test.input)
		if err != nil {
			if test.errWant == "" {
				t.Errorf("ParseTTL(%s) => %v, want %v", test.input, err, test.want)
				continue
			}
			if strings.HasSuffix(err.Error(), test.errWant) {
				continue
			}
			t.Errorf("ParseTTL(%s) => %q, want %q", test.input, err, test.errWant)
			continue
		}

		if !reflect.DeepEqual(triples, test.want) {
			t.Errorf("ParseTTL(%s) => %v, want %v", test.input, triples, test.want)
		}
	}
}

// ttlTestSuite is a representation of the official W3C test suite for Turtle
// which is found at: http://www.w3.org/2013/TurtleTests/
var ttlTestSuite = []struct {
	input   string
	errWant string
	want    []rdf.Triple
}{
	//# atomic tests
	//<#IRI_subject> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "IRI_subject" ;
	//   rdfs:comment "IRI subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <IRI_subject.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#IRI_with_four_digit_numeric_escape> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "IRI_with_four_digit_numeric_escape" ;
	//   rdfs:comment "IRI with four digit numeric escape (\\u)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <IRI_with_four_digit_numeric_escape.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`<http://a.example/\u0073> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#IRI_with_eight_digit_numeric_escape> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "IRI_with_eight_digit_numeric_escape" ;
	//   rdfs:comment "IRI with eight digit numeric escape (\\U)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <IRI_with_eight_digit_numeric_escape.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`<http://a.example/\U00000073> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#IRI_with_all_punctuation> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "IRI_with_all_punctuation" ;
	//   rdfs:comment "IRI with all punctuation" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <IRI_with_all_punctuation.ttl> ;
	//   mf:result    <IRI_with_all_punctuation.nt> ;
	//   .

	{`<scheme:!$%25&amp;'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "scheme:!$%25&amp;'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#bareword_a_predicate> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "bareword_a_predicate" ;
	//   rdfs:comment "bareword a predicate" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <bareword_a_predicate.ttl> ;
	//   mf:result    <bareword_a_predicate.nt> ;
	//   .

	{`<http://a.example/s> a <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#old_style_prefix> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "old_style_prefix" ;
	//   rdfs:comment "old-style prefix" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <old_style_prefix.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#SPARQL_style_prefix> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "SPARQL_style_prefix" ;
	//   rdfs:comment "SPARQL-style prefix" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <SPARQL_style_prefix.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`PREFIX p: <http://a.example/>
p:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefixed_IRI_predicate> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefixed_IRI_predicate" ;
	//   rdfs:comment "prefixed IRI predicate" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefixed_IRI_predicate.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
<http://a.example/s> p:p <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefixed_IRI_object> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefixed_IRI_object" ;
	//   rdfs:comment "prefixed IRI object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefixed_IRI_object.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
<http://a.example/s> <http://a.example/p> p:o .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefix_only_IRI> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefix_only_IRI" ;
	//   rdfs:comment "prefix-only IRI (p:)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefix_only_IRI.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix p: <http://a.example/s>.
p: <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefix_with_PN_CHARS_BASE_character_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefix_with_PN_CHARS_BASE_character_boundaries" ;
	//   rdfs:comment "prefix with PN CHARS BASE character boundaries (prefix: AZazÀÖØöø...:)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefix_with_PN_CHARS_BASE_character_boundaries.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿豈﷏ﷰ�𐀀󯿽: <http://a.example/> .
<http://a.example/s> <http://a.example/p> AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿豈﷏ﷰ�𐀀󯿽:o .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefix_with_non_leading_extras> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefix_with_non_leading_extras" ;
	//   rdfs:comment "prefix with_non_leading_extras (_:a·̀ͯ‿.⁀)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefix_with_non_leading_extras.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix a·̀ͯ‿.⁀: <http://a.example/>.
a·̀ͯ‿.⁀:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries" ;
	//   rdfs:comment "localName with assigned, NFC-normalized, basic-multilingual-plane PN CHARS BASE character boundaries (p:AZazÀÖØöø...)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries.ttl> ;
	//   mf:result    <localName_with_assigned_nfc_bmp_PN_CHARS_BASE_character_boundaries.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:AZazÀÖØöø˿Ͱͽ΄῾‌‍⁰↉Ⰰ⿕、ퟻ﨎ﷇﷰ￯ .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/AZazÀÖØöø˿Ͱͽ΄῾‌‍⁰↉Ⰰ⿕、ퟻ﨎ﷇﷰ￯"},
		},
	}},

	//<#localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries" ;
	//   rdfs:comment "localName with assigned, NFC-normalized PN CHARS BASE character boundaries (p:AZazÀÖØöø...)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries.ttl> ;
	//   mf:result    <localName_with_assigned_nfc_PN_CHARS_BASE_character_boundaries.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:AZazÀÖØöø˿Ͱͽ΄῾‌‍⁰↉Ⰰ⿕、ퟻ﨎ﷇﷰ￯𐀀󠇯 .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/AZazÀÖØöø˿Ͱͽ΄῾‌‍⁰↉Ⰰ⿕、ퟻ﨎ﷇﷰ￯𐀀󠇯"},
		},
	}},

	//<#localName_with_nfc_PN_CHARS_BASE_character_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_nfc_PN_CHARS_BASE_character_boundaries" ;
	//   rdfs:comment "localName with nfc-normalize PN CHARS BASE character boundaries (p:AZazÀÖØöø...)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_nfc_PN_CHARS_BASE_character_boundaries.ttl> ;
	//   mf:result    <localName_with_nfc_PN_CHARS_BASE_character_boundaries.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿﨎﷏ﷰ￯𐀀󯿽 .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿﨎﷏ﷰ￯𐀀󯿽"},
		},
	}},

	//<#default_namespace_IRI> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "default_namespace_IRI" ;
	//   rdfs:comment "default namespace IRI (:ln)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <default_namespace_IRI.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix : <http://a.example/>.
:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#prefix_reassigned_and_used> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefix_reassigned_and_used" ;
	//   rdfs:comment "prefix reassigned and used" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefix_reassigned_and_used.ttl> ;
	//   mf:result    <prefix_reassigned_and_used.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
@prefix p: <http://b.example/>.
p:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://b.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#reserved_escaped_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "reserved_escaped_localName" ;
	//   rdfs:comment "reserved-escaped local name" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <reserved_escaped_localName.ttl> ;
	//   mf:result    <reserved_escaped_localName.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:\_\~\.\-\!\$\&\'\(\)\*\+\,\;\=\/\?\#\@\%00 <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/_~.-!$&'()*+,;=/?#@%00`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#percent_escaped_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "percent_escaped_localName" ;
	//   rdfs:comment "percent-escaped local name" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <percent_escaped_localName.ttl> ;
	//   mf:result    <percent_escaped_localName.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:%25 <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/%25`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#HYPHEN_MINUS_in_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "HYPHEN_MINUS_in_localName" ;
	//   rdfs:comment "HYPHEN-MINUS in local name" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <HYPHEN_MINUS_in_localName.ttl> ;
	//   mf:result    <HYPHEN_MINUS_in_localName.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:s- <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s-`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#underscore_in_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "underscore_in_localName" ;
	//   rdfs:comment "underscore in local name" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <underscore_in_localName.ttl> ;
	//   mf:result    <underscore_in_localName.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:s_ <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s_`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#localname_with_COLON> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localname_with_COLON" ;
	//   rdfs:comment "localname with COLON" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localname_with_COLON.ttl> ;
	//   mf:result    <localname_with_COLON.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:s: <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s:`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#localName_with_leading_underscore> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_leading_underscore" ;
	//   rdfs:comment "localName with leading underscore (p:_)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_leading_underscore.ttl> ;
	//   mf:result    <localName_with_leading_underscore.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:_ <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/_`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#localName_with_leading_digit> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_leading_digit" ;
	//   rdfs:comment "localName with leading digit (p:_)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_leading_digit.ttl> ;
	//   mf:result    <localName_with_leading_digit.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:0 <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/0`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#localName_with_non_leading_extras> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "localName_with_non_leading_extras" ;
	//   rdfs:comment "localName with_non_leading_extras (_:a·̀ͯ‿.⁀)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <localName_with_non_leading_extras.ttl> ;
	//   mf:result    <localName_with_non_leading_extras.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
p:a·̀ͯ‿.⁀ <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/a·̀ͯ‿.⁀`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#old_style_base> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "old_style_base" ;
	//   rdfs:comment "old-style base" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <old_style_base.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@base <http://a.example/>.
<s> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#SPARQL_style_base> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "SPARQL_style_base" ;
	//   rdfs:comment "SPARQL-style base" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <SPARQL_style_base.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`BASE <http://a.example/>
<s> <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#labeled_blank_node_subject> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_subject" ;
	//   rdfs:comment "labeled blank node subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_subject.ttl> ;
	//   mf:result    <labeled_blank_node_subject.nt> ;
	//   .

	{`_:s <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.Blank{ID: "s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#labeled_blank_node_object> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_object" ;
	//   rdfs:comment "labeled blank node object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_object.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> _:o .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "o"},
		},
	}},

	//<#labeled_blank_node_with_PN_CHARS_BASE_character_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_with_PN_CHARS_BASE_character_boundaries" ;
	//   rdfs:comment "labeled blank node with PN_CHARS_BASE character boundaries (_:AZazÀÖØöø...)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_with_PN_CHARS_BASE_character_boundaries.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> _:AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿豈﷏ﷰ�𐀀󯿽 .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "AZazÀÖØöø˿ͰͽͿ῿‌‍⁰↏Ⰰ⿯、퟿豈﷏ﷰ�𐀀󯿽"},
		},
	}},

	//<#labeled_blank_node_with_leading_underscore> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_with_leading_underscore" ;
	//   rdfs:comment "labeled blank node with_leading_underscore (_:_)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_with_leading_underscore.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> _:_ .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "_"},
		},
	}},

	//<#labeled_blank_node_with_leading_digit> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_with_leading_digit" ;
	//   rdfs:comment "labeled blank node with_leading_digit (_:0)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_with_leading_digit.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> _:0 .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "0"},
		},
	}},

	//<#labeled_blank_node_with_non_leading_extras> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "labeled_blank_node_with_non_leading_extras" ;
	//   rdfs:comment "labeled blank node with_non_leading_extras (_:a·̀ͯ‿.⁀)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <labeled_blank_node_with_non_leading_extras.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> _:a·̀ͯ‿.⁀ .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: `http://a.example/s`},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "a·̀ͯ‿.⁀"},
		},
	}},

	//<#anonymous_blank_node_subject> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "anonymous_blank_node_subject" ;
	//   rdfs:comment "anonymous blank node subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <anonymous_blank_node_subject.ttl> ;
	//   mf:result    <labeled_blank_node_subject.nt> ;
	//   .

	{`[] <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.Blank{ID: "b1"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#anonymous_blank_node_object> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "anonymous_blank_node_object" ;
	//   rdfs:comment "anonymous blank node object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <anonymous_blank_node_object.ttl> ;
	//   mf:result    <labeled_blank_node_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> [] .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "b1"},
		},
	}},

	//<#sole_blankNodePropertyList> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "sole_blankNodePropertyList" ;
	//   rdfs:comment "sole blankNodePropertyList [ <p> <o> ] ." ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <sole_blankNodePropertyList.ttl> ;
	//   mf:result    <labeled_blank_node_subject.nt> ;
	//   .

	{`[ <http://a.example/p> <http://a.example/o> ] .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.Blank{ID: "b1"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
	}},

	//<#blankNodePropertyList_as_subject> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "blankNodePropertyList_as_subject" ;
	//   rdfs:comment "blankNodePropertyList as subject [ … ] <p> <o> ." ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <blankNodePropertyList_as_subject.ttl> ;
	//   mf:result    <blankNodePropertyList_as_subject.nt> ;
	//   .

	{`[ <http://a.example/p> <http://a.example/o> ] <http://a.example/p2> <http://a.example/o2> .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.Blank{ID: "b1"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.URI{URI: "http://a.example/o"},
		},
		rdf.Triple{
			Subj: rdf.Blank{ID: "b1"},
			Pred: rdf.URI{URI: "http://a.example/p2"},
			Obj:  rdf.URI{URI: "http://a.example/o2"},
		},
	}},

	//<#blankNodePropertyList_as_object> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "blankNodePropertyList_as_object" ;
	//   rdfs:comment "blankNodePropertyList as object <s> <p> [ … ] ." ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <blankNodePropertyList_as_object.ttl> ;
	//   mf:result    <blankNodePropertyList_as_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> [ <http://a.example/p2> <http://a.example/o2> ] .`, "", []rdf.Triple{
		rdf.Triple{
			Subj: rdf.URI{URI: "http://a.example/s"},
			Pred: rdf.URI{URI: "http://a.example/p"},
			Obj:  rdf.Blank{ID: "b1"},
		},
		rdf.Triple{
			Subj: rdf.Blank{ID: "b1"},
			Pred: rdf.URI{URI: "http://a.example/p2"},
			Obj:  rdf.URI{URI: "http://a.example/o2"},
		},
	}},

	//<#blankNodePropertyList_with_multiple_triples> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "blankNodePropertyList_with_multiple_triples" ;
	//   rdfs:comment "blankNodePropertyList with multiple triples [ <s> <p> ; <s2> <p2> ]" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <blankNodePropertyList_with_multiple_triples.ttl> ;
	//   mf:result    <blankNodePropertyList_with_multiple_triples.nt> ;
	//   .

	{`[ <http://a.example/p1> <http://a.example/o1> ; <http://a.example/p2> <http://a.example/o2> ] <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{}},

	//<#nested_blankNodePropertyLists> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "nested_blankNodePropertyLists" ;
	//   rdfs:comment "nested blankNodePropertyLists [ <p1> [ <p2> <o2> ] ; <p3> <o3> ]" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <nested_blankNodePropertyLists.ttl> ;
	//   mf:result    <nested_blankNodePropertyLists.nt> ;
	//   .

	{`[ <http://a.example/p1> [ <http://a.example/p2> <http://a.example/o2> ] ; <http://a.example/p> <http://a.example/o> ].`, "", []rdf.Triple{}},

	//<#blankNodePropertyList_containing_collection> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "blankNodePropertyList_containing_collection" ;
	//   rdfs:comment "blankNodePropertyList containing collection [ <p1> ( … ) ]" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <blankNodePropertyList_containing_collection.ttl> ;
	//   mf:result    <blankNodePropertyList_containing_collection.nt> ;
	//   .

	{`[ <http://a.example/p1> (1) ] .`, "", []rdf.Triple{}},

	//<#collection_subject> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "collection_subject" ;
	//   rdfs:comment "collection subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <collection_subject.ttl> ;
	//   mf:result    <collection_subject.nt> ;
	//   .

	{`(1) <http://a.example/p> <http://a.example/o> .`, "", []rdf.Triple{}},

	//<#collection_object> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "collection_object" ;
	//   rdfs:comment "collection object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <collection_object.ttl> ;
	//   mf:result    <collection_object.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> (1) .`, "", []rdf.Triple{}},

	//<#empty_collection> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "empty_collection" ;
	//   rdfs:comment "empty collection ()" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <empty_collection.ttl> ;
	//   mf:result    <empty_collection.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> () .`, "", []rdf.Triple{}},

	//<#nested_collection> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "nested_collection" ;
	//   rdfs:comment "nested collection (())" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <nested_collection.ttl> ;
	//   mf:result    <nested_collection.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> ((1)) .`, "", []rdf.Triple{}},

	//<#first> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "first" ;
	//   rdfs:comment "first, not last, non-empty nested collection" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <first.ttl> ;
	//   mf:result    <first.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> ((1) 2) .`, "", []rdf.Triple{}},

	//<#last> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "last" ;
	//   rdfs:comment "last, not first, non-empty nested collection" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <last.ttl> ;
	//   mf:result    <last.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> (1 (2)) .`, "", []rdf.Triple{}},

	//<#LITERAL1> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL1" ;
	//   rdfs:comment "LITERAL1 'x'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL1.ttl> ;
	//   mf:result    <LITERAL1.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 'x' .`, "", []rdf.Triple{}},

	//<#LITERAL1_ascii_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL1_ascii_boundaries" ;
	//   rdfs:comment "LITERAL1_ascii_boundaries '\\x00\\x09\\x0b\\x0c\\x0e\\x26\\x28...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL1_ascii_boundaries.ttl> ;
	//   mf:result    <LITERAL1_ascii_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\x00	&([]' .`, "", []rdf.Triple{}},

	//<#LITERAL1_with_UTF8_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL1_with_UTF8_boundaries" ;
	//   rdfs:comment "LITERAL1_with_UTF8_boundaries '\\x80\\x7ff\\x800\\xfff...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL1_with_UTF8_boundaries.ttl> ;
	//   mf:result    <LITERAL_with_UTF8_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '߿ࠀ࿿က쿿퀀퟿�𐀀𿿽񀀀󿿽􀀀􏿽' .`, "", []rdf.Triple{}},

	//<#LITERAL1_all_controls> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL1_all_controls" ;
	//   rdfs:comment "LITERAL1_all_controls '\\x00\\x01\\x02\\x03\\x04...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL1_all_controls.ttl> ;
	//   mf:result    <LITERAL1_all_controls.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "\u0000\u0001\u0002\u0003\u0004\u0005\u0006\u0007\u0008\t\u000B\u000C\u000E\u000F\u0010\u0011\u0012\u0013\u0014\u0015\u0016\u0017\u0018\u0019\u001A\u001B\u001C\u001D\u001E\u001F" .`, "", []rdf.Triple{}},

	//<#LITERAL1_all_punctuation> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL1_all_punctuation" ;
	//   rdfs:comment "LITERAL1_all_punctuation '!\"#$%&()...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL1_all_punctuation.ttl> ;
	//   mf:result    <LITERAL1_all_punctuation.nt> ;
	//   .

	{"<http://a.example/s> <http://a.example/p> ' !\"#$%&():;<=>?@[]^_`{|}~' .", "", []rdf.Triple{}},

	//<#LITERAL_LONG1> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG1" ;
	//   rdfs:comment "LITERAL_LONG1 '''x'''" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG1.ttl> ;
	//   mf:result    <LITERAL1.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''x''' .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG1_ascii_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG1_ascii_boundaries" ;
	//   rdfs:comment "LITERAL_LONG1_ascii_boundaries '\\x00\\x26\\x28...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG1_ascii_boundaries.ttl> ;
	//   mf:result    <LITERAL_LONG1_ascii_boundaries.nt> ;
	//   .

	{"<http://a.example/s> <http://a.example/p> '\x00&([]' .", "", []rdf.Triple{}},

	//<#LITERAL_LONG1_with_UTF8_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG1_with_UTF8_boundaries" ;
	//   rdfs:comment "LITERAL_LONG1_with_UTF8_boundaries '\\x80\\x7ff\\x800\\xfff...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG1_with_UTF8_boundaries.ttl> ;
	//   mf:result    <LITERAL_with_UTF8_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''߿ࠀ࿿က쿿퀀퟿�𐀀𿿽񀀀󿿽􀀀􏿽''' .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG1_with_1_squote> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG1_with_1_squote" ;
	//   rdfs:comment "LITERAL_LONG1 with 1 squote '''a'b'''" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG1_with_1_squote.ttl> ;
	//   mf:result    <LITERAL_LONG1_with_1_squote.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''x'y''' .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG1_with_2_squotes> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG1_with_2_squotes" ;
	//   rdfs:comment "LITERAL_LONG1 with 2 squotes '''a''b'''" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG1_with_2_squotes.ttl> ;
	//   mf:result    <LITERAL_LONG1_with_2_squotes.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''x''y''' .`, "", []rdf.Triple{}},

	//<#LITERAL2> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL2" ;
	//   rdfs:comment "LITERAL2 \"x\"" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL2.ttl> ;
	//   mf:result    <LITERAL1.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "x" .`, "", []rdf.Triple{}},

	//<#LITERAL2_ascii_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL2_ascii_boundaries" ;
	//   rdfs:comment "LITERAL2_ascii_boundaries '\\x00\\x09\\x0b\\x0c\\x0e\\x21\\x23...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL2_ascii_boundaries.ttl> ;
	//   mf:result    <LITERAL2_ascii_boundaries.nt> ;
	//   .

	{"<http://a.example/s> <http://a.example/p> \"\x00	!#[]\" .", "", []rdf.Triple{}},

	//<#LITERAL2_with_UTF8_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL2_with_UTF8_boundaries" ;
	//   rdfs:comment "LITERAL2_with_UTF8_boundaries '\\x80\\x7ff\\x800\\xfff...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL2_with_UTF8_boundaries.ttl> ;
	//   mf:result    <LITERAL_with_UTF8_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "߿ࠀ࿿က쿿퀀퟿�𐀀𿿽񀀀󿿽􀀀􏿽" .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG2> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG2" ;
	//   rdfs:comment "LITERAL_LONG2 \"\"\"x\"\"\"" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2.ttl> ;
	//   mf:result    <LITERAL1.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> """x""" .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG2_ascii_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG2_ascii_boundaries" ;
	//   rdfs:comment "LITERAL_LONG2_ascii_boundaries '\\x00\\x21\\x23...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2_ascii_boundaries.ttl> ;
	//   mf:result    <LITERAL_LONG2_ascii_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> \"\x00!#[]\" .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG2_with_UTF8_boundaries> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG2_with_UTF8_boundaries" ;
	//   rdfs:comment "LITERAL_LONG2_with_UTF8_boundaries '\\x80\\x7ff\\x800\\xfff...'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2_with_UTF8_boundaries.ttl> ;
	//   mf:result    <LITERAL_with_UTF8_boundaries.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> """߿ࠀ࿿က쿿퀀퟿�𐀀𿿽񀀀󿿽􀀀􏿽""" .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG2_with_1_squote> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG2_with_1_squote" ;
	//   rdfs:comment "LITERAL_LONG2 with 1 squote \"\"\"a\"b\"\"\"" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2_with_1_squote.ttl> ;
	//   mf:result    <LITERAL_LONG2_with_1_squote.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> """x"y""" .`, "", []rdf.Triple{}},

	//<#LITERAL_LONG2_with_2_squotes> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "LITERAL_LONG2_with_2_squotes" ;
	//   rdfs:comment "LITERAL_LONG2 with 2 squotes \"\"\"a\"\"b\"\"\"" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2_with_2_squotes.ttl> ;
	//   mf:result    <LITERAL_LONG2_with_2_squotes.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> """x""y""" .`, "", []rdf.Triple{}},

	//<#literal_with_CHARACTER_TABULATION> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_CHARACTER_TABULATION" ;
	//   rdfs:comment "literal with CHARACTER TABULATION" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_CHARACTER_TABULATION.ttl> ;
	//   mf:result    <literal_with_CHARACTER_TABULATION.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '	' .`, "", []rdf.Triple{}},

	//<#literal_with_BACKSPACE> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_BACKSPACE" ;
	//   rdfs:comment "literal with BACKSPACE" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_BACKSPACE.ttl> ;
	//   mf:result    <literal_with_BACKSPACE.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '' .`, "", []rdf.Triple{}},

	//<#literal_with_LINE_FEED> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_LINE_FEED" ;
	//   rdfs:comment "literal with LINE FEED" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_LINE_FEED.ttl> ;
	//   mf:result    <literal_with_LINE_FEED.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''
''' .`, "", []rdf.Triple{}},

	//<#literal_with_CARRIAGE_RETURN> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_CARRIAGE_RETURN" ;
	//   rdfs:comment "literal with CARRIAGE RETURN" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_CARRIAGE_RETURN.ttl> ;
	//   mf:result    <literal_with_CARRIAGE_RETURN.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '''
''' .`, "", []rdf.Triple{}},

	//<#literal_with_FORM_FEED> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_FORM_FEED" ;
	//   rdfs:comment "literal with FORM FEED" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_FORM_FEED.ttl> ;
	//   mf:result    <literal_with_FORM_FEED.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '' .`, "", []rdf.Triple{}},

	//<#literal_with_REVERSE_SOLIDUS> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_REVERSE_SOLIDUS" ;
	//   rdfs:comment "literal with REVERSE SOLIDUS" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_REVERSE_SOLIDUS.ttl> ;
	//   mf:result    <literal_with_REVERSE_SOLIDUS.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\\' .`, "", []rdf.Triple{}},

	//<#literal_with_escaped_CHARACTER_TABULATION> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_escaped_CHARACTER_TABULATION" ;
	//   rdfs:comment "literal with escaped CHARACTER TABULATION" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_escaped_CHARACTER_TABULATION.ttl> ;
	//   mf:result    <literal_with_CHARACTER_TABULATION.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\t' .`, "", []rdf.Triple{}},

	//<#literal_with_escaped_BACKSPACE> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_escaped_BACKSPACE" ;
	//   rdfs:comment "literal with escaped BACKSPACE" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_escaped_BACKSPACE.ttl> ;
	//   mf:result    <literal_with_BACKSPACE.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\b' .`, "", []rdf.Triple{}},

	//<#literal_with_escaped_LINE_FEED> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_escaped_LINE_FEED" ;
	//   rdfs:comment "literal with escaped LINE FEED" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_escaped_LINE_FEED.ttl> ;
	//   mf:result    <literal_with_LINE_FEED.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\n' .`, "", []rdf.Triple{}},

	//<#literal_with_escaped_CARRIAGE_RETURN> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_escaped_CARRIAGE_RETURN" ;
	//   rdfs:comment "literal with escaped CARRIAGE RETURN" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_escaped_CARRIAGE_RETURN.ttl> ;
	//   mf:result    <literal_with_CARRIAGE_RETURN.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\r' .`, "", []rdf.Triple{}},

	//<#literal_with_escaped_FORM_FEED> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_escaped_FORM_FEED" ;
	//   rdfs:comment "literal with escaped FORM FEED" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_escaped_FORM_FEED.ttl> ;
	//   mf:result    <literal_with_FORM_FEED.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\f' .`, "", []rdf.Triple{}},

	//<#literal_with_numeric_escape4> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_numeric_escape4" ;
	//   rdfs:comment "literal with numeric escape4 \\u" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_numeric_escape4.ttl> ;
	//   mf:result    <literal_with_numeric_escape4.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\u006F' .`, "", []rdf.Triple{}},

	//<#literal_with_numeric_escape8> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_with_numeric_escape8" ;
	//   rdfs:comment "literal with numeric escape8 \\U" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_with_numeric_escape8.ttl> ;
	//   mf:result    <literal_with_numeric_escape4.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> '\U0000006F' .`, "", []rdf.Triple{}},

	//<#IRIREF_datatype> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "IRIREF_datatype" ;
	//   rdfs:comment "IRIREF datatype \"\"^^<t>" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <IRIREF_datatype.ttl> ;
	//   mf:result    <IRIREF_datatype.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "1"^^<http://www.w3.org/2001/XMLSchema#integer> .`, "", []rdf.Triple{}},

	//<#prefixed_name_datatype> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "prefixed_name_datatype" ;
	//   rdfs:comment "prefixed name datatype \"\"^^p:t" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <prefixed_name_datatype.ttl> ;
	//   mf:result    <IRIREF_datatype.nt> ;
	//   .

	{`@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
<http://a.example/s> <http://a.example/p> "1"^^xsd:integer .`, "", []rdf.Triple{}},

	//<#bareword_integer> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "bareword_integer" ;
	//   rdfs:comment "bareword integer" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <bareword_integer.ttl> ;
	//   mf:result    <IRIREF_datatype.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 1 .`, "", []rdf.Triple{}},

	//<#bareword_decimal> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "bareword_decimal" ;
	//   rdfs:comment "bareword decimal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <bareword_decimal.ttl> ;
	//   mf:result    <bareword_decimal.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 1.0 .`, "", []rdf.Triple{}},

	//<#bareword_double> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "bareword_double" ;
	//   rdfs:comment "bareword double" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <bareword_double.ttl> ;
	//   mf:result    <bareword_double.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 1E0 .`, "", []rdf.Triple{}},

	//<#double_lower_case_e> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "double_lower_case_e" ;
	//   rdfs:comment "double lower case e" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <double_lower_case_e.ttl> ;
	//   mf:result    <double_lower_case_e.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 1e0 .`, "", []rdf.Triple{}},

	//<#negative_numeric> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "negative_numeric" ;
	//   rdfs:comment "negative numeric" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <negative_numeric.ttl> ;
	//   mf:result    <negative_numeric.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> -1 .`, "", []rdf.Triple{}},

	//<#positive_numeric> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "positive_numeric" ;
	//   rdfs:comment "positive numeric" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <positive_numeric.ttl> ;
	//   mf:result    <positive_numeric.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> +1 .`, "", []rdf.Triple{}},

	//<#numeric_with_leading_0> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "numeric_with_leading_0" ;
	//   rdfs:comment "numeric with leading 0" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <numeric_with_leading_0.ttl> ;
	//   mf:result    <numeric_with_leading_0.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> 01 .`, "", []rdf.Triple{}},

	//<#literal_true> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_true" ;
	//   rdfs:comment "literal true" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_true.ttl> ;
	//   mf:result    <literal_true.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> true .`, "", []rdf.Triple{}},

	//<#literal_false> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "literal_false" ;
	//   rdfs:comment "literal false" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <literal_false.ttl> ;
	//   mf:result    <literal_false.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> false .`, "", []rdf.Triple{}},

	//<#langtagged_non_LONG> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "langtagged_non_LONG" ;
	//   rdfs:comment "langtagged non-LONG \"x\"@en" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <langtagged_non_LONG.ttl> ;
	//   mf:result    <langtagged_non_LONG.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "chat"@en .`, "", []rdf.Triple{}},

	//<#langtagged_LONG> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "langtagged_LONG" ;
	//   rdfs:comment "langtagged LONG \"\"\"x\"\"\"@en" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <langtagged_LONG.ttl> ;
	//   mf:result    <langtagged_non_LONG.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> """chat"""@en .`, "", []rdf.Triple{}},

	//<#lantag_with_subtag> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "lantag_with_subtag" ;
	//   rdfs:comment "lantag with subtag \"x\"@en-us" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <lantag_with_subtag.ttl> ;
	//   mf:result    <lantag_with_subtag.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> "chat"@en-us .`, "", []rdf.Triple{}},

	//<#objectList_with_two_objects> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "objectList_with_two_objects" ;
	//   rdfs:comment "objectList with two objects … <o1>,<o2>" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <objectList_with_two_objects.ttl> ;
	//   mf:result    <objectList_with_two_objects.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p> <http://a.example/o1>, <http://a.example/o2> .`, "", []rdf.Triple{}},

	//<#predicateObjectList_with_two_objectLists> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "predicateObjectList_with_two_objectLists" ;
	//   rdfs:comment "predicateObjectList with two objectLists … <o1>,<o2>" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <predicateObjectList_with_two_objectLists.ttl> ;
	//   mf:result    <predicateObjectList_with_two_objectLists.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p1> <http://a.example/o1>; <http://a.example/p2> <http://a.example/o2> .`, "", []rdf.Triple{}},

	//<#repeated_semis_at_end> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "repeated_semis_at_end" ;
	//   rdfs:comment "repeated semis at end <s> <p> <o> ;; <p2> <o2> ." ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <repeated_semis_at_end.ttl> ;
	//   mf:result    <predicateObjectList_with_two_objectLists.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p1> <http://a.example/o1>;; <http://a.example/p2> <http://a.example/o2> .`, "", []rdf.Triple{}},

	//<#repeated_semis_not_at_end> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "repeated_semis_not_at_end" ;
	//   rdfs:comment "repeated semis not at end <s> <p> <o> ;;." ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <repeated_semis_not_at_end.ttl> ;
	//   mf:result    <repeated_semis_not_at_end.nt> ;
	//   .

	{`<http://a.example/s> <http://a.example/p1> <http://a.example/o1>;; .`, "", []rdf.Triple{}},

	//# original tests-ttl
	//<#turtle-syntax-file-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-file-01" ;
	//   rdfs:comment "Empty file" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-file-01.ttl> ;
	//   .

	{``, "", []rdf.Triple{}},

	//<#turtle-syntax-file-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-file-02" ;
	//   rdfs:comment "Only comment" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-file-02.ttl> ;
	//   .

	{`#Empty file.`, "", []rdf.Triple{}},

	//<#turtle-syntax-file-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-file-03" ;
	//   rdfs:comment "One comment, one empty line" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-file-03.ttl> ;
	//   .

	{`#One comment, one empty line.
`, "", []rdf.Triple{}},

	//<#turtle-syntax-uri-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-uri-01" ;
	//   rdfs:comment "Only IRIs" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-uri-01.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-uri-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-uri-02" ;
	//   rdfs:comment "IRIs with Unicode escape" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-uri-02.ttl> ;
	//   .

	{`# x53 is capital S
<http://www.w3.org/2013/TurtleTests/\u0053> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-uri-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-uri-03" ;
	//   rdfs:comment "IRIs with long Unicode escape" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-uri-03.ttl> ;
	//   .

	{`# x53 is capital S
<http://www.w3.org/2013/TurtleTests/\U00000053> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-uri-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-uri-04" ;
	//   rdfs:comment "Legal IRIs" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-uri-04.ttl> ;
	//   .

	{`# IRI with all chars in it.
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p>
<scheme:!$%25&'()*+,-./0123456789:/@ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz~?#> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-base-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-base-01" ;
	//   rdfs:comment "@base" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-base-01.ttl> ;
	//   .

	{`@base <http://www.w3.org/2013/TurtleTests/> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-base-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-base-02" ;
	//   rdfs:comment "BASE" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-base-02.ttl> ;
	//   .

	{`BASE <http://www.w3.org/2013/TurtleTests/>`, "", []rdf.Triple{}},

	//<#turtle-syntax-base-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-base-03" ;
	//   rdfs:comment "@base with relative IRIs" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-base-03.ttl> ;
	//   .

	{`@base <http://www.w3.org/2013/TurtleTests/> .
<s> <p> <o> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-base-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-base-04" ;
	//   rdfs:comment "base with relative IRIs" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-base-04.ttl> ;
	//   .

	{`base <http://www.w3.org/2013/TurtleTests/>
<s> <p> <o> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-01" ;
	//   rdfs:comment "@prefix" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-02" ;
	//   rdfs:comment "PreFIX" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-02.ttl> ;
	//   .

	{`PreFIX : <http://www.w3.org/2013/TurtleTests/>`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-03" ;
	//   rdfs:comment "Empty PREFIX" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-03.ttl> ;
	//   .

	{`PREFIX : <http://www.w3.org/2013/TurtleTests/>
:s :p :123 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-04" ;
	//   rdfs:comment "Empty @prefix with % escape" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :%20 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-05" ;
	//   rdfs:comment "@prefix with no suffix" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-05.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
: : : .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-06> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-06" ;
	//   rdfs:comment "colon is a legal pname character" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-06.ttl> ;
	//   .

	{`# colon is a legal pname character
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@prefix x: <http://www.w3.org/2013/TurtleTests/> .
:a:b:c  x:d:e:f :::: .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-07> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-07" ;
	//   rdfs:comment "dash is a legal pname character" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-07.ttl> ;
	//   .

	{`# dash is a legal pname character
@prefix x: <http://www.w3.org/2013/TurtleTests/> .
x:a-b-c  x:p x:o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-08> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-08" ;
	//   rdfs:comment "underscore is a legal pname character" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-08.ttl> ;
	//   .

	{`# underscore is a legal pname character
@prefix x: <http://www.w3.org/2013/TurtleTests/> .
x:_  x:p_1 x:o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-prefix-09> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-prefix-09" ;
	//   rdfs:comment "percents in pnames" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-prefix-09.ttl> ;
	//   .

	{`# percents
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@prefix x: <http://www.w3.org/2013/TurtleTests/> .
:a%3E  x:%25 :a%3Eb .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-01" ;
	//   rdfs:comment "string literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-01.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "string" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-02" ;
	//   rdfs:comment "langString literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-02.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "string"@en .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-03" ;
	//   rdfs:comment "langString literal with region" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-03.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "string"@en-uk .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-04" ;
	//   rdfs:comment "squote string literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-04.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 'string' .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-05" ;
	//   rdfs:comment "squote langString literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-05.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 'string'@en .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-06> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-06" ;
	//   rdfs:comment "squote langString literal with region" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-06.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 'string'@en-uk .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-07> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-07" ;
	//   rdfs:comment "long string literal with embedded single- and double-quotes" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-07.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> """abc""def''ghi""" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-08> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-08" ;
	//   rdfs:comment "long string literal with embedded newline" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-08.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> """abc
def""" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-09> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-09" ;
	//   rdfs:comment "squote long string literal with embedded single- and double-quotes" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-09.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> '''abc
def''' .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-10> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-10" ;
	//   rdfs:comment "long langString literal with embedded newline" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-10.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> """abc
def"""@en .`, "", []rdf.Triple{}},

	//<#turtle-syntax-string-11> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-string-11" ;
	//   rdfs:comment "squote long langString literal with embedded newline" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-string-11.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> '''abc
def'''@en .`, "", []rdf.Triple{}},

	//<#turtle-syntax-str-esc-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-str-esc-01" ;
	//   rdfs:comment "string literal with escaped newline" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-str-esc-01.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "a\n" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-str-esc-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-str-esc-02" ;
	//   rdfs:comment "string literal with Unicode escape" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-str-esc-02.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "a\u0020b" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-str-esc-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-str-esc-03" ;
	//   rdfs:comment "string literal with long Unicode escape" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-str-esc-03.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "a\U00000020b" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-pname-esc-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-pname-esc-01" ;
	//   rdfs:comment "pname with back-slash escapes" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-pname-esc-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :\~\.\-\!\$\&\'\(\)\*\+\,\;\=\/\?\#\@\_\%AA .`, "", []rdf.Triple{}},

	//<#turtle-syntax-pname-esc-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-pname-esc-02" ;
	//   rdfs:comment "pname with back-slash escapes (2)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-pname-esc-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :0123\~\.\-\!\$\&\'\(\)\*\+\,\;\=\/\?\#\@\_\%AA123 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-pname-esc-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-pname-esc-03" ;
	//   rdfs:comment "pname with back-slash escapes (3)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-pname-esc-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:xyz\~ :abc\.:  : .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-01" ;
	//   rdfs:comment "bnode subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
[] :p :o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-02" ;
	//   rdfs:comment "bnode object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p [] .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-03" ;
	//   rdfs:comment "bnode property list object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p [ :q :o ] .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-04" ;
	//   rdfs:comment "bnode property list object (2)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p [ :q1 :o1 ; :q2 :o2 ] .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-05" ;
	//   rdfs:comment "bnode property list subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-05.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
[ :q1 :o1 ; :q2 :o2 ] :p :o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-06> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-06" ;
	//   rdfs:comment "labeled bnode subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-06.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
_:a  :p :o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-07> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-07" ;
	//   rdfs:comment "labeled bnode subject and object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-07.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s  :p _:a .
_:a  :p :o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-08> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-08" ;
	//   rdfs:comment "bare bnode property list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-08.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
[ :p  :o ] .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-09> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-09" ;
	//   rdfs:comment "bnode property list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-09.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
[ :p  :o1,:2 ] .
:s :p :o  .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bnode-10> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-bnode-10" ;
	//   rdfs:comment "mixed bnode property list and triple" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bnode-10.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .

:s1 :p :o .
[ :p1  :o1 ; :p2 :o2 ] .
:s2 :p :o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-01" ;
	//   rdfs:comment "integer literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-01.ttl> ;
	//   .

	{`<s> <p> 123 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-02" ;
	//   rdfs:comment "negative integer literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-02.ttl> ;
	//   .

	{`<s> <p> -123 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-03" ;
	//   rdfs:comment "positive integer literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-03.ttl> ;
	//   .

	{`<s> <p> +123 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-04" ;
	//   rdfs:comment "decimal literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-04.ttl> ;
	//   .

	{`# This is a decimal.
<s> <p> 123.0 . `, "", []rdf.Triple{}},

	//<#turtle-syntax-number-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-05" ;
	//   rdfs:comment "decimal literal (no leading digits)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-05.ttl> ;
	//   .

	{`# This is a decimal.
<s> <p> .1 . `, "", []rdf.Triple{}},

	//<#turtle-syntax-number-06> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-06" ;
	//   rdfs:comment "negative decimal literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-06.ttl> ;
	//   .

	{`# This is a decimal.
<s> <p> -123.0 . `, "", []rdf.Triple{}},

	//<#turtle-syntax-number-07> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-07" ;
	//   rdfs:comment "positive decimal literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-07.ttl> ;
	//   .

	{`# This is a decimal.
<s> <p> +123.0 . `, "", []rdf.Triple{}},

	//<#turtle-syntax-number-08> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-08" ;
	//   rdfs:comment "integer literal with decimal lexical confusion" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-08.ttl> ;
	//   .

	{`# This is an integer
<s> <p> 123.`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-09> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-09" ;
	//   rdfs:comment "double literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-09.ttl> ;
	//   .

	{`<s> <p> 123.0e1 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-10> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-10" ;
	//   rdfs:comment "negative double literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-10.ttl> ;
	//   .

	{`<s> <p> -123e-1 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-number-11> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-number-11" ;
	//   rdfs:comment "double literal no fraction" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-number-11.ttl> ;
	//   .

	{`<s> <p> 123.E+1 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-datatypes-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-datatypes-01" ;
	//   rdfs:comment "xsd:byte literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-datatypes-01.ttl> ;
	//   .

	{`@prefix xsd:     <http://www.w3.org/2001/XMLSchema#> .
<s> <p> "123"^^xsd:byte .`, "", []rdf.Triple{}},

	//<#turtle-syntax-datatypes-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-datatypes-02" ;
	//   rdfs:comment "integer as xsd:string" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-datatypes-02.ttl> ;
	//   .

	{`@prefix rdf:     <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix xsd:     <http://www.w3.org/2001/XMLSchema#> .
<s> <p> "123"^^xsd:string .`, "", []rdf.Triple{}},

	//<#turtle-syntax-kw-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-kw-01" ;
	//   rdfs:comment "boolean literal (true)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-kw-01.ttl> ;
	//   .

	{`<s> <p> true .`, "", []rdf.Triple{}},

	//<#turtle-syntax-kw-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-kw-02" ;
	//   rdfs:comment "boolean literal (false)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-kw-02.ttl> ;
	//   .

	{`<s> <p> false .`, "", []rdf.Triple{}},

	//<#turtle-syntax-kw-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-kw-03" ;
	//   rdfs:comment "'a' as keyword" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-kw-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s a :C .`, "", []rdf.Triple{}},

	//<#turtle-syntax-struct-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-struct-01" ;
	//   rdfs:comment "object list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-struct-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :o1 , :o2 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-struct-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-struct-02" ;
	//   rdfs:comment "predicate list with object list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-struct-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p1 :o1 ;
   :p2 :o2 .`, "", []rdf.Triple{}},

	//<#turtle-syntax-struct-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-struct-03" ;
	//   rdfs:comment "predicate list with object list and dangling ';'" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-struct-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p1 :o1 ;
   :p2 :o2 ;
   .`, "", []rdf.Triple{}},

	//<#turtle-syntax-struct-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-struct-04" ;
	//   rdfs:comment "predicate list with multiple ;;" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-struct-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p1 :o1 ;;
   :p2 :o2 
   .`, "", []rdf.Triple{}},

	//<#turtle-syntax-struct-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-struct-05" ;
	//   rdfs:comment "predicate list with multiple ;;" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-struct-05.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p1 :o1 ;
   :p2 :o2 ;;
   .`, "", []rdf.Triple{}},

	//<#turtle-syntax-lists-01> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-lists-01" ;
	//   rdfs:comment "empty list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-lists-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p () .`, "", []rdf.Triple{}},

	//<#turtle-syntax-lists-02> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-lists-02" ;
	//   rdfs:comment "mixed list" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-lists-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p (1 "2" :o) .`, "", []rdf.Triple{}},

	//<#turtle-syntax-lists-03> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-lists-03" ;
	//   rdfs:comment "isomorphic list as subject and object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-lists-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
(1) :p (1) .`, "", []rdf.Triple{}},

	//<#turtle-syntax-lists-04> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-lists-04" ;
	//   rdfs:comment "lists of lists" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-lists-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
(()) :p (()) .`, "", []rdf.Triple{}},

	//<#turtle-syntax-lists-05> rdf:type rdft:TestTurtlePositiveSyntax ;
	//   mf:name    "turtle-syntax-lists-05" ;
	//   rdfs:comment "mixed lists with embedded lists" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-lists-05.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
(1 2 (1 2)) :p (( "a") "b" :o) .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bad-uri-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-uri-01" ;
	//   rdfs:comment "Bad IRI : space (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-uri-01.ttl> ;
	//   .

	{`# Bad IRI : space.
<http://www.w3.org/2013/TurtleTests/ space> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Bad IRI : space (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-uri-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-uri-02" ;
	//   rdfs:comment "Bad IRI : bad escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-uri-02.ttl> ;
	//   .

	{`# Bad IRI : bad escape
<http://www.w3.org/2013/TurtleTests/\u00ZZ11> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Bad IRI : bad escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-uri-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-uri-03" ;
	//   rdfs:comment "Bad IRI : bad long escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-uri-03.ttl> ;
	//   .

	{`# Bad IRI : bad escape
<http://www.w3.org/2013/TurtleTests/\U00ZZ1111> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Bad IRI : bad long escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-uri-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-uri-04" ;
	//   rdfs:comment "Bad IRI : character escapes not allowed (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-uri-04.ttl> ;
	//   .

	{`# Bad IRI : character escapes not allowed.
<http://www.w3.org/2013/TurtleTests/\n> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Bad IRI : character escapes not allowed (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-uri-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-uri-05" ;
	//   rdfs:comment "Bad IRI : character escapes not allowed (2) (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-uri-05.ttl> ;
	//   .

	{`# Bad IRI : character escapes not allowed.
<http://www.w3.org/2013/TurtleTests/\/> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Bad IRI : character escapes not allowed (2) (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-prefix-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-prefix-01" ;
	//   rdfs:comment "No prefix (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-prefix-01.ttl> ;
	//   .

	{`# No prefix
:s <http://www.w3.org/2013/TurtleTests/p> "x" .`, "No prefix (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-prefix-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-prefix-02" ;
	//   rdfs:comment "No prefix (2) (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-prefix-02.ttl> ;
	//   .

	{`# No prefix
@prefix rdf:     <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
<http://www.w3.org/2013/TurtleTests/s> rdf:type :C .`, "No prefix (2) (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-prefix-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-prefix-03" ;
	//   rdfs:comment "@prefix without URI (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-prefix-03.ttl> ;
	//   .

	{`# @prefix without URI.
@prefix ex: .`, "@prefix without URI (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-prefix-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-prefix-04" ;
	//   rdfs:comment "@prefix without prefix name (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-prefix-04.ttl> ;
	//   .

	{`# @prefix without prefix name .
@prefix <http://www.w3.org/2013/TurtleTests/> .`, "@prefix without prefix name (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-prefix-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-prefix-05" ;
	//   rdfs:comment "@prefix without ':' (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-prefix-05.ttl> ;
	//   .

	{`# @prefix without :
@prefix x <http://www.w3.org/2013/TurtleTests/> .`, "@prefix without ':' (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-base-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-base-01" ;
	//   rdfs:comment "@base without URI (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-base-01.ttl> ;
	//   .

	{`# @base without URI.
@base .`, "@base without URI (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-base-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-base-02" ;
	//   rdfs:comment "@base in wrong case (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-base-02.ttl> ;
	//   .

	{`# @base in wrong case.
@BASE <http://www.w3.org/2013/TurtleTests/> .`, "@base in wrong case (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-base-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-base-03" ;
	//   rdfs:comment "BASE without URI (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-base-03.ttl> ;
	//   .

	{`# FULL STOP used after SPARQL BASE
BASE <http://www.w3.org/2013/TurtleTests/> .
<s> <p> <o> .`, "BASE without URI (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-01" ;
	//   rdfs:comment "Turtle is not TriG (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-01.ttl> ;
	//   .

	{`# Turtle is not TriG
{ <http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> }`, "Turtle is not TriG (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-02" ;
	//   rdfs:comment "Turtle is not N3 (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-02.ttl> ;
	//   .

	{`# Turtle is not N3
<http://www.w3.org/2013/TurtleTests/s> = <http://www.w3.org/2013/TurtleTests/o> .`, "Turtle is not N3 (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-03" ;
	//   rdfs:comment "Turtle is not NQuads (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-03.ttl> ;
	//   .

	{`# Turtle is not NQuads
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> <http://www.w3.org/2013/TurtleTests/g> .`, "Turtle is not NQuads (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-04" ;
	//   rdfs:comment "Turtle does not allow literals-as-subjects (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-04.ttl> ;
	//   .

	{`# Turtle does not allow literals-as-subjects
"hello" <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "Turtle does not allow literals-as-subjects (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-05" ;
	//   rdfs:comment "Turtle does not allow literals-as-predicates (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-05.ttl> ;
	//   .

	{`# Turtle does not allow literals-as-predicates
<http://www.w3.org/2013/TurtleTests/s> "hello" <http://www.w3.org/2013/TurtleTests/o> .`, "Turtle does not allow literals-as-predicates (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-06> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-06" ;
	//   rdfs:comment "Turtle does not allow bnodes-as-predicates (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-06.ttl> ;
	//   .

	{`# Turtle does not allow bnodes-as-predicates
<http://www.w3.org/2013/TurtleTests/s> [] <http://www.w3.org/2013/TurtleTests/o> .`, "Turtle does not allow bnodes-as-predicates (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-07> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-07" ;
	//   rdfs:comment "Turtle does not allow labeled bnodes-as-predicates (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-07.ttl> ;
	//   .

	{`# Turtle does not allow bnodes-as-predicates
<http://www.w3.org/2013/TurtleTests/s> _:p <http://www.w3.org/2013/TurtleTests/o> .`, "Turtle does not allow labeled bnodes-as-predicates (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-kw-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-kw-01" ;
	//   rdfs:comment "'A' is not a keyword (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-kw-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s A :C .`, "'A' is not a keyword (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-kw-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-kw-02" ;
	//   rdfs:comment "'a' cannot be used as subject (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-kw-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
a :p :o .`, "'a' cannot be used as subject (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-kw-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-kw-03" ;
	//   rdfs:comment "'a' cannot be used as object (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-kw-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p a .`, "'a' cannot be used as object (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-kw-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-kw-04" ;
	//   rdfs:comment "'true' cannot be used as subject (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-kw-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
true :p :o .`, "'true' cannot be used as subject (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-kw-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-kw-05" ;
	//   rdfs:comment "'true' cannot be used as object (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-kw-05.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s true :o .`, "'true' cannot be used as object (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-01" ;
	//   rdfs:comment "{} fomulae not in Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-01.ttl> ;
	//   .

	{`# {} fomulae not in Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .

{ :a :q :c . } :p :z .
`, "{} fomulae not in Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-02" ;
	//   rdfs:comment "= is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-02.ttl> ;
	//   .

	{`# = is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .

:a = :b .`, "= is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-03" ;
	//   rdfs:comment "N3 paths not in Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-03.ttl> ;
	//   .

	{`# N3 paths
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@prefix ns: <http://www.w3.org/2013/TurtleTests/p#> .

:x.
  ns:p.
    ns:q :p :z .`, "N3 paths not in Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-04" ;
	//   rdfs:comment "N3 paths not in Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-04.ttl> ;
	//   .

	{`# N3 paths
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@prefix ns: <http://www.w3.org/2013/TurtleTests/p#> .

:x^ns:p :p :z .`, "N3 paths not in Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-05" ;
	//   rdfs:comment "N3 is...of not in Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-05.ttl> ;
	//   .

	{`# N3 is...of
@prefix : <http://www.w3.org/2013/TurtleTests/> .

:z is :p of :x .`, "N3 is...of not in Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-06> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-06" ;
	//   rdfs:comment "N3 paths not in Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-06.ttl> ;
	//   .

	{`# = is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .

:a.:b.:c .`, "N3 paths not in Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-07> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-07" ;
	//   rdfs:comment "@keywords is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-07.ttl> ;
	//   .

	{`# @keywords is not Turtle
@keywords a .
x a Item .`, "@keywords is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-08> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-08" ;
	//   rdfs:comment "@keywords is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-08.ttl> ;
	//   .

	{`# @keywords is not Turtle
@keywords a .
x a Item .`, "@keywords is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-09> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-09" ;
	//   rdfs:comment "=> is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-09.ttl> ;
	//   .

	{`# => is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s => :o .`, "=> is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-10> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-10" ;
	//   rdfs:comment "<= is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-10.ttl> ;
	//   .

	{`# <= is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s <= :o .`, "<= is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-11> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-11" ;
	//   rdfs:comment "@forSome is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-11.ttl> ;
	//   .

	{`# @forSome is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@forSome :x .`, "@forSome is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-12> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-12" ;
	//   rdfs:comment "@forAll is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-12.ttl> ;
	//   .

	{`# @forAll is not Turtle
@prefix : <http://www.w3.org/2013/TurtleTests/> .
@forAll :x .`, "@forAll is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-n3-extras-13> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-n3-extras-13" ;
	//   rdfs:comment "@keywords is not Turtle (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-n3-extras-13.ttl> ;
	//   .

	{`# @keywords is not Turtle
@keywords .
x @a Item .`, "@keywords is not Turtle (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-08> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-08" ;
	//   rdfs:comment "missing '.' (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-08.ttl> ;
	//   .

	{`# No DOT
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o>`, "missing '.' (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-09> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-09" ;
	//   rdfs:comment "extra '.' (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-09.ttl> ;
	//   .

	{`# Too many DOT
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> . .`, "extra '.' (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-10> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-10" ;
	//   rdfs:comment "extra '.' (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-10.ttl> ;
	//   .

	{`# Too many DOT
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> . .
<http://www.w3.org/2013/TurtleTests/s1> <http://www.w3.org/2013/TurtleTests/p1> <http://www.w3.org/2013/TurtleTests/o1> .`, "extra '.' (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-11> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-11" ;
	//   rdfs:comment "trailing ';' no '.' (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-11.ttl> ;
	//   .

	{`# Trailing ;
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> ;`, "trailing ';' no '.' (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-12> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-12" ;
	//   rdfs:comment "subject, predicate, no object (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-12.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> `, "subject, predicate, no object (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-13> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-13" ;
	//   rdfs:comment "subject, predicate, no object (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-13.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> `, "subject, predicate, no object (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-14> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-14" ;
	//   rdfs:comment "literal as subject (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-14.ttl> ;
	//   .

	{`# Literal as subject
"abc" <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/p>  .`, "literal as subject (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-15> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-15" ;
	//   rdfs:comment "literal as predicate (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-15.ttl> ;
	//   .

	{`# Literal as predicate
<http://www.w3.org/2013/TurtleTests/s> "abc" <http://www.w3.org/2013/TurtleTests/p>  .`, "literal as predicate (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-16> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-16" ;
	//   rdfs:comment "bnode as predicate (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-16.ttl> ;
	//   .

	{`# BNode as predicate
<http://www.w3.org/2013/TurtleTests/s> [] <http://www.w3.org/2013/TurtleTests/p>  .`, "bnode as predicate (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-struct-17> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-struct-17" ;
	//   rdfs:comment "labeled bnode as predicate (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-struct-17.ttl> ;
	//   .

	{`# BNode as predicate
<http://www.w3.org/2013/TurtleTests/s> _:a <http://www.w3.org/2013/TurtleTests/p>  .`, "labeled bnode as predicate (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-lang-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-lang-01" ;
	//   rdfs:comment "langString with bad lang (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-lang-01.ttl> ;
	//   .

	{`# Bad lang tag
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "string"@1 .`, "langString with bad lang (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-esc-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-esc-01" ;
	//   rdfs:comment "Bad string escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-esc-01.ttl> ;
	//   .

	{`# Bad string escape
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "a\zb" .`, "Bad string escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-esc-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-esc-02" ;
	//   rdfs:comment "Bad string escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-esc-02.ttl> ;
	//   .

	{`# Bad string escape
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "\uWXYZ" .`, "Bad string escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-esc-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-esc-03" ;
	//   rdfs:comment "Bad string escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-esc-03.ttl> ;
	//   .

	{`# Bad string escape
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "\U0000WXYZ" .`, "Bad string escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-esc-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-esc-04" ;
	//   rdfs:comment "Bad string escape (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-esc-04.ttl> ;
	//   .

	{`# Bad string escape
<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> "\U0000WXYZ" .`, "Bad string escape (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-pname-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-pname-01" ;
	//   rdfs:comment "'~' must be escaped in pname (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-pname-01.ttl> ;
	//   .

	{`# ~ must be escaped.
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:a~b :p :o .`, "'~' must be escaped in pname (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-pname-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-pname-02" ;
	//   rdfs:comment "Bad %-sequence in pname (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-pname-02.ttl> ;
	//   .

	{`# Bad %-sequence
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:a%2 :p :o .`, "Bad %-sequence in pname (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-pname-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-pname-03" ;
	//   rdfs:comment "Bad unicode escape in pname (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-pname-03.ttl> ;
	//   .

	{`# No \u (x39 is "9")
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:a\u0039 :p :o .`, "Bad unicode escape in pname (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-01" ;
	//   rdfs:comment "mismatching string literal open/close (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-01.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p "abc' .`, "mismatching string literal open/close (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-02" ;
	//   rdfs:comment "mismatching string literal open/close (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-02.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p 'abc" .`, "mismatching string literal open/close (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-03" ;
	//   rdfs:comment "mismatching string literal long/short (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-03.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p '''abc' .`, "mismatching string literal long/short (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-04" ;
	//   rdfs:comment "mismatching long string literal open/close (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-04.ttl> ;
	//   .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p """abc''' .`, "mismatching long string literal open/close (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-05" ;
	//   rdfs:comment "Long literal with missing end (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-05.ttl> ;
	//   .

	{`# Long literal with missing end
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p """abc
def`, "Long literal with missing end (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-06> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-06" ;
	//   rdfs:comment "Long literal with extra quote (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-06.ttl> ;
	//   .

	{`# Long literal with 4"
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p """abc""""@en .`, "Long literal with extra quote (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-string-07> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-string-07" ;
	//   rdfs:comment "Long literal with extra squote (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-string-07.ttl> ;
	//   .

	{`# Long literal with 4'
@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p '''abc''''@en .`, "Long literal with extra squote (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-num-01> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-01" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-num-01.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 123.abc .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-num-02> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-02" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-num-02.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 123e .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-num-03> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-03" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-num-03.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 123abc .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-num-04> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-04" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-num-04.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> 0x123 .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#turtle-syntax-bad-num-05> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-05" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-num-05.ttl> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> +-1 .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#turtle-eval-struct-01> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-eval-struct-01" ;
	//   rdfs:comment "triple with IRIs" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-struct-01.ttl> ;
	//   mf:result    <turtle-eval-struct-01.nt> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-eval-struct-02> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-eval-struct-02" ;
	//   rdfs:comment "triple with IRIs and embedded whitespace" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-struct-02.ttl> ;
	//   mf:result    <turtle-eval-struct-02.nt> ;
	//   .

	{`<http://www.w3.org/2013/TurtleTests/s> 
      <http://www.w3.org/2013/TurtleTests/p1> <http://www.w3.org/2013/TurtleTests/o1> ;
      <http://www.w3.org/2013/TurtleTests/p2> <http://www.w3.org/2013/TurtleTests/o2> ; 
      .`, "", []rdf.Triple{}},

	//<#turtle-subm-01> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-01" ;
	//   rdfs:comment "Blank subject" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-01.ttl> ;
	//   mf:result    <turtle-subm-01.nt> ;
	//   .

	{`@prefix : <#> .
[] :x :y .`, "", []rdf.Triple{}},

	//<#turtle-subm-02> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-02" ;
	//   rdfs:comment "@prefix and qnames" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-02.ttl> ;
	//   mf:result    <turtle-subm-02.nt> ;
	//   .

	{`# Test @prefix and qnames
@prefix :  <http://example.org/base1#> .
@prefix a: <http://example.org/base2#> .
@prefix b: <http://example.org/base3#> .
:a :b :c .
a:a a:b a:c .
:a a:a b:a .`, "", []rdf.Triple{}},

	//<#turtle-subm-03> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-03" ;
	//   rdfs:comment ", operator" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-03.ttl> ;
	//   mf:result    <turtle-subm-03.nt> ;
	//   .

	{`# Test , operator
@prefix : <http://example.org/base#> .
:a :b :c,
      :d,
      :e .`, "", []rdf.Triple{}},

	//<#turtle-subm-04> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-04" ;
	//   rdfs:comment "; operator" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-04.ttl> ;
	//   mf:result    <turtle-subm-04.nt> ;
	//   .

	{`# Test ; operator
@prefix : <http://example.org/base#> .
:a :b :c ;
   :d :e ;
   :f :g .`, "", []rdf.Triple{}},

	//<#turtle-subm-05> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-05" ;
	//   rdfs:comment "empty [] as subject and object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-05.ttl> ;
	//   mf:result    <turtle-subm-05.nt> ;
	//   .

	{`# Test empty [] operator; not allowed as predicate
@prefix : <http://example.org/base#> .
[] :a :b .
:c :d [] .`, "", []rdf.Triple{}},

	//<#turtle-subm-06> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-06" ;
	//   rdfs:comment "non-empty [] as subject and object" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-06.ttl> ;
	//   mf:result    <turtle-subm-06.nt> ;
	//   .

	{`# Test non empty [] operator; not allowed as predicate
@prefix : <http://example.org/base#> .
[ :a :b ] :c :d .
:e :f [ :g :h ] .`, "", []rdf.Triple{}},

	//<#turtle-subm-07> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-07" ;
	//   rdfs:comment "'a' as predicate" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-07.ttl> ;
	//   mf:result    <turtle-subm-07.nt> ;
	//   .

	{`# 'a' only allowed as a predicate
@prefix : <http://example.org/base#> .
:a a :b .`, "", []rdf.Triple{}},

	//<#turtle-subm-08> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-08" ;
	//   rdfs:comment "simple collection" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-08.ttl> ;
	//   mf:result    <turtle-subm-08.nt> ;
	//   .

	{`@prefix : <http://example.org/stuff/1.0/> .
:a :b ( "apple" "banana" ) .
`, "", []rdf.Triple{}},

	//<#turtle-subm-09> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-09" ;
	//   rdfs:comment "empty collection" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-09.ttl> ;
	//   mf:result    <turtle-subm-09.nt> ;
	//   .

	{`@prefix : <http://example.org/stuff/1.0/> .
:a :b ( ) .
`, "", []rdf.Triple{}},

	//<#turtle-subm-10> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-10" ;
	//   rdfs:comment "integer datatyped literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-10.ttl> ;
	//   mf:result    <turtle-subm-10.nt> ;
	//   .

	{`# Test integer datatyped literals using an OWL cardinality constraint
@prefix owl: <http://www.w3.org/2002/07/owl#> .

# based on examples in the OWL Reference

_:hasParent a owl:ObjectProperty .

[] a owl:Restriction ;
  owl:onProperty _:hasParent ;
  owl:maxCardinality 2 .`, "", []rdf.Triple{}},

	//<#turtle-subm-11> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-11" ;
	//   rdfs:comment "decimal integer canonicalization" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-11.ttl> ;
	//   mf:result    <turtle-subm-11.nt> ;
	//   .

	{`<http://example.org/res1> <http://example.org/prop1> 000000 .
<http://example.org/res2> <http://example.org/prop2> 0 .
<http://example.org/res3> <http://example.org/prop3> 000001 .
<http://example.org/res4> <http://example.org/prop4> 2 .
<http://example.org/res5> <http://example.org/prop5> 4 .`, "", []rdf.Triple{}},

	//<#turtle-subm-12> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-12" ;
	//   rdfs:comment "- and _ in names and qnames" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-12.ttl> ;
	//   mf:result    <turtle-subm-12.nt> ;
	//   .

	{`# Tests for - and _ in names, qnames
@prefix ex1: <http://example.org/ex1#> .
@prefix ex-2: <http://example.org/ex2#> .
@prefix ex3_: <http://example.org/ex3#> .
@prefix ex4-: <http://example.org/ex4#> .

ex1:foo-bar ex1:foo_bar "a" .
ex-2:foo-bar ex-2:foo_bar "b" .
ex3_:foo-bar ex3_:foo_bar "c" .
ex4-:foo-bar ex4-:foo_bar "d" .`, "", []rdf.Triple{}},

	//<#turtle-subm-13> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-13" ;
	//   rdfs:comment "tests for rdf:_<numbers> and other qnames starting with _" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-13.ttl> ;
	//   mf:result    <turtle-subm-13.nt> ;
	//   .

	{`# Tests for rdf:_<numbers> and other qnames starting with _
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .
@prefix ex:  <http://example.org/ex#> .
@prefix :    <http://example.org/myprop#> .

ex:foo rdf:_1 "1" .
ex:foo rdf:_2 "2" .
ex:foo :_abc "def" .
ex:foo :_345 "678" .`, "", []rdf.Triple{}},

	//<#turtle-subm-14> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-14" ;
	//   rdfs:comment "bare : allowed" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-14.ttl> ;
	//   mf:result    <turtle-subm-14.nt> ;
	//   .

	{`# Test for : allowed
@prefix :    <http://example.org/ron> .

[] : [] .

: : : .
`, "", []rdf.Triple{}},

	//<#turtle-subm-15> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-15" ;
	//   rdfs:comment "simple long literal" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-15.ttl> ;
	//   mf:result    <turtle-subm-15.nt> ;
	//   .

	{`# Test long literal
@prefix :  <http://example.org/ex#> .
:a :b """a long
	literal
with
newlines""" .`, "", []rdf.Triple{}},

	//<#turtle-subm-16> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-16" ;
	//   rdfs:comment "long literals with escapes" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-16.ttl> ;
	//   mf:result    <turtle-subm-16.nt> ;
	//   .

	{`@prefix : <http://example.org/foo#> .

## \U00015678 is a not a legal codepoint
## :a :b """\nthis \ris a \U00015678long\t
## literal\uABCD
## """ .
## 
## :d :e """\tThis \uABCDis\r \U00015678another\n
## one
## """ .

# \U00015678 is a not a legal codepoint
# \U00012451 in Cuneiform numeric ban 3
:a :b """\nthis \ris a \U00012451long\t
literal\uABCD
""" .

:d :e """\tThis \uABCDis\r \U00012451another\n
one
""" .`, "", []rdf.Triple{}},

	//<#turtle-subm-17> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-17" ;
	//   rdfs:comment "floating point number" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-17.ttl> ;
	//   mf:result    <turtle-subm-17.nt> ;
	//   .

	{`@prefix : <http://example.org/#> .

:a :b  1.0 .
`, "", []rdf.Triple{}},

	//<#turtle-subm-18> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-18" ;
	//   rdfs:comment "empty literals, normal and long variant" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-18.ttl> ;
	//   mf:result    <turtle-subm-18.nt> ;
	//   .

	{`@prefix : <http://example.org/#> .

:a :b "" .

:c :d """""" .
`, "", []rdf.Triple{}},

	//<#turtle-subm-19> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-19" ;
	//   rdfs:comment "positive integer, decimal and doubles" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-19.ttl> ;
	//   mf:result    <turtle-subm-19.nt> ;
	//   .

	{`@prefix : <http://example.org#> .
:a :b 1.0 .
:c :d 1 .
:e :f 1.0e0 .`, "", []rdf.Triple{}},

	//<#turtle-subm-20> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-20" ;
	//   rdfs:comment "negative integer, decimal and doubles" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-20.ttl> ;
	//   mf:result    <turtle-subm-20.nt> ;
	//   .

	{`@prefix : <http://example.org#> .
:a :b -1.0 .
:c :d -1 .
:e :f -1.0e0 .`, "", []rdf.Triple{}},

	//<#turtle-subm-21> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-21" ;
	//   rdfs:comment "long literal ending in double quote" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-21.ttl> ;
	//   mf:result    <turtle-subm-21.nt> ;
	//   .

	{`# Test long literal
@prefix :  <http://example.org/ex#> .
:a :b """John said: "Hello World!\"""" .`, "", []rdf.Triple{}},

	//<#turtle-subm-22> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-22" ;
	//   rdfs:comment "boolean literals" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-22.ttl> ;
	//   mf:result    <turtle-subm-22.nt> ;
	//   .

	{`@prefix : <http://example.org#> .
:a :b true .
:c :d false .`, "", []rdf.Triple{}},

	//<#turtle-subm-23> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-23" ;
	//   rdfs:comment "comments" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-23.ttl> ;
	//   mf:result    <turtle-subm-23.nt> ;
	//   .

	{`# comment test
@prefix : <http://example.org/#> .
:a :b :c . # end of line comment
:d # ignore me
  :e # and me
      :f # and me
        .
:g :h #ignore me
     :i,  # and me
     :j . # and me

:k :l :m ; #ignore me
   :n :o ; # and me
   :p :q . # and me`, "", []rdf.Triple{}},

	//<#turtle-subm-24> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-24" ;
	//   rdfs:comment "no final mewline" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-24.ttl> ;
	//   mf:result    <turtle-subm-24.nt> ;
	//   .

	{`# comment line with no final newline test
@prefix : <http://example.org/#> .
:a :b :c .
#foo`, "", []rdf.Triple{}},

	//<#turtle-subm-25> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-25" ;
	//   rdfs:comment "repeating a @prefix changes pname definition" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-25.ttl> ;
	//   mf:result    <turtle-subm-25.nt> ;
	//   .

	{`@prefix foo: <http://example.org/foo#>  .
@prefix foo: <http://example.org/bar#>  .

foo:blah foo:blah foo:blah .
`, "", []rdf.Triple{}},

	//<#turtle-subm-26> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-26" ;
	//   rdfs:comment "Variations on decimal canonicalization" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-26.ttl> ;
	//   mf:result    <turtle-subm-26.nt> ;
	//   .

	{`<http://example.org/foo> <http://example.org/bar> "2.345"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "1"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "1.0"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "1."^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "1.000000000"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.3"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.234000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.2340000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.23400000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.234000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.2340000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.23400000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.234000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.2340000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.23400000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.234000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.2340000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.23400000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.234000000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.2340000000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "2.23400000000000000000005"^^<http://www.w3.org/2001/XMLSchema#decimal> .
<http://example.org/foo> <http://example.org/bar> "1.2345678901234567890123457890"^^<http://www.w3.org/2001/XMLSchema#decimal> .`, "", []rdf.Triple{}},

	//<#turtle-subm-27> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "turtle-subm-27" ;
	//   rdfs:comment "Repeating @base changes base for relative IRI lookup" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-subm-27.ttl> ;
	//   mf:result    <turtle-subm-27.nt> ;
	//   .

	{`# In-scope base URI is <http://www.w3.org/2013/TurtleTests/turtle-subm-27.ttl> at this point
<a1> <b1> <c1> .
@base <http://example.org/ns/> .
# In-scope base URI is http://example.org/ns/ at this point
<a2> <http://example.org/ns/b2> <c2> .
@base <foo/> .
# In-scope base URI is http://example.org/ns/foo/ at this point
<a3> <b3> <c3> .
@prefix : <bar#> .
:a4 :b4 :c4 .
@prefix : <http://example.org/ns2#> .
:a5 :b5 :c5 .`, "", []rdf.Triple{}},

	//<#turtle-eval-bad-01> rdf:type rdft:TestTurtleNegativeEval ;
	//   mf:name    "turtle-eval-bad-01" ;
	//   rdfs:comment "Bad IRI : good escape, bad charcater (negative evaluation test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-bad-01.ttl> ;
	//   .

	{`# Bad IRI : good escape, bad charcater
<http://www.w3.org/2013/TurtleTests/\u0020> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-eval-bad-02> rdf:type rdft:TestTurtleNegativeEval ;
	//   mf:name    "turtle-eval-bad-02" ;
	//   rdfs:comment "Bad IRI : hex 3C is < (negative evaluation test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-bad-02.ttl> ;
	//   .

	{`# Bad IRI : hex 3C is <
<http://www.w3.org/2013/TurtleTests/\u003C> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-eval-bad-03> rdf:type rdft:TestTurtleNegativeEval ;
	//   mf:name    "turtle-eval-bad-03" ;
	//   rdfs:comment "Bad IRI : hex 3E is  (negative evaluation test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-bad-03.ttl> ;
	//   .

	{`# Bad IRI : hex 3E is >
<http://www.w3.org/2013/TurtleTests/\u003E> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//<#turtle-eval-bad-04> rdf:type rdft:TestTurtleNegativeEval ;
	//   mf:name    "turtle-eval-bad-04" ;
	//   rdfs:comment "Bad IRI : {abc} (negative evaluation test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-eval-bad-04.ttl> ;
	//   .

	{`# Bad IRI
<http://www.w3.org/2013/TurtleTests/{abc}> <http://www.w3.org/2013/TurtleTests/p> <http://www.w3.org/2013/TurtleTests/o> .`, "", []rdf.Triple{}},

	//# tests requested by Jeremy Carroll
	//# http://www.w3.org/2011/rdf-wg/wiki/Turtle_Candidate_Recommendation_Comments#c35
	//<#comment_following_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "comment_following_localName" ;
	//   rdfs:comment "comment following localName" ;
	//   rdft:approval rdft:Proposed ;
	//   mf:action    <comment_following_localName.ttl> ;
	//   mf:result    <IRI_spo.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:o#comment
.`, "", []rdf.Triple{}},

	//<#number_sign_following_localName> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "number_sign_following_localName" ;
	//   rdfs:comment "number sign following localName" ;
	//   rdft:approval rdft:Proposed ;
	//   mf:action    <number_sign_following_localName.ttl> ;
	//   mf:result    <number_sign_following_localName.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:o\#numbersign
.`, "", []rdf.Triple{}},

	//<#comment_following_PNAME_NS> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "comment_following_PNAME_NS" ;
	//   rdfs:comment "comment following PNAME_NS" ;
	//   rdft:approval rdft:Proposed ;
	//   mf:action    <comment_following_PNAME_NS.ttl> ;
	//   mf:result    <comment_following_PNAME_NS.nt> ;
	//   .

	{`@prefix p: <http://a.example/> .
<http://a.example/s> <http://a.example/p> p:#comment
.`, "", []rdf.Triple{}},

	//<#number_sign_following_PNAME_NS> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "number_sign_following_PNAME_NS" ;
	//   rdfs:comment "number sign following PNAME_NS" ;
	//   rdft:approval rdft:Proposed ;
	//   mf:action    <number_sign_following_PNAME_NS.ttl> ;
	//   mf:result    <number_sign_following_PNAME_NS.nt> ;
	//   .

	{`@prefix p: <http://a.example/>.
<http://a.example/s> <http://a.example/p> p:\#numbersign
.`, "", []rdf.Triple{}},

	//# tests from Dave Beckett
	//# http://www.w3.org/2011/rdf-wg/wiki/Turtle_Candidate_Recommendation_Comments#c28
	//<#LITERAL_LONG2_with_REVERSE_SOLIDUS> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "LITERAL_LONG2_with_REVERSE_SOLIDUS" ;
	//   rdfs:comment "REVERSE SOLIDUS at end of LITERAL_LONG2" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <LITERAL_LONG2_with_REVERSE_SOLIDUS.ttl> ;
	//   mf:result    <LITERAL_LONG2_with_REVERSE_SOLIDUS.nt> ;
	//   .

	{`@prefix : <http://example.org/ns#> .

:s :p1 """test-\\""" .`, "", []rdf.Triple{}},

	//<#turtle-syntax-bad-LITERAL2_with_langtag_and_datatype> rdf:type rdft:TestTurtleNegativeSyntax ;
	//   mf:name    "turtle-syntax-bad-num-05" ;
	//   rdfs:comment "Bad number format (negative test)" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <turtle-syntax-bad-LITERAL2_with_langtag_and_datatype.ttl> ;
	//   .

	{`<http://example.org/resource> <http://example.org#pred> "value"@en^^<http://www.w3.org/1999/02/22-rdf-syntax-ns#XMLLiteral> .`, "Bad number format (negative test)", []rdf.Triple{}},

	//<#two_LITERAL_LONG2s> rdf:type rdft:TestTurtleEval ;
	//   mf:name    "two_LITERAL_LONG2s" ;
	//   rdfs:comment "two LITERAL_LONG2s testing quote delimiter overrun" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <two_LITERAL_LONG2s.ttl> ;
	//   mf:result    <two_LITERAL_LONG2s.nt> ;
	//   .

	{`# Test long literal twice to ensure it does not over-quote
@prefix :  <http://example.org/ex#> .
:a :b """first long literal""" .
:c :d """second long literal""" .`, "", []rdf.Triple{}},

	//<#langtagged_LONG_with_subtag> rdf:type rdft:TestTurtleEval ;
	//   mf:name      "langtagged_LONG_with_subtag" ;
	//   rdfs:comment "langtagged LONG with subtag \"\"\"Cheers\"\"\"@en-UK" ;
	//   rdft:approval rdft:Approved ;
	//   mf:action    <langtagged_LONG_with_subtag.ttl> ;
	//   mf:result    <langtagged_LONG_with_subtag.nt> ;
	//   .

	{`# Test long literal with lang tag
@prefix :  <http://example.org/ex#> .
:a :b """Cheers"""@en-UK .`, "", []rdf.Triple{}},

	//# tests from David Robillard
	//# http://www.w3.org/2011/rdf-wg/wiki/Turtle_Candidate_Recommendation_Comments#c21
	//<#turtle-syntax-bad-blank-label-dot-end>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Blank node label must not end in dot" ;
	//	mf:name "turtle-syntax-bad-blank-label-dot-end" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-blank-label-dot-end.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
_:b1. :p :o .`, "Blank node label must not end in dot", []rdf.Triple{}},

	//<#turtle-syntax-bad-number-dot-in-anon>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Dot delimeter may not appear in anonymous nodes" ;
	//	mf:name "turtle-syntax-bad-number-dot-in-anon" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-number-dot-in-anon.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .

:s
	:p [
		:p1 27.
	] .`, "Dot delimeter may not appear in anonymous nodes", []rdf.Triple{}},

	//<#turtle-syntax-bad-ln-dash-start>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Local name must not begin with dash" ;
	//	mf:name "turtle-syntax-bad-ln-dash-start" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-ln-dash-start.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :-o .`, "Local name must not begin with dash", []rdf.Triple{}},

	//<#turtle-syntax-bad-ln-escape>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Bad hex escape in local name" ;
	//	mf:name "turtle-syntax-bad-ln-escape" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-ln-escape.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :o%2 .`, "Bad hex escape in local name", []rdf.Triple{}},

	//<#turtle-syntax-bad-ln-escape-start>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Bad hex escape at start of local name" ;
	//	mf:name "turtle-syntax-bad-ln-escape-start" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-ln-escape-start.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s :p :%2o .`, "Bad hex escape at start of local name", []rdf.Triple{}},

	//<#turtle-syntax-bad-ns-dot-end>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Prefix must not end in dot" ;
	//	mf:name "turtle-syntax-bad-ns-dot-end" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-ns-dot-end.ttl> .

	{`@prefix eg. : <http://www.w3.org/2013/TurtleTests/> .
eg.:s eg.:p eg.:o .`, "Prefix must not end in dot", []rdf.Triple{}},

	//<#turtle-syntax-bad-ns-dot-start>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Prefix must not start with dot" ;
	//	mf:name "turtle-syntax-bad-ns-dot-start" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-ns-dot-start.ttl> .

	{`@prefix .eg : <http://www.w3.org/2013/TurtleTests/> .
.eg:s .eg:p .eg:o .`, "Prefix must not start with dot", []rdf.Triple{}},

	//<#turtle-syntax-bad-missing-ns-dot-end>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Prefix must not end in dot (error in triple, not prefix directive like turtle-syntax-bad-ns-dot-end)" ;
	//	mf:name "turtle-syntax-bad-missing-ns-dot-end" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-missing-ns-dot-end.ttl> .

	{`valid:s valid:p invalid.:o .`, "Prefix must not end in dot (error in triple, not prefix directive like turtle-syntax-bad-ns-dot-end)", []rdf.Triple{}},

	//<#turtle-syntax-bad-missing-ns-dot-start>
	//	rdf:type rdft:TestTurtleNegativeSyntax ;
	//	rdfs:comment "Prefix must not start with dot (error in triple, not prefix directive like turtle-syntax-bad-ns-dot-end)" ;
	//	mf:name "turtle-syntax-bad-missing-ns-dot-start" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-bad-missing-ns-dot-start.ttl> .

	{`.undefined:s .undefined:p .undefined:o .`, "Prefix must not start with dot (error in triple, not prefix directive like turtle-syntax-bad-ns-dot-end)", []rdf.Triple{}},

	//<#turtle-syntax-ln-dots>
	//	rdf:type rdft:TestTurtlePositiveSyntax ;
	//	rdfs:comment "Dots in pname local names" ;
	//	mf:name "turtle-syntax-ln-dots" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-ln-dots.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s.1 :p.1 :o.1 .
:s..2 :p..2 :o..2.
:3.s :3.p :3.`, "", []rdf.Triple{}},

	//<#turtle-syntax-ln-colons>
	//	rdf:type rdft:TestTurtlePositiveSyntax ;
	//	rdfs:comment "Colons in pname local names" ;
	//	mf:name "turtle-syntax-ln-colons" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-ln-colons.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
:s:1 :p:1 :o:1 .
:s::2 :p::2 :o::2 .
:3:s :3:p :3 .
::s ::p ::o .
::s: ::p: ::o: .`, "", []rdf.Triple{}},

	//<#turtle-syntax-ns-dots>
	//	rdf:type rdft:TestTurtlePositiveSyntax ;
	//	rdfs:comment "Dots in namespace names" ;
	//	mf:name "turtle-syntax-ns-dots" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-ns-dots.ttl> .

	{`@prefix e.g: <http://www.w3.org/2013/TurtleTests/> .
e.g:s e.g:p e.g:o .`, "", []rdf.Triple{}},

	//<#turtle-syntax-blank-label>
	//	rdf:type rdft:TestTurtlePositiveSyntax ;
	//	rdfs:comment "Characters allowed in blank node labels" ;
	//	mf:name "turtle-syntax-blank-label" ;
	//        rdft:approval rdft:Approved ;
	//	mf:action <turtle-syntax-blank-label.ttl> .

	{`@prefix : <http://www.w3.org/2013/TurtleTests/> .
_:0b :p :o . # Starts with digit
_:_b :p :o . # Starts with underscore
_:b.0 :p :o . # Contains dot, ends with digit`, "", []rdf.Triple{}},
}