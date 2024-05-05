package misc

import sr "linear-db/pkg/structure"

func DatabaseExists(dbs *sr.Databases, db string) bool {
    for _, d := range dbs.Databases {
        if d.Name == db {
            return true
        }
    }
    return false
}

func IndexOf(e string, a *sr.Databases) int{
	for i, v := range(a.Databases){
		if e == v.Name {
			return i
		}
	}
	return -1
}