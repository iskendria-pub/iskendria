# alexandria
Scientific publications published on a Hyperledger Sawtooth blockchain. This project is part of Martijn Dirkse's internship at
The Privacy Factory, Urmond, Limburg. The project is based on the report
"Blockchain Technology for Publishing Scientific Papers", written as the final paper of the Blockchain Technology Consultant
course described at https://www.3estack.io/en/training/.

This application publishes hashes of scientific manuscripts on a Hyperledger Sawtooth blockchain. The blockchain thus offers
cryptographic proof that a manuscript has been published. The following systems exist:
* Blockchain: Holds hashes of scientific manuscripts along with status information. Data maintained in the blockchain is the
final truth.
* Client: Command-line application used by everyone who wants to edit the blockchain.
* Portal: Web application where users can find full-text manuscripts. This system is not the core business of the Alexandria
team, but it has to be maintained as long as other parties do not make portals themselves.
* Major: Command-line user interface used by the Alexandria team to maintain the system.

The following roles exist:
* Author: The author of manuscripts.
* Editor: Administrator of a scientific journal.
* Reviewer: Writes reviews about manuscripts, allowing editors to make informed decisions about them.
* Major: Member of the Alexandria team.

Note that Authors, Editors and Reviewers are all a Person.

# Detailed requirements

## Persons

AX-10: Each person in the system has the following mandatory attributes:
* UUID
* Public key.
* Private key.
* isMajor.
* Name.
* Email address.
* isSigned: Boolean indicating whether the Alexandra team knows the person.
* saldo.
* hasBiography (boolean).

AX-20: The Client should allow anyone to add a person to the system, with isMajor = false.

AX-30: A person should be identified by her UUID. The public key is not a good unique id because a user must be able to
update it, see AX-60.

AX-40: It is not allowed that a public key is shared between multiple persons.

AX-50: The private key of a person should not be stored on the blockchain and it should not be sent
over the network.

AX-60: The Client should allow each person to update her key pair, Name or Email address.

AX-70: With the Major tool a major should be able to set or reset the isSigned property and the isMajor property of each person.
When a person updates her personal data, these properties should not be reset. A Major can reset them if needed.

AX-80: The saldo is an integer value. The saldo is decreased when the person uses the system. She
should not be able to use the system when her saldo is zero. She should be able to buy new credits to continue usage of the
system.

AX-90: Each person has the following optional attributes:
* Institution.
* Telephone number.
* Address.
* Zip code.
* Country.
* Government ID.
* Biography.

The Client tool should allow her to add and edit these.

AX-95: The Client tool should allow a person to clear her biography. The hasBiography property should reflect whether there
is a biography.

AX-100: The Biography should be treated as a document, see AX-1000.

AX-110: The Client should be able to generate a public/private key pair.

AX-120: The Client tool should allow every person to see all of her information.

AX-130: The Major tool should allow every major to see all information of every person.

## Documents

AX-1000: The following kinds of documents exist in the system:
* Biography.
* Manuscript.
* Journal description.
* Review text.

AX-1010: Each document is hashed, only the hash is stored on the Blockchain.

AX-1020: The Portal should allow everyone to upload the full-text of a document, provided that its hash agrees with the hash
on the blockchain.

AX-1030: Two formats are allowed for documents: PDF or UTF-8 text. UTF-8 text is relevant because it can be rendered easily
within the Portal.

AX-1040: The following properties of a document should be maintained on the Blockchain:
* The hash.
* The document kind.
* The format.
* Optional context values.

AX-1050: Documents are identified by their hash; no two documents can exist with the same hash. When someone causes a hash
collision by uploading a document, the document should be refused by the Blockchain with an error message. The user is
required to apply a trivial modification to the document to have another hash. This is acceptable because the chance that a hash
collisions occurs is negligible.

AX-1060: UTF-8 text documents are trimmed (leading and trailing space characters removed) before they are hashed. The reason is
that text editors add trailing space characters sometimes as was observed with the editor vim.

AX-1070: UTF-8 text documents can sometimes be short, for example reviews with the text "Good work". Therefore UTF-8 texts
should be prepended by context values to make it unique. Here is an example: Assume that the user with UUID abcdef submits the
review text "Good work\n\n\n" about the manuscript with hash 01234567. Then first AX-1060 is applied, resulting in "Good work".
Then the text "01234567\nabcdef\nGood work" is hashed.

AX-1080: Here are the context values to prepend to a UTF-8 text before hashing:
* A Biography is prepended by the person id it is about.
* A Journal Description is prepended by the Journal id.
* A Review Text is prepended by the hash of the manuscript followed by the id of the reviewer.

AX-1090: Manuscripts are always PDF documents. Reviews, Journal Descriptions and Biographies are always UTF-8 text.

## Manuscript

AX-1500: A Manuscript has the following properties:
* Hash.
* Version number (one positive number).
* Commit message (not hashed, text goes on the blockchain).
* Hash of previous version (if version >= 2).
* Hash of next version.
* Title.
* List of AuthorInfo items.
* Status.
* Journal.
* List of Review.
* Volume.
* First page.
* Last page.

AX-1510: An AuthorInfo item has the following properties:
* The author.
* didSign: True if the author signed for being author.
* authorSeq: The order of the author list of a publication is very significant.

AX-1520: The status of a manuscript can be: INIT, NEW, REVIEWABLE, REJECTED, PUBLISHED or ASSIGNED. These mean:
* INIT: The list of authors is being established.
* NEW: The information is complete.
* REVIEWABLE: The editor declared that the manuscript is suited for the journal, provided it gets the right positive reviews.
* REJECTED: The editor rejected this manuscript for publication. This does not prevent anyone to submit a new version.
* PUBLISHED: The manuscript is published but has not been assigned to a volume.
* ASSIGNED: The manuscript is published in a volume.

AX-1530: The id of a Manuscript is its hash.

AX-1540: The Client tool should allow an existing person to submit a new manuscript. The following information should be
provided:
* The PDF text.
* The title.
* The list of authors.
* The journal.
* An optional commit message.

The system will then set remaining properties as follows:
* The Hash is calculated from the PDF text.
* The Version Number is set to 1.
* The Status is set to INIT.
* When the person signing the transaction is also in the list of authors, set didSign=true for that AuthorInfo.

AX-1550: The Client tool should allow every author of a manuscript to submit a new version, provided that the status is not PUBLISHED or ASSIGNED. The following information should be included:
* The PDF text.
* A mandatory commit message.
* The previous version.
* The title.
* The list of authors.

The system will set the remaining properties as follows:
* The hash is calculated from the PDF text.
* The Version Number is one more than the version number of the previous version.
* The status is set to INIT.
* the Journal equals the journal of the previous version.

The system will also assign the nextHash property of the previous version.

AX-1560: The Client tool should allow everyone in the list of authors to sign for being author. When every author has signed,
the status should go to NEW or REVIEWABLE.

AX-1570: The Client tool should allow an editor to change the manuscript status from NEW to REVIEWABLE. An editor applies
here to the journal of the manuscript. This change is propagated to previous and future versions.

AX-1580: The Client tool should allow reviewers to write reviews. This is allowed for documents that are not in state
INIT or NEW.

AX-1590: The Client tool should allow an editor to change the manuscript status from REVIEWABLE to REJECTED or PUBLISHED.
She is required to reference all the reviews that guided her decision.

AX-1600: The Client tool should allow an editor to assign a volume to a published manuscript. The following information
should be included:
* The volume id.
* The first page.
* The last page.

The system should update the manuscript state to ASSIGNED.

AX-1610: The state machine for manuscripts should be such that the order of manipulations is not important. Each sequence
of manipulations should result in the same bytes on the blockchain.

AX-1620: The Client should allow everyone to see metadata about every Manuscript.

## Journal

AX-2000: A Journal should have the following mandatory properties:
* UUID.
* Title.
* isSigned (boolean).
* hasDescription (boolean).
* Optional journal description, treated as document, see AX-1000.
* List of EditorInfo.
* List of reviewable manuscripts.
* List of published unassigned manuscripts.
* List of volumes.

AX-2010: An EditorInfo has properties Person and EditorState. The editor state can be ACCEPTED or PROPOSED.
These states allow the editors of a journal to resign and to assign colleagues, while each added editor
should have both her own signature and the signature of an existing editor.

AX-2020: A Journal is identified by its UUID.

AX-2030: The Client tool should allow every known person to create a journal. The following information should be provided:
* The title.
* An optional description.

The system should create a UUID and include it in the request to the Blockchain. The system should then set the provided
properties and make the signer the first editor of the journal, with state ACCEPTED.
The property hasDescription should be true iff a description was provided.

AX-2040: The Client tool should allow each editor to update the title or the journal description. She should
also be able to remove the journal description.

AX-2050: The Major tool should allow majors to set or reset the isSigned property of a journal. This indicates whether she 
approves the existence of the journal. This way a Major can ban rogue journals. This is not about banning bad science, but more
about protecting trademarks. An Alexandria journal titled "Nature" could be banned because the journal Nature already exists
outside Alexandria.

AX-2060: The Client tool should allow each editor of a Journal to resign. This just removes her from the list of editors.

AX-2070. The Client tool should allow each editor of a Journal to invite another editor. That other person is added to the
EditorInfo list with status PROPOSED. A proposed editor does not have the rights of an editor yet.

AX-2080: The Client tool should allow each PROPOSED editor to sign, taking the state to ACCEPTED.

AX-2090: All manuscripts that have state REVIEWABLE that are also the last version should occur in the reviewable
manuscripts list of their journal.

AX-2100: All manuscripts that have state PUBLISHED should occur in the published unassigned list of their journal.

AX-2110: The Client tool should allow each editor to create a Volume. A Volume has the following properties:
* UUID of journal.
* issue string.
* list of manuscripts.

AX-2120: A manuscript can only be in a volume when its state is ASSIGNED. An ASSIGNED journal should be in exactly
one Volume.

AX-2130: Volumes should not be editable.

## Reviews

AX-2500: A Review, see AX-1580, has the following properties:
* The manuscript it is about.
* One author.
* Hash of text.
* Judgement.
* isUsedByEditor (see AX-1590).

AX-2510: The Judgement in a review can be "ACCEPTED" or "REJECTED". There is no judgement for review requested,
because a new version is treated here as a new manuscript.

AX-2520: The order of reviews is not important. Each order of adding reviews should result in the same bytes stored on the
blockchain.

AX-2530: A review is not editable.

## Credits

AX-3000: The following actions should cost credit:
* Person edits her data.
* Author submits new Manuscript.
* Author submits new version of Manuscript.
* Reviewer submits review about Manuscript.
* Editor makes manuscripts eligible for reviews.
* Editor publishes manuscript.
* Editor assigns manuscript to volume.
* Editor starts new journal.
* Editor edits journal properties.
* Editor assigns other editor.
* Editor accepts editorship.
* Editor creates volume.

AX-3010: The Major tool should allow each major to adjust prices. Each action of AX-3000 should have its own price.

AX-3020: The Major tool should allow each major to add and withdraw credit to everyone.

## User interface

### Client tool

AX-4000: The Client tool is an interactive command-line tool. This means that it runs in a command shell and that you
do not exit after a command is completed. After a command is completed, a prompt is shown asking for a new command. There
should be a command to exit the Client tool.

AX-4010: The Client tool should have subcommands "general", "author", "reviewer" and "editor" that group all possible commands.
These should act as modes. For example, when you apply the "author" command, all following commands are interpreted as
author commands. The exit command should leave the subcommand.

AX-4020: The Client tool should have a help command that shows context sensitive help.

AX-4030: The Client should have a login command, allowing the user to impersonate a Person on the blockchain. The login command
should give feedback on who you are.

AX-4040: The Client tool should offer an Options within each of the four subcommands. The Options command shows what actions are
possible given the saldo of the logged-in user.

AX-4050: The Client tool should offer a command "cv", which shows the biography, submitted manuscripts and
editorships of journals of the logged-in person.

### Portal tool

AX-4300: The Portal tool has a screen that shows all journals by title. Journals with isSigned=true should stand out.

AX-4310: When the user clicks a journal in the Portal tool, a screen with journal details will appear. It shows the properties
of the journal, see AX-2000. It has separate link for each volume and additional links pointing to reviewable manuscripts and
published unassigned manuscripts.

AX-4320: The journal details screen shows the description. There are three cases:
* There is a journal description hash on the blockchain, but the corresponding text is not known to the portal. In this case,
there is a Submit button allowing the user to upload the text.
* There is a journal description hash on the blockchain and there is a corresponding text. Then the text is shown and
there is also a "Verify" button.
* There is no journal description hash on the blockchain. Then the screen indicates something like "not available".

AX-4330: Every editor, author or reviewer mentioned in Portal is a link. Clicking such a link should produce a CV of the person.

AX-4335: Editors, Authors and Reviewers who are signed persons should stand out.

AX-4340: The CV screen shows person properties, see AX-10 and AX-90.

AX-4350: There are three cases for the person Biography:
* There is a person Biography on the blockchain (hasBiography = true), but the text is not known to Portal. In this case
a Submit button is shown.
* There is a person Biography and the text is known. In this case Portal shows a Verify button.
* There is no Biography (hasBiography = false). Then a message like "not available" is shown, no buttons.

AX-4355: The CV screen should include manuscripts submitted and editorships held.

AX-4360: The manuscript-related links in the journal details screen point to a table of contents. A table of contents shows
a list of manuscripts. Each entry in the list shows at least the title and the first few authors.

AX-4370: The Portal tool has a manuscript detail screen. It shows all Manuscript properties.

AX-4380: The Manuscript detail screen has two cases related to the full-text:
* The full text is not known to the Portal. In this case there is a Submit button.
* The full text is known. In this case there is a Verify button.

AX-4390: The Manuscript details screen has links to reviews and previous versions.

AX-4400: Portal has a reviews list screen. For each review, there are two cases:
* The full text of the review is known. In this case, there is a verify button.
* The full text of the review is not known. In this case, there is a Submit button.

AX-4410: The reviews list screen should highlight reviews that the editor used to take her decision. The
decision whether to reject or publish (AX-1590).

AX-4420: The reviews list screen should include all versions and their reviews.

### Major tool

AX-4700: The Major tool is an interactive command-line application.

AX-4710: When the blockchain is empty, the Major tool should allow everyone to bootstrap the blockchain. The user should provide
all price levels mentioned in AX-3000. The key that signed is the bootstrap key.

AX-4720: The Major can use the Client tool to create a Person for herself that has the bootstrap key as key.
Then she can use the Major tool to declare that person Major. She signs that with the bootstrap key.

AX-4730: The Major tool should allow any Major to disable the bootstrap key. After that, the bootstrap key is still usable
as long as it is the key of a person having isMajor = true.

## Blockchain

AX-5000: Every update of a value on the blockchain should mention the current value. When the real current value differs
from the current value according to the request, then the transaction should be rejected. This allows the blockchain to
detect conflicts.

AX-5010: A Hyperledger Sawtooth transaction processor can return an error. Care should be taken to return the right error,
because otherwise the Validator will cause useless retries. In the xo transaction processor in the Hyperledger Sawtooth
docs, a processor.InvalidTransactionError is used.

AX-5020: When an ordered list on the blockchain (e.g. all versions of a Manuscript) is appended, then the request should
reference the current last item. When the last current item is not correct anymore, then the transaction should be rejected.

AX-5030: Unordered lists (e.g. the reviews of a Manuscript) should be stored with an artificial sort order. This way the
same inserts produce the same list, irrespective of the order in which the inserts are applied.

AX-5040: Each transaction that costs credits should include the price. This way, the blockchain can resolve conflicts
between ordinary transactions and price changes.
