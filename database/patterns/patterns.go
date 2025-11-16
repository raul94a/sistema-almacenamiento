package patterns



type _DatabaseStringConnectionPatterns struct {
	MySQL string
}

var  patterns *_DatabaseStringConnectionPatterns


func CreateDSNPatterns() *_DatabaseStringConnectionPatterns{
	if patterns != nil {
		return patterns
	}
    mysql := `^[^:\s]+:[^@\s]+@(?:tcp|udp)\([^)]+\)/[^?\s]+(?:\?[^=&\s]+=[^&\s]+(?:&[^=&\s]+=[^&\s]+)*)?$`

	_patterns := &_DatabaseStringConnectionPatterns{
		MySQL: mysql,
	}

	patterns = _patterns
	return patterns;
}