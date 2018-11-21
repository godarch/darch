package staging

// ByAge implements sort.Interface for []Person based on
// the Age field.
type sortStageImageNamed []StagedImageNamed

func (a sortStageImageNamed) Len() int           { return len(a) }
func (a sortStageImageNamed) Less(i, j int) bool { return a[i].Ref.FullName() < a[j].Ref.FullName() }
func (a sortStageImageNamed) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
