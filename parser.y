%{
package main

import "fmt"

func setTotal(l yyLexer, v Total) {
  l.(*lexer).total = v
}
%}

%union{
  result Total

  ch byte

  any interface{}
  mapstr map[string]interface{}
  list []interface{}
  null *interface{}

  str string
  integer int
  float float64
  Boolean bool
}

%token <ch> OP CL COLON
%token <any> VALUE NUMBER WORD INTEGER
%token <str> OLT CLT
%token <Boolean> TRUE FALSE
%token <null> NULL

%type <result> full
%type <mapstr> block pair
%type <any> value long_text
%type <list> list values

%start main

%%

main: full
  {
    setTotal(yylex, $1)
  }
;

full: WORD block {
	$$ = Total{
		docName: $1.(string),
		data: $2,
	}
};

block: 	   OP 	   CL { $$ = map[string]interface{}{} }
	|  OP pair CL { $$ = $2 };

pair: 	  WORD COLON value { $$ = map[string]interface{}{$1.(string): $3} }
	| WORD COLON value pair
	{
		$4[$1.(string)] = $3
		$$ = $4
	}

value: INTEGER { $$ = $1; }
	| NULL { $$ = nil; }
	| WORD { $$ = $1; }
	| TRUE { $$ = $1; }
	| FALSE { $$ = $1; }
	| list { $$ = $1; }
	| long_text { $$ = $1; }
	| block { $$ = $1; };

list: '[' values ']'
{
	fmt.Println($2)
	$$ = $2
}

long_text: OLT value CLT { $$ = $2; }

values: value {$$ = []interface{}{$1}}
	| value values
	{
		fmt.Println($2)
		$$ = append($2, $1)
	}