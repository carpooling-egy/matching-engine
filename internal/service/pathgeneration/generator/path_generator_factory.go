package generator

func NewPathGenerator() PathGenerator {
	pathGeneratorType := getPathGeneratorType()
	switch pathGeneratorType {
	case "insertion":
		return NewInsertionPathGenerator()
	case "random_topological":
		return NewRandomTopologicalGenerator()
	default:
		return NewInsertionPathGenerator()
	}
}
