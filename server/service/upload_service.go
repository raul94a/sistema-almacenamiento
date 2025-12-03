package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/storage-system/server/models"
	repo "github.com/storage-system/server/repositories"
)


type UploadService interface {
	Upload(uploadObject models.UploadObject, data []byte, jwt *string) (string, error)
}


type UploadObjectService struct {
	Repository *repo.ObjectRepository
}

//TODO: Idempotency should be present in HTTP handler
func (u *UploadObjectService) Upload(
	uploadObject models.UploadObject, 
	data []byte, 
	jwt *string,) (string, error){
	
	// Should create and ID for this file
	uuid, err := uuid.NewV6()
	if err != nil {
		// TODO: Remove
		log.Printf("ErrorService - Upload() %s",err.Error())
		return "", nil
	}
	id := uuid.String()
	// TODO: Remove
	log.Printf("Service - CreateObjectId %s",id)
	
	object := uploadObject.MapToObject(id)

	// what fields should we add at this point

	// TODO: Check userowner - if jwt userID is not the same as the one passed. Or maybe we can use the jwt


	// TODO: Disk -  go get github.com/shirou/gopsutil/v3/disk
	//				 https://pkg.go.dev/github.com/shirou/gopsutil/disk
	// TODO: What Strategy to use?
	//	The project can be configured to be used by tenant, user or whatever.
	//  We can use different strategies for storage management.
	//  for example, if a company like dropbox was using this app,
	//. maybe it was interesting for them to have a Disk per person, yes.
	//. So when the request comes here, with the uuid we can know where to store
	//  the data. (a user is associated to a disk usage).

	//  however, if I was a company that just want to storage a documents from a department
	//  and I had one machine with 10 disks, maybe I would choose a strategy where first
	//. a disk is filled, and then we can change to the next one.

	// Also, if I had two disks maybe a roundrobin will be just fine.
	// Maybe, we may open the possibility to delegate to the developer the
	// capability of choosing the disk where store the file. So, maybe an app that
	// is using this backend can compute previously `where` to place the file (in what disk),
	// so this information is forwarded to this system and we don't need to compute it.

	// Nonetheless, this place a big issue. If the storage location is not the same machine
	// where this app is running, the bytes must be forwarded to the storage location. This
	// can be extreme challenging. There're possibilities like stablishing a FTP tunnel

	// 1. Start with LocalDisk + TenantAffinitySelector using consistent hashing
	// 2. Usage tracking in PostgreSQL (tenant_id + disk_id + bytes_used)
	// 3. Configurable number of replicas (1 or 2)
	// 4. Add admin API later to repin tenants if needed
	// 5. Add "fill_sequential" as alternative mode for "appliance" deployments
	// 4. TenantAffinitySelector 
	//    - Each tenant gets pinned to one or more disks
	//    - Uses consistent hashing or simple map[tenantID -> diskID[]]
	/*type TenantAffinitySelector struct {
		tenantToDisks sync.Map // tenantID -> []string
	}*/

	object.Disk = ""
	// TODO: Location - Should we use uuid + filename?
	// TODO: Location - What should we do if there's a clash of uuid + filename (highly unlikely)
	object.Location = ""
	
	// Deleted to false
	object.Deleted = false


	err = u.Repository.UploadFile(&object)

	if err != nil {
		return "", err
	}
	
	return id,nil

}