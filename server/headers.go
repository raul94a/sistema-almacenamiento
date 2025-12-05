package server


// Headers Especiales usadas en la App
// 
// El objetivo es securizar y autenticar las peticiones:

// a) Quién hace el request -Usuario
//	  Se puede hacer a través de un token JWT que pertenezca a la app principal
//    No obstante, para hacerlo funcionar se necesita que ese backend implemente un JWK
//	  ¿Puede este usuario acceder al recurso?
//	   X-Access-User

// b) Desde donde hace el request -Dominio
//    Para autorización de dominios. ¿Está el dominio incluido en el whitelist?
//	   X-Site

// c) ¿Estoy autorizado a pedir recursos a este servidor?
//.   Puede ser que un usuario de otra aplicación que esté permitido y desde
//.   un dominio permitido (apartados a y b cumplidos) sea hackeado. Al final
//.   los headers también podrían hackearse. Solo son strings.
//.   En tal caso, ¿cómo poner una capa más de seguridad? Con un token especial
//.   creado desde la propia app, un JWK, a través de una clave publica.
//.   Esta clave debe ser compartida con el servidor y ser utilizada para firmar un
//.   secreto almacenado en el servidor desde el que se hace la petición.
//.   la clave privada de este sistema puede validar la firma de dicho secreto.
//.   X-Signature
//.
//.   webservices necesarios:

//.   PUBLIC ACCESS
//.   * GET /security/public-key
//		Retorna la clave pública del servidor de almacenamiento de ficheros 
//
//.   RESTRINGED ACCESS: SAME-ORIGIN
// 	  * POST /security/asymmetric/create
//		check: https://stackoverflow.com/questions/5244129/use-rsa-private-key-to-generate-public-key
//	    Metadata
//		 - Headers + CORS
//		 - Body: {
//		          	"secret": "MY_SECRET",
//					"size": 1024/2048/4096
//				 }
//.
//. 
//
//
//
//