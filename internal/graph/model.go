package graph

// Snapshot is a deterministic structural view of a repository.
// It contains observed/derived facts and explicitly labeled high-confidence
// inferences. Low-confidence hypotheses belong in a separate layer and must
// never be mixed into hard enforcement by default.
type Snapshot struct {
	Version int      `json:"version"`
	Root    string   `json:"root"`
	Git     GitState `json:"git"`
	Nodes   []Node   `json:"nodes"`
	Edges   []Edge   `json:"edges"`
}

type GitState struct {
	Available    bool     `json:"available"`
	ChangedFiles []string `json:"changed_files,omitempty"`
}

type Confidence string

const (
	ConfidenceObserved     Confidence = "observed"
	ConfidenceDeclared     Confidence = "declared"
	ConfidenceDerived      Confidence = "derived"
	ConfidenceInferredHigh Confidence = "inferred_high"
	ConfidenceCandidate    Confidence = "candidate"
)

type Node struct {
	ID         string            `json:"id"`
	Kind       NodeKind          `json:"kind"`
	Name       string            `json:"name"`
	Path       string            `json:"path,omitempty"`
	Source     string            `json:"source"`
	Confidence Confidence        `json:"confidence"`
	Evidence   []string          `json:"evidence,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Edge struct {
	From       string     `json:"from"`
	To         string     `json:"to"`
	Kind       EdgeKind   `json:"kind"`
	Source     string     `json:"source"`
	Confidence Confidence `json:"confidence"`
	Evidence   []string   `json:"evidence,omitempty"`
}

type NodeKind string

const (
	NodeRepository NodeKind = "repository"
	NodeProject    NodeKind = "project"
	NodeModule     NodeKind = "module"
	NodeFile       NodeKind = "file"
	NodeSymbol     NodeKind = "symbol"
)

type EdgeKind string

const (
	EdgeContains  EdgeKind = "contains"
	EdgeDependsOn EdgeKind = "depends_on"
)
