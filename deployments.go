package now

import (
	"time"
)

// IncompleteDeployment is the contents of a deploy object before upload
type IncompleteDeployment struct {
	ID        string   `json:"deploymentID"`
	URL       string   `json:"url"`
	TotalSize int      `json:"totalSize"`
	Missing   []string `json:"missing"`
	Warnings  []string `json:"warnings"`
}

// Deployment is the contents of a deploy object
type Deployment struct {
	UID            string     `json:"uid"`
	Host           string     `json:"host"`
	State          string     `json:"state"`
	StateTimestamp *time.Time `json:"stateTs,omitempty"`
}

// DeploymentContentType represents a DeploymentContent type string
type DeploymentContentType string

// DeploymentContentTypes
const (
	TypeDir  DeploymentContentType = "directory"
	TypeFile DeploymentContentType = "file"
)

// DeploymentContent represents a file or directory for deploy
type DeploymentContent interface {
	GetType() DeploymentContentType
	GetName() string
}

// DeploymentDir is the contents of a directory object for deploy
type DeploymentDir struct {
	Type     DeploymentContentType `json:"type"`
	Name     string                `json:"name"`
	Children []DeploymentContent   `json:"children"`
}

// GetName implements the DeploymentContent interface
func (d DeploymentDir) GetName() string {
	return d.Name
}

// GetType implements the DeploymentContent interface
func (d DeploymentDir) GetType() DeploymentContentType {
	return d.Type
}

// DeploymentFile is the contents of a file object for deploy
type DeploymentFile struct {
	Type         DeploymentContentType `json:"type"`
	Name         string                `json:"name"`
	UID          string                `json:"uid,omitempty"`
	Scripts      map[string]string     `json:"scripts,omitempty"`
	Dependencies map[string]string     `json:"dependencies,omitempty"`
	Version      string                `json:"version,omitempty"`
	Description  string                `json:"description,omitempty"`
}

// GetName implements the DeploymentContent interface
func (d DeploymentFile) GetName() string {
	return d.Name
}

// GetType implements the DeploymentContent interface
func (d DeploymentFile) GetType() DeploymentContentType {
	return d.Type
}

// Alias represents a deployment alias object
type Alias struct {
	UID     string     `json:"uid,omitempty"`
	OldUID  string     `json:"oldId,omitempty"`
	Alias   string     `json:"alias"`
	Created *time.Time `json:"created,omitempty"`
}
