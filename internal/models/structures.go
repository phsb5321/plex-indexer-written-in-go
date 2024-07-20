package models

type File struct {
	Name     string `json:"name"`
	FullPath string `json:"fullPath"`
}

type Folder struct {
	Name       string   `json:"name"`
	FullPath   string   `json:"fullPath"`
	Files      []File   `json:"files,omitempty"`
	SubFolders []Folder `json:"subFolders,omitempty"`
}

type Episode struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

type Season struct {
	Name     string    `json:"name"`
	Episodes []Episode `json:"episodes"`
}

type SeriesStructure struct {
	Seasons []Season `json:"seasons"`
}
