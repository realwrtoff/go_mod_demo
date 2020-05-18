package cache

import (
	"fmt"
	"sync"
	"testing"
)

type Person struct {
	Name string
	Age int
	Desc string
}


func TestNewMemKv(t *testing.T) {
	m := NewMemKv()
	p := &Person{
		Name: "jim",
		Age: 34,
		Desc: "666",
	}
	m.Set(p.Name, p)

	var wg = sync.WaitGroup{}
	for i:=0;i<10000;i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, kv *MemKv, i int) {
			defer wg.Done()
			v, _ := kv.Get("jim")
			p1 := v.(*Person)
			p1.Age = i
			p1.Desc = fmt.Sprintf("write %d", p1.Age)
		}(&wg, m, i)
	}

	for i:=0;i<10;i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, kv *MemKv, i int) {
			defer wg.Done()
			v, _ := kv.Get("jim")
			p2 := v.(*Person)
			t.Logf("jim age %d, desc %s", p2.Age, p2.Desc)
		}(&wg, m, i)
	}

	v, _ := m.Get("jim")
	pe := v.(*Person)
	t.Logf("end jim age %d, desc %s", pe.Age, pe.Desc)
}