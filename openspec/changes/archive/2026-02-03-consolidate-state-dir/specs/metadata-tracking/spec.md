## MODIFIED Requirements

### Requirement: JSON serialization
The system SHALL serialize metadata to JSON at `<state_dir>/run_metadata.json`.

#### Scenario: Metadata save after each iteration
- **WHEN** Each iteration completes
- **THEN** System writes updated RunMetadata to `<state_dir>/run_metadata.json` with proper ISO8601 timestamps

#### Scenario: Backward compatible JSON format
- **WHEN** System serializes metadata
- **THEN** JSON structure matches Python ciccio format exactly for backward compatibility

### Requirement: SaveMetadata receives config
The system SHALL pass full config to SaveMetadata for state directory access.

#### Scenario: SaveMetadata signature
- **WHEN** SaveMetadata is called
- **THEN** Function receives projectRoot, *config.Config, and metadata parameters
