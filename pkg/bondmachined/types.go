package bondmachined

type PullArtifactRequest struct {
	ArtifactName string `json:"artifact_name"`
	BoardModel   string `json:"board_model"`
	Force        bool   `json:"force,omitempty"`
}
