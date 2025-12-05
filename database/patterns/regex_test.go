package patterns

import (
	"regexp"
	"testing"

)

func Test_Mysql_Regexp(t *testing.T) {
    // Regex optimizado para Go (RE2)
	_patterns := CreateDSNPatterns()
    mysql_pattern := _patterns.MySQL

    re := regexp.MustCompile(mysql_pattern)

    shouldFailTests := []string{
		"raul@tcp(host)/db",
        "raul:pass@http(host)/db",
        "raul:pass@tcp()/db",
        "raul:pass@tcp(host)/db?",
        "raul:pass@tcp(host)/db?key",
        "raul:pass@tcp(host)/db?k=v&w",
        "root:root@tcp(localhost:3388)/",
        "admin:p@ssw0rd@udp(mysql.raul94a.es:3306)/tienda?env=prod&debug=false",


    }

    shouldPassTests := []string {
		"root:root@tcp(localhost:3388)/database_test",
        "raul:secreto123@tcp(192.168.1.100)/miapp",
        "admin:passw0rd@udp(mysql.raul94a.es:3306)/tienda?env=prod&debug=false",
        "user:pass@tcp(localhost)/testdb?modo=dev",
        "bob:muy-seguro_2025@tcp(10.0.0.5)/backup?ssl=true&timeout=30",
        "test:testpwd@tcp(dev.api.es:44444)/api_database",
    }

    for _, s := range shouldPassTests {
        if !re.MatchString(s) {
            t.Fatalf("Match MySQL Regex failed: `shouldPassTests` %s", s)
        } 
    }

    for _, s := range shouldFailTests {
        if re.MatchString(s) {
            t.Fatalf("Do not match MySQL Regex failed: `shouldFailTests` %s", s)
        } 
    }
}