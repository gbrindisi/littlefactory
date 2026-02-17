## MODIFIED Requirements

### Requirement: Project root detection via Factoryfile
The system SHALL locate project root by searching for Factoryfile.

#### Scenario: Factoryfile in current directory
- **WHEN** Current working directory contains Factoryfile
- **THEN** System uses current directory as project root

#### Scenario: Factoryfile in parent directory
- **WHEN** Current directory does not have Factoryfile but parent directory does
- **THEN** System walks up directory tree and uses first directory containing Factoryfile as project root

#### Scenario: No Factoryfile found
- **WHEN** No Factoryfile found in current or parent directories
- **THEN** System returns error indicating no littlefactory project found

## REMOVED Requirements

### Requirement: Project root detection via .beads directory
**Reason**: Replaced by Factoryfile detection - decouples project detection from task backend
**Migration**: Ensure Factoryfile exists in project root
