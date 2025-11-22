package models

type Team struct {
	ID      string
	Name    string
	Members []TeamMember
}

type TeamMember struct {
	UserID   string
	Username string
	IsActive bool
}
