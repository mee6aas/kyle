package assigns

var (
	//                     invkID
	assignments = make(map[string]assign)
)

type assign struct {
	id     string           // ID of the assignment
	holder chan interface{} // channel to pass result
}
