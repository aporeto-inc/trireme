package crypto

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestComputeVerifyHMAC tests the compute and verify of HMAC functions
func TestComputeVerifyHMAC(t *testing.T) {
	Convey("Given a token and a key", t, func() {

		token := make([]byte, 256)
		for i := uint8(0); i < 255; i++ {
			token[i] = byte(i)
		}

		key := make([]byte, 32)
		for i := uint8(0); i < 32; i++ {
			key[i] = i
		}

		Convey("When I sign the token with the key", func() {
			expectedMac, err := ComputeHmac256(token, key)
			So(err, ShouldBeNil)
			So(expectedMac, ShouldNotBeNil)

			Convey("I should be able to verify the token with the same key", func() {
				verified := VerifyHmac(token, expectedMac, key)
				So(verified, ShouldBeTrue)
			})

			Convey("If I provide the worng key, I should fail verification", func() {
				fakeKey := make([]byte, 32)
				verified := VerifyHmac(token, expectedMac, fakeKey)
				So(verified, ShouldBeFalse)
			})

			Convey("If I provide the wright key, but I havae the wrong signature, I should fail verification", func() {
				failedMac := make([]byte, 32)
				verified := VerifyHmac(token, failedMac, key)
				So(verified, ShouldBeFalse)
			})
		})
	})
}

// TestRandomString tests the random string generation function and the random byte generation
func TestRandomString(t *testing.T) {
	Convey("Given a string length of 16", t, func() {
		length := 16
		Convey("I should be able to generate two random strings that are not equal", func() {
			string1, err1 := GenerateRandomString(length)
			string2, err2 := GenerateRandomString(length)
			So(err1, ShouldBeNil)
			So(err2, ShouldBeNil)
			So(string1, ShouldNotEqual, string2)
		})
	})
	Convey("Given a string length of 32", t, func() {
		length := 16
		Convey("I should be able to generate two random strings that are not equal", func() {
			string1, err1 := GenerateRandomString(length)
			string2, err2 := GenerateRandomString(length)
			So(err1, ShouldBeNil)
			So(err2, ShouldBeNil)
			So(string1, ShouldNotEqual, string2)
		})
	})
}

// TestFuncLoadEllipticCurve
func TestFuncLoadEllipticCurve(t *testing.T) {
	Convey("Given a valid EC key", t, func() {
		keyPEM := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPkiHqtH372JJdAG/IxJlE1gv03cdwa8Lhg2b3m/HmbyoAoGCCqGSM49
AwEHoUQDQgAEAfAL+AfPj/DnxrU6tUkEyzEyCxnflOWxhouy1bdzhJ7vxMb1vQ31
8ZbW/WvMN/ojIXqXYrEpISoojznj46w64w==
-----END EC PRIVATE KEY-----`
		Convey("I should be able to load the key", func() {
			key, err := LoadEllipticCurveKey([]byte(keyPEM))
			So(key, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})
	})

	Convey("Given an invalid PEM BLOCK", t, func() {
		keyPEM := ""
		Convey("I should get an error", func() {
			_, err := LoadEllipticCurveKey([]byte(keyPEM))
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given an invalid Key file", t, func() {
		keyPEM := `-----BEGIN EC PRIVATE KEY-----
-----END EC PRIVATE KEY-----`
		Convey("I should get an error", func() {
			_, err := LoadEllipticCurveKey([]byte(keyPEM))
			So(err, ShouldNotBeNil)
		})
	})
}

// TestFuncLoadRootCertificates test the loading of root certs in a cert pool
func TestFuncLoadRootCertificates(t *testing.T) {
	Convey("Given a valid certificate chain", t, func() {
		caPool := `-----BEGIN CERTIFICATE-----
MIIBhTCCASwCCQC8b53yGlcQazAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABJxneTUqhbtgEIwpKUUzwz3h92SqcOdIw3mfQkMjg3Vobvr6JKlpXYe9xhsN
rygJmLhMAN9gjF9qM9ybdbe+m3owCgYIKoZIzj0EAwIDRwAwRAIgC1fVMqdBy/o3
jNUje/Hx0fZF9VDyUK4ld+K/wF3QdK4CID1ONj/Kqinrq2OpjYdkgIjEPuXoOoR1
tCym8dnq4wtH
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIB3jCCAYOgAwIBAgIJALsW7pyC2ERQMAoGCCqGSM49BAMCMEsxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTEMMAoGA1UEBwwDU0pDMRAwDgYDVQQKDAdUcmlyZW1l
MQ8wDQYDVQQDDAZ1YnVudHUwHhcNMTYwOTI3MjI0OTAwWhcNMjYwOTI1MjI0OTAw
WjBLMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4G
A1UECgwHVHJpcmVtZTEPMA0GA1UEAwwGdWJ1bnR1MFkwEwYHKoZIzj0CAQYIKoZI
zj0DAQcDQgAE4c2Fd7XeIB1Vfs51fWwREfLLDa55J+NBalV12CH7YEAnEXjl47aV
cmNqcAtdMUpf2oz9nFVI81bgO+OSudr3CqNQME4wHQYDVR0OBBYEFOBftuI09mmu
rXjqDyIta1gT8lqvMB8GA1UdIwQYMBaAFOBftuI09mmurXjqDyIta1gT8lqvMAwG
A1UdEwQFMAMBAf8wCgYIKoZIzj0EAwIDSQAwRgIhAMylAHhbFA0KqhXIFiXNpEbH
JKaELL6UXXdeQ5yup8q+AiEAh5laB9rbgTymjaANcZ2YzEZH4VFS3CKoSdVqgnwC
dW4=
-----END CERTIFICATE-----`

		Convey("I should be able to get a valid certificate chain", func() {
			roots := LoadRootCertificates([]byte(caPool))
			So(roots, ShouldNotBeNil)
			So(len(roots.Subjects()), ShouldEqual, 2)
		})
	})
}

// TestFuncLoadAndVerifyCertificate
func TestLoadAndVerifyCertificate(t *testing.T) {
	Convey("Given a valid certificate chain", t, func() {
		caPool := `-----BEGIN CERTIFICATE-----
MIIBhTCCASwCCQC8b53yGlcQazAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABJxneTUqhbtgEIwpKUUzwz3h92SqcOdIw3mfQkMjg3Vobvr6JKlpXYe9xhsN
rygJmLhMAN9gjF9qM9ybdbe+m3owCgYIKoZIzj0EAwIDRwAwRAIgC1fVMqdBy/o3
jNUje/Hx0fZF9VDyUK4ld+K/wF3QdK4CID1ONj/Kqinrq2OpjYdkgIjEPuXoOoR1
tCym8dnq4wtH
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIB3jCCAYOgAwIBAgIJALsW7pyC2ERQMAoGCCqGSM49BAMCMEsxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTEMMAoGA1UEBwwDU0pDMRAwDgYDVQQKDAdUcmlyZW1l
MQ8wDQYDVQQDDAZ1YnVudHUwHhcNMTYwOTI3MjI0OTAwWhcNMjYwOTI1MjI0OTAw
WjBLMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4G
A1UECgwHVHJpcmVtZTEPMA0GA1UEAwwGdWJ1bnR1MFkwEwYHKoZIzj0CAQYIKoZI
zj0DAQcDQgAE4c2Fd7XeIB1Vfs51fWwREfLLDa55J+NBalV12CH7YEAnEXjl47aV
cmNqcAtdMUpf2oz9nFVI81bgO+OSudr3CqNQME4wHQYDVR0OBBYEFOBftuI09mmu
rXjqDyIta1gT8lqvMB8GA1UdIwQYMBaAFOBftuI09mmurXjqDyIta1gT8lqvMAwG
A1UdEwQFMAMBAf8wCgYIKoZIzj0EAwIDSQAwRgIhAMylAHhbFA0KqhXIFiXNpEbH
JKaELL6UXXdeQ5yup8q+AiEAh5laB9rbgTymjaANcZ2YzEZH4VFS3CKoSdVqgnwC
dW4=
-----END CERTIFICATE-----`
		roots := LoadRootCertificates([]byte(caPool))
		So(roots, ShouldNotBeNil)
		So(len(roots.Subjects()), ShouldEqual, 2)

		Convey("Given a certificate signed by the intermediatery", func() {
			certPEM := `-----BEGIN CERTIFICATE-----
MIIBhjCCASwCCQCPCdgp39gHJTAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABAHwC/gHz4/w58a1OrVJBMsxMgsZ35TlsYaLstW3c4Se78TG9b0N9fGW1v1r
zDf6IyF6l2KxKSEqKI854+OsOuMwCgYIKoZIzj0EAwIDSAAwRQIgQwQn0jnK/XvD
KxgQd/0pW5FOAaB41cMcw4/XVlphO1oCIQDlGie+WlOMjCzrV0Xz+XqIIi1pIgPT
IG7Nv+YlTVp5qA==
-----END CERTIFICATE-----`
			Convey("I should be able to load and verify the certificate", func() {
				cert, err := LoadAndVerifyCertificate([]byte(certPEM), roots)
				So(cert, ShouldNotBeNil)
				So(err, ShouldBeNil)
			})
		})
	})

	Convey("Given the root CA certificate only ", t, func() {
		caPool := `
-----BEGIN CERTIFICATE-----
MIIB3jCCAYOgAwIBAgIJALsW7pyC2ERQMAoGCCqGSM49BAMCMEsxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTEMMAoGA1UEBwwDU0pDMRAwDgYDVQQKDAdUcmlyZW1l
MQ8wDQYDVQQDDAZ1YnVudHUwHhcNMTYwOTI3MjI0OTAwWhcNMjYwOTI1MjI0OTAw
WjBLMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4G
A1UECgwHVHJpcmVtZTEPMA0GA1UEAwwGdWJ1bnR1MFkwEwYHKoZIzj0CAQYIKoZI
zj0DAQcDQgAE4c2Fd7XeIB1Vfs51fWwREfLLDa55J+NBalV12CH7YEAnEXjl47aV
cmNqcAtdMUpf2oz9nFVI81bgO+OSudr3CqNQME4wHQYDVR0OBBYEFOBftuI09mmu
rXjqDyIta1gT8lqvMB8GA1UdIwQYMBaAFOBftuI09mmurXjqDyIta1gT8lqvMAwG
A1UdEwQFMAMBAf8wCgYIKoZIzj0EAwIDSQAwRgIhAMylAHhbFA0KqhXIFiXNpEbH
JKaELL6UXXdeQ5yup8q+AiEAh5laB9rbgTymjaANcZ2YzEZH4VFS3CKoSdVqgnwC
dW4=
-----END CERTIFICATE-----`
		roots := LoadRootCertificates([]byte(caPool))
		So(roots, ShouldNotBeNil)
		So(len(roots.Subjects()), ShouldEqual, 1)

		Convey("Given a certificate signed by the intermediatery", func() {
			certPEM := `-----BEGIN CERTIFICATE-----
MIIBhjCCASwCCQCPCdgp39gHJTAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABAHwC/gHz4/w58a1OrVJBMsxMgsZ35TlsYaLstW3c4Se78TG9b0N9fGW1v1r
zDf6IyF6l2KxKSEqKI854+OsOuMwCgYIKoZIzj0EAwIDSAAwRQIgQwQn0jnK/XvD
KxgQd/0pW5FOAaB41cMcw4/XVlphO1oCIQDlGie+WlOMjCzrV0Xz+XqIIi1pIgPT
IG7Nv+YlTVp5qA==
-----END CERTIFICATE-----`
			Convey("I should be able to fail to verify the certificate ", func() {
				cert, err := LoadAndVerifyCertificate([]byte(certPEM), roots)
				So(cert, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given a good CA ", t, func() {
		caPool := `-----BEGIN CERTIFICATE-----
MIIBhTCCASwCCQC8b53yGlcQazAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABJxneTUqhbtgEIwpKUUzwz3h92SqcOdIw3mfQkMjg3Vobvr6JKlpXYe9xhsN
rygJmLhMAN9gjF9qM9ybdbe+m3owCgYIKoZIzj0EAwIDRwAwRAIgC1fVMqdBy/o3
jNUje/Hx0fZF9VDyUK4ld+K/wF3QdK4CID1ONj/Kqinrq2OpjYdkgIjEPuXoOoR1
tCym8dnq4wtH
-----END CERTIFICATE-----`

		roots := LoadRootCertificates([]byte(caPool))
		So(roots, ShouldNotBeNil)
		So(len(roots.Subjects()), ShouldEqual, 1)

		Convey("Given a bad certificate ", func() {
			certPEM := `-----BEGIN CERTIFICATE-----
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABAHwC/gHz4/w58a1OrVJBMsxMgsZ35TlsYaLstW3c4Se78TG9b0N9fGW1v1r
zDf6IyF6l2KxKSEqKI854+OsOuMwCgYIKoZIzj0EAwIDSAAwRQIgQwQn0jnK/XvD
KxgQd/0pW5FOAaB41cMcw4/XVlphO1oCIQDlGie+WlOMjCzrV0Xz+XqIIi1pIgPT
IG7Nv+YlTVp5qA==
-----END CERTIFICATE-----`
			Convey("I should be able to fail to verify the certificate ", func() {
				cert, err := LoadAndVerifyCertificate([]byte(certPEM), roots)
				So(cert, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})

	Convey("Given bad certificate   ", t, func() {

		Convey("Where the certificate block is bad  ", func() {
			certPEM := ``
			Convey("I should be able to fail to verify the certificate ", func() {
				cert, err := LoadAndVerifyCertificate([]byte(certPEM), nil)
				So(cert, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

// TestLoadAndVerifyECSecrets
func TestLoadAndVerifyECSecrets(t *testing.T) {

	Convey("Given a valid EC key", t, func() {
		keyPEM := `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIPkiHqtH372JJdAG/IxJlE1gv03cdwa8Lhg2b3m/HmbyoAoGCCqGSM49
AwEHoUQDQgAEAfAL+AfPj/DnxrU6tUkEyzEyCxnflOWxhouy1bdzhJ7vxMb1vQ31
8ZbW/WvMN/ojIXqXYrEpISoojznj46w64w==
-----END EC PRIVATE KEY-----`
		Convey("I should be able to load the key", func() {
			key, err := LoadEllipticCurveKey([]byte(keyPEM))
			So(key, ShouldNotBeNil)
			So(err, ShouldBeNil)
		})

		Convey("Given a valid certificate chain", func() {
			caPool := `-----BEGIN CERTIFICATE-----
MIIBhTCCASwCCQC8b53yGlcQazAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABJxneTUqhbtgEIwpKUUzwz3h92SqcOdIw3mfQkMjg3Vobvr6JKlpXYe9xhsN
rygJmLhMAN9gjF9qM9ybdbe+m3owCgYIKoZIzj0EAwIDRwAwRAIgC1fVMqdBy/o3
jNUje/Hx0fZF9VDyUK4ld+K/wF3QdK4CID1ONj/Kqinrq2OpjYdkgIjEPuXoOoR1
tCym8dnq4wtH
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIB3jCCAYOgAwIBAgIJALsW7pyC2ERQMAoGCCqGSM49BAMCMEsxCzAJBgNVBAYT
AlVTMQswCQYDVQQIDAJDQTEMMAoGA1UEBwwDU0pDMRAwDgYDVQQKDAdUcmlyZW1l
MQ8wDQYDVQQDDAZ1YnVudHUwHhcNMTYwOTI3MjI0OTAwWhcNMjYwOTI1MjI0OTAw
WjBLMQswCQYDVQQGEwJVUzELMAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4G
A1UECgwHVHJpcmVtZTEPMA0GA1UEAwwGdWJ1bnR1MFkwEwYHKoZIzj0CAQYIKoZI
zj0DAQcDQgAE4c2Fd7XeIB1Vfs51fWwREfLLDa55J+NBalV12CH7YEAnEXjl47aV
cmNqcAtdMUpf2oz9nFVI81bgO+OSudr3CqNQME4wHQYDVR0OBBYEFOBftuI09mmu
rXjqDyIta1gT8lqvMB8GA1UdIwQYMBaAFOBftuI09mmurXjqDyIta1gT8lqvMAwG
A1UdEwQFMAMBAf8wCgYIKoZIzj0EAwIDSQAwRgIhAMylAHhbFA0KqhXIFiXNpEbH
JKaELL6UXXdeQ5yup8q+AiEAh5laB9rbgTymjaANcZ2YzEZH4VFS3CKoSdVqgnwC
dW4=
-----END CERTIFICATE-----`
			roots := LoadRootCertificates([]byte(caPool))
			So(roots, ShouldNotBeNil)
			So(len(roots.Subjects()), ShouldEqual, 2)

			Convey("Given a certificate signed by the intermediatery", func() {
				certPEM := `-----BEGIN CERTIFICATE-----
MIIBhjCCASwCCQCPCdgp39gHJTAKBggqhkjOPQQDAjBLMQswCQYDVQQGEwJVUzEL
MAkGA1UECAwCQ0ExDDAKBgNVBAcMA1NKQzEQMA4GA1UECgwHVHJpcmVtZTEPMA0G
A1UEAwwGdWJ1bnR1MB4XDTE2MDkyNzIyNDkwMFoXDTI2MDkyNTIyNDkwMFowSzEL
MAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMQwwCgYDVQQHDANTSkMxEDAOBgNVBAoM
B1RyaXJlbWUxDzANBgNVBAMMBnVidW50dTBZMBMGByqGSM49AgEGCCqGSM49AwEH
A0IABAHwC/gHz4/w58a1OrVJBMsxMgsZ35TlsYaLstW3c4Se78TG9b0N9fGW1v1r
zDf6IyF6l2KxKSEqKI854+OsOuMwCgYIKoZIzj0EAwIDSAAwRQIgQwQn0jnK/XvD
KxgQd/0pW5FOAaB41cMcw4/XVlphO1oCIQDlGie+WlOMjCzrV0Xz+XqIIi1pIgPT
IG7Nv+YlTVp5qA==
-----END CERTIFICATE-----`
				Convey("I should be able to load and verify the certificate", func() {
					cert, err := LoadAndVerifyCertificate([]byte(certPEM), roots)
					So(cert, ShouldNotBeNil)
					So(err, ShouldBeNil)
				})

				Convey("Given I have valid EC key, certificate chain and signed certificate", func() {
					key, cert, certPool, err := LoadAndVerifyECSecrets([]byte(keyPEM), []byte(certPEM), []byte(caPool))

					Convey("I should be able to load and verify all the certificates and keys in the right data structures", func() {
						So(key, ShouldNotBeNil)
						So(cert, ShouldNotBeNil)
						So(certPool, ShouldNotBeNil)
						So(err, ShouldBeNil)
					})
				})

				Convey("Given I have invalid EC key, valid certificate chain and signed certificate", func() {
					invalidKeyPEM := `-----BEGIN EC PRIVATE KEY-----
			-----END EC PRIVATE KEY-----`
					key, cert, certPool, err := LoadAndVerifyECSecrets([]byte(invalidKeyPEM), []byte(certPEM), []byte(caPool))

					Convey("I should be able to fail verifying the EC key", func() {
						So(key, ShouldBeNil)
						So(cert, ShouldBeNil)
						So(certPool, ShouldBeNil)
						So(err, ShouldResemble, fmt.Errorf("unable to parse pem block: -----BEGIN EC PRIVATE KEY-----\n\t\t\t-----END EC PRIVATE KEY-----"))
					})
				})

				Convey("Given I have valid EC key, invalid certificate chain and valid signed certificate", func() {
					invalidCaPool := `-----BEGIN CERTIFICATE-----
			-----END CERTIFICATE-----`
					key, cert, certPool, err := LoadAndVerifyECSecrets([]byte(keyPEM), []byte(certPEM), []byte(invalidCaPool))

					Convey("I should be able to fail loading the certificate pool", func() {
						So(key, ShouldBeNil)
						So(cert, ShouldBeNil)
						So(certPool, ShouldBeNil)
						So(err, ShouldResemble, fmt.Errorf("unable to load root certificate pool"))
					})
				})

				Convey("Given I have valid EC key, certificate chain and invalid signed certificate", func() {
					invalidCertPEM := `-----BEGIN CERTIFICATE-----
		-----END CERTIFICATE-----`
					key, cert, certPool, err := LoadAndVerifyECSecrets([]byte(keyPEM), []byte(invalidCertPEM), []byte(caPool))

					Convey("I should be able to fail verifying certificate (bad certificate)", func() {
						So(key, ShouldBeNil)
						So(cert, ShouldBeNil)
						So(certPool, ShouldBeNil)
						So(err, ShouldResemble, fmt.Errorf("unable to decode pem block: -----BEGIN CERTIFICATE-----\n\t\t-----END CERTIFICATE-----"))
					})
				})
			})
		})
	})
}
