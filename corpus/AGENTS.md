# The corpus

* A corpus is a folder filled with human- and machine-readable .md files.
* The corpus should reflect the current state of the apps/services/libraries in the repo.
* "Current state" here means architecture, user roles and most importantly business rules.
* User roles are important: it can be physical human users, or other consuming services and systems.
* The level of detail should be such, that a product manager can read the files and verify that the repo does what it's supposed to do.
* The purpose of the corpus is to serve as _verification_. It does not replace tests and formal verification, it is the connection to human brains and big picture thinking.
* However, there may be some big picture technical architectural decisions as well. Mainly ones that are important for the secure and performant functioning of the software.

## Using the corpus

* When software engineers make changes to the code, relevant parts of the corpus should be updated, so the changed behavior of the software can be verified on a high level.
* Product managers, customer support etc can use an LLM (or their own brain) to answer questions they might have about the software product. They may also spot errors in the corpus, which may be corrected either only in the corpus (if the code already behaves correctly), or the change might require also fixing the code.
* It is up to LLMs to compare code and corpus and spot discrepancies and provide suggestions for improvement.

## Format of the corpus

* Avoid long markdown files with a ton of detail. Split up the text into logical chunks in their own files. Make an INDEX.md that provides a general overview of the folder (and keep it up to date if you change the file layout).
* Use markdown tables and mermaid diagrams as needed, but only when they truly add value.
* Don't start numbering features or decisions etc, loose references to other .md files is enough.
* Don't reference spec files or their numbering.
