# project-detection

## What It Does
Detects the littlefactory project root by walking up the directory tree looking for a Factoryfile. This is the entry point for all commands that need to know which project they are operating on.

## Requirements
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

## Boundaries

## Gotchas
