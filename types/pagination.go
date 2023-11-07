package types

type PaginationQuery struct {
	Limit  uint
	Offset uint
}

type PaginationReturn struct {
	TotalResults uint
	CurrentPage  uint
	NextPage     uint
	TotalPages   uint
}
