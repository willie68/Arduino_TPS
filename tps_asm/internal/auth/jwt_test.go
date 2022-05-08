package auth

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

// const testTokenSignature = "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
// const testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
const testToken = "Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJOVmlESFlDalNwNVZNV0J0bUhXZUoyV1JmMThKdGlTenBMemdhckhuOWxBIn0.eyJleHAiOjE2MTk4NTM3ODMsImlhdCI6MTYxOTg1MzQ4MywiYXV0aF90aW1lIjoxNjE5ODUzNDgyLCJqdGkiOiI1MTNkMTU4NC04Y2QyLTQ0NTUtODQ2NS0xNGRlZTc4MTViYzkiLCJpc3MiOiJodHRwczovL29uZS1kbS1kZXYuZWFzeTAyLnByb2FjdGNsb3VkLmRlL2F1dGgvcmVhbG1zL2MxMzNkajU1YmN1bGplMm5qMmcwIiwiYXVkIjpbImJyb2tlciIsImFjY291bnQiXSwic3ViIjoiODNlOTQ2NzItOTRmOC00NzYwLWE2M2YtY2UwZjA2OWExMzUxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiZXMtZG0iLCJub25jZSI6ImRmZDM2NjU3LTBhYzktNGQyNi05NGQ4LWE2MjBkZWY0ZDc2MyIsInNlc3Npb25fc3RhdGUiOiIxOTc2YmU1ZS0wNTdkLTQ3ODUtOTk0NS1jMTA4ODExZGYxYzAiLCJhY3IiOiIxIiwiYWxsb3dlZC1vcmlnaW5zIjpbImh0dHA6Ly9sb2NhbGhvc3Q6ODA4MS8qIiwiaHR0cDovL2xvY2FsaG9zdDo4MDgwLyoiLCJodHRwczovL29uZS1kbS1kZXYuZWFzeTAyLnByb2FjdGNsb3VkLmRlLyoiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ0ZW5hbnRhZG1pbi1jMTMzZGtkNWJjdWxqZTJuajJpMCIsInRlbmFudGFkbWluLWMxMzNkdmQ1YmN1bGplMm5qMmpnIiwidGVuYW50dXNlci1jMTMzZHZkNWJjdWxqZTJuajJqZyIsIm9mZmxpbmVfYWNjZXNzIiwidGVuYW50dXNlci1jMTMzZGtkNWJjdWxqZTJuajJpMCIsInRlbmFudHVzZXItYzEzM2R1NTViY3VsamUybmoyaWciLCJ0ZW5hbnRhZG1pbi1jMTMzZHU1NWJjdWxqZTJuajJpZyIsInVtYV9hdXRob3JpemF0aW9uIiwidGVuYW50YWRtaW4tYzEzM2R1dDViY3VsamUybmoyajAiLCJ0ZW5hbnR1c2VyLWMxMzNkdXQ1YmN1bGplMm5qMmowIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYnJva2VyIjp7InJvbGVzIjpbInJlYWQtdG9rZW4iXX0sImFjY291bnQiOnsicm9sZXMiOlsibWFuYWdlLWFjY291bnQiLCJtYW5hZ2UtYWNjb3VudC1saW5rcyIsInZpZXctcHJvZmlsZSJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJ0ZW5hbnRzIjpbImMxMzNkdmQ1YmN1bGplMm5qMmpnIiwiYzEzM2RrZDViY3VsamUybmoyaTAiLCJjMTMzZHV0NWJjdWxqZTJuajJqMCIsImMxMzNkdTU1YmN1bGplMm5qMmlnIl0sImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwibmFtZSI6IldpbGZyaWVkZCBLbGFhcyIsInN5c3RlbWlkZW50aXRpZXMiOlt7InN5c3RlbSI6ImMxMzNmdHJiMTMyaWJqZ25lbHZnIiwidGVuYW50IjoiYzEzM2R2ZDViY3VsamUybmoyamciLCJ1c2VyaWQiOiJydXBwIn1dLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJ3aWxmcmllZC5rbGFhc0BlYXN5LmRlIiwiZ2l2ZW5fbmFtZSI6IldpbGZyaWVkZCIsImZhbWlseV9uYW1lIjoiS2xhYXMiLCJlbWFpbCI6IndpbGZyaWVkLmtsYWFzQGVhc3kuZGUifQ.DfqG7q0jhutoK_gAE_2UNnlc1b5YxdBiLd2dz5Kf4jmPSrr11oAWzcp8gEvO1fTX6f3f_iGv5v6Ds4ULzjdEL6pLzZvsR5eWSzaanHnB8BNnXF_hIKn8Ac2HI7nqCFfkd38ENt0n6gM2y2Jh0u-WkLKonQYjcxR8co-KIZwsrTdNmydC6njFLbawGSKwk2RKGgrku4LDIjrZJG4LDj_BucWF79GxQHYSyW91d-oKI4ckoUi6hP4mbRFtl7DQGYns_UPHj1Hz5wjMm_JbmOLiG1hsxhCRSzPKofs-G93j6SYLa5pDBNSOf3Erv3to-k6gMjhn9BZfP8JvuSl9Q2LigQ"
const testTokenSignature = "DfqG7q0jhutoK_gAE_2UNnlc1b5YxdBiLd2dz5Kf4jmPSrr11oAWzcp8gEvO1fTX6f3f_iGv5v6Ds4ULzjdEL6pLzZvsR5eWSzaanHnB8BNnXF_hIKn8Ac2HI7nqCFfkd38ENt0n6gM2y2Jh0u-WkLKonQYjcxR8co-KIZwsrTdNmydC6njFLbawGSKwk2RKGgrku4LDIjrZJG4LDj_BucWF79GxQHYSyW91d-oKI4ckoUi6hP4mbRFtl7DQGYns_UPHj1Hz5wjMm_JbmOLiG1hsxhCRSzPKofs-G93j6SYLa5pDBNSOf3Erv3to-k6gMjhn9BZfP8JvuSl9Q2LigQ"

func TestToken(t *testing.T) {
	ast := assert.New(t)
	jwt, err := DecodeJWT(testToken)
	if err != nil {
		t.Logf("error occured: %v", err)
		t.Fail()
	}
	ast.NotNil(jwt)
	ast.NotNil(jwt.Token)

	header := jwt.Header
	ast.NotNil(header)
	ast.Equal("RS256", header["alg"])
	ast.Equal("JWT", header["typ"])

	payload := jwt.Payload
	ast.NotNil(payload)
	ast.Equal("83e94672-94f8-4760-a63f-ce0f069a1351", payload["sub"])
	ast.Equal("Wilfriedd Klaas", payload["name"])
	ast.Equal(float64(1619853483), payload["iat"])

	sig := jwt.Signature
	ast.NotNil(sig)
	ast.Equal(testTokenSignature, sig)
}

func TestRegex(t *testing.T) {
	resourceName := "/a!b\"c§d$e%f&g/"
	m := regexp.MustCompile("[^a-zA-Z_]")
	resourceName = m.ReplaceAllString(resourceName, "_")
	fmt.Println(resourceName)
}
