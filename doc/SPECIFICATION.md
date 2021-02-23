# Poly

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
and [RAML](https://raml.org/). Instead, Poly is designed based on succint
syntatical constructs that are already familiar to developers. This enables
the use of operators, scope delimiters, documentation engines and other
constructs that are either not possible or not natural with typical data
serialization formats.

Similarly to its counterparts, Poly focuses on the Hypertext Transfer Protocol
(HTTP), but it additionally considers other alternative protocol schemes from
design. For example, OpenAPI and RAML do not support the association of indexes
with models and procedure parameters, both typical when using binary encodings.
Thefore, the likes of [protobuf](https://github.com/protocolbuffers/protobuf)
and [Thrift](https://thrift.apache.org/) cannot be described by those
alternatives. Poly, instead, ignores any declarations related to transport and
encoding schemes, making it agnostic to representational concerns and networking
stacks. This enables the likes of content negotiation, or alternative protocols.
Neither of those things is possible with OpenAPI or RAML.

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

We borrow the conventions from [RFC 2234](https://tools.ietf.org/html/rfc2234)
and other utilities from standards. We do not, however, fully compromise to
following them strictly in this version. Expect minor differences.

All of the mechanisms specified in this document are described in both prose
and an augmented Backus-Naur Form (BNF) defined in RFC 2234. Section 6.1 of
RFC 2234 defines a set of core rules that are used by this specification,
and not repeated here. Implementers need to be familiar with the notation and
content of RFC 2234 in order to understand this specification.

## Basic Rules

The following rules are used throughout this specification to describe basic
parsing constructs (e.g. terminals). Unless stated otherwise, all constructs
assume an UTF-8 encoding.

```
    LETTER      = ...   ; One underscore (U+005F) or category L character
    DECIMAL     = ...   ; One code point of 0 (U+0030) throught 9 (U+0039)
    NUMBER      = ...   ; One character in the Unicode category N
    SPACE       = ...   ; One character, as per Unicode's White Space property
```

From this list, it should be clear that Poly aims at maximizing support for
Unicode, but the proposed use of space is questionable. This should be reviewed
in future versions of this document.

The following non-terminals are also considered, for a veriety of purposes that
will be made clear for each context:

```
    empty           = ; Empty
    slug            = LETTER / LETTER slug
    type            = slug
    name            = slug
    number          = NUMBER / NUMBER number
    integer         = DIGIT / DIGIT integer
```

## Specification

### Syntax

```
syntax      = "syntax" integer ";"
```

The Syntax declaration indicates the version of the Poly IDL that the file is
using. The practical consequence is having the parser select the appropriate
engine to process and interpert the input. This declaration is required and
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
the right place for documentation such as API versioning information. In the
course of this analysis, the following tags were considered:

| Name             | Tag      | Description                                |
|------------------|----------|--------------------------------------------|
| Version          | @version | The API version of what is being declared. |
| Contact          | n/a      | Not supported.                             |
| Terms of Service | n/a      | Not supported.                             |
| License          | n/a      | To be determined.                          |

With a few notes:

* The Version tag should not be confused with the Syntax declaration. The former
declares the API version for the services being defined in the files, while the
latter declares the version for the Poly IDL processor that interprets them.

* Contact information is not supported based on the premise that including
contact information at a file level is bad practice, since such information
is usually project-level metadata.

* Terms of service are not covered by the specification because there’s no
service. The IDL defines the interface, which can be implemented by several
services. It’s notable that different implementations of the declaration may
have different terms of service.

* Licensing is still being considered. If the license is given as a copyright
notice at the top of the document, the engine will not be capable of parsing it
deterministically. Javadoc, for example, only supports defining the copyright
from the Command Line Interface (CLI). The final decision may depend on the
choice of a documentation engine.

### Primitives

```
native              = "native" native-decl ";"
native-decl         = type / native-array-decl
native-array-decl   = "[" type "," "..." "]"
```

Primitives are data types that must be natively supported by the engine. A
primitive is declared with the “native” keyword, and is referred to as a Native
declaration.

If some declared primitive is not understood by the engine, the engine must fail
the process with an error. Although primitives can be declared in any file (e.g.
they are syntactically valid), declaring primitives that are not supported by
the engine will cause general failure, and thus it is not recommended. If the
same primitive name is declared more than once, the engine should also fail with
an error.

Primitives are not declared on every file, but on a standard declaration space
that is included implicitly by the engine. How that space is declared is up to
the implementation, and interpreters may decide to have a file physically
allocated on disk space, but are not required to do so. However, whatever
declarations are found in that space should be effective; that is, if a
declaration is not present, then the engine must not recognize it.

The supported native types are summarized as follows:

| Type         | Description                                 |
|--------------|---------------------------------------------|
| int32        | Signed 32-bit integer.                      |
| int64        | Signed 64-bit integer.                      |
| float        | Floating point value.                       |
| double       | Floating point value with double precision. |
| string       | To be determined.                           |
| wstring      | To be determined.                           |
| [array, ...] | An heterogeneous collection of entities.    |

Annotated as follows:

* It's still being considered whether the integer types have shorter versions
abbreviated as `i32` and `i64`. Unsigned integers are not supported because (1)
Poly is not capable of enforcing validation and (2) those are not supported by
every platform.

* The string and wstring types are still under analysis. The general idea is to
support Unicode, but neither of these types is coupled with that concept. The
main idea is to make `string` an abstract metatype that can represent any
collection of characters, which would be well aligned with the integer
considerations above.

* The name `array` in the standard declaration doesn't mean anything, and it is
not a reserved keyword. `[T, ...]`, for example, is also valid and the two are
semantically equivalent.

* Arrays are heterogenous but only if declared that way. A `[Pet]` declaration,
for example, denotes "an array of `Pet` objects", but nothing else. The ellipsis
symbol signifies a repetition, meaning that declarations such as `[Pet, Error]`
are possible, denoting a mixed array of `Pet` and `Error` objects.

### Fields

```
field               = type name modifier-list
modifier-list       = modifier modifier-list / empty
modifier            = name / "?" / "!"
```

A Field declares a variable, but not by itself (e.g. there's no such thing as
global variables), but rather as a part of other constructs, such as Model and
Endpoint declarations. The purpose of such declarations is only made clear when
found within the scope of such contexts.

Fields have a type, a name, and zero or more modifiers. The type must correspond
to either a Primitive or a Model, declared somewhere in the declaration space,
causing an error, that not being the case. These declarations, thus, associate
a data type with a given name, reserving it under the scope of the declaration.

A modifier can be a Slug, a question mark (`?`), or an exclamation mark (`!`). A
modifier causes some form of semantic change to the declaration, such as marking
it as "sensitive". The question mark (`?`) and the exclamation mark (`!`) are
special constructs that mean "optional" and "required", respectively, for the
sake of simplicity and code brevity. The following table summarizes the
available modifiers:

| Modifier   | Description                 |
|------------|-----------------------------|
| sensitive  | Flags a field as sensitive. |
| deprecated | To be determined.           |
| required   | Flags a field as required.  |
| optional   | Flags a field as optional.  |
| !          | Alias for "required".       |
| ?          | Alias for "optional".       |

Some comments are in order:

* The `sensitive` keyword tells the interpreter to flag the field as being of a
sensitive nature, such as passwords. This will instruct UI generators to use an
obfuscation control, for example. More importantly, it can tell the deployment
enviroment not to log this attribute, or not to broadcast it on the network.

* Modifiers are not reserved keywords, but rather any symbol that respects the
`name` rule. This is motivated by future compatibility, since modifiers can be
added to the language.

* The `required` and `optional` modifiers are mutually exclusively, as are their
alias counterparts. The semantic value of the action is to flag the field as
being of a certain nature, and that nature cannot contradict itself.

* By default, fields are required. The explicit presence of the `required`
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

It’s notable that Field declarations are not top-level declarations, and thus
are not terminated by a semicolon. Also, given that the declaration ends with a
list of slugs, other grammatical constructs should be careful enough to include
some form of separator, such as a comma, in order to avoid creating ambiguities.

### Tagged Fields

```
tagged-field    = tag ":" field
tag             = number
```

A Tagged Field introduces a tag number to a Field. The tag creates a static
order of elements, to be used in Field Lists, associating each field with a
numeric identifier. The purpose of these identifiers is to provide a numeric
ordering, often used by binary encodings. Protobuf, for example, requires these
tags in models, while Thrift also requires them in parameter declarations. One
requirement that these two implementations share is the fact that assigning a
tag is a permanent action, in the sense that the tag number becomes reserved.
Future revisions of the API must not reutilize such tags, and instead add new
ones in case of need and comment out the ones that become deprecated. For the
sake of compatibility, Poly adheres to this principle.

### Tagged Field Lists

```
tagged-field-list   =  tagged-field
tagged-field-list   =/ tagged-field ","
tagged-field-list   =/ tagged-field "," tagged-field-list
```

A Tagged Field List is a non-empty sequence of Tagged Fields separated by
commas, and optionally terminated by an additional comma (for convenience).
It's noticeable that elements may appear out of order, and that the numeric
identifiers are merely semantic.

### Location

```
location    =  "path"
location    =/ "query"
location    =/ "header"
location    =/ "cookie"
```

The Location attribute refers to where a value is expected to appear in the
context of the HTTP protocol. For example, a `header` qualifier indicates that
the value is expected to be passed in the headers. The following table lists
available location modifiers:

| Location | Description                                                 |
|----------|-------------------------------------------------------------|
| body     | The parameter is passed in the body.                        |
| path     | The parameter is passed in the URL, in the path component.  |
| query    | The parameter is passed in the URL, in the query component. |
| header   | The parameter is passed as a header.                        |
| cookie   | The argument is passed as a cookie.                         |

Some additional comments should be added:

* `path` applies to parameterized endpoints, where the value of the attribute
will replace some template variable. An example would be to replace `itemId`
in `items/{itemId}`.

* `query` values are encoded according to the `application/x-www-form-urlencoded`,
MIME type, as defined by [RFC 1866](https://tools.ietf.org/html/rfc1866).

Although there are several references, in these last comments, to encoding
schemes, those are not enforced by Poly, and thus it cannot be said that Poly
is introducing any concerns regarding encoding. Rather, these notes indicate
what is required by the specific RFC specifications that respect each item.
Therefore, these notes are not normative, but basically an observation of a
normalization that has already occurred elsewhere.

### Parameters

```
parameter   = location field

```

Parameters extend on Fields by specifying the Location in which the field is to
appear. This means that Parameters are constructs that are meant for use within
the context of some API related declaration, such as Endpoints. The Location
appears first in the construct.

### Tagged Parameters

```
tagged-parameter    = tag ":" parameter
```

A Tagged Parameter introduces a tag number to a Parameter. Similarly to Tagged
Fields, the tag creates a static order of elements, associating each field with
a numeric identifier. The purpose is also the same, to provide a numeric
semantic ordering that helps maintain backwards compatibility in the context of
binary encoding protocols. The main difference between the two is that a Tagged
Parameter list uses Parameters, instead of Fields, adding a Location attribute
to the tagged element.

### Tagged Parameter Lists

```
tagged-parameter-list   =  tagged-parameter
tagged-parameter-list   =/ tagged-parameter ","
tagged-parameter-list   =/ tagged-parameter "," tagged-parameter-list
```

A Tagged Parameter List declaration is identical to a Tagged Field List in
every way but one: that each element of the list is a Parameter, and thus
introduces a Location attribute. All other considerations apply, including
that the numeric identifiers are merely semantic and may appear out of order.

### Models

```
model   = "model" name "{" tagged-field-list "}"
```

A Model declaration defines a complex data type based on Primitive types or
other Models. Models can be used as input and output for endpoints, as well
as other situations where a data type is appropriate. When generating code,
engines may create classes and other data structures from models, as they
are meant to represent structures that hold data. As per Poly's promise, models
do not represent any form of encoding, and are rather an abstraction of data
relationships and model graphs.

Models do not support inheritance. Although that would probably improve code
reuse, there’s no viable way in which the order of the identifiers could be
maintained. For example, if a parent model is changed, how does that reflect
on the order of the child’s attributes?

One idea was to have the parent attributes being counted backwards, down from
zero. This, however, brings two problems, one related with how to encode the
negative numbers while still being space efficient, and the other with
inheritance at multiple levels; that is, if the numbering for the parent model
is reversed, how is the numbering for the grandparent model? Either way, this
feature would not be supported by most engines, and thus encapsulation is
encouraged instead.

### Verbs

```
verb    =  "get"
verb    =/ "put"
verb    =/ "post"
verb    =/ "delete"
verb    =/ "options"
verb    =/ "head"
verb    =/ "patch"
verb    =/ "trace"
```

### Paths

```
path        =  directory
path        =/ directory "/" path
path        =/ empty
directory   = 
```

Should we optionally allow leading and trailing slashes?

### Endpoints

```
endpoint    = verb path "(" tagged-parameter-list ")" type
```


Endpoint = Verb Path “(“ ParameterList “)” Type

Path = Directory | Directory “/” Path | ε
Directory = [a-zA-Z0-9$_@.&!*”’():;, ] | “-” | Escape | Hex
Escape = “\%” HexDigit HexDigit
HexDigit = [a-fA-F0-9]
ParameterList = Parameter “,” ParameterList | Parameter | ε









### Principles and Guidelines



1.3 Terminology
1.4 Overall Operation



2 Notational Conventions and Generic Grammar

2.1 Augmented BNF







# Future Work

* I think that I have already considered the RFC as to what are the valid
characters in a URL, but I'm not sure. The two should match.
