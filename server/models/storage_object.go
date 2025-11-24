package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Object representa la entidad de la base de datos mostrada en el diagrama.
type Object struct {
	// PK: id: uuid-v6
	// Usamos la librería google/uuid. Se generará en el hook BeforeCreate.
	ID uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`

	// Campos de Texto Básicos
	Description string `gorm:"type:text" json:"description"`
	Disk        string `gorm:"type:varchar(255)" json:"disk"`
	Location    string `gorm:"type:varchar(255)" json:"location"`

	// Foreign Keys / Campos Opcionales (Text?)
	// Usamos punteros (*string) para permitir valores NULL en la BBDD
	Bucket *string `gorm:"type:varchar(255);index" json:"bucket,omitempty"` // FK marcada en diagrama
	Region *string `gorm:"type:varchar(100)" json:"region,omitempty"`
	Parent *string `gorm:"type:varchar(255);index" json:"parent,omitempty"`

	// Encriptación y Seguridad
	Encrypted        bool    `gorm:"not null;default:false" json:"encrypted"`
	EncryptionMethod *string `gorm:"type:varchar(100)" json:"encryption_method,omitempty"`
	Public           bool    `gorm:"not null;default:false" json:"public"`

	// Detalles del Archivo
	Filename  string  `gorm:"type:varchar(255)" json:"filename"`
	Extension string  `gorm:"type:varchar(50)" json:"extension"`
	Hash      string  `gorm:"type:varchar(255)" json:"hash"` // Útil para verificar integridad
	Size      float64 `gorm:"type:float" json:"size"`        // float64 según diagrama
	Unit      string  `gorm:"type:varchar(20)" json:"unit"`

	// Control de Usuario (FK)
	UserOwner string `gorm:"column:user_owner;type:varchar(255);index" json:"user_owner"`

	// Control de Borrado Lógico (Doble verificación según diagrama)
	Deleted bool           `gorm:"not null;default:false" json:"deleted"`
	
	// Timestamps y Soft Delete nativo de GORM
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName permite sobrescribir el nombre de la tabla si no quieres que sea "objects"
func (Object) TableName() string {
	return "objects"
}

// BeforeCreate es un hook de GORM para asegurar que el UUID sea v6 antes de insertar.
func (o *Object) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		// Generamos UUID v6 (Time-ordered), ideal para claves primarias de BBDD
		// Nota: Requiere "github.com/google/uuid" v1.6.0+
		o.ID, err = uuid.NewV6()
		if err != nil {
			return err
		}
	}
	return
}