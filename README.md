This project, Alexandria, publishes scientific publications on a Hyperledger Sawtooth blockchain. It has started as Martijn Dirkse's internship at PrivacyO Technologies B.V., Urmond, Limburg. The project is based on the report "Blockchain Technology for Publishing Scientific Papers", written as the final paper of the Blockchain Technology Consultant course described at https://www.3estack.io/en/training/.

Alexandria hashes of scientific manuscripts on a Hyperledger Sawtooth blockchain. The blockchain thus offers cryptographic proof that a manuscript has been published. The following systems exist:

* Blockchain: Holds hashes of scientific manuscripts along with status information. Data maintained in the blockchain is the
final truth.
* Client: Command-line application used by everyone who wants to edit the blockchain.
* Portal: Web application where users can find full-text manuscripts. This system is not the core business of the Alexandria
team, but it has to be maintained as long as other parties do not make portals themselves.
* Major Tool: Command-line user interface used by the Alexandria team to maintain the system.

The following roles exist:

* Author: The author of manuscripts.
* Editor: Administrator of a scientific journal.
* Reviewer: Writes reviews about manuscripts, allowing editors to make informed decisions about them.
* Major: Member of the Alexandria team.

Note that Authors, Editors and Reviewers are all a Person.
