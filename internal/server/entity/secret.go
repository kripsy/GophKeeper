package entity

type Secret struct {
	ID          int    `json:"id,omitempty"`
	Type        string `json:"type" validate:"required,oneof=text binary card login_password"`
	Data        []byte `json:"data" validate:"required"`
	Meta        string `json:"meta,omitempty"`
	ChunkNum    int    `json:"chunk_num,omitempty"`
	TotalChunks int    `json:"total_chunks,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
}
