# Change Log
All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [1.0.0] - 2015-10-28
### Changed
- Ran golint on project resulting in minor code changes.
- Restructured tests into internal and external. internal tests test un-exported 
functionality, external tests test client facing functionality.
- Updated README to clarify installation instructions.

### Added
- Created method IsAuthed to check whether user is authenticated.
- Library is now versioned pinboard.Version().
