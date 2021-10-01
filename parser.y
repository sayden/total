%{
package main

func setTotal(l yyLexer, v Total) {
  l.(*lexer).total = v
}
%}

%union{
	Total Total

	Value *Value
	Pair *Pair
	Object Object

	Char rune
	Float float64
	Integer int
	String string
	Boolean bool
	List Values
	Nulltype *interface{}
}

// Separators
%token <Char> OP CL COLON OB CB
%token <String> OLT CLT

// Tokens
%token <Value> VALUE
%token <Float> FLOAT
%token <Integer> INTEGER
%token <String> TEXT WORD
%token <Boolean> BOOLEAN
%token <Nulltype> NULLTYPE
%token <List> LIST
%token <Object> OBJECT

// Types
%type <List> list_values
%type <Value> value list
%type <String> long_text
%type <Pair> pair
%type <Object> object
%type <Total> full
%type <Object> block

%start main

%%

main: full
  {
    setTotal(yylex, $1)
  }
;

full: WORD block {
	$$ = Total{
		docName: $1,
		data: $2,
	}
};

block: OP 	     CL { $$ = newObject() }
	|  OP object CL { $$ = $2 }
	;

object: pair { $$ = newObject($1) }
	| object pair { $$ = append($1, $2) }
	;

pair: WORD COLON value { $$ = &Pair{name: $1, value: $3} }
	| WORD COLON block { $$ = &Pair{name: $1, value: &Value{kind: OBJECT, data: $3}} }
	;

value: INTEGER { $$ = &Value{kind: INTEGER, data: $1} }
	| FLOAT { $$ = &Value{kind: FLOAT, data: $1} }
	| NULLTYPE { $$ = &Value{kind: NULLTYPE, data: nil} }
	| WORD { $$ = &Value{kind: WORD, data:$1 } }
	| BOOLEAN { $$ = &Value{kind: BOOLEAN, data: $1 } }
	| long_text { $$ = &Value{ kind: TEXT, data: $1 } }
	| list {$$ = &Value {kind: LIST, data: $1 }}
	;


long_text: OLT CLT { $$ = "" }
	| OLT value CLT
	;

list: OB CB { $$ = &Value{kind: LIST, data: &Values{} } }
	| OB list_values CB { $$ = &Value{kind: LIST, data: $2 } }
	;

list_values: value { $$ = Values{$1} }
	| value list_values
	{
		$$ = append($2, $1	)
	}
	;