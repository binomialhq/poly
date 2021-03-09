# Poly Language Specification

## Copyright Notice

Copyright (C) 2021 Binomial. All Rights Reserved.

## Abstract

Poly is an Interface Description Language (IDL) proposed as a simpler and less
verbose alternative to the likes of OpenAPI and RAML. It serves the purpose of
describing Application Programming Interfaces (API) in a way that is concise
and yet complete to enable the modeling of APIs and code generation. The main
rationale is to avoid the excessively verbose constructs of other alternatives,
which often rely on data serialization formats, such as YAML and JSON, in order
to achieve the same effect, if not more.

## Purpose

Poly is an Interface Description Language (IDL) designed with the intent of
providing a simple and effective way of describing the interfaces of web
services and components, with focus on the Hypertext Transfer Protocol (HTTP).
The main motivation is driven by the excessively verbose nature of other
current-practice alternatives, namely [OpenAPI](https://www.openapis.org/)
and [RAML](https://raml.org/). Instead, Poly is inspired on succint syntatical
constructs that are already familiar to developers. This enables the use of
operators, scope delimiters, documentation engines and other constructs that
are either not possible or not natural with typical data serialization formats.

Similarly to its counterparts, Poly focuses on HTTP, yet additionally
considering other alternative protocol schemes by design. For example, OpenAPI
and RAML do not support the association of indexes with models and procedure
parameters, both typical when using binary encodings. Thefore, the likes of
[protobuf](https://github.com/protocolbuffers/protobuf) and
[Thrift](https://thrift.apache.org/) cannot be described by those alternatives.
Poly, instead, ignores any declarations related to transport and encoding
schemes, making it agnostic to representational concerns and networking stacks.
This enables the likes of content negotiation, or alternative protocols. Neither
of those things is possible with OpenAPI or RAML.

It must be said that Poly supports a disjoint (yet overlapping) set of features
from those of RAML and OpenAPI, in the sense that neither is strictly a superset
of the other two. In other words, Poly does not support all features offered by
OpenAPI or RAML, nor do those support all features proposed by Poly. The former
is explained by observing that not all constructs offered by the alternatives
directly relate to API specifications, such as `servers` and `baseUri`. These
markers, provided by OpenAPI and RAML, respectively, mix, in our opinion,
concepts that should not be mixed, that of an API specification and the location
of some specific implementation of it. Poly proposes a strict separation of
these concerns. The fact that Poly implements features that are not available
in the alternatives, is rather explained by its consideration for a wider
variety of scenarios, as well as purpose-specific language constructs.

From a syntatic point of view, Poly more closely resembles protobuf and Apache
Thrift. In fact, many linguistic constructs are inspired by those. The key
differentiation factor is that Poly does not implement any form of transport or
encoding scheme, as those do. This makes Poly portable, in the sense that it is
expressive enough to represent services that operate over any type of network
stack. With that in mind, Poly can be used as a tool to mediate communication
links between services that make use of differing stacks, something that neither
of the alternatives can claim. Instead, the language focuses on specifing API
endpoints, and not much more. The declarative representation of communication
stacks, however, is not neglectable; rather, Poly suggests that those are
represented separatelly, and the use of different tools for different purposes.

Finally, because Poly does not implement an encoding scheme of its own, it fits
better within the context of OpenAPI and RAML than it does protobuf or Thrift.
All of them are used for code generation, being differentiated only by the fact
that the latter implement custom encoding schemes. Poly can be used for code
generation, but also for service matchmaking, the main purpose of its inception.
Poly, thus, can enable the communication of services that implement APIs that
would be compatible, had they not implemented incompatible encodings and other
varying transport-related concerns. If two services, for example, implement
compatible APIs but rely on different encoding schemes, Poly can be used as a
translator to mediate the communication, easily enabling an integration where
it would otherwise be very difficult.

## Philosophy

Poly was designed with several principals in mind:

1. Code brevity, clarity, and elegance, in the sense that declarations are as
short as they can be without loss of expressivity, but also while maintaining
human-friendly characteristics;

2. Familiarity, expressed as a design principle based on linguistic constructs
that are already familiar to developers;

3. Agnosticism, meaning that the declarative language does not bind an API to
any particular language, protocol, or encoding scheme;

4. Oriented to best-practices, by introducing language constructs that expand
into well-composed API declarations (e.g. the `rest` keyword);

5. Expressiveness, in its ability to (1) represent different architecural styles
(e.g. RPC vs REST) and (2) a wide variety of API design formalisms;

6. Machine-readability, as the ability to be parsed and processed by computer
programs;

7. Determinism, meaning that all declarative constructs can be interpreted
without ambiguity.

## Origins

The name "Poly" comes from "polyglot", in representation of its cross-language
and cross-stack design principles. Poly's main purpose is to enable services
with incompatible stacks to communicate, by implementing the necessary means
to mediate integrations. As with other tools of its kind, Poly can be used to
generate code in a wide diversity of languages. For these reasons, we like to
say that Poly "thinks" and "talks" in many languages, in reference to code
generation and protocol stack translation, respectively.

## Purpose of This Document

This document is meant as a formal description of the Poly language, by means
of a context-free grammar expressed in Augmented Backus-Naur Form (ABNF). It is
not meant as a tutorial or gentle introduction to Poly, neither is it meant as
a formal Request For Comment (RFC) or similar documents. This specification is
not expected to formally define all requirements, in contradiction to the
purpose of such documents, and therefore will need future revisions to further
improve on its formalisms. We do not exclude the idea of turning it into a
formal RFC in the future and submitting a standardization proposal.

## Notational Conventions

All grammatical constructs specified in this documented are formally described
resorting to an Augmented Backus-Naur Form (ABNF) syntax, as defined by
[RFC 5234](https://tools.ietf.org/html/rfc5234). Such constructs appear in
lowercase.

Terminal symbols are described by regular expressions with Unicode support, as
described by [TR18](https://www.unicode.org/reports/tr18/). Terminals are
expressed in all uppercase.

We do not, however, fully compromise to following them strictly in this version.
Expect minor differences.

## Specification

### Case Sensitivity

Poly is case-sensitive.

### Basic Rules

The following rules are used in this specification to describe reusable parsing
constructs. No specific form of encoding is assumed.

```
LETTER      = [a-zA-Z]      ; US-ASCII letters
DECIMAL     = [0-9]         ; US-ASCII numbrers 0 through 9
SPACE       = \s            ; US-ASCII space [\u0085\t\r\n\v\f]
HEXDIGIT    = [0-9A-F]      ; Hexadecimal digit
ULETTER     = \p{L}         ; Any Unicode L-category code point
UDECIMAL    = \p{Nd}        ; Any Unicode Nd-category code point
USPACE      = \p{Z}         ; Any Unicode Z-category code point
```

Implementations may decide to support L-category characters in combination with
M-category markers. Although the `ULETTER` declaration does recognize diacritics
when encoded together as a single code point, it also fails to recognize them
if they are encoded separately. For example, `à` can be encoded as `U+00E0`
or as `U+0061 U+0300`. The former encodes `à` as a single code point, while the
latter encodes `a` followed by the accent modifier `U+0300`. In order to support
modifiers, the engine may replace the `ULETTER` declaration above:

```
ULETTER     = \p{L}\p{M}*   ; Any L-category code point and M-category mark
```

As will be made clear throughout the specification, some names can also include
a set of extra characters. This set includes the Unicode characters `U+005F`
(Low Line, more oftenly known as "underscore") and `U+2010` (Hyphen), enabling
symbols such as `User-Agent` and `_tcp` to be defined.

```
EXTRA       = [\u005F\u2010]
```

Several non-terminals are also defined for reuse throughout this specification.

```
letter      = ULETTER / EXTRA
decimal     = UDECIMAL
word        = letter *(letter / decimal)
```

In case the engine is not meant to support Unicode, for whatever reason, some of
the rules above may be replaced with their US-ASCII counterparts:

```
letter      = LETTER / EXTRA
decimal     = DECIMAL
```

An `empty` non-terminal is also defined for the sake of defining empty rules
that do not consume any input:

```
empty   = ""            ; Empty string, does not consume input
```

### Treatment of Space

```
space   = +SPACE
```

Space is generally ignored by the engine unless indicated otherwise. However,
it's notable that "space" is defined as "US-ASCII space", and not the broader
`USPACE` declaration with Unicode support. In fact, it's recommended that the
engine triggers warnings for code points that match `USPACE` but not `SPACE`,
when in a context where space is to be ignored.

### Symbols

```
DQUOTE      = ["]                   ; Single U+0022 (")
NOTDQUOTE   = [^"\\]                ; Anything but U+0022 (") or U+005C (\)
ESCAPED     = "\\" .                ; U+005C (\) followed by one of anything
```

```
symbol          = symbol-simple / symbol-quoted
symbol-simple   = letter *( letter / decimal )
symbol-quoted   = DQUOTE *( *NOTDQUOTE [ ESCAPED ] ) DQUOTE
```

A Symbol is a name that appears as the identifier of some declaration. For
example, in the declaration `model Pet`, `Pet` is the symbol, or name, that can
be used to refer to the model being declared. Several constructs accept symbols
as a means of identification.

Symbols may also be quoted, in which case they can also be refered to as Quoted
Symbols. Quoted Symbols enable characters that do not fit in the Symbol
specification to be used. For example, the declaration `model "ASN.1"` is valid,
and declares a symbol called `ASN.1`, discarding the quotes.

Any reference to a Quoted Symbol must be provided exactly as the original
declaration, including Space. This means that Quoted Symbols may be hard to
read, and thus their usage is discouraged, unless under extraordinary and
justified circumstances. How to handle multi-line strings is an open issue.

### Scopes

Several Poly declarations introduce new symbols to a symbol space. How and where
these symbols are declared bears significance as to where the symbols are made
visible. The symbol space is composed of a stack of scopes, each of which
defining its own symbols. For the purposes of name resolution, Poly defines two
default scopes and a system for introducing new ones.

The first scope in the stack is the bottommost and is introduced by default. It
is called the Gobal Space. All declarations that appear outside any other
declaration are inserted here. In the following example, the `Pet` model is
being declared in the Global Space, but `name` and `tag` are not.

```
model Pet {
    1: string name,
    2: string tag
}
```

Declarations may further push a scope onto the stack, effectively overloading
all scopes bellow it. Therefore, in the example above, the `model` declaration
creates a new scope that will contain only the symbols therein defined. When
this push happens, that scope becomes the Active Space.

The symbol lookup system, then, begins by looking for a symbol in the Active
Space. Only if one is not found there, does the implementation look for the
same symbol in the scope that preceeds it in the stack. Therefore, scopes are
searched in reverse depth-order, starting with the topmost entry. This creates
an overloading effect.

When a name is not found in any scope as a result of a lookup, the engine must
fail with an error. Symbols that are declared and yet already found in the same
scope must also fail with an error, with an indication that the declaration is
duplicated.

Finally, the Primitive Space is one in which symbols are always visible and
cannot be overloaded. This is the ideal space to declare names that are reserved
by the system. If the name of a primitive is redeclared, in any context, the
system must fail with an error. If the name of a primitive is referenced, then
it can only refer to the primitive with that name. This can be implemented by
always starting the lookup and symbol declaration processes with the Primitive
Space.

### Value Categories

Poly defines two value categories, according to whether they produce new symbols
or not. An l-value, thus, refers to an expression that persists by introducing
some symbol into the Active Scope. On the other hand, r-values do not introduce
any new symbols to any scope, and thus do not persist beyond the point where
they are declared.

It's notable that each construct, small as it might be, produces values of
either one or the other category. Consider three constructs in the example
below. The quoted symbol "ASN.1" and the scope delimited by curly braces are
r-values since they do not introduce any symbols. This may seem counterintuitive
because the quoted symbol is a Symbol, but it is in fact the `model` expression
that declares it, and thus the `model` expression is the one that is an l-value.

```
model "ASN.1" {

}
```

### Comments

```
ENDLINE             = .*$           ; Anything up to EOL
MULTILINE           = (?!\*\/)*     ; Anything except "*/", multiline
```

```
comment             = comment-single / comment-multiple
comment-single      = "//" ENDLINE
comment-multiple    = "/*" MULTILINE "*/"
```

Poly supports C-style comments in both single-line and multi-line forms. That
is, comments can be started by a double forward slash (`//`) and span until the
end of the line, or by the multi-line start sequence (`/*`) and span multiple
lines until the multi-line end sequence (`*/`). In either case, the span area
is ignored by the interpreter. Nested comments are not allowed.

### Object Mappings

Throughout the document, some sections refer to Top-Level Declarations and
Key-Value Mappings. These concepts are related to whether objects are inline
when passed to or returned by a procedure, meaning that they are either encoded
as a key-value mapping, or as a top-level object.

Key-Value Mappings are encoded as keys of some top-level object map. For
example, given a model `Pet` with a field `name`, and the declaration `Pet pet`,
the `pet` value is encoded as a key-value mapping for that top-level object;
that is, if the object was to be encoded in JSON, it would look like this:

```
{ "pet": { "name": "Poly" } }
```

By contrast, when referenced as a Top-Level Object, the declaration is not named
and only the type (e.g. `Pet`) is used. In that case, the object is encoded
according to the following, if it were to be encoded in JSON as well:

```
{ "name": "Poly" }
```

### Syntax

```
syntax      = "syntax" +decimal
```

The Syntax declaration indicates the version of the Poly IDL that the file is
using. The practical consequence is having the engine select the appropriate
parser and rule set to process and interpret the input. This declaration is
required and must be the first non-empty and non-comment line in the file. It
must not appear more than once per file either.

There's only one version number currently supported, and that's `0` (zero). This
version does not guarantee the continuity of any of its constructs, and thus
future revisions might not be compatible. The following table summarizes version
numbers, for future reference:

| API Version | Description                            |
|-------------|----------------------------------------|
| 0           | The syntax described in this document. |

The version number may have future implications with regard to importing other
files and external declarations. When a file imports another, their versions can
only compare in one of three ways:

* The imported file declares a higher syntax version. In this case, if the
engine supports the version, it should be interpreted according to that. If not,
the engine must fail with an error (e.g. "unsupported version");

* The imported file declares a lower syntax version. In this case, the
interpreter must support the version, the implication being that each engine
is capable of interpreting all of the previous versions;

* The versions match, in which case no issue exists.

It’s notable that versions need not have backwards compatibility, other than
this declaration being the first on the file. That’s a perpetual requirement.

### Primitives

```
native              = "native" native-decl
native-decl         = word / native-array-decl
native-array-decl   = "[" word "," "..." "]"
```

Primitives are data types that must be natively supported by the engine. A
primitive is declared with the `native` keyword, if declared at all. If some
declared primitive is not understood by the engine, the engine must fail the
process with an error. Although primitives can be declared in any file (e.g.
they are syntactically valid), declaring primitives that are not supported by
the engine will cause general failure, and thus doing so is not recommended.

Primitive declarations introduce symbols in the Primitive Space. If the symbol
already exists in the Primitive Space (e.g the declaration is repeated), the
engine must fail with an error.

Primitives are not declared on every file, but on a standard declaration space
that is included implicitly by the engine. How that space is declared is up to
the implementation, and interpreters may decide to have a file physically
allocated on persistent memory, but are not required to do so. However, whatever
declarations are found there should be effective; that is, if a declaration is
not present, then the engine must not recognize it.

The supported primitives are summarized as follows:

| Primitive    | Description                                 |
|--------------|---------------------------------------------|
| i32, int32   | Signed 32-bit integer.                      |
| i64, int64   | Signed 64-bit integer.                      |
| float        | Floating point value.                       |
| double       | Floating point value with double precision. |
| string       | To be determined.                           |
| wstring      | To be determined.                           |
| void         | No type.                                    |
| [array, ...] | An heterogeneous collection of entities.    |

Unsigned integers are not supported because (1) Poly does not introduce any
concerns regarding validation and (2) those are not supported by every platform.

The string and wstring types are still under analysis. The general idea is to
support Unicode, but neither of these types is coupled with that concept. The
main idea is to make `string` an abstract metatype that can represent any
collection of characters.

The name `array` in the standard declaration doesn't mean anything, and is
not a reserved keyword. `[T, ...]`, for example, is also valid and the two are
semantically equivalent. Engines must not make a distinction of the choice of
this name.

Arrays are heterogenous but only if declared that way. A `[Pet]` declaration,
for example, denotes "an array of `Pet` objects", but nothing else. The ellipsis
symbol signifies a repetition, meaning that declarations such as `[Pet, Error]`
are possible, denoting a mixed array of `Pet` and `Error` objects.

The `void` primitive cannot be used to declare Symbols. Rather, it's sometimes
used to override some Reference, indicating that the Reference does not apply
in a certain context. Attempts to declare a `void` type must result in a
semantic error.

### References

```
reference       =  [ "." ] reference-decl
reference-decl  =  symbol
reference-decl  =/ symbol "." reference-decl
```

References consist of a way to refer to symbols whatever the scope that declares
them. For example, if a model in the Global Space `Pet` declares a field `name`,
the reference `Pet.name` expands into the same declarative construct. Given that
the same semantic entity is used, the referrer inherits all properties of the
reference, including type, modifiers, and any other attributes. This enables
the reuse of declarations without having to explicitly rewrite all the details.

The lookup for the first component of a reference follows the same process as
described before. This means that the symbol is looked up in the Active Space
first, and then successfully traced by the other scopes on the stack. This also
applies if that's the only symbol in the list. For symbols that appear on the
list after that, the lookup does not backtrace, and instead only looks for the
symbol in the scope that is referenced up to that point. Therefore, the lookup
`Service.Model.Field` looks for `Service` using the previous method and
`Model.Field` only within the context of that first scope, if any.

It's notable that this scheme may create ambiguities. The example that follows
shows an ambiguous lookup for `PetStore`, which happens because the symbol
exists both in the Active Space and the Global Space scopes. This means that
both declarations can be referenced in the same way, hence the ambiguity. In
this case, the `Veterinary.PetStore` construct is the right choice for the
engine, since it lives closer to the scope where the symbol is being referenced.

```
service PetStore {

    model Pet {
        1: int32 id,
        2: string name
    }
}

service Veterinary {

    model PetStore {
        1: int32 storeId
    }

    model PetOwner {
        1: PetStore.Pet             // Ambiguous
    }
}
```

In order to refer to the `PetStore` declaration in the Global Space, the entire
expression may be preceeded by the name resolution operator `.`. That is, in
the example above, the correct way to refer to `PetStore.Pet` from the Global
Space would be `.PetStore.Pet`. This forces the engine to begin the lookup from
the bottommost scope. Thus, a dot at the beginning triggers full backtracking
before the lookup is performed, down to the Global Space.

### Type

```
type                = type-simple / type-array
type-simple         = reference
type-array          = "[" type-array-list "]"
type-array-list     = type
type-array-list     = type "," type-array-list
```

The Type construct is used in declarations that refer to a type that already
exists, be that a Symbol or a Primitive. The Type construct must not be used
in the context of a Symbol for a declaration that is introducing that Symbol.

Types include arrays, which can be nested and refer to Quoted Symbols. The
declaration `[.PetStore.Pet, [int32], "ASN.1"]`, for example, illustrates all
three scenarios.

### Annotations

```
annotation  = +decimal ":"
```

Annotations associate expressions with a given numeric identifier, and can be
used to create static indexes of elements. Annotations must be valid integers,
meaning that leading zeros are removed, if they exist. Since other engines do
not support zero as an index value, Poly doesn't either, and thus annotations
consisting of all zeros constitute a semantic error. Annotations are, therefore,
one-indexed.

The annotations do not necessarily dictate order, in the sense that the
identifiers need not be sorted. In fact, the annotations are not processed 
Poly, since they are meant for binary encodings.

Protobuf, for example, requires these in model fields, while Thrift also
requires them in parameter declarations. One requirement that these two
implementations share is the fact that assigning an annotation is a permanent
action, in the sense that the number becomes reserved. Future revisions of the
API must not reutilize them, and instead add new ones in case of need, while
commenting out the ones that become deprecated. For the sake of compatibility,
Poly adheres to this principle.

When annotations appear in lists, specifically Field Lists and Parameter Lists,
a mix of annotated and non-annotated entities is not allowed. Although such
constructs are valid from a syntatic perspective, such scenario must not pass
semantic validation, otherwise provoking an error.

Currently, annotations only support the Hindu-Arabic numeral system. This is
justified by one's own lack of knowledge of other numeral systems, and their
use in a context where the numerals bear meaning. Future revisions should
change this behaviour to span other numeral systems as well.

### Modifiers

```
modifier            = word / "?" / "!"
modifier-list       = modifier modifier-list / empty
```

Modifiers can be `word`s, question marks (`?`), or exclamation marks (`!`), and
cause some form of semantic change to the declaration. The `deprecated` keyword,
for example, causes the field to be flagged as being deprecated. The question
(`?`) and the exclamation (`!`) marks are special sugary constructs that mean
`optional` and `required`, respectively. The following table summarizes all
recognized modifiers:

| Modifier   | Description                        |
|------------|------------------------------------|
| sensitive  | Flags a field as sensitive.        |
| deprecated | Flags a field as being deprecated. |
| required   | Flags a field as required.         |
| optional   | Flags a field as optional.         |
| !          | Alias for "required".              |
| ?          | Alias for "optional".              |

The `sensitive` keyword tells the interpreter to flag the field as being of a
sensitive nature, applying to the likes of passwords. This will instruct UI
generators to use an obfuscation control, for example. More importantly, it can
tell the deployment enviroment not to log this attribute, or not to broadcast
it on the network.

The `deprecated` modifier can be used by code generators and other engines to
bear semantic value in varying and unspecified ways. In the context of Poly, in
a more direct sense, it serves documentation purposes.

The `required` and `optional` modifiers are mutually exclusively, as are their
alias counterparts. The semantic value of the action is to flag an entity as
being of a certain nature, and that nature cannot contradict itself. Therefore,
`required` and `optional`, in their literal or abreviated forms, must not appear
in the same expression, otherwise causing the system to error.

### Fields

```
field       =  type symbol modifier-list / reference
```

A Field declares an attribute of some construct in which it appears, meaning
that the purpose of a Field expression is only made clear when found within
the scope of some other declaration. They are, however, l-value expressions,
and introduce the Symbol they declare in the context in which they appear.

Fields have a Type, a name, and zero or more Modifiers. The Type must correspond
to either a Primitive or a Model, declared somewhere in the symbol space,
causing an error that not being the case. These declarations, thus, associate a
data type and modifiers with a given name, the same name that is to be used for
symbol lookups.

By default, Fields are required. The explicit presence of the `required`
modifier is redundant, unless some option is passed to the interpreter that
changes the default (e.g. command line options). In that sense, omitting the
declaration doesn't mean "required", but rather "whatever is the currently
active default". For the most part, thought, that is supposed to be `required`.

For the sake of readability, it’s recommended that the special exclamation and
question marks appear before any other declaration and right next to the field
name. Not doing so is not syntactically invalid (nor should it be), but it does
make it harder to read. Consider the following example, and whether its clear
that the question mark refers to `password`, not `sensitive`:

```
string password sensitive ?
```

### Field Lists

```
field-list          = "{" field-list-items "}" 
field-list-items    =  field-list-decl
field-list-items    =/ field-list-decl ","
field-list-items    =/ field-list-decl "," field-list-items
field-list-decl     =  [ annotation ] field
```

A Field List is a non-empty sequence of annotated Fields separated by commas,
and optionally terminated by an additional comma (for convenience). It's also
possible that the annotations are not present at all, with the implication that
the declaration might not be supported by some tools. Field Lists push a new
Scope into the symbol space.

When a Reference appears instead of a `type` and `symbol` pair, the Reference
must refer to a Field declaration. It cannot refer to a Parameter, Model, or
other types of constructs, otherwise provoking an engine error.

### Location

```
location    = "body" / "path" / "query" / "header" / "cookie"
```

The Location attribute refers to where a value is expected to appear in the
context of the HTTP protocol. For example, a `header` qualifier indicates that
the value is expected to be passed in the headers. The following table lists
all available location modifiers:

| Location | Description                                                 |
|----------|-------------------------------------------------------------|
| body     | The parameter is passed in the body of the request.         |
| path     | The parameter is passed in the URL, in the path component.  |
| query    | The parameter is passed in the URL, in the query component. |
| header   | The parameter is passed as a header.                        |
| cookie   | The parameter is passed as a cookie.                        |

`path` applies to parameterized endpoints, where the value of the attribute
will replace some Path Template variable. For example, if an expression declares
some URL endpoint such as `items/{itemId}`, the Path Template `itemId` would be
replaced by some value `itemId` with a `path` Location attribute.

The semantic value of the Location attribute under the context of non-HTTP
protocols is yet to be determined.

### Parameters

```
parameter   =  [ location ] field
```

Parameters extend on Fields by specifying the Location in which the field is to
appear. When omitted, `body` is assumed as the default Location, saving some
typing. The exception to this is when the field is declared by Reference, and
that Reference is to a Parameter. In that case, the Location is whatever is
specified in the original Parameter construct, and it cannot be overriden. In
fact, any attempt to so results in an error. Locations can still be specified if
the Reference is to a Field.

### Parameter Lists

```
parameter-list          =  "{" parameter-list-items "}"
parameter-list-items    =  parameter-list-decl
parameter-list-items    =/ parameter-list-decl ","
parameter-list-items    =/ parameter-list-decl "," parameter-list-items
parameter-list-decl     =/ [ annotation ] parameter
```

Parameter Lists are similar to Field Lists, except that each element of the list
is a Parameter and that when References are used, they must refer to a Field or
Parameter declaration. Referring to any other type of declaration must result in
an error. All other considerations regarding Field Lists apply, including that
the construct pushes a new stack into the symbol space.

### Prototyping

```
prototype   = "from" reference
```

When inheriting from a prototype, declarations implement all symbols declared
by that prototype. Depending on the situation, that can be done explicitly or
implicitly. Multiple inheritance is not supported, motivated by the Diamond
Problem and the lack of support in several languages.

When prototyping and using annotated declarations, the annotations are not
inherited, and all fields must be explicitly listed in order to avoid conflicts.
The example below, of Model prototyping, illustrates what would happen if
annotations were to be inherited, creating a conflict that the engine would not
be capable of resolving.

```
model NewPet {
    1: string name required,
    2: string tag
}

model Pet from NewPet {
    1: int32 id required        // Conflict! Annotations match
}
```

For that reason, all parent declarations (e.g. fields and parameters) must be
explicitly referenced by the child, reassigning the annotations, otherwise
incurring in the penalty of an error being indicated by the engine. The example
below shows an error due to the annotations being inherited but not explicitly
reassigned, changing the nature and the location of the error, in comparison
to the previous example.

```
model NewPet {
    1: string name required,
    2: string tag
}

model Pet from NewPet {             // Error! "name" and "tag" are not declared
    1: int32 id required
}
```

The example that follows illustrates the correct way for prototyping with
annotations. Notably, all fields are listed, but only declared as References.
Types, modifiers, and location attributes are inherited implicitly, and cannot
be overloaded. Thus, the child Field or Parameter holds all properties of the
corresponding parent, with the exception of the annotation itself. This enforces
that changes made to the parent also reflect on the child, including `sensitive`
and `deprecated` modifiers.

```
model NewPet {
    1: string name required,
    2: string tag required
}

model Pet from NewPet {
    1: int32 id required,
    2: NewPet.name,
    3: NewPet.tag
}
```

From a code generation perspective, generators may implement these constructs
using strict inheritance, since the declarations respect the "is-a" relationship
and, thus, the Liskov Substitution Principle. In fact, any context in which a
parent is declared and a child is given in its place, the extra Fields or
Parameters declared by the child should be discarded, and the entity processed
according to the substitution principles. Not doing so is a violation of the
declaration.

Since annotations can be omitted from Field and Parameter lists, and for the
sake of convenience, those that are inherited may be omitted if the child does
not use annotations. It still holds that such declarations are not compatible
with certain processors, and thus this approach is not recommended. It is,
however, convenient.

```
model NewPet {
    string name required,
    string tag required
}

model Pet from NewPet {             // Valid, since fields are not tagged
    int32 id required,
}
```

In fact, this is valid regardless of whether the parent uses annotations. If the
parent does use annotations, then the child simply ignores them, even relaxing
the need for them to be declared.

```
model NewPet {
    1: string name required,
    2: string tag required
}

model Pet from NewPet {             // Valid, annotations are ignored
    int32 id required
}
```

In the reverse scenario, with the parent not using annotations, the child may
declare them explicitly, since they would be reassigned anyway.

```
model NewPet {
    string name required,
    string tag required
}

model Pet from NewPet {             // Valid, all fields are annotated
    1: int32 id required,
    2: NewPet.name,
    3: NewPet.tag
}
```

Another key aspect of prototype inheritance is that if one parent fields is
declared explicitly, then they all have to be. This works as a semantic
guarantee that all fields are properly handled, and that updates to the parent
are visibly reflected on the child, rather than silently absorved.

None of this rationale applies to Service declarations, however, since they
don't implement annotations at all. For that reason, Services don't need to
explicitly declare inherited constructs, as the other prototypable entities do. 

The main difference between these entities is that `model`, `in`, and `out` do
not inherit the scope of their parents, having the declarations copied instead.
If the declarations are not annotated, it just so happens that the copy can be
automatic. `service`, on the other hand, inherits a copy of the parent's scope.
To be precise, when `model`, `in`, or `out` prototyping occurs, the scope must
be redeclared with references to the same symbols, even if automatically, while
`service` keeps a copy of the parent scope as a scope of its own.

The practical consequence of this is that Services cannot overload declarations,
including methods. This happens because Poly enforces that child Services be
merely an extension of their parents, without any disregard for the contract
that they promote.

One final observation is that no declaration can inherit a prototype from a
different type of declaration. That is, Models inherit from Models, Services
inherit from Services, and so on, and the prototype chain must not be mixed.
The only exception to this rule is with Groups, which can only be prototyped
from Templates, not other Groups.

### Models

```
model           = "model" model-decl
model-decl      = symbol [ prototype ] field-list
```

A Model declaration defines a complex data type constructed from Primitive types
and other Models. Models can be used as input and ouput for procedures, as well
as other situations where a data type is appropriate.

When generating code, engines may generate classes and other data structures
from models, as they are meant to represent structures that hold data. Models
are not used to represent any form of encoding, and are rather an abstract
representation of data graphs.

### Input

```
input                   =  "in" input-decl
input-decl              =  symbol input-decl-options
input-decl-options      =  [ annotation ] reference [ prototype ] [ parameter-list ]
input-decl-options      =/ prototype [ parameter-list ]
input-decl-options      =/ parameter-list
```

The `in` keyword is used to declare constructs that are used as input to
Procedures. The expression accepts a Reference that is used as a Top-Level
Declaration and a Parameter List that specifies a Key-Value Mapping. The
Top-Level declaration does not accept a Location specifier, and always refers
to `body`.

It's notable that both the Top-Level Declaration Reference and the Parameter
List can be specified together, in which case the Parameter List adds to the
`body` declaration introduced by the Top-Level Declaration. When that happens,
any Parameter with the `body` modifier is illegal in the context of the
Parameter List, since the body is already explicitly defined. In the example
below, the `Authorization` declaration further specifies that the Input
declaration is to accept an `Authorization` header, besides the `Pet` argument
that is already passed in the body. The `owner_id` declaration, however, is
illegal, since it might create several symbol conflicts.

```
in PetIn Pet {
    header string Authorization,
    body int32 owner_id                 // Illegal, body already fully declared
}
```

An optional Annotation can also be given for the Top-Level declaration. In that
case, if a matching Annotation is given in the Parameter List, the engine must
fail with an error, since the Annotation has already been specified.

```
in PetIn 1:Pet {
    1: header string Authorization      // Error, annotation already specified
}
```

### Exception Annotations

```
exception-annotation    = +decimal [ "xx" ] ":"
```

An Exception Annotation is a type of annotation where the identifiers may
include `xx` after an initial digit. This enables annotations such as `2xx`, to
signify "any success HTTP response code", or `4xx`, meaning "any client error".
Exception Annotations can be used to list response constructs that apply
according to return codes, even specific codes, such as `201`.

Although any non-empty sequence of numbers followed or not by `xx` is
syntatically valid, only annotations that correspond in form to an HTTP status
code, if `xx` was to be replaced by any two numeric digits, are semantically
accepted by the interpreter. Any violation of this rule causes an error, with
the indication that the annotation is not valid.

### Exceptions

```
exception           =  exception-annotation exception-decl
exception-decl      =  [ annotation ] type [ parameter-list ]
exception-decl      =/ parameter-list
```

Exceptions consist of declarations of alternatives for Parameter Lists,
according to the context of a given HTTP status code. As with the Input
declaration, the expression accepts a Type, representing a Top-Level
Declaration, and a Parameter List, specifying a Key-Value Mapping. The
Top-Level Declaration's Location attribute cannot be modified, and always
refers to the body. The example that follows illustrates one such construct.

```
2xx: 1:[Pet] {
    2: header string Authorization
}
```

Other considerations for the Input construct also apply. When a Reference is
present, `body` declarations in the Parameter List are illegal, but, if omitted,
any number of `body` Parameters may appear, in which case the attributes are
passed as Key-Value Mappings instead.

### Exception List

```
exception-list          =  "{" exception-list-items "}"
exception-list-items    =  exception
exception-list-items    =/ exception ","
exception-list-items    =/ exception "," exception-list-items
```

Exception Lists introduce lists of Exceptions. As is common with lists of the
same kind, they push a new Scope into the symbol space.

### Output

```
output                  =  "out" output-decl
output-decl             =  symbol output-decl-options
output-decl-options     =  [ annotation ] type [ prototype ] [ exception-list ]
output-decl-options     =/ prototype [ exception-list ]
output-decl-options     =/ exception-list
```

Ouput declarations define output models with associated HTTP return codes. The
declaration allows for a Reference to define the Top-Level Model that is to be
sent with the body for success responses, while also listing Exceptions for
other specific HTTP status codes. In the example below, the declaration defines
`Pet` as the top-level body parameter and a `ServerError` as a top-level entity
for responses that indicate server errors.

```
out PetOut Pet {
    5xx: .ServerError
}
```

When a Reference declaration is given, it defaults to a `body` Location under a
`2xx` Exception, meaning that the two declarations in the examples that follow
are equivalent. In fact, the two must not be specified together, since that
creates a conflicting declaration, given that the body is specified twice.

```
out PetOut Pet
```

```
out PetOut {                    // Equivalent
    2xx: Pet
}
```

```
out PetOut Pet {
    2xx: Pet                    // Error, body is already specified
}
```

```
out PetOut Pet {
    2xx: {
        body Pet pet            // Error, body is already specified
    }
}
```

When declaring specific success status codes, however, the declaration can be
overriden. In the example that follows, all success codes return a `Pet` object,
while the `204` status code (No Content) returns nothing. This is valid because
the Exception is more specific than the Top-Level declaration.

```
out PetOut Pet {
    204: void
}
```

The same rationale applies to Top-Level Annotations. Given that the Top-Level
declaration specifies a construct for the `2xx` Exception, the Annotation does
not apply to other Exceptions, and thus the same Annotation can be repeated.

```
out PetOut 1:Pet {
    204: void,                  // Annotation omitted for void
    4xx: 1:.ClientError         // Annotations do not conflict
}
```

### Path Templates

```
path-template   = "{" word "}"
```

Path Templates parameterize Universal Resource Identifiers (URI). For example,
in the path `pets/{petId}`, `petId` is a Path Template argument that is supposed
to be replaced by a value in the path component of the URI.

All values that are valid as per the `segment` specification in
[RFC 3986](https://tools.ietf.org/html/rfc3986), including empty strings, are
also valid as a replacement for templates. Notably, values that match `path`,
as per the same RFC, are not valid, since those would potentially result in
conflicting declarations. For example, if `petId` in, say, `/pets/{petId}` was
to be replaced by `poly/attr`, it would create a conflict if a second rule was
defined as `/pets/{petId}/attr`, in which case the substitution of `petId` with
`poly` results in the same URI path component: `/pets/poly/attr`.

### Path

```
path                    = *( "/" [ parameterized-segment ] )
parameterized-segment   = *pchar [ path-template ] *pchar
segment                 = *pchar
pchar                   = unreserved / pct-encoded / sub-delims / ":" / "@"
unreserved              = LETTER / DECIMAL / "-" / "." / "_" / "~"
pct-encoded             = "%" HEXDIGIT HEXDIGIT
sub-delims              = "!" / "$" / "&" / "'" / "(" / ")" / "*" / "+" / "," / ";" / "="
```

A Path corresponds to a URI `path` component, as per section 3.3 of the
[RFC 3986](https://tools.ietf.org/html/rfc3986) specification, with the
exception that paths must start with a forward slash (`/`), unlike the
original RFC, which allows a path to begin with a segment. Paths can also
be empty, in which case they correspond to the root document, without a
trailing slash. If a single forward slash is given, then the path represents
the root document with a trailing slash.

It's notable that Poly makes a distinction (as it should) between paths that
end with a trailing slash and those that don't. Therefore, the URIs `/pets`
and `/pets/` correspond to different server endpoints.

Not unlike the original RFC specification, paths consisting of a sequence of
forward slashes, such as `///`, are valid paths. This is especially true given
that paths can be parametrized by Path Templates, since Path Templates can be
replaced by empty values.

### Verbs

```
verb    = "GET" / "POST" / "PUT" / "DELETE" / "PATCH"
```

A verb corresponds to the HTTP definition of a method, as defined by the
HTTP/1.1 specification in [RFC 7231](https://tools.ietf.org/html/rfc7231),
and indicates the purpose of a request. The following verbs are supported:

* `GET`, `POST`, `PUT`, and `DELETE`
([RFC 7231](https://tools.ietf.org/html/rfc7231));

* `PATCH` ([RFC 5789](https://tools.ietf.org/html/rfc5789)).

Several other verbs were left out of the specification, mostly because they
are deprecated, not used, or not relevant in the context of this specification.
However, some of such decisions are questionable. The `CONNECT` verb, for
example, may bear relevancy in many cases. Therefore, it's reasonable to assume
that the set of verbs excluded from the specification will change in future
revisions. The excluded verbs are the following:

* `OPTIONS`, `HEAD`, `TRACE`, and `CONNECT`
([RFC 7231](https://tools.ietf.org/html/rfc7231));

* `LINK` and `UNLINK` ([RFC 2068](https://tools.ietf.org/html/rfc2068));

* The HTTP/2 `PRI` verb ([RFC 7540](https://tools.ietf.org/html/rfc7540));

* All WebDAV extensions.

### Procedures

```
procedure       = verb path procedure-in procedure-out / reference
prodecure-in    = [ annotation ] type
procedure-out   = [ annotation ] type [ exception-list ]
```

Procedures represent server endpoints, and are associated with an HTTP method
(or Verb) that indicates the type of action to take on a given resource. The
resource is indicated by a Path, possibly parameterized, and the method by a
Verb. The following example illustrates one such declaration:

```
GET /pets void 1:[Pet] {
    4xx: ClientError,
    5xx: ServerError
}
```

If either the `procedure-in` or the `procedure-out` constructs correspond to a
Reference to Input or Output declarations, respectively, then those must not be
annotated, otherwise causing the engine to raise an error.

If either the `procedure-in` or the `procedure-out` constructs correspond to a
Model Reference instead, then the Reference can be annotated, in which case it
defines the Top-Level declaration for either the Input or Output declarations,
respectively.

A Procedure declaration accepts an Exception List as the last construct, in
which case the list of Exceptions is added to the Output specification given.
This behaves the same way as the Output declaration when produced with the `out`
keyword, with the exception that the construct must not use Prototypes, and it
does not introduce any symbol in the symbol space.

Procedures are r-values, and cannot be referenced directly. They can, however,
be referenced when annotated, when in the context of a Group Annotation. In
that case, the given Reference must be to a Procedure, otherwise resulting in a
semantic error.

### Group Annotations

```
group-annotation    = word ":"
```

Group Annotations work as other types of Annotations, in the sense that they
provide an identifier for the declaration that follows them. The key difference
is that Group Annotations are not numeric, but rather `word`s. Annotations
introduce the construct that follows them associated with the symbol that they
define on the Active Scope, and are thus l-values.

### Annotated Procedure List

```
annotated-procedure-list        =  "{" annotated-procedure-list-items "}"
annotated-procedure-list-items  =  annotated-procedure-list-decl
annotated-procedure-list-items  =/ annotated-procedure-list-decl ","
annotated-procedure-list-items  =/ annotated-procedure-list-decl "," annotated-procedure-list-items
annotated-procedure-list-decl   =  group-annotation procedure
```

Annotated Procedure Lists are like any other list, in the sense that they
consist of comma-sperated values delimited by curly brackets. The list
introduces Group Annotations with corresponding Procedures or procedure
References.

A Group Annotation must not appear repeated in the list, otherwise causing an
error.

### Templates

```
template        = "template" template-decl
template-decl   = symbol template-entity-list annotated-procedure-list

template-entity-list        = "(" template-entity-list-items ")"
template-entity-list-items  = symbol
template-entity-list-items  = symbol "," template-entity-list-items
```

Templates define abstract constructs that can be used as a means for code reuse.
The code below, for example, illustrates how one could roughly represent the
REST concept, by specifying a group of methods that respects to the paradigm.

```
template PetStoreTemplate (Entity, NewEntity) {

    list: GET void [Entity],

    create: PUT Entity NewEntity,

    read: GET /{id} void Entity,

    update: POST /{id} Entity Entity,

    delete: DELETE /{id} void void
}
```

Templates define a list of Procedures and another of abstract entities. Such
entities can be replaced by specific Models when the Template is being
instantied.

Templates themselves can be Prototyped, in which case they inherit the
declarations from the prototype parent. The annotations, however, are not
overriden, and instead must not match, otherwise triggering an error.

Templates are l-values, and any of the Symbols that they declare can be referred
by References.






### Groups

```
group       = "group" group-decl
group-decl  = symbol path [ prototype ] annotated-procedure-list
```

Groups are constructs used to group Procedures that relate over some semantic
aspect or their base path. That is, they can, for example, refer to the same or
similar entities, or otherwise share the same base URI.

Groups introduce a Group Annotation list, in which the annotated elements can
only be Procedures. This creates a list of annotated Procedures that lives and
is referenceable under the same entity. 

In the example that follows, the base path `/pets` is inherited by the `list`
Procedure, since the Path component is empty. As for the `read` method, the
component is attached, resulting in `/pets/{petId}`.

```
group PetStoreGroup /pets {

    list: GET void [Pet],

    read: GET /{petId} void Pet
}
```

A trailing slash may appear in the Group URI, but that is not recommended. This
is because the Path components must already begin with a slash, meaning that
a duplication would occur. If URI for the Group in the example above was to be
changed to `/pets/`, then the `read` method would live under `/pets//{petdId}`,
with two slashes.

The Group Annotations can be used to create References to Procedures. In that
sense, the method for `read`, above, can be referenced as `PetStoreGroup.read`.

Groups support Prototyping, but from Templates, rather than from other Groups.
As such, Groups don't support real inheritance, but rather an implementation of
Templates, similar to the Java `interface` concept.

Groups are r-values, introducing the Symbol that they define in the Active
Scope.









In the example that follows, `PetStoreGroup` inherits all constructs from
`PetStoreTemplate`, while defining `Entity` and `NewEntity` as `Pet` and
`NewPet`, respectively.

```
group PetStoreGroup /pets from PetStoreTemplate (Pet, NewPet)
```

Further defining other Procedures is optional, but, in any case, those still
need to be annotated. In case the annotations overlap with the Template, the
Procedure is not overloaded, and rather triggers a system error.

```
group PetStoreGroup /pets from PetStoreTemplate (Pet, NewPet) {

    create: PUT Entity Entity,      // Error, "create" is already defined
}
```


The Paths that appear on the Template will be resolved according to the base
Path given by the Groups that implement it. That is, the `read` Procedure for
the example above, would yield `/pets/{id}`for the instantiation that
follows it.









### Future Work

* How can the language support arbitrary key-value mappings? Currently, there's
no way to specify an arbitrary object, in the sense that the keys are not well
defined;

* Are theyre any situations in which `void` can be annotated? If so, the
description about annotations must be reviewed;

* It should be said that declarations that are not present in an API
specification are not a violation of the declaration. That is, services may
implement methods beyond those declared;

* The language should support abstract models, being of little consequence to
the declarations, but while giving an indication to code generators that the
model is not an actual model;

* Multiple inheritance would greatly improve code reuse. Although many languages
do not support multiple inheritance natively, there are ways in which code
generators can work around this limitation;

* Alternatives are not supported for I/O declarations. This means that the
language cannot express constructs such as "I accept either A or B";

* Consider validators as a way to ellaborate on what a given Template is
supposed to accept. Validators can be regular or context-free;

* The Location attribute is deeply coupled with the HTTP protocol, and thus
support for other transport schemes is compromised. Future revisions should
put this into consideration;

* It's not known how to handle Unicode strings, or their respective encodings.
As it currently stands, it would appear that the declaration forces services to
support any kind of string, whatever the encoding, and that might not be too
realistic;

* This document is not final and several constructs are still being added to
this specification. The major ones being considered are `group`, `rest`,
`socket`, and `websocket`.
