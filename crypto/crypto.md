# Crypto proposal

For keeping the app secure, we're providing 3 layers of security

A) JWT Authorization. 

B) Hybrid RSA/AES Encryption

C) Domain whitelist

-----
A) JWT Authorization is not directly implemented by the storage system.
    it relies on a third-party that have to expose its JWKS.
    With these JWKS the JWT Token can be decyphered and used to
    Authenticate the user. Special fields will be allowed to be included
    into the JWT through an admin user interface.

B) Hybrid RSA/AES Encryption. Provided by the Storage System. This will
    be used as a method to verify the requests. The final output of
    this technique will be a encrypted string, the `Signature`, included
    in the headers. The secrets and keys can be created/modified through the
    admin user interface. Also, some endpoints will be provided for automation
    of key rotations.

C) Domain whitelist. Registration of the domains that can talk to our system.
    If the user lack of firewall knowledge this will be pretty useful.
    The Storage system is able to accept domains that are allowed to talk
    to the Storage System.

All of the above security layers can be used in any combination of them. Also,
any of these layers can be disabled through the admin user interface.

Needs to this system can be successful:

a.SDKs
    We do need to build SDKs that manages the request to the storage system.
    These will manage all the complex encryption needed for the Signature.

The endpoints will be available with simple methods.


# Kotlin


```kotlin
// @param key: Storage System Key
// @param secret: Storage System Secret
data class StorageSystemConfig(private val key: String? = null, private val secret: String? = null, private val url: String) {
    
}
public class StorageSystemClient(val config: StorageSystemConfig) {

    private val headerBuilder = HeaderBuilder()

    private val remoteConfig = RemoteConfig()

    init {
        remoteConfig.fetchPublicKey(
            config.url + endpoint.public_key
        )
    }


}  

private class HeaderBuilder {

}
private class RemoteConfig {
    private var pk: String = ""
    fun FetchPublicKey(url: String){
        // okHttp or Ktor
       this.pk =  http.request('GET', 'URL') 
    }

    fun getPk(): String {
        return this.pk
    }
}
```