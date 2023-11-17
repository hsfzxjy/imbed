package lazy

import (
	"sync"

	"github.com/hsfzxjy/imbed/core/ref"
)

type Object interface {
	GetData() ([]byte, error)
	GetSHA() (ref.Sha256, error)
}

type MustSHAObject interface {
	Object
	MustSHA() ref.Sha256
}

type MustDataSHAObject interface {
	Object
	MustData() []byte
	MustSHAObject
}

type shaCache struct {
	once sync.Once
	sha  ref.Sha256
}

type ConstData []byte
type constData = ConstData

func (d ConstData) GetData() ([]byte, error) {
	return d, nil
}

func (d ConstData) MustData() []byte {
	return d
}

type ConstSHA ref.Sha256
type constSHA = ConstSHA

func (s ConstSHA) GetSHA() (ref.Sha256, error) {
	return ref.Sha256(s), nil
}

func (s ConstSHA) MustSHA() ref.Sha256 {
	return ref.Sha256(s)
}

type SHAF struct {
	Fn func() ref.Sha256
	shaCache
}
type shaF = SHAF

func (s *SHAF) GetSHA() (ref.Sha256, error) {
	return s.MustSHA(), nil
}

func (s *SHAF) MustSHA() ref.Sha256 {
	s.once.Do(func() {
		s.sha = s.Fn()
	})
	return s.sha
}

type dataF struct {
	once sync.Once
	p    []byte
	fn   func() []byte
}

func (d *dataF) GetData() ([]byte, error) {
	return d.MustData(), nil
}

func (d *dataF) MustData() []byte {
	d.once.Do(func() {
		d.p = d.fn()
	})
	return d.p
}

type data struct {
	constData
	shaCache
}

func (d *data) GetSHA() (ref.Sha256, error) {
	return d.MustSHA(), nil
}

func (d *data) MustSHA() ref.Sha256 {
	d.once.Do(func() {
		d.sha = ref.Sha256HashSum(d.constData)
	})
	return d.sha
}

func Data(p []byte) *data {
	return &data{constData: p}
}

func _() { var _ MustDataSHAObject = &data{} }

type dataSHA struct {
	constData
	constSHA
}

func DataSHA(p []byte, sha ref.Sha256) *dataSHA {
	return &dataSHA{constData: p, constSHA: constSHA(sha)}
}

func _() { var _ MustDataSHAObject = &dataSHA{} }

type dataSHAFunc struct {
	constData
	shaF
}

func DataSHAFunc(p []byte, sha func() ref.Sha256) *dataSHAFunc {
	return &dataSHAFunc{constData: p, shaF: shaF{Fn: sha}}
}

func _() { var _ MustDataSHAObject = &dataSHAFunc{} }

type dataFunc struct {
	once     sync.Once
	p        []byte
	sha      ref.Sha256
	dataFunc func() []byte
}

func (d *dataFunc) MustData() []byte {
	d.compute()
	return d.p
}

func (d *dataFunc) MustSHA() ref.Sha256 {
	d.compute()
	return d.sha
}

func (d *dataFunc) compute() {
	d.once.Do(func() {
		d.p = d.dataFunc()
		d.sha = ref.Sha256HashSum(d.p)
	})
}

func (d *dataFunc) GetData() ([]byte, error) {
	d.compute()
	return d.p, nil
}

func (d *dataFunc) GetSHA() (ref.Sha256, error) {
	d.compute()
	return d.sha, nil
}

func DataFunc(data func() []byte) *dataFunc {
	return &dataFunc{dataFunc: data}
}

func _() { var _ MustDataSHAObject = &dataFunc{} }

type dataFuncSHA struct {
	dataF
	constSHA
}

func DataFuncSHA(data func() []byte, sha ref.Sha256) *dataFuncSHA {
	return &dataFuncSHA{dataF: dataF{fn: data}, constSHA: constSHA(sha)}
}

func _() { var _ MustDataSHAObject = &dataFuncSHA{} }

type dataFuncSHAFunc struct {
	dataF
	shaF
}

func DataFuncSHAFunc(data func() []byte, sha func() ref.Sha256) *dataFuncSHAFunc {
	return &dataFuncSHAFunc{dataF: dataF{fn: data}, shaF: shaF{Fn: sha}}
}

func _() { var _ MustDataSHAObject = &dataFuncSHAFunc{} }

type dataFunc2 struct {
	once     sync.Once
	p        []byte
	sha      ref.Sha256
	err      error
	dataFunc func() ([]byte, error)
}

func (d *dataFunc2) compute() {
	d.once.Do(func() {
		p, err := d.dataFunc()
		if err != nil {
			d.err = err
			return
		}
		d.p = p
		d.sha = ref.Sha256HashSum(p)
	})
}

func (d *dataFunc2) GetData() ([]byte, error) {
	d.compute()
	return d.p, d.err
}

func (d *dataFunc2) GetSHA() (ref.Sha256, error) {
	d.compute()
	return d.sha, d.err
}

func DataFunc2(data func() ([]byte, error)) *dataFunc2 {
	return &dataFunc2{dataFunc: data}
}

func _() { var _ Object = &dataFunc2{} }

type dataF2 struct {
	once sync.Once
	fn   func() ([]byte, error)
	p    []byte
	err  error
}

func (d *dataF2) GetData() ([]byte, error) {
	d.once.Do(func() {
		d.p, d.err = d.fn()
	})
	return d.p, d.err
}

type dataFunc2SHA struct {
	dataF2
	constSHA
}

func DataFunc2SHA(data func() ([]byte, error), sha ref.Sha256) *dataFunc2SHA {
	return &dataFunc2SHA{dataF2: dataF2{fn: data}, constSHA: constSHA(sha)}
}

func _() { var _ MustSHAObject = &dataFunc2SHA{} }

type dataFunc2SHAFunc struct {
	dataF2
	shaF
}

func DataFunc2SHAFunc(data func() ([]byte, error), sha func() ref.Sha256) *dataFunc2SHAFunc {
	return &dataFunc2SHAFunc{dataF2: dataF2{fn: data}, shaF: shaF{Fn: sha}}
}

func _() { var _ MustSHAObject = &dataFunc2SHAFunc{} }
