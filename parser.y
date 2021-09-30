%{
package main

import "fmt"

func setScl(l yyLexer, v Scl) {
  l.(*MyHandmadeLexer).scl = v
}
%}

%union{
  result Scl

  ch byte

  any interface{}
  mapstr map[string]interface{}
  list []interface{}
  null *interface{}

  str string
  integer int
  float float64
}

%token <integer> INTEGER
%token <ch> OP CL COLON
%token <any> VALUE NUMBER BOOL WORD
%token <null> NULL

%type <result> full
%type <mapstr> block pair
%type <any> value
%type <list> list values

%start main

%%

main: full
  {
    setScl(yylex, $1)
  }
;

full: WORD block {
	$$ = Scl{
		docName: $1.(string),
		data: $2,
	}
};

block: OP 	   CL { $$ = map[string]interface{}{} }
	|  OP pair CL { $$ = $2 };

pair: WORD COLON value { $$ = map[string]interface{}{$1.(string): $2} }
	| WORD COLON value pair
	{
		$4[$1.(string)] = $3
		$$ = $4
	}

value: INTEGER {$$ = $1} | BOOL {$$ = $1}  | NULL {$$ = nil} | WORD {$$ = $1}  | list {$$ = $1} | block {$$ = $1};

list: '[' values ']'
{
	fmt.Println($2)
	$$ = $2
}

values: value {$$ = []interface{}{$1}}
	| value values
	{
		fmt.Println($2)
		$$ = append($2, $1)
	}