This project, Iskendria, publishes scientific publications on a Hyperledger Sawtooth blockchain. The programming language is Golang. It has started as Martijn Dirkse's internship at PrivacyO Technologies B.V., Urmond, Limburg. The project is based on the report "Blockchain Technology for Publishing Scientific Papers", written as the final paper of the Blockchain Technology Consultant course described at https://www.3estack.io/en/training/. For more information and for a download link for the report, see www.iskendria.org.

Iskendria hashes scientific manuscripts on a Hyperledger Sawtooth blockchain. The blockchain thus offers cryptographic proof that a manuscript has been published. The following systems exist:

* Blockchain: Holds hashes of scientific manuscripts along with status information. Data maintained in the blockchain is the final truth.
* Client: Command-line application used by everyone who wants to edit the blockchain.
* Portal: Web application where users can find full-text manuscripts. This system is not the core business of the Iskendria team, but it has to be maintained as long as other parties do not make portals themselves.
* Major Tool: Command-line user interface used by the Iskendria team to maintain the system.

The following roles exist:

* Author: The author of manuscripts.
* Editor: Administrator of a scientific journal.
* Reviewer: Writes reviews about manuscripts, allowing editors to make informed decisions about them.
* Major: Member of the Iskendria team.

Note that Authors, Editors and Reviewers are all a Person.

More information about the functionality of this software can be found in the Markdown files (.md) in the root directory. These files can be formatted using LaTeX, a howto for doing this may be included later. The resulting .pdf files are checked in in the directory pdfdoc. Note that the information is a bit outdated; it was written before coding started.

On August 31 2019 I ported this project from our private gitlab repository to the present GitHub project. We also had Jira issues. First I ported these using https://github.com/hbrands/jira-issues-importer. This copied the issue names and preserved the issue numbers. The number of the Jira issue was put in the first comment of the corresponding GitHub issue. Then I manually put the Jira description as the second comment and copied the Jira comments. The creation times of the issues were preserved, but the creation times of the comments were lost.

Some issues are done in the sense that they are implemented, but have important comments. These issues are not closed on GitHub, but can be distinguished by the tag "documentation". So open issues with tag "documentation" do not require any coding, but require that the comments are properly saved.
