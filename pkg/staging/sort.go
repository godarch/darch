package staging

// sortStageImageNamedByName implements sort.Interface for []StagedImageNamed
// based on the FullName field in an ascending order.
type sortStagedImageNamedByName []StagedImageNamed

func (a sortStagedImageNamedByName) Len() int { return len(a) }
func (a sortStagedImageNamedByName) Less(i, j int) bool {
	return a[i].Ref.FullName() < a[j].Ref.FullName()
}
func (a sortStagedImageNamedByName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// sortStageImageNamedByNameDesc implements sort.Interface for []StagedImageNamed
// based on the FullName field in an descending order.
type sortStagedImageNamedByNameDesc []StagedImageNamed

func (a sortStagedImageNamedByNameDesc) Len() int { return len(a) }
func (a sortStagedImageNamedByNameDesc) Less(i, j int) bool {
	return a[i].Ref.FullName() > a[j].Ref.FullName()
}
func (a sortStagedImageNamedByNameDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// sortStagedImageNamedByAge implements sort.Interface for []StagedImageNamed
// based on the CreationTime field in an ascending order.
type sortStagedImageNamedByAge []StagedImageNamed

func (a sortStagedImageNamedByAge) Len() int { return len(a) }
func (a sortStagedImageNamedByAge) Less(i, j int) bool {
	return a[i].CreationTime.Before(a[j].CreationTime)
}
func (a sortStagedImageNamedByAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// sortStagedImageNamedByAgeDesc implements sort.Interface for []StagedImageNamed
// based on the CreationTime field in an descending order.
type sortStagedImageNamedByAgeDesc []StagedImageNamed

func (a sortStagedImageNamedByAgeDesc) Len() int { return len(a) }
func (a sortStagedImageNamedByAgeDesc) Less(i, j int) bool {
	return a[i].CreationTime.After(a[j].CreationTime)
}
func (a sortStagedImageNamedByAgeDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
