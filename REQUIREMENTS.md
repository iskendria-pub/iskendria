# Alexandria Detailed Requirements

This document presents detailed requirements for Alexandria. Please read README.md first.

## Persons

AX-10: The following properties should be maintained about each person:

* id.
* creation timestamp (AX-5520).
* modification timestamp (AX-5520).
* public key.
* private key.
* name.
* email address.
* isMajor.
* isSigned: Boolean indicating whether the Alexandra team knows the person.
* balance.
* hasBiography (boolean).

AX-20: The Major Tool should allow only majors to create new persons.

AX-30: A person should be identified by her id. The public key is not a good unique id because a user must be able to update that, see AX-60. See AX-5060 about id fields.

AX-50: The private key of a person should not be stored on the blockchain and it should not be sent over the network.

AX-60: The Client should allow each person to update her key pair, Name or Email address.

AX-70: With the Major tool a major should be able to set or reset the isSigned property and the isMajor property of each person. When a person updates her personal data, these properties should not be reset. A Major can reset them if needed.

AX-80: The balance is an integer value. The balance is decreased when the person uses the system. She should not be able to use the system when her balance is zero.

AX-85: Each major should be allowed to add balance to every person.

AX-90: Each person has the following optional properties:

* organization.
* telephone number.
* address.
* postal code.
* country.
* biography.
* biography format, see AX-1030.
* extraInfo.

The Client tool should allow her to add and edit these.

AX-95: The Client tool should allow a person to clear her biography. The hasBiography property should reflect whether there is a biography.

AX-100: The Biography should be treated as a document, see AX-1000.

AX-110: The Client tool should be able to generate a public/private key pair.

AX-120: The Client tool should allow every person to see all of her information.

AX-130: The Major tool should allow every major to see all information of every person.

AX-140: The modification timestamp of a person should be updated when one or more of the properties of AX-10 or AX-90 is updated.

AX-150: Although the Client Tool and the Portal tool provide curriculum vitae information (AX-4050 and AX-4340), a person does not have a list of manuscripts or a list of editorships. CV functionality should be provided by listening to events coming from the blockchain, not by keeping the information on the blockchain itself.

## Documents

AX-1000: The following kinds of documents exist in the system:

* Biography.
* Manuscript.
* Journal description.
* Review text.

AX-1010: Each document is hashed, only the hash is stored on the Blockchain. The hash algorithm is SHA-512.

AX-1020: The Portal should allow everyone to upload the full-text of a document, provided that its hash agrees with the hash on the blockchain.

AX-1030: Two formats are allowed for documents: PDF or UTF-8 text. UTF-8 text is relevant because it can be rendered easily within the Portal.

AX-1040: The following properties of a document should be maintained on the Blockchain:

* The hash.
* The format.

AX-1050: Documents themselves are not identified, but only the entity it is part of. For example, a biography is treated as a property of a person.

AX-1090: Manuscripts are always PDF documents. Reviews, Journal Descriptions and Biographies are always UTF-8 text.

## Manuscript

AX-1500: A Manuscript has the following properties:

* id.
* creation timestamp.
* modification timestamp.
* hash.
* manuscript format, see AX-1030.
* manuscript thread id.
* version number (one positive number).
* commit message (not hashed, text goes on the blockchain).
* title.
* list of AuthorInfo items.
* status.
* journalId.
* volume id.
* first page.
* last page.

AX-1510: An AuthorInfo item has the following properties:

* The author.
* didSign: True if the author signed for being author.
* authorSeq: The order of the author list of a publication is very significant.

AX-1513: Each manuscript is part of a manuscript thread. A manuscript thread has the following properties:

* id
* manuscripts
* isReviewable

AX-1516: The manuscripts property of a manuscript thread is an ordered list of manuscript ids. The order corresponds to the submission order.

AX-1520: The status of a manuscript can be: INIT, NEW, REVIEWABLE, REJECTED, PUBLISHED or ASSIGNED. These mean:

* INIT: Note that the list of authors is specified when a manuscript is created. All authors have to sign for being author. This state indicates that some authors did not sign. 
* NEW: All authors signed for being author, but there is no editor signature yet that the manuscript is worth reviewing.
* REVIEWABLE: The journal editor should do a quick check to see whether the manuscript is worth reviewing. This state indicates that the editor signed for this.
* REJECTED: The editor rejected this manuscript for publication. This does not prevent anyone to submit a new version.
* PUBLISHED: The manuscript is published but has not been assigned to a volume.
* ASSIGNED: The manuscript is published in a volume.

AX-1540: The Client tool should allow an existing person to submit a new manuscript. The following information should be provided:

* The PDF text.
* The title.
* The list of authors.
* The journal.
* An optional commit message. This message allows authors, editors and reviewers to reference different versions in a meaningful way.

The system will then set remaining properties as follows:

* The id is generated.
* The creation timestamp and the modification timestamp are set to the current time according to the host preparing the command to be sent to the blockchain. See AX-5500.
* The hash is calculated from the PDF text.
* The manuscript format should be set to PDF.
* A new manuscript thread is created with the manuscript included as the first item.
* The version number is set to 1.
* The status is set to INIT.
* When the person signing the transaction is also in the list of authors, set didSign=true for that AuthorInfo.
* The volume id, the first page and the last page are made empty.

AX-1550: The Client tool should allow every author of a manuscript to submit a new version, provided that the status is not PUBLISHED or ASSIGNED. The following information should be included:

* The PDF text.
* A mandatory commit message.
* The previous version.
* The title.
* The list of authors.

The system will set the remaining properties as follows:

* The creation timestamp and the modification timestamp are set to the current time according to the host that prepares the command to be sent to the blockchain.
* The hash is calculated from the PDF text.
* The manuscript format should be set to PDF.
* The version number is one more than the version number of the previous version.
* The new manuscript is added to the manuscript thread.
* The status is set to INIT.
* the journal equals the journal of the previous version.
* The volume, the first page and the last page are made empty.

AX-1560: The Client tool should allow everyone in the list of authors to sign for being author. When every author has signed, the status should go to NEW or REVIEWABLE. The state will be NEW when the manuscript thread has isReviewable = false. Otherwise, the state will be REVIEWABLE.

AX-1570: The Client tool should allow an editor to change the manuscript status from NEW to REVIEWABLE. An editor applies here to the journal of the manuscript. This change is applied to all manuscripts in the manuscript thread. All manuscripts in the thread with status NEW should be updated to REVIEWABLE. The isReviewable property of the thread should be set (AX-1513).

AX-1580: The Client tool should allow reviewers to write reviews. This is allowed for documents that are not in state INIT or NEW.

AX-1590: The Client tool should allow an editor to change the manuscript status from REVIEWABLE to REJECTED or PUBLISHED. She is required to reference all the reviews that guided her decision.

AX-1600: The Client tool should allow an editor to assign a volume to a published manuscript. The following information should be included:

* The volume id.
* The first page.
* The last page.

The system should update the manuscript state to ASSIGNED.

AX-1610: The state machine for manuscripts should be such that the order of manipulations is not important.

AX-1620: The Client should allow everyone to see metadata about every Manuscript.

AX-1630: The modification timestamp of a manuscript should be updated when:

* An author signs for being author (AX-1560).
* An editor changes the state (AX-1570, AX-1590, AX-1600)

## Journal

AX-2000: A Journal should have the following mandatory properties:

* id.
* creation timestamp (see AX-5520).
* modification timestamp (see AX-5520).
* title.
* isSigned (boolean).
* hasDescription (boolean).
* descriptionHash.
* descriptionFormat (see AX-1030).
* List of EditorInfo.

AX-2010: An EditorInfo has properties Person and EditorState. The editor state can be ACCEPTED or PROPOSED. These states allow the editors of a journal to resign and to assign colleagues, while each added editor should have both her own signature and the signature of an existing editor.

AX-2030: The Client tool should allow every known person to create a journal. The following information should be provided:

* The title.
* An optional description.

AX-2040: The Client tool should allow each editor to update the title or the journal description. She should also be able to remove the journal description.

AX-2050: The Major tool should allow majors to set or reset the isSigned property of a journal. This indicates whether she approves the existence of the journal. This way a Major can ban rogue journals. This is not about banning bad science, but more about protecting trademarks. An Alexandria journal titled "Nature" could be banned because the journal Nature already exists outside Alexandria.

AX-2060: The Client tool should allow each editor of a Journal to resign. This just removes her from the list of editors.

AX-2070. The Client tool should allow each editor of a Journal to invite another editor. That other person is added to the EditorInfo list with status PROPOSED. A proposed editor does not have the rights of an editor yet.

AX-2080: The Client tool should allow each PROPOSED editor to sign, taking the state to ACCEPTED.

AX-2110: The Client tool should allow each editor to create a Volume. A Volume has the following properties:

* id.
* creation timestamp.
* id of journal.
* issue string.

AX-2120: A manuscript can only be in a volume when its state is ASSIGNED. An ASSIGNED journal should be in exactly one Volume.

AX-2130: Volumes should not be editable.

AX-2140: The modification timestamp of a journal should be updated when:

* The title is changed.
* The isSigned property is changed.
* The description is set, cleared or changed.
* One or more of the EditorInfo items is changed.

## Reviews

AX-2500: A Review, see AX-1580, has the following properties:

* id.
* creation timestamp.
* The manuscript id it is about.
* One author.
* Hash of text.
* Format of text.
* Judgement.
* isUsedByEditor (see AX-1590).

AX-2510: The Judgement in a review can be "ACCEPTED" or "REJECTED". There is no judgement for review requested, because a new version is treated here as a new manuscript.

AX-2530: A review is not editable, except the isUsedByEditor property.

## Credits

AX-3000: The following actions should cost credit:

* Major edits settings.
* Major changes person authorization.
* Major changes journal authorization.
* Major creates person.
* Person edits her data.
* Author submits new Manuscript.
* Author submits new version of Manuscript.
* Author signs for being author.
* Reviewer submits review about Manuscript.
* Editor makes manuscripts eligible for reviews.
* Editor rejects manuscript.
* Editor publishes manuscript.
* Editor assigns manuscript to volume.
* Editor starts new journal.
* Editor edits journal properties.
* Editor invites other editor.
* Editor accepts editorship.
* Editor creates volume.

AX-3010: The Major tool should allow each major to adjust prices. Each action of AX-3000 should have its own price.

AX-3015: The price list has an associated creation timestamp and an associated modification timestamp.

## User interface

### Client tool

AX-4000: The Client tool is an interactive command-line tool. This means that it runs in a command shell and that you do not exit after a command is completed. After a command is completed, a prompt is shown asking for a new command. There should be a command to exit the Client tool.

AX-4010: The Client tool should have subcommands "general", "author", "reviewer" and "editor" that group all possible commands. These should act as modes. For example, when you apply the "author" command, all following commands are interpreted as author commands. The exit command should leave the subcommand.

AX-4020: The Client tool should have a help command that shows context sensitive help.

AX-4030: The Client should have a login command, allowing the user to impersonate a Person on the blockchain. The login command should give feedback on who you are.

AX-4050: The Client tool should offer a command "cv", which shows the biography, submitted manuscripts and editorships of journals of the logged-in person.

### Portal tool

AX-4300: The Portal tool has a screen that shows all journals by title. Journals with isSigned=true should stand out.

AX-4310: When the user clicks a journal in the Portal tool, a screen with journal details will appear. It shows the properties of the journal, see AX-2000. It has separate link for each volume and additional links pointing to reviewable manuscripts and published unassigned manuscripts.

AX-4320: The journal details screen shows the journal description. There are three cases:

* There is a journal description hash on the blockchain, but the corresponding text is not known to the portal. In this case, there is a Submit button allowing the user to upload the text.
* There is a journal description hash on the blockchain and there is a corresponding text. Then the text is shown and there is also a "Verify" button.
* There is no journal description hash on the blockchain. Then the screen indicates something like "not available".

AX-4330: Every editor, author or reviewer mentioned in Portal is a link. Clicking such a link should produce a CV of the person.

AX-4335: Editors, Authors and Reviewers who are signed persons should stand out.

AX-4340: The CV screen shows person properties, see AX-10 and AX-90.

AX-4350: There are three cases for the person Biography:

* There is a person Biography on the blockchain (hasBiography = true), but the text is not known to Portal. In this case a Submit button is shown.
* There is a person Biography and the text is known. In this case Portal shows a Verify button.
* There is no Biography (hasBiography = false). Then a message like "not available" is shown, no buttons.

AX-4355: The CV screen should include manuscripts submitted and editorships held.

AX-4360: The manuscript-related links in the journal details screen point to a table of contents. A table of contents shows a list of manuscripts. Each entry in the list shows at least the title and the first few authors.

AX-4370: The Portal tool has a manuscript detail screen. It shows all Manuscript properties.

AX-4380: The Manuscript detail screen has two cases related to the full-text:

* The full text is not known to the Portal. In this case there is a Submit button.
* The full text is known. In this case there is a Verify button.

AX-4390: The Manuscript details screen has links to reviews and previous versions.

AX-4400: Portal has a reviews list screen. For each review, there are two cases:

* The full text of the review is known. In this case, there is a verify button.
* The full text of the review is not known. In this case, there is a Submit button.

AX-4410: The reviews list screen should highlight reviews that the editor used to take her decision. The decision whether to reject or publish (AX-1590).

AX-4420: The reviews list screen should include all versions and their reviews.

### Major tool

AX-4700: The Major tool is an interactive command-line application.

AX-4710: When the blockchain is empty, the Major tool should allow everyone to bootstrap the blockchain. The user should provide all price levels mentioned in AX-3000. The user should also provide person create information, see AX-10 and AX-90. This will result in a bootstrapped blockchain with one person who is major. The key of the person is the key that signed the bootstrap request.

## Blockchain

AX-5000: Every update of a value on the blockchain should mention the current value. For example, when a person edits her name then the expected current name is in the request. When the current name does not match the expected value, then the transaction should be rejected. This allows the blockchain to detect conflicts.

AX-5010: A Hyperledger Sawtooth transaction processor can return an error. Care should be taken to return the right error, because otherwise the Validator will cause useless retries. In the xo transaction processor in the Hyperledger Sawtooth docs, a processor.InvalidTransactionError is used.

AX-5020: When an ordered list on the blockchain (e.g. all versions of a Manuscript) is appended, then the request should reference the current last item. When the last current item is not correct anymore, then the transaction should be rejected. We do not allow multiple successors of the same item in a list.

AX-5040: Each transaction that costs credits should include the price. This way, the blockchain can resolve conflicts between ordinary transactions and price changes.

AX-5050: The following items are identified with property named "id" that has a unique value:

* Person
* Manuscript
* Manuscript thread
* Journal
* Volume
* Review

 The id-value is used as an address in Sawtooth and should therefore be globally unique. For example a person id should not be used to identify an manuscript, a manuscript thread or anything else. Editors, authors and reviewers reference ids to manipulate the blockchain.

AX-5060: When a user creates an object with an id, she does not provide that id herself. The client tool or the major tool is responsible for generating the id. The id should be the blockchain address where the object is stored.

AX-5070: Each Hyperledger Sawtooth address should contain only one item. In theory, generating an address for a new item can result in an address collision. This is solved by basing generated addresses on a uuid. When a collision occurs, a new uuid can be generated resulting in a new address. This can be done until the generated address is free.

AX-5080: There is a fixed address that is filled when the blockchain is bootstrapped. This address will hold the list of prices, AX-3000. Using a fixed address allows the blockchain to check whether there is bootstrap information.

## Timestamps

AX-5500: Creation and modification timestamps are generated by the host that sends requests to the blockchain.

AX-5510: At the moment, the blockchain is not required to verify timestamps. In a later version of the software, the blockchain will check that client-generated timestamps are close to the real time.

AX-5520: Timestamps are sent and stored as seconds since UNIX epoch. This way, the time zone used by a host does not influence the value sent.
