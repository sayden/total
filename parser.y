%{
package total
%}

%union{
	Total total

	value *value
	kv *keyValue
	object object

	char rune
	float float64
	integer int
	string string
	boolean bool
	list values
	nulltype interface{}
	any interface{}
}

// Separators
%token <char> OP CL COLON OB CB NL
%token <string> OLT CLT

// Tokens
%token <value> VALUE
%token <float> FLOAT
%token <integer> INTEGER
%token <string> TEXT WORD
%token <boolean> BOOLEAN
%token <nulltype> NULLTYPE
%token <list> LIST
%token <object> OBJECT

// Types
%type <list> list_values list
%type <value> value full
%type <string> long_text
%type <kv> kv
%type <object> object block
%type <Total> main

%start main

%%

main: WORD COLON full
  {
    yylex.(*myscanner).total = total{
                           		docName: $1,
                           		data: $3,
                           	}
  }
;

full: block	{ $$ = &value{kind:OBJECT, data:$1} }
	| list 	{ $$ = &value{kind:LIST, data:$1} }
	;

block: OP 	     CL  { $$ = newObject() }
	|  OP NL 	 CL  { $$ = newObject() }
	|  OP NL object CL  { $$ = $3 }
	;

object: kv { $$ = object{$1} }
	| object kv { $$ = append($1, $2) }
	;

kv:   WORD COLON value NL { $$ = &keyValue{name: $1, value: $3} }
	;

value: INTEGER { $$ = &value{kind: INTEGER, data: $1} }
	| FLOAT { $$ = &value{kind: FLOAT, data: $1} }
	| NULLTYPE { $$ = &value{kind: NULLTYPE, data: nil} }
	| BOOLEAN { $$ = &value{kind: BOOLEAN, data: $1} }
	| WORD { $$ = &value{kind: WORD, data:$1} }
	| long_text { $$ = &value{ kind: TEXT, data: $1} }
	| list { $$ = &value{kind: LIST, data: $1} }
	| block { $$ = &value{kind: OBJECT, data: $1} }
	;

long_text: OLT CLT 				{ $$ = "" }
	| OLT TEXT CLT 				{ $$ = $2 }
	;

list: OB 			 	CB 		{ $$ = values{} }
	| OB NL 	 	 	CB  	{ $$ = values{} }
	| OB NL list_values NL CB      { $$ = $3 }
	| OB list_values 	CB 		{ $$ = $2 }
	;

list_values: list_values value 	{ $$ = append($1, $2) }
	| value 					{ $$ = values{$1} }
	;