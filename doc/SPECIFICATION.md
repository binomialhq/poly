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
is easily explained by the fact that Poly attempts to consider a wider variety
of scenarios, as well as purpose-specific language constructs. The latter, on
the other hand, can be explained by observing that not all constructs offered
by the alternatives directly relate to API specifications, such as `servers`
and `baseUri`. These markers, provided by OpenAPI and RAML, respectively, mix,
in our opinion, concepts that should not be mixed, that of an API specification
and the location of some specific implementation of it. Poly proposes a strict
separation of these concerns.

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
of a context-free grammar expressed in augmented Backus-Naur Form (BNF). It is
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
ULETTER     = \p{L}         ; Any Unicode L-category code point
UDECIMAL    = \p{Nd}        ; Any Unicode Nd-category code point
USPACE      = \p{Z}         ; Any Unicode Z-category code point
```

Implementations may decide to support L-category characters in combination with
M-category markers. Although the `ULETTER` declaration does recognize diacritics
when encoded together as a single code point, it also fails to recognize them
if the two are encoded separately. For example, `à` can be encoded as U+00E0 or
as U+0061 U+0300. The former encodes `à` as a single code point, while the
latter encodes `a` followed by the accent modifier. In order to support
modifiers, the engine may replace the `ULETTER` declaration above:

```
ULETTER     = \p{L}\p{M}*   ; Any L-category code point and M-category mark
```

As will be made clear throughout the specification, some names can also include
a set of extra characters. This set includes the Unicode characters U+005F (Low
Line, more oftenly known as "underscore") and U+2010 (Hyphen), enabling symbols
such as `User-Agent` and `_tcp` to be defined.

```
EXTRA   = [\u005F\u2010]
```

Several non-terminals are also defined for reuse throughout this specification.
This set specifies the basic glyphs that can be processed by the engine. As per
this specification, engines should support Unicode. The notable exception is
the `number` declaration, which discards numeral systems other than the
Hindu-Arabic. This is most unfortunately left out of the specification due to
one's own lack of knowledge in other numeral systems, an issue that should be
addressed in future revisions.

```
letter      = ULETTER / EXTRA
decimal     = UDECIMAL
number      = DECIMAL
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
space   = SPACE
```

Space is generally ignored by the engine unless indicated otherwise. However,
it's notable that "space" is defined as "US-ASCII space", and not the broader
`USPACE` declaration with Unicode support. In fact, it's recommended that the
engine triggers warnings for code points that match `USPACE` but not `SPACE`,
when not in a context that indicates that space is not to be ignored.

### Symbols

```
DQUOTE      = ["]                   ; Single U+0022 (")
NOTDQUOTE   = [^"\\]                ; Anything but U+0022 (") or U+005C (\)
ESCAPED     = "\\" .                ; U+005C (\) followed by anything
```

```
symbol          = symbol-simple / symbol-quoted
symbol-simple   = letter *( letter / decimal )
symbol-quoted   = DQUOTE *( *NOTDQUOTE [ ESCAPED ] ) DQUOTE
```

A Symbol is a name that appears as the name of some declaration. For example,
in the declaration `model Pet`, `Pet` is the symbol, or name, that can be used
to refer to the model being declared. Several constructs accept symbols as a
means of identification.

Symbols may also be quoted, so that characters that do not fit this definition
can be used. For example, the declaration `model "ASN.1"` is valid, and declares
a symbol called `ASN.1`, discarding the quotes.

Any reference to a quoted symbol must be provided exactly as the original
declaration, including Space. This means that quoted symbols may be hard to
read, and thus their usage is discouraged, unless under extraordinary and
justified circumstances. How to handle multi-line strings is an open issue.

### Scopes

Several Poly declarations introduce new symbols to the symbol space. How and
where these symbols are declared bears significance as to where the symbols are
made visible. The symbol space is composed of a stack of scopes, each of which
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

The symbol lookup system, then, begins by looking for a symbol in the current
scope. Only if one is not found there, does the implementation look for the
same symbol in the scope that preceeds it in the stack. Therefore, scopes are
searched in reverse depth-order, starting with the topmost entry. This creates
an overloading effect.

When a name is not found in any scope as a result of a lookup, the engine must
fail with an error. Symbols that are declared and yet already found on the same
scope must also fail with an error, with an indication that the declaration is
duplicated.

Finally, the Primitive Space is one in which symbols are always visible and
cannot be overloaded. Symbols declared in this space, thus, may not be
overloaded, making it ideal to declare names that are reserved by the system.
If the name of a primitive is redeclared, in any context, the system must fail
with an error. If the name of a primitive is referenced, then it can only refer
to the primitive with that name. This can be implemented by always starting the
lookup process with the Primitive Space.

By default, all constructs that produce an identified declaration (e.g. not
anonymous) introduce a symbol to the scope where they are defined. Such
constructs are properly identified throught the document.

### References

```
reference           =  ["."] symbol
reference           =/ symbol "." reference
```

References consist of a way to refer to symbols whatever the scope that declares
them. For example, if a model in the Global Space `Pet` declares a field `name`,
the reference `Pet.name` expands into the same declarative construct. Given that
the same semantic entity is used, the referrer inherits all properties of the
reference, including type, modifiers, and any other attributes. This enables
the reuse of declarations without having to explicitly rewrite all the details.

The lookup for the first component of a reference list follows the same process
as described before. This means that the symbol is looked up first in the Active
Space, and then successfully traced by the other scopes on the stack. This also
applies if that's the only symbol in the list. For symbols that appear on the
list after that, the lookup does not backtrace, and instead only looks for the
symbol in the scope that is referenced up to that point. Therefore, the lookup
`Service.Model.Field` looks for `Service` using the previous method and
`Model.Field` only within the context of that first scope.

It's notable that this scheme may create ambiguities. The example that follows
shows an ambiguous lookup for `PetStore`, which happens because the symbol
exists both in the Active and the Global Space scopes. This means that both
declarations can be referenced in the same way, hence the ambiguity. In this
case, the `Veterinary.PetStore` construct is the right choice for the engine,
since it lives closer to the scope where the symbol is being referenced.

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

### Identification

An expression is said to be _identified_ when it introduces a new symbol that is
directly associated with the rule. `model Pet`, for example, is identified,
since the `Pet` symbol can be used to refer back to the declaration. Poly also
supports _anonymous_ expressions, which in turn produce a rule, but do not
associate it with any symbol. The example below illustrates an anonymous `model`
declaration, classified as such for not being associated with any symbol.
Throughout this specification, each section makes clear which rules can be
identified or anonymous.

```
model {
    1: int32 id,
    2: string name
}
```

### Value Categories

Poly defines two value categories, according to whether they produce new symbols
or not. An l-value, thus, refers to an expression that persists by introducing
some symbol into the Active Scope. On the other hand, r-values do not introduce
any new symbols to any scope, and thus do not persist beyond the point where
they are declared. All identified declarations are l-values, as are assignments,
while references and anonymous declarations are r-values.

It's notable that each construct, small as it might be, produces values of
either one or the other category. Consider three constructs in the example
below. The quoted symbol "ASN.1" and the scope delimited by curly braces are
r-values since they do not introduce any symbols. This may seem counterintuitive
because the quoted symbol is a symbol, but it is in fact the `model` expression
that declares it, and thus the `model` expression is the one that is an l-value.

```
model "ASN.1" {

}
```

### Comments

```
ENDLINE             = .*$           ; Anything up to EOL
MULTILINE           = (?!\*\/)*     ; Anything except "*/"
```

```
comment                 = comment-single / comment-multiple
comment-single          = "//" ENDLINE
comment-multiple        = "/*" MULTILINE "*/"
```

Poly supports C-style comments in both single-line and multi-line forms. That
is, comments can be started by a double forward slash (`//`) and span until the
end of the line, or by the multi-line start sequence (`/*`) and span multiple
lines until the multi-line end sequence (`*/`). In either case, the span area
is ignored by the interpreter.

When either form of commenting immediately preceeds an expression, it will be
considered associated with the expression that it preceeds. Specifically, for
comments that are associated with l-values, the comment block is processed by
the documentation engine, if one exists.

### Syntax

```
syntax      = "syntax" integer
```

The Syntax declaration indicates the version of the Poly IDL that the file is
using. The practical consequence is having the parser select the appropriate
engine to process and interpret the input. This declaration is required and
must be the first non-empty and non-comment line in the file.

There's only one version number currently supported, and that's "0" (zero). This
version does not guarantee the continuity of any of its constructs, and thus
future revisions may not be compatible. The following table summarizes version
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

A comment block that applies to the Syntax declaration is interpreted by the
engine as file-level documentation. By applying to the document, this makes it
the right place for documentation such as API versioning information.

### Primitives

```
native              = "native" native-decl
native-decl         = word / native-array-decl
native-array-decl   = "[" word "," "..." "]"
```

Primitives are data types that must be natively supported by the engine. A
primitive is declared with the `native` keyword. If some declared primitive is
not understood by the engine, the engine must fail the process with an error.
Although primitives can be declared in any file (e.g. they are syntactically
valid), declaring primitives that are not supported by the engine will cause
general failure, and thus that is not recommended.

Primitive declarations introduce symbols in the Primitive Space, as declared by
the `word` grammatical rule. If the symbol already exists in the Primitive Space
(e.g the declaration is repeated), the engine must fail with an error.

Primitives are not declared on every file, but on a standard declaration space
that is included implicitly by the engine. How that space is declared is up to
the implementation, and interpreters may decide to have a file physically
allocated on persistent memory, but are not required to do so. However, whatever
declarations are found there should be effective; that is, if a declaration is
not present, then the engine must not recognize it.

Primitives can also be assigned as any other form of expression, since any
reference to them just returns their declaration. As per the example below,
the assignment creates an alias for the primitive type, but the alias is not
itself a primitive. The main difference is that type alias can be overriden by
deeper declaration spaces, while primitives cannot. In other words, type alias
are declared in the scope where they appear, and not in the Primitive Space.

```
i32 = int32;
```

The supported primitives are summarized as follows:

| Primitive    | Description                                 |
|--------------|---------------------------------------------|
| int32        | Signed 32-bit integer.                      |
| int64        | Signed 64-bit integer.                      |
| float        | Floating point value.                       |
| double       | Floating point value with double precision. |
| string       | To be determined.                           |
| wstring      | To be determined.                           |
| [array, ...] | An heterogeneous collection of entities.    |

It's still being considered whether the integer types have shorter versions
abbreviated as `i32` and `i64`. Unsigned integers are not supported because (1)
Poly is not capable of enforcing validation and (2) those are not supported by
every platform.

The string and wstring types are still under analysis. The general idea is to
support Unicode, but neither of these types is coupled with that concept. The
main idea is to make `string` an abstract metatype that can represent any
collection of characters, which would be well aligned with the considerations
above.

The name `array` in the standard declaration doesn't mean anything, and it is
not a reserved keyword. `[T, ...]`, for example, is also valid and the two are
semantically equivalent.

Arrays are heterogenous but only if declared that way. A `[Pet]` declaration,
for example, denotes "an array of `Pet` objects", but nothing else. The ellipsis
symbol signifies a repetition, meaning that declarations such as `[Pet, Error]`
are possible, denoting a mixed array of `Pet` and `Error` objects.

### Annotations

```
annotation  = +number ":"
```

Annotations associate expressions with a given numeric identifier, and can be
used to create static indexes of elements. The annotations do not necessarily
dictate order, however, in the sense that the identifiers need not be sorted.
In fact, the annotations are not processed by Poly, since they are meant for
binary encodings.

Protobuf, for example, requires these in models, while Thrift also requires them
in parameter declarations. One requirement that these two implementations share
is the fact that assigning an annotation is a permanent action, in the sense
that the number becomes reserved. Future revisions of the API must not reutilize
them, and instead add new ones in case of need, while commenting out the ones
that become deprecated. For the sake of compatibility, Poly adheres to this
principle.

When annotations appear in lists, specifically Field Lists and Parameter Lists,
a mix of annotated and non-annotated entities is not allowed. Although such
constructs are valid from a syntatic perspective, such scenario must not pass
semantic validation, provoking an error that base being verified.

As was clarified before, numbers with semantic meaning are only supported for
the Hindu-Arabic numeral system, which is the case of annotations. Future
revisions should change this behaviour to span other numeric systems as well.

### Modifiers

```
modifier            = word / "?" / "!"
modifier-list       = modifier modifier-list / empty
```

Modifiers can be `word`s, question marks (`?`), or exclamation marks (`!`), and
cause some form of semantic change to the declaration. The `deprecated` keyword,
for example, causes the field to be flagged as being deprecated. The question
(`?`) and the exclamation (`!`) marks are special sugary constructs that mean
`optional` and `required`, respectively. The following table summarizes the
recognized modifiers:

| Modifier   | Description                 |
|------------|-----------------------------|
| sensitive  | Flags a field as sensitive. |
| deprecated | To be determined.           |
| required   | Flags a field as required.  |
| optional   | Flags a field as optional.  |
| !          | Alias for "required".       |
| ?          | Alias for "optional".       |

The `sensitive` keyword tells the interpreter to flag the field as being of a
sensitive nature, applying to the likes of passwords. This will instruct UI
generators to use an obfuscation control, for example. More importantly, it can
tell the deployment enviroment not to log this attribute, or not to broadcast
it on the network.

Modifiers are not reserved keywords, but rather any string that respects the
`word` rule. This is motivated by future compatibility, since modifiers can be
added to the language. However, this bears the consequence that modifiers can
appear as quoted identifiers, which is syntatically valid, but not very useful.
Regardless, if a modifier is known to the system it is also recognizable in its
quoted version, and thus `sensitive` and `"sensitive"` are equivalent.

The `required` and `optional` modifiers are mutually exclusively, as are their
alias counterparts. The semantic value of the action is to flag an entity as
being of a certain nature, and that nature cannot contradict itself. Therefore,
`required` and `optional`, as well as `!` and `?`, must not appear in the same
expression, otherwise causing the system to error.

### Fields

```
field               =  field-type field-name modifier-list
field-type          =  symbol
field-name          =  symbol
```

A Field declares an attribute of some construct in which it appears, meaning
that the purpose of a Field expression is only made clear when found within
the scope of some other declaration. They are, however, l-value expressions,
and introduce the symbol that they declare in the context in which they appear.

Fields have a type, a name, and zero or more modifiers. The type must correspond
to either a Primitive or a Model, declared somewhere in the symbol space,
causing an error that not being the case. These declarations, thus, associate a
data type and modifiers with a given name, the same name that is to be used for
symbol lookups.

By default, fields are required. The explicit presence of the `required`
modifier is redundant, unless some option is passed to the interpreter that
changes the default (e.g. command line options).

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
field-list          =  field-list-decl
field-list          =/ field-list-decl ","
field-list          =/ field-list-decl "," field-list
field-list-decl     =  [ annotation ] field
```

A Field List is a non-empty sequence of annotated Fields separated by commas,
and optionally terminated by an additional comma (for convenience). It's also
possible that the are not present at all, with the implication that the
declaration might not be supported by some tools.

### Location

```
location    = word
```

The Location attribute refers to where a value is expected to appear in the
context of the HTTP protocol. For example, a `header` qualifier indicates that
the value is expected to be passed in the headers. The following table lists
all available location modifiers:

| Location | Description                                                 |
|----------|-------------------------------------------------------------|
| body     | The parameter is passed in the body.                        |
| path     | The parameter is passed in the URL, in the path component.  |
| query    | The parameter is passed in the URL, in the query component. |
| header   | The parameter is passed as a header.                        |
| cookie   | The parameter is passed as a cookie.                        |

`path` applies to parameterized endpoints, where the value of the attribute
will replace some template variable. For example, if an expression declares
some URL endpoint such as `items/{itemId}`, the template `itemId` would be
replaced by the value in `path`. All other options should be clear.

### Parameters

```
parameter   =  location field / field
```

Parameters extend on Fields by specifying the Location in which the field is to
appear. This means that Parameters are constructs that are meant for use within
the context of some API related declaration, such as Prodcedures. The Location
appears first in the construct. When omitted, `body` is assumed as the default,
saving some typing.

### Parameter Lists

```
parameter-list          =  parameter-list-decl
parameter-list          =/ parameter-list-decl ","
parameter-list          =/ parameter-list-decl "," parameter-list
parameter-list-decl     =/ [ annotation ] parameter
```

A Parameter List declaration is identical to a Field List in every way but one:
that each element of the list is a Parameter, and thus introduces a Location
attribute. All other considerations apply.

### Prototyping

```
prototype   = ":" symbol
```

Poly implements four declarative constructs that support prototype inheritance,
namelly `service`, `model`, `in`, and `out`. When inheriting from a prototype,
declarations implement all symbols declared by that prototype. Depending on the
situation, that can be done explicitly or implicitly. Multiple inheritance is
not supported, motivated by the Diamond Problem and the lack of support in
several languages.

When prototyping `model`, `in`, or `out` and using annotated declarations, the
annotations are not inherited, and all fields must be explicitly listed in order
to avoid conflicts. The example below, of model prototyping, illustrates what
would happen if annotations were to be inherited, creating a conflict that the
engine would not be capable of resolving.

```
model NewPet {
    1: string name required,
    2: string tag
}

model Pet: NewPet {
    1: int32 id required        // Conflict! Annotations match
}
```

For that reason, all parent declarations (e.g. fields and parameters) must be
explicitly referenced by the child, reassigning the annotations, otherwise
incurring in the penalty of an error being indicated by the engine. The example
below shows an error where the annotations are inherited but not explicitly
reassigned, changing the nature and the location of the error, in comparison
to the previous example.

```
model NewPet {
    1: string name required,
    2: string tag
}

model Pet: NewPet {             // Error! "name" and "tag" are not declared
    1: int32 id required
}
```

The example that follows illustrates the correct way for prototyping with
annotations. Notably, all fields are listed, but only declared as References.
Types, modifiers, and location attributes are inherited implicitly, and cannot
be overloaded. Thus, the child field or parameter holds all properties of the
corresponding parent, with the exception of the annotation itself. This enforces
that changes made to the parent also reflect on the child, including `sensitive`
and `deprecated` modifiers.

```
model NewPet {
    1: string name required,
    2: string tag required
}

model Pet: NewPet {
    1: int32 id required,
    2: NewPet.name,
    3: NewPet.tag
}
```

From a code generation perspective, generators may implement these constructs
using strict inheritance, since the declarations respect the "is-a" relationship
and, thus, the Liskov Substitution Principle. In fact, any context in which a
parent is declared and a child is given in its place, the extra fields or
parameters declared by the child should be discarded, and the entity processed
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

model Pet: NewPet {        // Valid, since fields are not tagged
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

model Pet: NewPet {         // Valid, annotations are ignored
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

model Pet: NewPet {         // Valid, all fields are annotated
    1: int32 id required,
    2: NewPet.name,
    3: NewPet.tag
}
```

Another key aspect of prototype inheritance is that if one parent field or
parameter is declared explicitly then they all have to be. This works as a
semantic construct to guarantee that all fields are properly handled, and that
updates to the parent are visibly reflected in the child, rather than silently
absorved.

None of this rationale applies to `service` declarations, however, since they
don't implement annotations at all. For that reason, services don't need to
explicitly declare inherited constructs, as the other prototypable entities do. 

The main difference between these entities is that `model`, `in`, and `out` do
not inherit the scope of their parents, while `service` inherits a copy. To be
precise, when `model`, `in`, or `out` prototyping occurs, the scope must be
redeclared with references to the same symbols, while `service` keeps a copy of
the parent scope as a scope of its own.

The practical consequence of this is that services cannot overload declarations,
including methods. This happens because Poly enforces that child services be
merely an extension of their parents, without any disregard for the contract
that they promote.

Anonymous prototyping is also supported. In this scenario, all of the same
principles apply, except that the child entity is not named. Although this
may appear useless, it's notable that the entity may still be an r-value for
an assignment, meaning that it still has the pontential to become visible in
the scope.

```
model : NewPet {
    1: int32 id required,
    2: NewPet.name,
    3: NewPet.tag
}
```

One final observation is that no declaration can inherit a prototype from a
different type of declaration. That is, models inherit from models, services
inherit from services, and so on, and the prototype chain must not be mixed.

### Models

```
model                   = "model" model-decl
model-decl              = [ symbol ] [ prototype ] "{" field-list "}"
```

A Model declaration defines a complex data type constructed over Primitive types
or other Models. Models can be used as input and ouput for procedures, as well
as other situations where a data type is appropriate. Specifically, models can
be used anywhere were a `symbol` appears in the grammar.

When generating code, engines may generate classes and other data structures
from models, as they are meant to represent structures that hold data. Models
are not used to represent any form of encoding, and are rather an abstract
representation of data graphs.

Models can be identified or anonymous. When identified, a model explicitly
indicates the symbol that becomes associated with the declaration, and, when
anonymous, that association either does not exist or is achieved through an
assignment. The following declarations are, thus, equivalent:

```
model Pet {
    1: string name,
    2: string tag
}

Pet = model {
    1: string name,
    2: string tag
}
```

Therefore, model declarations return the entity that they declare, which is
introduced in the active scope. This entity can be assigned to a symbol (e.g. a
constant). If an identified model is also assigned, then the model can be
referred in either way, but the names must not match, otherwise incurring in a
violation of a repeated declaration for the same scope.

```
PetModel = model Pet {      // PetModel and Pet are equivalent
    1: string name,
    2: string tag
}
```

A third option is to declare an anonymous model without any form of assignment,
as examplified below. In this case, the model declaration is ineffective, and
may be discarded by the engine (e.g. garbage collected), since it cannot be
referrenced. More specifically, as the model is to be introduced in the scope,
the scope discards it because it doesn't have a symbol for it, but it does not
produce an error. This type of construct may be used during the design and
development stages, or to flag models that are intented for future reviews of
the API.

```
model {
    1: string name,
    2: string tag
}
```

The use of curly braces is indicative of a declaration space that is pushed onto
the declaration space stack. Therefore, names that are declared in this space do
not conflict with names from the parent scope, with the notable exception of the
Primitive Space. The example below shows how three declarations of the same
symbol `name` do not conflict with each other, since they live in different
declaration spaces.

```
name = ...;             // Some declaration...

model Pet {
    1: string name      // No conflict
}

model NewPet {
    1: string name      // Still no conflict
}
```

### Verbs

```
verb    =  word
```

A verb corresponds to the HTTP definition of a method, as defined by the
HTTP/1.1 specification in [RFC 7231](https://tools.ietf.org/html/rfc7231),
and indicates the purpose of a request. Verbs can be expressed in all lower
or all upper case, but not in mixed case. The following verbs are supported:

* `GET`, `POST`, `PUT`, and `DELETE`
([RFC 7231](https://tools.ietf.org/html/rfc7231));

* `PATCH` ([RFC 5789](https://tools.ietf.org/html/rfc5789)).

Several other verbs were left out of the specification, mostly because they
are deprecated, not used, or not relevant in the context of the purpose of
this specification. However, some of such decisions are questionable. The
`CONNECT` verb, for example, may bear relevancy in many cases. Therefore, it's
reasonable to assume that the set of verbs excluded from the specification will
be reviewed in future revisions. The excluded verbs are the following:

* `OPTIONS`, `HEAD`, `TRACE`, and `CONNECT`
([RFC 7231](https://tools.ietf.org/html/rfc7231));

* `LINK` and `UNLINK` ([RFC 2068](https://tools.ietf.org/html/rfc2068));

* The HTTP/2 `PRI` verb ([RFC 7540](https://tools.ietf.org/html/rfc7540));

* All WebDAV extensions.

### Input

```
input           = "in" input-decl
input-decl      = [ symbol ] [ prototype ] "{" parameter-list "}"
```

The `in` keyword can be used to specify input parameters for a given construct,
mostly Prodecures.


"Normal" (or explicit) declaration:

```
in PetIn {
    1: Pet pet
}
```

"Short" version declares top-level:

```
in PetIn Pet;
```

### Output

```
output                  = "out" output-decl
output-decl             = output-decl-toplevel / output-decl-explicit
output-decl-toplevel    = symbol symbol
output-decl-explicit    = [ symbol ] [ prototype ] "{" parameter-list "}"
```

### Procedures

```
```

Normal declaration:

```
out PetOut {
    1: Pet pet
}
```

Short version declares top-level:

```
out PetOut Pet;
```
