package dag

type Route struct {
	route    []string
	routeMap map[string]bool
}

func (r *Route) Push(node string) {
	r.route = append(r.route, node)
	if r.routeMap == nil {
		r.routeMap = make(map[string]bool)
	}
	r.routeMap[node] = true
}

func (r *Route) Pop() string {
	if len(r.route) == 0 {
		return ""
	}
	node := r.route[len(r.route)-1]
	delete(r.routeMap, node)
	r.route = r.route[:len(r.route)-1]
	return node
}

func (r *Route) Contains(node string) bool {
	return r.routeMap[node]
}
