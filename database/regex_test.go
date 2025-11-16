package database

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/storage-system/database/patterns"
)

func Test_Mysql_Regexp(*testing.T) {
    // Regex optimizado para Go (RE2)
	_patterns := patterns.CreateDSNPatterns()
    mysql_pattern := _patterns.MySQL

    re := regexp.MustCompile(mysql_pattern)

    tests := []string{
		"root:root@tcp(localhost:3388)/",
		"root:root@tcp(localhost:3388)/database_test",
        "raul:secreto123@tcp(192.168.1.100)/miapp",
        "admin:p@ssw0rd@udp(mysql.raul94a.es:3306)/tienda?env=prod&debug=false",
        "user:pass@tcp(localhost)/testdb?modo=dev",
        "bob:muy-seguro_2025@tcp(10.0.0.5)/backup?ssl=true&timeout=30",
        // Fails
        "raul@tcp(host)/db",
        "raul:pass@http(host)/db",
        "raul:pass@tcp()/db",
        "raul:pass@tcp(host)/db?",
        "raul:pass@tcp(host)/db?key",
        "raul:pass@tcp(host)/db?k=v&w",
    }

    for _, s := range tests {
        if re.MatchString(s) {
            fmt.Printf("PASS → %s\n", s)
        } else {
            fmt.Printf("FAIL → %s\n", s)
        }
    }
}