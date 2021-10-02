%{
package total
%}

%union{
	Total Total

	value *value
	kv *Kv
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
%token <Char> OP CL COLON OB CB
%token <String> OLT CLT

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
%type <list> list_values
%type <value> value
%type <list> list
%type <string> long_text
%type <kv> kv
%type <object> object block
%type <any> full
%type <Total> main

%start main

%%

main: WORD full
  {
    yylex.(*lexer).total = Total{
                           		docName: $1,
                           		data: $2,
                           	}
  }
;

full: block	{ $$ = &value{kind:OBJECT, data:$1} }
	| list 	{ $$ = &value{kind:LIST, data:$1} };

block: OP 	     CL  { $$ = newObject() }
	|  OP object CL { $$ = $2 }
	;

object: kv { $$ = object{$1} }
	| object kv { $$ = append($1, $2) }
	;

kv:   WORD COLON value { $$ = &Kv{name: $1, value: $3} }
	| WORD       block { $$ = &Kv{name: $1, value: &value{kind: OBJECT, data: $2}} }
	;

value: INTEGER { $$ = &value{kind: INTEGER, data: $1} }
	| FLOAT { $$ = &value{kind: FLOAT, data: $1} }
	| NULLTYPE { $$ = &value{kind: NULLTYPE, data: nil} }
	| WORD { $$ = &value{kind: WORD, data:$1} }
	| BOOLEAN { $$ = &value{kind: BOOLEAN, data: $1} }
	| long_text { $$ = &value{ kind: TEXT, data: $1} }
	| list { $$ = &value{kind: LIST, data: $1} }
	| block { $$ = &value{kind: OBJECT, data: $1} }
	;


long_text: OLT CLT 				{ $$ = "" }
	| OLT TEXT CLT 				{ $$ = $2 }
	;

list: OB 			 CB 		{ $$ = values{} }
	| OB list_values CB 		{ $$ = $2 }
	;

list_values: list_values value 	{ $$ = append($1, $2) }
	| value 					{ $$ = values{$1} }
	;